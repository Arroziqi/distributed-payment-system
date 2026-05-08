<script setup lang="ts">
import { useRoute } from 'vue-router';
import BaseIcon from '@/components/atoms/BaseIcon.vue';

const route = useRoute();

const navigation = [
  { name: 'Dashboard', href: '/', icon: 'LayoutDashboard' },
  { name: 'Wallet', href: '/wallet', icon: 'Wallet' },
  { name: 'Transactions', href: '/transactions', icon: 'ArrowRightLeft' },
  { name: 'Notifications', href: '/notifications', icon: 'Bell' },
  { name: 'Observability', href: '/observability', icon: 'LineChart' },
  { name: 'Analytics', href: '/analytics', icon: 'PieChart' },
  { name: 'Settings', href: '/settings', icon: 'Settings' },
];

const isActive = (path: string) => {
  if (path === '/' && route.path !== '/') return false;
  return route.path.startsWith(path);
};
</script>

<template>
  <aside class="flex h-full w-64 flex-col border-r bg-card transition-all">
    <div class="flex h-16 items-center border-b px-6">
      <BaseIcon name="Zap" class="h-6 w-6 text-primary mr-2" />
      <span class="text-lg font-bold tracking-tight">PayFlow</span>
    </div>
    <nav class="flex-1 space-y-1 px-3 py-4">
      <router-link
        v-for="item in navigation"
        :key="item.name"
        :to="item.href"
        :class="[
          'group flex items-center rounded-md px-3 py-2 text-sm font-medium transition-colors',
          isActive(item.href)
            ? 'bg-primary text-primary-foreground'
            : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
        ]"
      >
        <BaseIcon
          :name="item.icon"
          :class="['mr-3 h-5 w-5', isActive(item.href) ? 'text-primary-foreground' : 'text-muted-foreground group-hover:text-accent-foreground']"
        />
        {{ item.name }}
      </router-link>
    </nav>
    <div class="p-4 border-t">
      <div class="rounded-lg bg-muted/50 p-3">
        <p class="text-xs font-medium text-muted-foreground mb-1">Status</p>
        <div class="flex items-center">
          <div class="h-2 w-2 rounded-full bg-green-500 mr-2" />
          <span class="text-[10px] font-medium">System Online</span>
        </div>
      </div>
    </div>
  </aside>
</template>
