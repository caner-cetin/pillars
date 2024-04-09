<template>
    <q-dialog v-model="showIf">
        <q-card style="width: 600px; max-width: 80vw; padding:20px">
            <q-form @submit="feedSingleEarthquake" @reset="onReset" class="q-gutter-md">
                <q-input v-model="place" label="Place" hint="Short and concise location of the place" lazy-rules
                    :rules="[val => val && val.length > 0 || 'Required!', val => val.length <= 50 || 'Maximum 50 characters!']" />
                <q-input v-model="mag" label="Magnitude" lazy-rules
                    :rules="[val => val && val.length > 0 || 'Required!', val => val >= 0 && val <= 10 || 'Magnitude should be between 0 and 10!']" />
                <q-input v-model="lat" label="Latitude" lazy-rules
                    :rules="[val => val && val.length > 0 || 'Required!', val => val >= -90 && val <= 90 || 'Latitude should be between -90 and 90!']" />
                <q-input v-model="lon" label="Longitude" lazy-rules
                    :rules="[val => val && val.length > 0 || 'Required!', val => val >= -180 && val <= 180 || 'Longitude should be between -180 and 180!']" />
                <div align="right">
                    <q-btn label="Reset" type="reset" color="primary" flat class="q-ml-sm" />
                    <q-btn label="Submit" type="submit" color="primary" />
                </div>
            </q-form>
        </q-card>
    </q-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useQuasar, QSpinnerGears } from 'quasar'
import { ObjectId } from 'bson'
import { httpClient } from '../axios'
import { Earthquake } from '../types/Earthquakes'

const qua = useQuasar()
const showIf = defineModel('showIf', { type: Boolean, required: true })
const place = ref('')
const mag = ref('')
const lat = ref('')
const lon = ref('')
function onReset() {
    place.value = ''
    mag.value = ''
    lat.value = ''
    lon.value = ''
}
function feedSingleEarthquake() {
    const submit_dismiss = qua.notify({
        spinner: QSpinnerGears,
        message: 'Submitting...',
        timeout: 2000
    })
    const data: Earthquake = {
        _id: new ObjectId(),
        mag: parseFloat(mag.value),
        place: place.value,
        time: new Date(),
        location: {
            type: "Point",
            coordinates: [parseFloat(lon.value), parseFloat(lat.value)]
        }
    }
    httpClient.put('earthquakes/insert', data).then(response => {
        submit_dismiss()
        setTimeout(() => {
            qua.notify({
                message: 'Earthquake submitted.',
                color: 'green',
                position: 'bottom'
            })
        }, 2000)
        showIf.value = false
        onReset()
    }).catch((error) => {
        if (error.response) {
            console.log(error.response.data);
            console.log(error.response.status);
            console.log(error.response.headers);
        } else if (error.request) {
            console.log(error.request);
        } else {
            console.log('Error', error.message);
        }
        console.log(error.config);
        submit_dismiss()
        qua.notify({
            message: 'Failed to submit earthquake.',
            color: 'red',
            position: 'bottom'
        })
    })
}
</script>