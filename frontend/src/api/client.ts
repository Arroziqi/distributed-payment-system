import axios from 'axios';
import type { AxiosRequestConfig, AxiosResponse } from 'axios';

const AUTH_API = import.meta.env.VITE_AUTH_API || 'http://localhost:8081';
const WALLET_API = import.meta.env.VITE_WALLET_API || 'http://localhost:8082';
const TRANSACTION_API = import.meta.env.VITE_TRANSACTION_API || 'http://localhost:8083';
const NOTIFICATION_API = import.meta.env.VITE_NOTIFICATION_API || 'http://localhost:8084';

export const AXIOS_INSTANCE = axios.create();

import { useAuthStore } from '@/stores/auth.store';

let isRefreshing = false;
let failedQueue: any[] = [];

const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

AXIOS_INSTANCE.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore();
    const token = authStore.accessToken;
    
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // Dynamic Base URL based on path
    if (config.url?.startsWith('/auth')) {
      config.baseURL = AUTH_API;
    } else if (config.url?.startsWith('/wallet') || config.url?.startsWith('/wallets')) {
      config.baseURL = WALLET_API;
    } else if (config.url?.startsWith('/api/v1/users')) {
      config.baseURL = AUTH_API;
    } else if (config.url?.startsWith('/transactions')) {
      config.baseURL = TRANSACTION_API;
    } else if (config.url?.startsWith('/notifications')) {
      config.baseURL = NOTIFICATION_API;
    }

    return config;
  },
  (error) => Promise.reject(error)
);

AXIOS_INSTANCE.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    const authStore = useAuthStore();

    const isAuthRequest = originalRequest.url?.includes('/auth/login') || 
                         originalRequest.url?.includes('/auth/register') ||
                         originalRequest.url?.includes('/auth/refresh');
    
    if (error.response?.status === 401 && !isAuthRequest && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then((token) => {
            originalRequest.headers.Authorization = `Bearer ${token}`;
            return AXIOS_INSTANCE(originalRequest);
          })
          .catch((err) => Promise.reject(err));
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        const newToken = await authStore.refreshToken();
        processQueue(null, newToken);
        originalRequest.headers.Authorization = `Bearer ${newToken}`;
        return AXIOS_INSTANCE(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError, null);
        authStore.logout();
        
        if (window.location.pathname !== '/login') {
          window.location.href = '/login';
        }
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }
    return Promise.reject(error);
  }
);

export const customInstance = <T>(config: AxiosRequestConfig): Promise<T> => {
  return AXIOS_INSTANCE(config).then((res: AxiosResponse<T>) => res.data);
};

export default customInstance;
