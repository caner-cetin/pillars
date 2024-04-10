<template>
  <q-layout view="hHh Lpr fFf">
    <q-page-container>
      <q-page>
        <div class="row" style="height: 52.5vh;">
          <GoogleMap style="width: 100%; height: 100%" :center="center" :zoom="3">
            <Marker v-for="(location, index) in positions" :key="index"
              :options="{ position: { lat: location.lat, lng: location.lng } }" @click="showMarkerDetails(location)" />
          </GoogleMap>
        </div>
        <div class="row" style="height: 47.5vh; display: flex;">
          <div class="column" style="width: 30%; height: 100%; ">
            <div class="row justify-center items-center" style="width: 100%; height: 10%">
              <q-btn glossy label="Submit Earthquakes" @click="activateScriptDialog = !activateScriptDialog" />
            </div>
            <div class="row justify-center items-center" style="width: 100%; height: 10%; padding-top: 10px;">
              <q-btn glossy label="Wipe Earthquakes"
                @click="activateWipeEarthquakesDialog = !activateWipeEarthquakesDialog" />
            </div>
            <div class="row justify-center items-center" style="width: 100%; height: 70%;">
              <div v-if="currentRowEarthquakes">
                <EarthquakeInfoCard v-show="renderMarker" :earthquake="cachedEarthquakes[selectedMarkerIndex]" />
              </div>
            </div>
            <div class="row justify-center items-center" style="width: 100%; height: 10%;">
              <div v-if="connectedToWs">
                <q-icon name="sym_o_rss_feed" /> Listening for changes...
              </div>
              <div v-else>
                <q-icon name="sym_o_signal_disconnected" style="padding-right: 10px;" />Live updates are disabled.<q-btn
                  push color="red-14" label="Reconnect" @click='reconnectWS' size="sm" />
              </div>
            </div>
          </div>
          <div class="column" style="width: 70%; height: 100%;">
            <div class="text-h6 text-center" style="width: 100%;">
              <q-table title="Latest Earthquakes" :rows="currentRowEarthquakes" :columns="latestEarthquakesColumns"
                binary-state-sort dense style="height: 100%; width: 100%;" :row-key="isUniqueTableKey"
                v-model:pagination="earthquakeTablePagination" @request="onLatestEarthquakeTableScroll"
                :rows-per-page-options="[5]" no-data-label="No earthquakes submitted." :filter="latestEarthquakeFilter"
                @row-click="onLatestEarthquakeRowClick" :loading="latestEarthquakesTableLoading">
                <template v-slot:top-right>
                  <q-checkbox v-model="onlyShowEarthquakesInTheTableAtMap" @click="adjustMarkers()" /> <a
                    class="text-subtitle2 text-weight-light" style="padding-right: 20px;">Only show earthquakes in the
                    table at the
                    map.</a>
                  <SearchBar v-model:search="latestEarthquakeFilter" />
                </template>
              </q-table>
            </div>
          </div>
        </div>
        <q-dialog v-model='activateScriptDialog'>
          <q-card>
            <q-card-section class="row items-center q-pb-none">
              <div class="text-h6">Available Scripts</div>
              <q-space />
              <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>
            <q-card-section>
              <div class="row">
                <q-btn glossy label="Feed a Single Earthquake"
                  @click="activateFirstScriptDialog = !activateFirstScriptDialog" />
              </div>
              <div class="row">
                <q-btn glossy label="Feed Earthquakes from USGS Dumps"
                  @click="activateSecondScriptDialog = !activateSecondScriptDialog" />
              </div>
            </q-card-section>
          </q-card>
        </q-dialog>

        <SingleEarthquakeInsertDialogue v-model:show-if="activateFirstScriptDialog" />
        <BulkEarthquakeInsertDialogue v-model:show-if="activateSecondScriptDialog" />
        <WipeAllEarthquakesDialogue v-model:show-if="activateWipeEarthquakesDialog"
          v-model:current-row-earthquakes="currentRowEarthquakes" v-model:cached-earthquakes="cachedEarthquakes" />
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { GoogleMap, Marker } from 'vue3-google-map'
import { wsURL } from '../axios'
import { BSON, EJSON, ObjectId } from 'bson';
import { ref, Ref, render } from 'vue'
import { QTableProps } from 'quasar'
import { Earthquake, LatLng, EarthquakeCount } from '../types/Earthquakes'
import { httpClient } from '../axios'
import { request } from 'http';

