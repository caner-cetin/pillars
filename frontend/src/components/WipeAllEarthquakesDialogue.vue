<template>
    <q-dialog v-model="showIf">
        <q-card>
            <q-card-section>
                <q-card-section>
                    <div class="text-h6">Wipe All Earthquakes</div>
                </q-card-section>
                <q-separator dark />
                <q-card-section>
                    Are you sure you want to wipe all earthquakes?
                </q-card-section>
                <q-card-actions align="right">
                    <q-btn label="Cancel" color="primary" @click="showIf = false" />
                    <q-btn label="Wipe" color="negative" @click="wipeAllEarthquakes" />
                </q-card-actions>
            </q-card-section>
        </q-card>
    </q-dialog>
</template>

<script setup lang="ts">
import { defineModel, PropType } from 'vue';
import { httpClient } from 'src/axios';
import { useQuasar } from 'quasar';
import { Earthquake } from 'src/types/Earthquakes';
const showIf = defineModel("showIf", { type: Boolean, required: true });
const cachedEarthquakes = defineModel("cachedEarthquakes", { type: Array as PropType<Earthquake[]>, required: true });
const currentRowEarthquakes = defineModel("currentRowEarthquakes", { type: Array as PropType<Earthquake[]>, required: true });
let qua = useQuasar();
function wipeAllEarthquakes() {
    httpClient.delete('/earthquakes/delete').then((response) => {
        setTimeout(() => {
            qua.notify({
                message: 'Wiped all earthquakes.',
                color: 'positive',
                position: 'bottom'
            });
        }, 1000);
        showIf.value = false;
        cachedEarthquakes.value = [];
        currentRowEarthquakes.value = [];
    }
    ).catch((error) => {
        setTimeout(() => {
            qua.notify({
                message: 'Failed to wipe all earthquakes.',
                color: 'negative',
                position: 'bottom'
            });
        }, 1000);
        console.error(error);
    });
}
</script>