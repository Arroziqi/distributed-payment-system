<script setup lang="ts">
import * as LucideIcons from 'lucide-vue-next';
import { computed } from 'vue';

interface Props {
  name: string;
  size?: number | string;
  color?: string;
  strokeWidth?: number | string;
  class?: any;
}

const props = withDefaults(defineProps<Props>(), {
  size: 20,
  strokeWidth: 2,
});

const icon = computed(() => {
  // @ts-ignore
  return LucideIcons[props.name];
});
</script>

<template>
  <component
    :is="icon"
    v-if="icon"
    :size="props.size"
    :color="props.color"
    :stroke-width="props.strokeWidth"
    :class="props.class"
  />
</template>
