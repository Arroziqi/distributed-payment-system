import { createRouter, createWebHistory } from 'vue-router';
import AuthLayout from '@/layouts/AuthLayout.vue';
import MainLayout from '@/layouts/MainLayout.vue';

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

router.beforeEach((to, _from, next) => {
  const isAuthenticated = !!localStorage.getItem('access_token');

  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login');
  } else if (to.meta.guest && isAuthenticated) {
    next('/');
  } else {
    next();
  }
});

export default router;
