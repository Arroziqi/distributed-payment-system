<template>
  <q-card class="base-card q-pa-md">
    <q-card-section v-if="$slots.header || title" class="q-pb-none">
      <slot name="header">
        <div class="text-h6">{{ title }}</div>
      </slot>
    </q-card-section>

    <q-card-section>
      <slot></slot>
    </q-card-section>

    <q-card-actions v-if="$slots.actions" align="right">
      <slot name="actions"></slot>
    </q-card-actions>
  </q-card>
</template>

<script setup lang="ts">
defineProps<{
  title?: string;
}>();
</script>

<style scoped lang="sass">
.base-card
  border-radius: 8px
  box-shadow: 0 4px 12px rgba(0,0,0,0.05)
</style>
