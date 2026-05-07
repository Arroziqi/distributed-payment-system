import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { usePostAuthLogin, usePostAuthRegister, usePostAuthLogout } from '@/api/generated/auth/auth';
import { toast } from 'vue-sonner';
import { jwtDecode } from 'jwt-decode';

export const useAuth = () => {
  const router = useRouter();
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'));
  const isAuthenticated = computed(() => !!user.value);

  const loginMutation = usePostAuthLogin();
  const registerMutation = usePostAuthRegister();
  const logoutMutation = usePostAuthLogout();

  const login = async (data: any) => {
    try {
      const tokens = await loginMutation.mutateAsync({ data });
      if (tokens && tokens.access_token) {
        localStorage.setItem('access_token', tokens.access_token || '');
        localStorage.setItem('refresh_token', tokens.refresh_token || '');
        
        // Decode token to get user ID
        const decoded: any = jwtDecode(tokens.access_token || '');
        const userId = decoded.sub || decoded.user_id || (tokens as any).user_id;
        
        const userData = { 
          id: userId,
          email: data.email,
          name: data.email.split('@')[0]
        };
        
        user.value = userData;
        localStorage.setItem('user', JSON.stringify(userData));
        localStorage.setItem('user_id', userId);
        
        toast.success('Login successful');
        router.push('/');
      }
    } catch (error: any) {
      const errorMsg = error.response?.data?.error || error.message || 'Login failed';
      toast.error(errorMsg);
    }
  };

  const register = async (data: any) => {
    try {
      const response: any = await registerMutation.mutateAsync({ data });
      // Depending on backend, response might be the body or { data }
      if (response) {
        toast.success('Registration successful. Please login.');
        router.push('/login');
      }
    } catch (error: any) {
      const errorMsg = error.response?.data?.error || error.message || 'Registration failed';
      toast.error(errorMsg);
    }
  };

  const logout = async () => {
    const refreshToken = localStorage.getItem('refresh_token');
    try {
      if (refreshToken) {
        await logoutMutation.mutateAsync({ data: { refresh_token: refreshToken } });
      }
    } finally {
      localStorage.clear();
      user.value = null;
      router.push('/login');
    }
  };

  return {
    user,
    isAuthenticated,
    login,
    register,
    logout,
    loading: loginMutation.isPending || registerMutation.isPending,
  };
};
