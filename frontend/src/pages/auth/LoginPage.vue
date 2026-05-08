<script setup lang="ts">
import { ref } from 'vue';
import LoginForm from '@/components/molecules/LoginForm.vue';
import { useAuth } from '@/features/auth/useAuth';

const { login, loading } = useAuth();
const error = ref<string | null>(null);

const handleLogin = async (data: any) => {
  error.value = null;
  try {
    await login(data);
  } catch (err: any) {
    error.value = err.response?.data?.error || err.message || 'Login failed';
  }
};
</script>

<template>
  <div class="space-y-6">
    <div class="space-y-2 text-center">
      <h1 class="text-2xl font-semibold tracking-tight">Login to your account</h1>
      <p class="text-sm text-muted-foreground">
        Enter your email below to login to your account
      </p>
    </div>
    
    <LoginForm :loading="loading" :error="error" @submit="handleLogin" />
    
    <div class="text-center text-sm">
      Don't have an account?
      <router-link to="/register" class="underline underline-offset-4 hover:text-primary">
        Sign up
      </router-link>
    </div>
  </div>
</template>
