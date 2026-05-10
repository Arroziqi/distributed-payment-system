import { createRouter, createWebHistory } from 'vue-router';
import AuthLayout from '@/layouts/AuthLayout.vue';
import MainLayout from '@/layouts/MainLayout.vue';
import { jwtDecode } from 'jwt-decode';

const routes = [
  {
    path: '/',
    component: MainLayout,
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/pages/dashboard/DashboardPage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'wallet',
        name: 'Wallet',
        component: () => import('@/pages/wallet/WalletPage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'wallet/transfer',
        name: 'Transfer',
        component: () => import('@/pages/wallet/TransferPage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'wallet/topup',
        name: 'TopUp',
        component: () => import('@/pages/wallet/TopUpPage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'transactions',
        name: 'Transactions',
        component: () => import('@/pages/transactions/TransactionsPage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'notifications',
        name: 'Notifications',
        component: () => import('@/pages/notifications/NotificationsPage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'observability',
        name: 'Observability',
        component: () => import('@/pages/observability/ObservabilityPage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/pages/profile/ProfilePage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('@/pages/settings/SettingsPage.vue'),
        meta: { requiresAuth: true }
      }
    ]
  },
  {
    path: '/auth',
    component: AuthLayout,
    children: [
      {
        path: 'login',
        name: 'Login',
        alias: '/login',
        component: () => import('@/pages/auth/LoginPage.vue'),
        meta: { guest: true }
      },
      {
        path: 'register',
        name: 'Register',
        alias: '/register',
        component: () => import('@/pages/auth/RegisterPage.vue'),
        meta: { guest: true }
      }
    ]
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

import { useAuthStore } from '@/stores/auth.store';

router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore();
  
  // Ensure auth is initialized (check tokens, maybe fetch user)
  if (authStore.accessToken && !authStore.currentUser) {
    await authStore.initializeAuth();
  }

  const isAuthenticated = authStore.isAuthenticated;

  if (to.meta.requiresAuth && !isAuthenticated) {
    authStore.logout(); // Ensure state is clean
    next('/login');
  } else if (to.meta.guest && isAuthenticated) {
    next('/');
  } else {
    next();
  }
});

export default router;
