<script setup lang="ts">
import { Badge } from '@/components/ui/badge';
import type { BadgeVariants } from '@/components/ui/badge';

interface Props {
  variant?: BadgeVariants['variant'];
  class?: any;
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
});
</script>

<template>
  <Badge :variant="props.variant" :class="props.class">
    <slot />
  </Badge>
</template>
