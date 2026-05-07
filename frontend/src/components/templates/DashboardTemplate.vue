<template>
  <q-layout view="hHh Lpr lFf">
    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
        <q-btn flat dense round icon="menu" aria-label="Menu" @click="leftDrawerOpen = !leftDrawerOpen" />
        <q-toolbar-title>
          Distributed Payment System
        </q-toolbar-title>
      </q-toolbar>
    </q-header>

    <SidebarNavigation v-model="leftDrawerOpen" @logout="$emit('logout')" />

    <q-page-container>
      <div class="q-pa-md">
        <slot></slot>
      </div>
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import SidebarNavigation from '@/components/organisms/SidebarNavigation.vue';

const leftDrawerOpen = ref(true);

defineEmits(['logout']);
</script>
