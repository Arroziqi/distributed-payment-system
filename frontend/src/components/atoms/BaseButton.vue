<script setup lang="ts">
import { Button } from '@/components/ui/button';
import type { ButtonVariants } from '@/components/ui/button';

interface Props {
  variant?: ButtonVariants['variant'];
  size?: ButtonVariants['size'];
  class?: any;
  loading?: boolean;
  disabled?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  size: 'default',
});
</script>

<template>
  <Button
    :variant="variant"
    :size="size"
    :class="props.class"
    :disabled="disabled || loading"
  >
    <span v-if="loading" class="mr-2 h-4 w-4 animate-spin border-2 border-current border-t-transparent rounded-full" />
    <slot />
  </Button>
</template>
