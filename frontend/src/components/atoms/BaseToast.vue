<script setup lang="ts">
import { toast } from 'vue-sonner';

const show = (message: string, options?: any) => {
  toast(message, options);
};

const success = (message: string, options?: any) => {
  toast.success(message, options);
};

const error = (message: string, options?: any) => {
  toast.error(message, options);
};

const warning = (message: string, options?: any) => {
  toast.warning(message, options);
};

const info = (message: string, options?: any) => {
  toast.info(message, options);
};

defineExpose({
  show,
  success,
  error,
  warning,
  info,
});
</script>

<template>
  <!-- No visual element, just logic wrapper for sonner -->
</template>
