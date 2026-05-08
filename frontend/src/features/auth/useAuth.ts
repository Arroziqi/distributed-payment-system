import { ref, computed, reactive, toRefs } from 'vue';
import { useRouter } from 'vue-router';
import { usePostAuthLogin, usePostAuthRegister, usePostAuthLogout, usePutAuthMe } from '@/api/generated/auth/auth';
import { toast } from 'vue-sonner';
import { jwtDecode } from 'jwt-decode';
import { AXIOS_INSTANCE } from '@/api/client';

const state = reactive({
  user: JSON.parse(localStorage.getItem('user') || 'null'),
});

export const useAuth = () => {
  const router = useRouter();
  const { user } = toRefs(state);
  
  const isTokenExpired = (token: string) => {
    try {
      const decoded: any = jwtDecode(token);
      const currentTime = Date.now() / 1000;
      return decoded.exp < currentTime;
    } catch {
      return true;
    }
  };

  const isAuthenticated = computed(() => {
    const token = localStorage.getItem('access_token');
    return !!user.value && !!token && !isTokenExpired(token);
  });

  const loginMutation = usePostAuthLogin();
  const registerMutation = usePostAuthRegister();
  const logoutMutation = usePostAuthLogout();
  const updateMeMutation = usePutAuthMe();

  const login = async (data: any) => {
    try {
      const tokens: any = await loginMutation.mutateAsync({ data });
      if (tokens && tokens.access_token) {
        localStorage.setItem('access_token', tokens.access_token || '');
        localStorage.setItem('refresh_token', tokens.refresh_token || '');
        
        // Decode token to get user ID
        const decoded: any = jwtDecode(tokens.access_token || '');
        const userId = decoded.sub || decoded.user_id || tokens.user_id;
        
        localStorage.setItem('user_id', userId);

        // Fetch actual profile
        const profile: any = await AXIOS_INSTANCE.get(`${import.meta.env.VITE_AUTH_API || 'http://localhost:8081'}/auth/me`, {
          headers: { Authorization: `Bearer ${tokens.access_token}` }
        }).then(res => res.data);

        const userData = { 
          id: userId,
          email: profile.email,
          name: profile.name
        };
        
        user.value = userData;
        localStorage.setItem('user', JSON.stringify(userData));
        
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
      // Depending on backend, response might be the body or { data }
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
    const refreshToken = localStorage.getItem('refresh_token');
    try {
      if (refreshToken) {
        await logoutMutation.mutateAsync({ data: { refresh_token: refreshToken } });
      }
    } catch (error) {
      console.error('Logout API call failed', error);
    } finally {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      localStorage.removeItem('user');
      localStorage.removeItem('user_id');
      user.value = null;
      router.push('/login');
    }
  };

  const updateProfile = async (data: { name?: string; email?: string }) => {
    try {
      const updatedUser: any = await updateMeMutation.mutateAsync({ data });
      const userData = { 
        ...user.value,
        name: updatedUser.name,
        email: updatedUser.email
      };
      state.user = userData;
      localStorage.setItem('user', JSON.stringify(userData));
      
      return updatedUser;
    } catch (error: any) {
      const errorMsg = error.response?.data?.error || error.message || 'Update failed';
      toast.error(errorMsg);
      throw error;
    }
  };

  return {
    user,
    isAuthenticated,
    login,
    register,
    logout,
    updateProfile,
    loading: computed(() => loginMutation.isPending.value || registerMutation.isPending.value || updateMeMutation.isPending.value),
  };
}
