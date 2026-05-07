<script setup lang="ts">
import { Loader2 } from 'lucide-vue-next';
import { computed } from 'vue';

interface Props {
  size?: number | string;
  class?: any;
}

const props = withDefaults(defineProps<Props>(), {
  size: 24,
});

const numericSize = computed(() => {
  return typeof props.size === 'string' ? parseInt(props.size, 10) : props.size;
});
</script>

<template>
  <Loader2 :size="numericSize" :class="['animate-spin', props.class]" />
</template>
