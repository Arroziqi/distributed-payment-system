import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { usePostAuthLogin, usePostAuthRegister, usePutAuthMe } from '@/api/generated/auth/auth';
import { toast } from 'vue-sonner';
import { useAuthStore } from '@/stores/auth.store';

export const useAuth = () => {
  const router = useRouter();
  const authStore = useAuthStore();
  
  const loginMutation = usePostAuthLogin();
  const registerMutation = usePostAuthRegister();
  const updateMeMutation = usePutAuthMe();

  const login = async (data: any) => {
    try {
      const tokens: any = await loginMutation.mutateAsync({ data });
      if (tokens && tokens.access_token) {
        await authStore.login(tokens);
        toast.success('Login successful');
        router.push('/');
      }
    } catch (error: any) {
      const errorMsg = error.response?.data?.error || error.message || 'Login failed';
      toast.error(errorMsg);
      throw error;
    }
  };

  const register = async (data: any) => {
    try {
      const response: any = await registerMutation.mutateAsync({ data });
      if (response) {
        toast.success('Registration successful. Please login.');
        router.push('/login');
      }
    } catch (error: any) {
      const errorMsg = error.response?.data?.error || error.message || 'Registration failed';
      toast.error(errorMsg);
      throw error;
    }
  };

  const logout = async () => {
    try {
      await authStore.logout();
    } catch (error) {
      console.error('Logout failed', error);
    } finally {
      router.push('/login');
    }
  };

  const updateProfile = async (data: { name?: string; email?: string }) => {
    try {
      const updatedUser: any = await updateMeMutation.mutateAsync({ data });
      await authStore.fetchCurrentUser(); // Refresh store data
      return updatedUser;
    } catch (error: any) {
      const errorMsg = error.response?.data?.error || error.message || 'Update failed';
      toast.error(errorMsg);
      throw error;
    }
  };

  return {
    user: computed(() => authStore.currentUser),
    isAuthenticated: computed(() => authStore.isAuthenticated),
    login,
    register,
    logout,
    updateProfile,
    loading: computed(() => loginMutation.isPending.value || registerMutation.isPending.value || updateMeMutation.isPending.value || authStore.loading),
  };
}
