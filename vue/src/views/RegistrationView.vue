<template>
  <div class="auth-page">
    <div class="auth-card">
      <h2>Регистрация</h2>
      <div class="auth-form">
        <div class="form-group">
          <label>Имя</label>
          <input v-model="name" type="text" placeholder="Ваше имя" />
        </div>
        <div class="form-group">
          <label>Email</label>
          <input v-model="email" type="email" placeholder="email@example.com" />
        </div>
        <div class="form-group">
          <label>Телефон</label>
          <input v-model="phone" type="tel" placeholder="+7 (999) 123-45-67" />
        </div>
        <div class="form-group">
          <label>Пароль</label>
          <input v-model="password" type="password" placeholder="Минимум 6 символов" />
        </div>
        <button class="btn-primary" @click="register" :disabled="loading">
          {{ loading ? 'Регистрация...' : 'Зарегистрироваться' }}
        </button>
      </div>
      <div class="auth-links">
        <router-link to="/login">Уже есть аккаунт? Войти</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { useAuthStore } from '../auth/stores/auth.store';
import AuthService from '../auth/services/auth.service';

export default {
  name: 'RegistrationView',
  data() {
    return {
      name: '',
      email: '',
      phone: '',
      password: '',
      loading: false,
    };
  },
  methods: {
    async register() {
      this.loading = true;
      try {
        const store = useAuthStore();
        const response = await AuthService.register({
          name: this.name,
          email: this.email,
          phone: this.phone,
          password: this.password,
        });
        store.user.token = response.token;
        store.user.name = response.user.name;
        store.user.email = response.user.email;
        store.user.id = response.user.id;
        store.isAuthenticated = true;
        localStorage.setItem('token', response.token);
        localStorage.setItem('user', JSON.stringify(response.user));
        this.$router.push('/');
      } catch (e) {
        alert('Ошибка регистрации');
      } finally {
        this.loading = false;
      }
    },
  },
};
</script>

<style scoped>
.auth-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60vh;
  width: 100%;
}
.auth-card {
  background: var(--light-dark-background-color);
  border-radius: 16px;
  padding: 32px;
  max-width: 400px;
  width: 100%;
}
.auth-card h2 { text-align: center; margin-bottom: 24px; }
.auth-form { display: flex; flex-direction: column; gap: 12px; }
.form-group { display: flex; flex-direction: column; gap: 4px; }
.form-group label { font-weight: 600; font-size: 0.9em; }
.form-group input {
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 14px;
}
.btn-primary {
  padding: 12px;
  background: #1976d2;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
}
.btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
.auth-links { text-align: center; margin-top: 16px; }
.auth-links a { color: #1976d2; text-decoration: none; font-size: 0.9em; }
</style>
