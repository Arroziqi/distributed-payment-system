<script setup lang="ts">
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';

interface Props {
  text: string;
  side?: 'top' | 'right' | 'bottom' | 'left';
  class?: any;
}

const props = withDefaults(defineProps<Props>(), {
  side: 'top',
});
</script>

<template>
  <TooltipProvider>
    <Tooltip>
      <TooltipTrigger as-child>
        <slot />
      </TooltipTrigger>
      <TooltipContent :side="props.side" :class="props.class">
        <p>{{ props.text }}</p>
      </TooltipContent>
    </Tooltip>
  </TooltipProvider>
</template>
