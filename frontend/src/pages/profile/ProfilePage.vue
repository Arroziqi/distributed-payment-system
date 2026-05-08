<script setup lang="ts">
import { reactive, onMounted, watch, ref } from 'vue';
import { useAuth } from '@/features/auth/useAuth';
import BaseIcon from '@/components/atoms/BaseIcon.vue';
import BaseButton from '@/components/atoms/BaseButton.vue';
import BaseAvatar from '@/components/atoms/BaseAvatar.vue';
import BaseAlert from '@/components/atoms/BaseAlert.vue';

const { user, updateProfile, loading } = useAuth();
const showSuccess = ref(false);
const successMessage = ref('');

const form = reactive({
  name: '',
  email: '',
});

onMounted(() => {
  if (user.value) {
    form.name = user.value.name || '';
    form.email = user.value.email || '';
  }
});

watch(user, (newUser) => {
  if (newUser) {
    form.name = newUser.name || '';
    form.email = newUser.email || '';
  }
}, { deep: true });

const handleSubmit = async () => {
  try {
    await updateProfile({
      name: form.name,
      email: form.email,
    });
    successMessage.value = 'Profile updated successfully';
    showSuccess.value = true;
  } catch (error) {
    // Error is handled by toast in useAuth
  }
};
</script>

<template>
  <div class="max-w-4xl mx-auto py-8 relative">
    <!-- Success Alert (Floating) -->
    <Transition
      enter-active-class="transform transition ease-out duration-300"
      enter-from-class="-translate-y-4 opacity-0"
      enter-to-class="translate-y-0 opacity-100"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="translate-y-0 opacity-100"
      leave-to-class="-translate-y-4 opacity-0"
    >
      <div v-if="showSuccess" class="absolute top-0 left-1/2 -translate-x-1/2 z-50 w-full max-w-md px-4 mt-4">
        <BaseAlert
          variant="success"
          title="Success"
          :description="successMessage"
          show-progress
          :duration="5000"
          show-close
          @close="showSuccess = false"
          class="shadow-2xl border-green-500/30 backdrop-blur-md bg-white/90 dark:bg-slate-900/90"
        />
      </div>
    </Transition>

    <div class="mb-8">
      <h1 class="text-3xl font-bold tracking-tight">Account Settings</h1>
      <p class="text-muted-foreground mt-1">Manage your account information and preferences.</p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
      <!-- Profile Preview -->
      <div class="space-y-6">
        <div class="bg-card border rounded-xl p-6 shadow-sm flex flex-col items-center text-center">
          <div class="relative group">
            <BaseAvatar 
              :src="'https://github.com/shadcn.png'" 
              :fallback="user?.name?.charAt(0) || 'U'" 
              class="h-32 w-32 border-4 border-background shadow-xl" 
            />
            <button class="absolute bottom-0 right-0 p-2 bg-primary text-primary-foreground rounded-full shadow-lg hover:bg-primary/90 transition-colors">
              <BaseIcon name="Camera" class="h-4 w-4" />
            </button>
          </div>
          
          <div class="mt-4">
            <h2 class="text-xl font-semibold">{{ user?.name }}</h2>
            <p class="text-sm text-muted-foreground">{{ user?.email }}</p>
          </div>
          
          <div class="mt-6 w-full pt-6 border-t flex justify-around text-center">
            <div>
              <p class="text-xs text-muted-foreground font-medium uppercase tracking-wider">Status</p>
              <p class="text-sm font-semibold text-green-500">Active</p>
            </div>
            <div class="border-x px-4">
              <p class="text-xs text-muted-foreground font-medium uppercase tracking-wider">Plan</p>
              <p class="text-sm font-semibold">Premium</p>
            </div>
          </div>
        </div>

        <div class="bg-muted/30 border border-dashed rounded-xl p-6">
          <div class="flex items-center gap-3 mb-3">
            <BaseIcon name="ShieldCheck" class="h-5 w-5 text-primary" />
            <h3 class="font-semibold">Security Note</h3>
          </div>
          <p class="text-xs text-muted-foreground leading-relaxed">
            Changing your email address will require you to use the new email for future logins. Make sure you have access to it.
          </p>
        </div>
      </div>

      <!-- Edit Form -->
      <div class="md:col-span-2 space-y-6">
        <div class="bg-card border rounded-xl shadow-sm overflow-hidden">
          <div class="px-6 py-4 border-b bg-muted/50">
            <h3 class="font-semibold">Personal Information</h3>
          </div>
          
          <form @submit.prevent="handleSubmit" class="p-6 space-y-6">
            <div class="grid grid-cols-1 gap-6">
              <div class="space-y-2">
                <label for="name" class="text-sm font-medium">Full Name</label>
                <div class="relative">
                  <BaseIcon name="User" class="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                  <input
                    id="name"
                    v-model="form.name"
                    type="text"
                    placeholder="John Doe"
                    class="flex h-10 w-full rounded-md border border-input bg-background px-10 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                    required
                  />
                </div>
              </div>

              <div class="space-y-2">
                <label for="email" class="text-sm font-medium">Email Address</label>
                <div class="relative">
                  <BaseIcon name="Mail" class="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                  <input
                    id="email"
                    v-model="form.email"
                    type="email"
                    placeholder="john@example.com"
                    class="flex h-10 w-full rounded-md border border-input bg-background px-10 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                    required
                  />
                </div>
              </div>
            </div>

            <div class="flex justify-end pt-4">
              <BaseButton type="submit" :disabled="loading" class="min-w-[120px]">
                <BaseIcon v-if="loading" name="Loader2" class="mr-2 h-4 w-4 animate-spin" />
                <span>{{ loading ? 'Saving...' : 'Save Changes' }}</span>
              </BaseButton>
            </div>
          </form>
        </div>

        <div class="bg-card border rounded-xl shadow-sm overflow-hidden border-destructive/20">
          <div class="px-6 py-4 border-b bg-destructive/5">
            <h3 class="font-semibold text-destructive">Danger Zone</h3>
          </div>
          <div class="p-6 flex items-center justify-between">
            <div>
              <p class="font-medium">Delete Account</p>
              <p class="text-sm text-muted-foreground">Permanently remove your account and all of your data.</p>
            </div>
            <BaseButton variant="destructive" size="sm">
              Delete
            </BaseButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
