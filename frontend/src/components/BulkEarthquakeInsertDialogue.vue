<template>
    <q-dialog v-model="showIf">
        <q-uploader :factory="factoryFn" flat bordered style="max-width: 600px" accept=".json, .geojson"
            label="Select / Drag&Drop Dumps" color="deep-purple-7" text-color="white" no-thumbnails multiple
            field-name="usgs-dump" />
    </q-dialog>
</template>


<script setup lang="ts">
import { defineModel } from 'vue'
import { restURL } from '../axios'

const showIf = defineModel('showIf', { type: Boolean, required: true })
function factoryFn(files: readonly File[]) {
    return {
        url: restURL + 'earthquakes/bulk-insert',
        method: 'PUT'
    }
}
</script>
