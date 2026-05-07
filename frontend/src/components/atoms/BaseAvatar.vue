<script setup lang="ts">
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';

interface Props {
  src?: string;
  alt?: string;
  fallback?: string;
  class?: any;
}

const props = defineProps<Props>();
</script>

<template>
  <Avatar :class="props.class">
    <AvatarImage v-if="props.src" :src="props.src" :alt="props.alt" />
    <AvatarFallback>{{ props.fallback || '?' }}</AvatarFallback>
  </Avatar>
</template>
