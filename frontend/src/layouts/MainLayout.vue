<script setup lang="ts">
import { ref, computed } from 'vue';
import SidebarNavigation from '@/components/organisms/SidebarNavigation.vue';
import TopNavbar from '@/components/organisms/TopNavbar.vue';
import { useAuth } from '@/features/auth/useAuth';

const isSidebarOpen = ref(true);
const { logout, user: authUser } = useAuth();

const user = computed(() => ({
  name: authUser.value?.name || 'User',
  email: authUser.value?.email || '',
  avatar: 'https://github.com/shadcn.png'
}));

const handleToggleSidebar = () => {
  isSidebarOpen.value = !isSidebarOpen.value;
};

const handleLogout = async () => {
  await logout();
};
</script>

<template>
  <div class="flex h-screen bg-background text-foreground overflow-hidden">
    <!-- Sidebar -->
    <SidebarNavigation
      :class="['fixed inset-y-0 left-0 z-50 transform md:relative md:translate-x-0 transition-transform duration-300 ease-in-out', !isSidebarOpen ? '-translate-x-full' : 'translate-x-0']"
    />
    
    <!-- Main Content -->
    <div class="flex flex-1 flex-col overflow-hidden">
      <TopNavbar
        :user="user"
        @toggle-sidebar="handleToggleSidebar"
        @logout="handleLogout"
      />
      
      <main class="flex-1 overflow-y-auto p-4 md:p-6">
        <div class="mx-auto max-w-7xl">
          <router-view />
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped>
/* Responsive tweaks if needed */
</style>
