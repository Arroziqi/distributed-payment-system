<script setup lang="ts">
import { Input } from '@/components/ui/input';

interface Props {
  modelValue?: string | number;
  type?: string;
  placeholder?: string;
  disabled?: boolean;
  class?: any;
  error?: boolean;
}

const props = defineProps<Props>();
const emit = defineEmits(['update:modelValue']);

const onInput = (event: Event) => {
  const target = event.target as HTMLInputElement;
  emit('update:modelValue', target.value);
};
</script>

<template>
  <Input
    :type="props.type"
    :value="props.modelValue"
    :placeholder="props.placeholder"
    :disabled="props.disabled"
    :class="[
      props.class,
      props.error && 'border-destructive ring-destructive'
    ]"
    @input="onInput"
  />
</template>
