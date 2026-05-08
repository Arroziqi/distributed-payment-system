<script setup lang="ts">
import { ref } from 'vue';
import RegisterForm from '@/components/molecules/RegisterForm.vue';
import { useAuth } from '@/features/auth/useAuth';

const { register, loading } = useAuth();
const error = ref<string | null>(null);

const handleRegister = async (data: any) => {
  error.value = null;
  if (data.password !== data.confirmPassword) {
    error.value = "Passwords do not match";
    return;
  }
  try {
    await register({
      name: data.name,
      email: data.email,
      password: data.password
    });
  } catch (err: any) {
    error.value = err.response?.data?.error || err.message || 'Registration failed';
  }
};
</script>

<template>
  <div class="space-y-6">
    <div class="space-y-2 text-center">
      <h1 class="text-2xl font-semibold tracking-tight">Create an account</h1>
      <p class="text-sm text-muted-foreground">
        Enter your details below to create your account
      </p>
    </div>
    
    <RegisterForm :loading="loading" :error="error" @submit="handleRegister" />
    
    <div class="text-center text-sm">
      Already have an account?
      <router-link to="/login" class="underline underline-offset-4 hover:text-primary">
        Login
      </router-link>
    </div>
  </div>
</template>
