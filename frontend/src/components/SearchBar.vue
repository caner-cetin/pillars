<script setup lang="ts">
import { defineModel } from 'vue'
const search = defineModel('search', { type: String, required: true })
</script>

<template>
    <q-input borderless dense debounce="300" v-model="search" placeholder="Search in any column">
        <template v-slot:append>
            <q-icon name="search" />
        </template>
    </q-input>
</template>