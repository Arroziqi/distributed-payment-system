<script setup lang="ts">
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import BaseIcon from '@/components/atoms/BaseIcon.vue';

interface Props {
  title: string;
  value: string | number;
  description?: string;
  icon: string;
  trend?: {
    value: number;
    isPositive: boolean;
  };
}

defineProps<Props>();
</script>

<template>
  <Card>
    <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
      <CardTitle class="text-sm font-medium">{{ title }}</CardTitle>
      <BaseIcon :name="icon" class="h-4 w-4 text-muted-foreground" />
    </CardHeader>
    <CardContent>
      <div class="text-2xl font-bold">{{ value }}</div>
      <p v-if="description" class="text-xs text-muted-foreground">
        <span v-if="trend" :class="trend.isPositive ? 'text-green-500' : 'text-red-500'">
          {{ trend.isPositive ? '+' : '' }}{{ trend.value }}%
        </span>
        {{ description }}
      </p>
    </CardContent>
  </Card>
</template>
