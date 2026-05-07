<script setup lang="ts">
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { AlertCircle, CheckCircle2, Info, AlertTriangle } from 'lucide-vue-next';
import { computed } from 'vue';

interface Props {
  variant?: 'default' | 'destructive' | 'success' | 'warning';
  title?: string;
  description?: string;
  class?: any;
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
});

const icon = computed(() => {
  switch (props.variant) {
    case 'destructive': return AlertCircle;
    case 'success': return CheckCircle2;
    case 'warning': return AlertTriangle;
    default: return Info;
  }
});

const variantClass = computed(() => {
  if (props.variant === 'success') return 'border-green-500 text-green-700 bg-green-50 dark:bg-green-900/10';
  if (props.variant === 'warning') return 'border-yellow-500 text-yellow-700 bg-yellow-50 dark:bg-yellow-900/10';
  return '';
});
</script>

<template>
  <Alert :variant="props.variant === 'destructive' ? 'destructive' : 'default'" :class="[variantClass, props.class]">
    <component :is="icon" class="h-4 w-4" />
    <AlertTitle v-if="props.title">{{ props.title }}</AlertTitle>
    <AlertDescription v-if="props.description || $slots.default">
      <slot>{{ props.description }}</slot>
    </AlertDescription>
  </Alert>
</template>
