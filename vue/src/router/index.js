import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../auth/stores/auth.store';

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('../views/HomeView.vue'),
  },
  {
    path: '/about',
    name: 'about',
    component: () => import('../views/AboutView.vue'),
  },
  {
    path: '/profile',
    name: 'profile',
    component: () => import('../views/ProfileView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/support',
    name: 'support',
    component: () => import('../views/SupportView.vue'),
  },
  {
    path: '/feetback',
    name: 'feetback',
    component: () => import('../views/FeetbackView.vue'),
  },
  {
    path: '/results',
    name: 'results',
    component: () => import('../views/ResultsView.vue'),
  },
  {
    path: '/settings',
    name: 'settings',
    component: () => import('../views/SettingsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/registration',
    name: 'registration',
    component: () => import('../views/RegistrationView.vue')
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/LogInView.vue')
  },
  {
    path: '/restObject',
    name: 'restObject',
    component: () => import('../views/RestObjectView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/surveys',
    name: 'surveys',
    component: () => import('../views/surveys/SurveyListView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/surveys/new',
    name: 'survey-new',
    component: () => import('../views/surveys/SurveyBuilderView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/surveys/:id/edit',
    name: 'survey-edit',
    component: () => import('../views/surveys/SurveyBuilderView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/survey/:id',
    name: 'survey-take',
    component: () => import('../views/surveys/SurveyTakeView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/analytics/:id',
    name: 'survey-analytics',
    component: () => import('../views/surveys/SurveyAnalyticsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/notifications/settings',
    name: 'notification-settings',
    component: () => import('../views/surveys/NotificationSettingsView.vue'),
    meta: { requiresAuth: true }
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore();
  await authStore.checkAuth();

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login');
  } else {
    next();
  }
});

export default router
