<script setup lang="ts">
import { Checkbox } from '@/components/ui/checkbox';

interface Props {
  checked?: boolean;
  disabled?: boolean;
  class?: any;
}

const props = defineProps<Props>();
const emit = defineEmits(['update:checked']);
</script>

<template>
  <Checkbox
    :checked="props.checked"
    :disabled="props.disabled"
    :class="props.class"
    @update:checked="emit('update:checked', $event)"
  />
</template>
