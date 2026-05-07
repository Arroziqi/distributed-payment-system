<script setup lang="ts">
import BaseInput from '@/components/atoms/BaseInput.vue';
import BaseIcon from '@/components/atoms/BaseIcon.vue';

interface Props {
  modelValue?: string;
  placeholder?: string;
  class?: any;
}

const props = defineProps<Props>();
defineEmits(['update:modelValue']);
</script>

<template>
  <div :class="['relative w-full max-w-sm items-center', props.class]">
    <BaseInput
      :model-value="props.modelValue"
      :placeholder="props.placeholder || 'Search...'"
      class="pl-10"
      @update:model-value="$emit('update:modelValue', $event)"
    />
    <span class="absolute start-0 inset-y-0 flex items-center justify-center px-3">
      <BaseIcon name="Search" class="h-4 w-4 text-muted-foreground" />
    </span>
  </div>
</template>
