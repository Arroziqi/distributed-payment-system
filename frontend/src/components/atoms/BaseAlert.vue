<script setup lang="ts">
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { AlertCircle, CheckCircle2, Info, AlertTriangle, X } from 'lucide-vue-next';
import { computed, onMounted, onUnmounted, ref, watch } from 'vue';

interface Props {
  variant?: 'default' | 'destructive' | 'success' | 'warning';
  title?: string;
  description?: string;
  class?: any;
  showProgress?: boolean;
  duration?: number; // in milliseconds
  showClose?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  duration: 5000,
  showClose: false,
});

const emit = defineEmits(['close']);

const progress = ref(100);
let timer: any = null;
let progressTimer: any = null;

const icon = computed(() => {
  switch (props.variant) {
    case 'destructive': return AlertCircle;
    case 'success': return CheckCircle2;
    case 'warning': return AlertTriangle;
    default: return Info;
  }
});

const variantClass = computed(() => {
  if (props.variant === 'success') return 'border-green-500/50 text-green-700 bg-green-50 dark:bg-green-900/10 dark:text-green-400';
  if (props.variant === 'warning') return 'border-yellow-500/50 text-yellow-700 bg-yellow-50 dark:bg-yellow-900/10 dark:text-yellow-400';
  return '';
});

const progressBarClass = computed(() => {
  switch (props.variant) {
    case 'success': return 'bg-green-500';
    case 'warning': return 'bg-yellow-500';
    case 'destructive': return 'bg-destructive';
    default: return 'bg-primary';
  }
});

const startTimer = () => {
  if (props.duration > 0) {
    const startTime = Date.now();
    const endTime = startTime + props.duration;

    progressTimer = setInterval(() => {
      const now = Date.now();
      const remaining = Math.max(0, endTime - now);
      progress.value = (remaining / props.duration) * 100;
      
      if (remaining <= 0) {
        clearInterval(progressTimer);
        emit('close');
      }
    }, 10);
  }
};

onMounted(() => {
  if (props.showProgress || props.duration > 0) {
    startTimer();
  }
});

onUnmounted(() => {
  if (progressTimer) clearInterval(progressTimer);
});
</script>

<template>
  <Alert 
    :variant="props.variant === 'destructive' ? 'destructive' : 'default'" 
    :class="['relative overflow-hidden transition-all duration-300', variantClass, props.class]"
  >
    <div class="flex items-start gap-4">
      <component :is="icon" class="h-5 w-5 mt-0.5 shrink-0" />
      <div class="flex-1">
        <AlertTitle v-if="props.title" class="font-bold tracking-tight">{{ props.title }}</AlertTitle>
        <AlertDescription v-if="props.description || $slots.default" class="text-sm opacity-90 leading-relaxed">
          <slot>{{ props.description }}</slot>
        </AlertDescription>
      </div>
      <button 
        v-if="showClose" 
        @click="emit('close')"
        class="shrink-0 p-1 hover:bg-black/5 dark:hover:bg-white/5 rounded-md transition-colors"
      >
        <X class="h-4 w-4" />
      </button>
    </div>

    <!-- Progress Bar -->
    <div 
      v-if="showProgress" 
      class="absolute bottom-0 left-0 h-1 transition-all duration-[10ms] ease-linear"
      :class="progressBarClass"
      :style="{ width: `${progress}%` }"
    />
  </Alert>
</template>
