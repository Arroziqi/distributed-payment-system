import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { AXIOS_INSTANCE } from '@/api/client';
import { jwtDecode } from 'jwt-decode';

interface User {
  id: string;
  name: string;
  email: string;
}

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string | null>(null);
  const currentUser = ref<User | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const isAuthenticated = computed(() => {
    if (!accessToken.value) return false;
    try {
      const decoded: any = jwtDecode(accessToken.value);
      return decoded.exp > Date.now() / 1000;
    } catch {
      return false;
    }
  });

  const isLoggedIn = computed(() => !!currentUser.value && isAuthenticated.value);
  const userName = computed(() => currentUser.value?.name || '');
  const userEmail = computed(() => currentUser.value?.email || '');

  const login = async (tokens: { access_token: string, user_id?: string }) => {
    accessToken.value = tokens.access_token;
    await fetchCurrentUser();
  };

  const fetchCurrentUser = async () => {
    loading.value = true;
    error.value = null;
    try {
      // Using the new v1 endpoint
      const response = await AXIOS_INSTANCE.get('/api/v1/users/me');
      currentUser.value = response.data;
    } catch (err: any) {
      error.value = err.response?.data?.error || err.message || 'Failed to fetch user profile';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const refreshToken = async () => {
    try {
      // The refresh token is now in an HttpOnly cookie, so we don't need to send it in the body
      // but the backend still supports JSON body as fallback.
      // We'll just call the endpoint.
      const response = await AXIOS_INSTANCE.post('/auth/refresh', {}, {
        withCredentials: true // Important for cookies
      });
      accessToken.value = response.data.access_token;
      return response.data.access_token;
    } catch (err) {
      logout();
      throw err;
    }
  };

  const logout = async () => {
    try {
      await AXIOS_INSTANCE.post('/auth/logout', {}, { withCredentials: true });
    } catch (err) {
      console.error('Logout API failed', err);
    } finally {
      accessToken.value = null;
      currentUser.value = null;
      // We don't clear localStorage here, pinia-plugin-persistedstate will sync it
    }
  };

  const initializeAuth = async () => {
    if (accessToken.value && !isAuthenticated.value) {
      try {
        await refreshToken();
      } catch {
        logout();
      }
    } else if (accessToken.value && !currentUser.value) {
      await fetchCurrentUser();
    }
  };

  return {
    accessToken,
    currentUser,
    isAuthenticated,
    loading,
    error,
    isLoggedIn,
    userName,
    userEmail,
    login,
    logout,
    fetchCurrentUser,
    refreshToken,
    initializeAuth,
  };
}, {
  persist: {
    pick: ['accessToken', 'currentUser'],
  },
});
