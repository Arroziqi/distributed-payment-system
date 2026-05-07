<script setup lang="ts">
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import BaseAvatar from '@/components/atoms/BaseAvatar.vue';
import BaseIcon from '@/components/atoms/BaseIcon.vue';

interface User {
  name: string;
  email: string;
  avatar?: string;
}

interface Props {
  user: User;
}

defineProps<Props>();
defineEmits(['logout', 'profile', 'settings']);
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <button class="flex items-center space-x-2 focus:outline-none">
        <BaseAvatar :src="user.avatar" :fallback="user.name.charAt(0)" class="h-8 w-8" />
        <div class="hidden md:block text-left">
          <p class="text-xs font-medium">{{ user.name }}</p>
          <p class="text-[10px] text-muted-foreground">{{ user.email }}</p>
        </div>
        <BaseIcon name="ChevronDown" class="h-4 w-4 text-muted-foreground" />
      </button>
    </DropdownMenuTrigger>
    <DropdownMenuContent class="w-56" align="end">
      <DropdownMenuLabel>My Account</DropdownMenuLabel>
      <DropdownMenuSeparator />
      <DropdownMenuGroup>
        <DropdownMenuItem @click="$emit('profile')">
          <BaseIcon name="User" class="mr-2 h-4 w-4" />
          <span>Profile</span>
        </DropdownMenuItem>
        <DropdownMenuItem @click="$emit('settings')">
          <BaseIcon name="Settings" class="mr-2 h-4 w-4" />
          <span>Settings</span>
        </DropdownMenuItem>
      </DropdownMenuGroup>
      <DropdownMenuSeparator />
      <DropdownMenuItem @click="$emit('logout')" class="text-destructive focus:text-destructive">
        <BaseIcon name="LogOut" class="mr-2 h-4 w-4" />
        <span>Log out</span>
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
