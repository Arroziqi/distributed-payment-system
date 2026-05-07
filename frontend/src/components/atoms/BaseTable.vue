<script setup lang="ts">
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';

interface Column {
  key: string;
  label: string;
  class?: string;
}

interface Props {
  columns: Column[];
  data: any[];
  caption?: string;
  loading?: boolean;
  class?: any;
}

const props = defineProps<Props>();
</script>

<template>
  <div class="rounded-md border">
    <Table :class="props.class">
      <TableCaption v-if="props.caption">{{ props.caption }}</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead v-for="col in props.columns" :key="col.key" :class="col.class">
            {{ col.label }}
          </TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <template v-if="props.loading">
          <TableRow v-for="i in 5" :key="i">
            <TableCell v-for="col in props.columns" :key="col.key">
              <div class="h-4 w-full animate-pulse bg-muted rounded" />
            </TableCell>
          </TableRow>
        </template>
        <template v-else-if="props.data.length === 0">
          <TableRow>
            <TableCell :colspan="props.columns.length" class="h-24 text-center">
              No results.
            </TableCell>
          </TableRow>
        </template>
        <template v-else>
          <TableRow v-for="(row, idx) in props.data" :key="idx">
            <TableCell v-for="col in props.columns" :key="col.key" :class="col.class">
              <slot :name="`cell(${col.key})`" :row="row" :value="row[col.key]">
                {{ row[col.key] }}
              </slot>
            </TableCell>
          </TableRow>
        </template>
      </TableBody>
    </Table>
  </div>
</template>
