<template>
    <q-card class="earthquake-card">
        <q-card-section>
            <div v-if="earthquake">
                <div class="text-h4">
                    <q-icon name="sym_o_earthquake">
                    </q-icon>
                    {{ earthquake.mag }}
                </div>
                <div class="text-subtitle2">
                    {{ earthquake.place }}
                    <q-separator spaced inset />
                    <q-icon name="sym_o_location_on">

                    </q-icon>
                    {{ earthquake.location.coordinates[1] }} , {{ earthquake.location.coordinates[0] }}
                    <q-separator spaced inset />
                    <div class="text-caption text-grey">
                        {{ new Date(earthquake.time).toLocaleString() }}
                    </div>
                </div>
            </div>
        </q-card-section>

    </q-card>
</template>


<script lang="ts">
import { PropType, Ref } from 'vue'
import { Earthquake } from '../types/Earthquakes'
import { QCard, QCardSection, QSeparator } from 'quasar'
export default {
    props: {
        earthquake: {
            type: Object as PropType<Earthquake>,
            required: false
        },
    },
    components: {
        QCard,
        QCardSection,
        QSeparator,
    },
}
</script>