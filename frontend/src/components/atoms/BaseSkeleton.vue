<script setup lang="ts">
import { Skeleton } from '@/components/ui/skeleton';

interface Props {
  class?: any;
}

const props = defineProps<Props>();
</script>

<template>
  <Skeleton :class="props.class" />
</template>
