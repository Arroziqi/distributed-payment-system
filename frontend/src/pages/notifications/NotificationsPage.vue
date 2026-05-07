<script setup lang="ts">
import NotificationItem from '@/components/molecules/NotificationItem.vue';
import BaseButton from '@/components/atoms/BaseButton.vue';
import BaseIcon from '@/components/atoms/BaseIcon.vue';

import { ref } from 'vue';

const notifications = ref<{
  id: number;
  title: string;
  message: string;
  time: string;
  type: 'info' | 'success' | 'warning' | 'error';
  unread: boolean;
}[]>([
  { id: 1, title: 'Payment Success', message: 'You have successfully paid $50 to Starbucks.', time: '2 mins ago', type: 'success', unread: true },
  { id: 2, title: 'Low Balance', message: 'Your wallet balance is below $10. Please top up.', time: '1 hour ago', type: 'warning', unread: true },
  { id: 3, title: 'Security Alert', message: 'New login detected from Safari on Mac.', time: '5 hours ago', type: 'info', unread: false },
  { id: 4, title: 'Top up Failed', message: 'Your top up request was rejected by the bank.', time: 'Yesterday', type: 'error', unread: false },
]);
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-3xl font-bold tracking-tight">Notifications</h2>
        <p class="text-muted-foreground">Stay updated with your account activity.</p>
      </div>
      <BaseButton variant="outline">
        <BaseIcon name="CheckCheck" class="mr-2 h-4 w-4" />
        Mark all as read
      </BaseButton>
    </div>

    <div class="space-y-4 max-w-3xl">
      <NotificationItem
        v-for="note in notifications"
        :key="note.id"
        v-bind="note"
      />
    </div>
    
    <div v-if="notifications.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
      <div class="rounded-full bg-muted p-4 mb-4">
        <BaseIcon name="BellOff" class="h-8 w-8 text-muted-foreground" />
      </div>
      <h3 class="text-lg font-semibold">No notifications</h3>
      <p class="text-muted-foreground">We'll notify you when something happens.</p>
    </div>
  </div>
</template>
