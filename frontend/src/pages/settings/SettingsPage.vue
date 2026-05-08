<script setup lang="ts">
import BaseIcon from '@/components/atoms/BaseIcon.vue';
import BaseButton from '@/components/atoms/BaseButton.vue';
import { useRouter } from 'vue-router';

const router = useRouter();
</script>

<template>
  <div class="flex flex-col items-center justify-center min-h-[60vh] text-center px-4">
    <div class="relative mb-8">
      <div class="absolute inset-0 bg-primary/20 blur-3xl rounded-full scale-150 animate-pulse"></div>
      <div class="relative bg-card border rounded-2xl p-8 shadow-2xl backdrop-blur-sm">
        <BaseIcon name="Settings" class="h-20 w-20 text-primary mx-auto animate-spin-slow" :stroke-width="1.5" />
      </div>
    </div>
    
    <h1 class="text-4xl font-extrabold tracking-tight mb-4 bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
      Settings Coming Soon
    </h1>
    
    <p class="text-muted-foreground max-w-md mx-auto mb-8 text-lg">
      We're working hard to bring you a comprehensive settings dashboard. Stay tuned for updates on security, preferences, and more!
    </p>
    
    <div class="flex flex-col sm:flex-row gap-4">
      <BaseButton size="lg" @click="router.push('/')">
        <BaseIcon name="LayoutDashboard" class="mr-2 h-5 w-5" />
        Back to Dashboard
      </BaseButton>
      <BaseButton variant="outline" size="lg" @click="router.push('/profile')">
        <BaseIcon name="User" class="mr-2 h-5 w-5" />
        View Profile
      </BaseButton>
    </div>
  </div>
</template>

<style scoped>
.animate-spin-slow {
  animation: spin 8s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
