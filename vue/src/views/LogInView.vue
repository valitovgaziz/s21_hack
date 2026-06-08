<template>
  <div class="auth-page">
    <div class="auth-card">
      <h2>Вход в PulseHR</h2>

      <div v-if="step === 'phone'" class="auth-form">
        <div class="form-group">
          <label>Номер телефона</label>
          <input v-model="phone" type="tel" placeholder="+7 (999) 123-45-67" />
        </div>
        <button class="btn-primary" @click="sendOTP" :disabled="sending">
          {{ sending ? 'Отправка...' : 'Получить код' }}
        </button>
        <div class="auth-divider">или</div>
        <button class="btn-secondary" @click="step = 'email'">Войти по email</button>
      </div>

      <div v-else-if="step === 'otp'" class="auth-form">
        <p>Код отправлен на {{ phone }}</p>
        <div class="form-group">
          <label>Код из SMS</label>
          <input v-model="otpCode" type="text" maxlength="6" placeholder="000000" />
        </div>
        <button class="btn-primary" @click="verifyOTP" :disabled="verifying">
          {{ verifying ? 'Проверка...' : 'Подтвердить' }}
        </button>
      </div>

      <div v-else class="auth-form">
        <div class="form-group">
          <label>Email</label>
          <input v-model="email" type="email" placeholder="email@example.com" />
        </div>
        <div class="form-group">
          <label>Пароль</label>
          <input v-model="password" type="password" placeholder="Пароль" />
        </div>
        <button class="btn-primary" @click="loginEmail" :disabled="loading">
          {{ loading ? 'Вход...' : 'Войти' }}
        </button>
        <div class="auth-divider">или</div>
        <button class="btn-secondary" @click="step = 'phone'">Войти по телефону</button>
      </div>

      <div class="auth-links">
        <router-link to="/registration">Нет аккаунта? Зарегистрироваться</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { useAuthStore } from '../auth/stores/auth.store';
import AuthService from '../auth/services/auth.service';

export default {
  name: 'LoginView',
  data() {
    return {
      step: 'phone',
      phone: '',
      otpCode: '',
      email: '',
      password: '',
      sending: false,
      verifying: false,
      loading: false,
    };
  },
  methods: {
    async sendOTP() {
      if (!this.phone) return;
      this.sending = true;
      try {
        await AuthService.loginPhone(this.phone);
        this.step = 'otp';
      } catch (e) {
        alert('Ошибка отправки кода');
      } finally {
        this.sending = false;
      }
    },
    async verifyOTP() {
      if (!this.otpCode || this.otpCode.length !== 6) return;
      this.verifying = true;
      try {
        const response = await AuthService.verifyOTP(this.phone, this.otpCode);
        const store = useAuthStore();
        store.user.token = response.token;
        store.user.name = response.user.name;
        store.user.email = response.user.email;
        store.user.id = response.user.id;
        store.isAuthenticated = true;
        localStorage.setItem('token', response.token);
        localStorage.setItem('user', JSON.stringify(response.user));
        this.$router.push('/');
      } catch (e) {
        alert('Неверный код');
      } finally {
        this.verifying = false;
      }
    },
    async loginEmail() {
      this.loading = true;
      try {
        const store = useAuthStore();
        await store.login({ email: this.email, password: this.password });
        this.$router.push('/');
      } catch (e) {
        alert('Ошибка входа');
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
.btn-secondary {
  padding: 12px;
  background: #e0e0e0;
  color: #333;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
}
.btn-primary:disabled, .btn-secondary:disabled { opacity: 0.6; cursor: not-allowed; }
.auth-divider { text-align: center; color: #888; font-size: 0.9em; }
.auth-links { text-align: center; margin-top: 16px; }
.auth-links a { color: #1976d2; text-decoration: none; font-size: 0.9em; }
</style>
