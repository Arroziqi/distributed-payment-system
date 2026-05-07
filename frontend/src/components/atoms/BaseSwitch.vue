<script setup lang="ts">
import { Switch } from '@/components/ui/switch';

interface Props {
  checked?: boolean;
  disabled?: boolean;
  class?: any;
}

const props = defineProps<Props>();
const emit = defineEmits(['update:checked']);
</script>

<template>
  <Switch
    :checked="props.checked"
    :disabled="props.disabled"
    :class="props.class"
    @update:checked="emit('update:checked', $event)"
  />
</template>
