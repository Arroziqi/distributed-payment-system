<script setup lang="ts">
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

interface Option {
  label: string;
  value: string;
}

interface Props {
  modelValue?: string;
  placeholder?: string;
  options: Option[];
  label?: string;
  class?: any;
}

const props = defineProps<Props>();
const emit = defineEmits(['update:modelValue']);
</script>

<template>
  <Select :model-value="props.modelValue" @update:model-value="emit('update:modelValue', $event)">
    <SelectTrigger :class="props.class">
      <SelectValue :placeholder="props.placeholder" />
    </SelectTrigger>
    <SelectContent>
      <SelectGroup>
        <SelectLabel v-if="props.label">{{ props.label }}</SelectLabel>
        <SelectItem v-for="opt in props.options" :key="opt.value" :value="opt.value">
          {{ opt.label }}
        </SelectItem>
      </SelectGroup>
    </SelectContent>
  </Select>
</template>
