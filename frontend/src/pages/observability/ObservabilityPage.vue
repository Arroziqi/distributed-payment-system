<script setup lang="ts">
import BaseBadge from '@/components/atoms/BaseBadge.vue';
import StatsCard from '@/components/molecules/StatsCard.vue';

const services = [
  { name: 'Auth Service', status: 'online', latency: '45ms', uptime: '99.9%' },
  { name: 'Wallet Service', status: 'online', latency: '12ms', uptime: '99.9%' },
  { name: 'Transaction Service', status: 'online', latency: '89ms', uptime: '99.8%' },
  { name: 'Notification Service', status: 'online', latency: '23ms', uptime: '99.9%' },
];
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-3xl font-bold tracking-tight">System Health</h2>
    </div>

    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <StatsCard title="CPU Usage" value="24%" icon="Cpu" description="Average across nodes" />
      <StatsCard title="Memory" value="1.2GB" icon="Database" description="Total allocated" />
      <StatsCard title="Requests/s" value="450" icon="Activity" :trend="{ value: 8, isPositive: true }" />
      <StatsCard title="Error Rate" value="0.02%" icon="AlertCircle" :trend="{ value: 1, isPositive: false }" />
    </div>

    <div class="rounded-xl border bg-card p-6">
      <h3 class="text-lg font-semibold mb-4">Service Status</h3>
      <div class="space-y-4">
        <div v-for="svc in services" :key="svc.name" class="flex items-center justify-between p-3 border rounded-lg">
          <div class="flex items-center">
            <div class="h-2 w-2 rounded-full bg-green-500 mr-3" />
            <span class="font-medium">{{ svc.name }}</span>
          </div>
          <div class="flex items-center space-x-6 text-sm text-muted-foreground">
            <span>Latency: {{ svc.latency }}</span>
            <span>Uptime: {{ svc.uptime }}</span>
            <BaseBadge>Healthy</BaseBadge>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