import EarthquakeInfoCard from '../components/EarthquakeInfoCard.vue'
import SearchBar from '../components/SearchBar.vue'
import WipeAllEarthquakesDialogue from '../components/WipeAllEarthquakesDialogue.vue'
import SingleEarthquakeInsertDialogue from '../components/SingleEarthquakeInsertDialogue.vue'
import BulkEarthquakeInsertDialogue from '../components/BulkEarthquakeInsertDialogue.vue'

const positions = ref([] as LatLng[])

const activateScriptDialog = ref(false)
const activateFirstScriptDialog = ref(false)
const activateSecondScriptDialog = ref(false)
const activateWipeEarthquakesDialog = ref(false)

const unshiftIndexEarthquakeTable = ref(0)
const latestEarthquakesTableLoading = ref(false)

const selectedMarkerIndex = ref(0)
const renderMarker = ref(false)
const center = { lat: 40.689247, lng: -74.044502 }

const currentRowEarthquakes: Ref<Earthquake[]> = ref([])
const cachedEarthquakes: Ref<Earthquake[]> = ref([])

const latestEarthquakesColumns: QTableProps['columns'] = [
  {
    name: 'place',
    required: true,
    label: 'Place',
    align: 'left',
    field: 'place',
    sortable: true,
    classes: "text-weight-light"
  },
  {
    name: 'mag',
    required: true,
    label: 'Magnitude',
    align: 'center',
    field: 'mag',
    sortable: true,
    classes: "text-weight-light"
  },
  {
    name: 'time',
    required: true,
    label: 'Time',
    align: 'center',
    field: 'time',
    format: (val: number | string | Date) => new Date(val).toLocaleString(),
    sortable: true,
    classes: "text-weight-light"
  }
]
const earthquakeTablePagination = ref({
  sortBy: 'time',
  descending: false,
  page: 1,
  rowsPerPage: 5,
  rowsNumber: 0,
})


let ws = new WebSocket(wsURL + "earthquakes/stream");
const connectedToWs = ref(false)
ws.onopen = () => {
  connectedToWs.value = true
}
ws.onclose = () => {
  connectedToWs.value = false
}
ws.onmessage = (event) => {
  try {
    const byte_blob = event.data as Blob
    byte_blob.text().then(text => {
      const data = JSON.parse(text) as Earthquake
      // if the earthquake s id is in cached earthquakes, then do nothing
      positions.value.push({ lat: data.location.coordinates[1], lng: data.location.coordinates[0] })
      if (currentRowEarthquakes.value.length == earthquakeTablePagination.value.rowsPerPage) {
        // pop the value at currentRowEarthquakes.length unshiftIndexEarthquakeTable
        let current_value = currentRowEarthquakes.value.at(unshiftIndexEarthquakeTable.value)
        if (current_value) {
          cachedEarthquakes.value.splice(cachedEarthquakes.value.indexOf(current_value), 1)
          currentRowEarthquakes.value.splice(unshiftIndexEarthquakeTable.value, 1)
        }
        unshiftIndexEarthquakeTable.value++
        if (unshiftIndexEarthquakeTable.value >= 5) {
          unshiftIndexEarthquakeTable.value = 0
        }
      }
      currentRowEarthquakes.value.unshift(data)
      cachedEarthquakes.value.unshift(data)
      earthquakeTablePagination.value.rowsNumber++
    }).catch(function (error) {
      console.error('An error occurred:', error)
    })
  }
  catch (error) {
    console.error('An error occurred:', error)
  }
}

getInitialEarthquakePage()

function reconnectWS() {
  ws = new WebSocket(wsURL + "earthquakes/stream");
  connectedToWs.value = true
}

