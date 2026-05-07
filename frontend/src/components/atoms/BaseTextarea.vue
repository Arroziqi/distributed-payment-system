<script setup lang="ts">
import { Textarea } from '@/components/ui/textarea';

interface Props {
  modelValue?: string;
  placeholder?: string;
  disabled?: boolean;
  class?: any;
}

const props = defineProps<Props>();
const emit = defineEmits(['update:modelValue']);
</script>

<template>
  <Textarea
    :model-value="props.modelValue"
    :placeholder="props.placeholder"
    :disabled="props.disabled"
    :class="props.class"
    @input="emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
  />
</template>
