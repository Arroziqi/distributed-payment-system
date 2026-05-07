<script setup lang="ts">
import {
  Pagination,
  PaginationEllipsis,
  PaginationFirst,
  PaginationLast,
  PaginationItem,
  PaginationNext,
  PaginationPrevious,
  PaginationContent,
} from '@/components/ui/pagination';

interface Props {
  total?: number;
  itemsPerPage?: number;
  siblingCount?: number;
  showEdges?: boolean;
  defaultPage?: number;
  class?: any;
}

const props = withDefaults(defineProps<Props>(), {
  total: 0,
  itemsPerPage: 10,
  siblingCount: 1,
  showEdges: true,
  defaultPage: 1,
});
</script>

<template>
  <Pagination
    :total="props.total"
    :items-per-page="props.itemsPerPage"
    :sibling-count="props.siblingCount"
    :show-edges="props.showEdges"
    :default-page="props.defaultPage"
  >
    <PaginationContent v-slot="{ items }" class="flex items-center gap-1">
      <PaginationFirst />
      <PaginationPrevious />

      <template v-for="(item, index) in items">
        <PaginationItem v-if="item.type === 'page'" :key="index" :value="item.value" :is-active="item.value === props.defaultPage">
          {{ item.value }}
        </PaginationItem>
        <PaginationEllipsis v-else :key="item.type" :index="index" />
      </template>

      <PaginationNext />
      <PaginationLast />
    </PaginationContent>
  </Pagination>
</template>
