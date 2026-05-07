<script setup lang="ts">
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';

interface Tab {
  label: string;
  value: string;
}

interface Props {
  tabs: Tab[];
  defaultValue?: string;
  class?: any;
}

const props = defineProps<Props>();
</script>

<template>
  <Tabs :default-value="props.defaultValue" :class="props.class">
    <TabsList>
      <TabsTrigger v-for="tab in props.tabs" :key="tab.value" :value="tab.value">
        {{ tab.label }}
      </TabsTrigger>
    </TabsList>
    <template v-for="tab in props.tabs" :key="tab.value">
      <TabsContent :value="tab.value">
        <slot :name="tab.value" />
      </TabsContent>
    </template>
  </Tabs>
</template>