const latestEarthquakeFilter = ref('')
const onlyShowEarthquakesInTheTableAtMap = ref(false)
const showMarkerDetails = (location: LatLng) => {
  renderMarker.value = true
  // search for the location in cachedEarthquakes
  for (let i = 0; i < cachedEarthquakes.value.length; i++) {
    if (cachedEarthquakes.value[i].location.coordinates[1] === location.lat && cachedEarthquakes.value[i].location.coordinates[0] === location.lng) {
      selectedMarkerIndex.value = i
      break
    }
  }
}
function onLatestEarthquakeRowClick(evt: Event, row: Earthquake) {
  renderMarker.value = true
  selectedMarkerIndex.value = cachedEarthquakes.value.indexOf(row)
}

function getInitialEarthquakePage() {
  latestEarthquakesTableLoading.value = true
  httpClient.get('earthquakes/count').then(response => {
    const data = response.data as EarthquakeCount
    earthquakeTablePagination.value.rowsNumber = data.count
  }).catch((error) => {
    console.error(error)
  })
  httpClient.get('earthquakes/get').then(response => {
    let earthquake = EJSON.deserialize(response.data) as Earthquake[]
    if (earthquake) {
      const earthquakeSet = new Set(cachedEarthquakes.value.map(data => data._id));
      earthquake.forEach(data => {
        if (!earthquakeSet.has(data._id)) {
          cachedEarthquakes.value.push(data);
          earthquakeSet.add(data._id);
        }
      });
      currentRowEarthquakes.value.splice(0, currentRowEarthquakes.value.length, ...earthquake)
      earthquake.forEach((data) => {
        positions.value.push({ lat: data.location.coordinates[1], lng: data.location.coordinates[0] })
      })
    }
  }).catch((error) => {
    console.error(error)
  })
  latestEarthquakesTableLoading.value = false
}
function adjustMarkers() {
  if (onlyShowEarthquakesInTheTableAtMap.value) {
    positions.value.splice(0, positions.value.length)
    currentRowEarthquakes.value.forEach((data) => {
      if (!positions.value.some((position) => position.lat === data.location.coordinates[1] && position.lng === data.location.coordinates[0])) {
        positions.value.push({ lat: data.location.coordinates[1], lng: data.location.coordinates[0] })
      }
    })
  } else {
    positions.value.splice(0, positions.value.length)
  }
}
function onLatestEarthquakeTableScroll(requestProp: {
  pagination: {
    sortBy: string;
    descending: boolean;
    page: number;
    rowsPerPage: number;
  };
  filter?: any;
  getCellValue: (col: any, row: any) => any;
}) {
  const { page, rowsPerPage, sortBy, descending } = requestProp.pagination
  const filter = requestProp.filter
  let limit = rowsPerPage
  let offset = (page - 1) * rowsPerPage
  latestEarthquakesTableLoading.value = true
  httpClient.get('earthquakes/count' + '?filter' + filter).then(response => {
    const data = response.data as EarthquakeCount
    console.log(data)
    earthquakeTablePagination.value.rowsNumber = data.count
  }).catch((error) => {
    console.error(error)
  })
  httpClient.get('earthquakes/get' + '?limit=' + limit + '&offset=' + offset + '&filter=' + filter).then(response => {
    let earthquake = EJSON.deserialize(response.data) as Earthquake[]
    currentRowEarthquakes.value.splice(0, currentRowEarthquakes.value.length, ...earthquake)
    if (onlyShowEarthquakesInTheTableAtMap.value) {
      positions.value.splice(0, positions.value.length)
    }
    earthquake.forEach((data) => {
      positions.value.push({ lat: data.location.coordinates[1], lng: data.location.coordinates[0] })
      cachedEarthquakes.value.push(data)
    })
  }).catch((error) => {
    console.error(error)
  })
  latestEarthquakesTableLoading.value = false
  earthquakeTablePagination.value.page = page
  earthquakeTablePagination.value.rowsPerPage = rowsPerPage
  earthquakeTablePagination.value.sortBy = sortBy
  earthquakeTablePagination.value.descending = descending
}
function isUniqueTableKey(row: Earthquake) {
  currentRowEarthquakes.value.forEach((data) => {
    if (data._id === row._id) {
      return false
    }
  })
}

</script>
