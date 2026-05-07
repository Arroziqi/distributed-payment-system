<script setup lang="ts">
import { Label } from '@/components/ui/label';

interface Props {
  for?: string;
  class?: any;
}

const props = defineProps<Props>();
</script>

<template>
  <Label :for="props.for" :class="props.class">
    <slot />
  </Label>
</template>
