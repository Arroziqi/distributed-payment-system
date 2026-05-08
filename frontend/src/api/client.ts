import axios from 'axios';
import type { AxiosRequestConfig, AxiosResponse } from 'axios';

const AUTH_API = import.meta.env.VITE_AUTH_API || 'http://localhost:8081';
const WALLET_API = import.meta.env.VITE_WALLET_API || 'http://localhost:8082';
const TRANSACTION_API = import.meta.env.VITE_TRANSACTION_API || 'http://localhost:8083';
const NOTIFICATION_API = import.meta.env.VITE_NOTIFICATION_API || 'http://localhost:8084';

export const AXIOS_INSTANCE = axios.create();

AXIOS_INSTANCE.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // Dynamic Base URL based on path
    if (config.url?.startsWith('/auth')) {
      config.baseURL = AUTH_API;
    } else if (config.url?.startsWith('/wallet') || config.url?.startsWith('/wallets')) {
      config.baseURL = WALLET_API;
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
  (error) => {
    const isAuthRequest = error.config?.url?.includes('/auth/login') || error.config?.url?.includes('/auth/register');
    
    if (error.response?.status === 401 && !isAuthRequest) {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      localStorage.removeItem('user');
      localStorage.removeItem('user_id');
      
      // Only redirect if not already on login page to avoid instant reloads
      if (window.location.pathname !== '/login') {
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

export const customInstance = <T>(config: AxiosRequestConfig): Promise<T> => {
  return AXIOS_INSTANCE(config).then((res: AxiosResponse<T>) => res.data);
};

export default customInstance;
