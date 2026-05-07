<script setup lang="ts">
import BaseIcon from '@/components/atoms/BaseIcon.vue';

interface Props {
  title: string;
  message: string;
  time: string;
  type: 'info' | 'success' | 'warning' | 'error';
  unread?: boolean;
}

defineProps<Props>();

const iconMap = {
  info: 'Info',
  success: 'CheckCircle2',
  warning: 'AlertTriangle',
  error: 'AlertCircle',
};

const colorMap = {
  info: 'text-blue-500',
  success: 'text-green-500',
  warning: 'text-yellow-500',
  error: 'text-red-500',
};
</script>

<template>
  <div :class="['flex items-start space-x-4 p-4 rounded-lg transition-colors hover:bg-muted/50', unread && 'bg-primary/5']">
    <div :class="['mt-1', colorMap[type]]">
      <BaseIcon :name="iconMap[type]" class="h-5 w-5" />
    </div>
    <div class="flex-1 space-y-1">
      <div class="flex items-center justify-between">
        <p class="text-sm font-medium leading-none">{{ title }}</p>
        <p class="text-xs text-muted-foreground">{{ time }}</p>
      </div>
      <p class="text-xs text-muted-foreground line-clamp-2">
        {{ message }}
      </p>
    </div>
    <div v-if="unread" class="h-2 w-2 rounded-full bg-primary mt-2" />
  </div>
</template>
