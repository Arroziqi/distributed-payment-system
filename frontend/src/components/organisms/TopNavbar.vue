<script setup lang="ts">
import { ref } from 'vue';
import ProfileDropdown from '@/components/molecules/ProfileDropdown.vue';
import SearchToolbar from '@/components/molecules/SearchToolbar.vue';
import BaseIcon from '@/components/atoms/BaseIcon.vue';
import BaseButton from '@/components/atoms/BaseButton.vue';
import { useRouter } from 'vue-router';

interface Props {
  user: {
    name: string;
    email: string;
    avatar?: string;
  };
}

defineProps<Props>();
defineEmits(['toggle-sidebar', 'logout']);

const router = useRouter();
const searchQuery = ref('');

const handleProfile = () => {
  router.push('/profile');
};

const handleSettings = () => {
  router.push('/settings');
};
</script>

<template>
  <header class="sticky top-0 z-30 flex h-16 w-full items-center justify-between border-b bg-background/95 px-4 backdrop-blur md:px-6">
    <div class="flex items-center gap-4">
      <BaseButton variant="ghost" size="icon" class="md:hidden" @click="$emit('toggle-sidebar')">
        <BaseIcon name="Menu" class="h-5 w-5" />
      </BaseButton>
      <SearchToolbar v-model="searchQuery" class="hidden sm:flex" />
    </div>
    
    <div class="flex items-center space-x-4">
      <BaseButton variant="ghost" size="icon" class="relative">
        <BaseIcon name="Bell" class="h-5 w-5" />
        <span class="absolute top-2 right-2 h-2 w-2 rounded-full bg-primary" />
      </BaseButton>
      
      <div class="h-6 w-px bg-border mx-2" />
      
      <ProfileDropdown
        :user="user"
        @logout="$emit('logout')"
        @profile="handleProfile"
        @settings="handleSettings"
      />
    </div>
  </header>
</template>
