<script setup lang="ts">
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/components/ui/drawer';

interface Props {
  title?: string;
  description?: string;
  class?: any;
}

const props = defineProps<Props>();
</script>

<template>
  <Drawer>
    <DrawerTrigger as-child>
      <slot name="trigger" />
    </DrawerTrigger>
    <DrawerContent :class="props.class">
      <div class="mx-auto w-full max-w-sm">
        <DrawerHeader>
          <DrawerTitle v-if="props.title">{{ props.title }}</DrawerTitle>
          <DrawerDescription v-if="props.description">
            {{ props.description }}
          </DrawerDescription>
        </DrawerHeader>
        <slot />
        <DrawerFooter v-if="$slots.footer">
          <slot name="footer" />
          <DrawerClose as-child>
            <slot name="close-trigger" />
          </DrawerClose>
        </DrawerFooter>
      </div>
    </DrawerContent>
  </Drawer>
</template>
