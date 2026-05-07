<script setup lang="ts">
import BaseIcon from '@/components/atoms/BaseIcon.vue';
</script>

<template>
  <div class="min-h-screen grid lg:grid-cols-2">
    <div class="hidden lg:flex flex-col justify-between p-12 bg-zinc-900 text-white">
      <div class="flex items-center gap-2">
        <div class="bg-primary rounded-lg p-1">
          <BaseIcon name="Zap" class="h-6 w-6 text-white" />
        </div>
        <span class="text-xl font-bold">PayFlow</span>
      </div>
      
      <div class="space-y-6">
        <h1 class="text-4xl font-bold leading-tight">
          Manage your global payments with confidence.
        </h1>
        <p class="text-zinc-400 text-lg max-w-md">
          The ultimate distributed payment system for modern teams. Secure, fast, and transparent.
        </p>
      </div>
      
      <div class="flex items-center gap-4">
        <div class="flex -space-x-3">
          <div v-for="i in 4" :key="i" class="h-10 w-10 rounded-full border-2 border-zinc-900 bg-zinc-800" />
        </div>
        <p class="text-sm text-zinc-400">
          Joined by over 10,000+ businesses worldwide.
        </p>
      </div>
    </div>
    
    <div class="flex items-center justify-center p-8 bg-background">
      <div class="w-full max-w-[400px] space-y-8">
        <div class="lg:hidden flex justify-center mb-8">
          <div class="flex items-center gap-2">
            <BaseIcon name="Zap" class="h-8 w-8 text-primary" />
            <span class="text-2xl font-bold">PayFlow</span>
          </div>
        </div>
        <router-view />
      </div>
    </div>
  </div>
</template>
