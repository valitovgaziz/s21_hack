import { defineStore } from 'pinia';
import { reactive, ref } from 'vue';
import AuthService from '../services/auth.service.js';
import { jwtDecode } from 'jwt-decode';

export const useAuthStore = defineStore('auth', () => {
  const user = reactive({name: '', email: '', id: 0, token: '', role: '', phone: ''});
  const isAuthenticated = ref(false);

  const initFromStorage = () => {
    const storedToken = localStorage.getItem('token');
    const storedUser = localStorage.getItem('user');

    if (storedToken && storedUser) {
      try {
        user.token = storedToken;
        const userData = JSON.parse(storedUser);
        user.name = userData.name || '';
        user.email = userData.email || '';
        user.id = userData.id || 0;
        user.role = userData.role || '';
        user.phone = userData.phone || '';
        isAuthenticated.value = true;
      } catch (error) {
        logout();
      }
    }
  };

  initFromStorage();

  const register = async (userData) => {
    try {
      const response = await AuthService.register(userData);
      if (response.token) {
        const decodedToken = jwtDecode(response.token);
        user.name = response.user?.name || userData.name;
        user.id = response.user?.id || 0;
        user.email = response.user?.email || userData.email;
        user.role = response.user?.role || 'employee';
        user.phone = response.user?.phone || '';
        isAuthenticated.value = true;
        user.token = response.token;

        localStorage.setItem('token', response.token);
        localStorage.setItem('user', JSON.stringify(response.user));
      }
      return response;
    } catch (error) {
      throw error;
    }
  };

  const login = async (credentials) => {
    try {
      const response = await AuthService.login(credentials);
      user.name = response.user?.name || '';
      user.id = response.user?.id || 0;
      user.email = response.user?.email || credentials.email;
      user.role = response.user?.role || 'employee';
      user.phone = response.user?.phone || '';
      isAuthenticated.value = true;
      user.token = response.token;

      localStorage.setItem('token', response.token);
      localStorage.setItem('user', JSON.stringify(response.user));
    } catch (error) {
      throw error;
    }
  };

  const logout = () => {
    isAuthenticated.value = false;
    user.name = '';
    user.token = '';
    user.email = '';
    user.id = 0;
    user.role = '';
    user.phone = '';

    localStorage.removeItem('token');
    localStorage.removeItem('user');
  };

  const checkAuth = async () => {
    try {
      const token = user.token || localStorage.getItem('token');
      if (token) {
        try {
          const response = await AuthService.checkAuth(token);
          user.name = response.user?.name || user.name;
          user.id = response.user?.id || user.id;
          user.email = response.user?.email || user.email;
          user.role = response.user?.role || user.role;
          user.phone = response.user?.phone || user.phone;
          isAuthenticated.value = true;
        } catch (error) {
          logout();
        }
      }
    } catch (error) {
      throw error;
    }
  };

  return { user, isAuthenticated, register, login, logout, checkAuth };
});
