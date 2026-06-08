<template>
  <div class="notif-settings">
    <h2>Настройки уведомлений</h2>

    <div v-if="loading" class="loading">Загрузка...</div>

    <div v-else class="settings-form">
      <div class="section">
        <h3>Каналы уведомлений</h3>
        <div class="channel-toggle">
          <label>
            <input type="checkbox" v-model="prefs.webPushEnabled" />
            Web Push (браузер)
          </label>
        </div>
        <div class="channel-toggle">
          <label>
            <input type="checkbox" v-model="prefs.smsEnabled" />
            SMS
          </label>
        </div>
        <div class="channel-toggle">
          <label>
            <input type="checkbox" v-model="prefs.emailEnabled" />
            E-mail
          </label>
        </div>
      </div>

      <div class="section">
        <h3>Предпочтительное время</h3>
        <select v-model="prefs.preferredTime">
          <option value="morning">Утром (9:00-12:00)</option>
          <option value="day">Днём (12:00-18:00)</option>
          <option value="evening">Вечером (18:00-21:00)</option>
        </select>
      </div>

      <div class="section">
        <h3>Подписки на push-уведомления</h3>
        <div v-if="subscriptions.length === 0" class="no-subs">Нет активных подписок</div>
        <div v-for="sub in subscriptions" :key="sub.id" class="sub-item">
          <span>{{ sub.deviceName || 'Устройство' }}</span>
          <button class="btn-danger" @click="removeSubscription(sub.id)">Отписаться</button>
        </div>
      </div>

      <div class="section">
        <h3>Push-уведомления</h3>
        <p v-if="pushSupported && !pushSubscribed">
          <button class="btn-primary" @click="subscribePush">Включить push-уведомления</button>
        </p>
        <p v-else-if="pushSubscribed" class="push-active">
          Push-уведомления активны
        </p>
        <p v-else class="push-unsupported">
          Push-уведомления не поддерживаются вашим браузером
        </p>
      </div>

      <button class="btn-primary" @click="savePreferences" :disabled="saving">
        {{ saving ? 'Сохранение...' : 'Сохранить настройки' }}
      </button>
    </div>
  </div>
</template>

<script>
import { notificationAPI } from '../../api/notifications';

export default {
  name: 'NotificationSettingsView',
  data() {
    return {
      loading: true,
      saving: false,
      prefs: {
        webPushEnabled: true,
        smsEnabled: true,
        emailEnabled: true,
        preferredTime: 'day',
      },
      subscriptions: [],
      pushSupported: false,
      pushSubscribed: false,
    };
  },
  async mounted() {
    this.pushSupported = 'Notification' in window && 'serviceWorker' in navigator;
    await this.load();
  },
  methods: {
    async load() {
      try {
        const prefsRes = await notificationAPI.getPreferences();
        this.prefs = prefsRes.data;
        const subsRes = await notificationAPI.getSubscriptions();
        this.subscriptions = subsRes.data;
      } catch (e) {
        console.error('Failed to load notification settings', e);
      } finally {
        this.loading = false;
      }
    },
    async savePreferences() {
      this.saving = true;
      try {
        await notificationAPI.updatePreferences(this.prefs);
        alert('Настройки сохранены');
      } catch (e) {
        alert('Ошибка сохранения');
      } finally {
        this.saving = false;
      }
    },
    async subscribePush() {
      try {
        const permission = await Notification.requestPermission();
        if (permission !== 'granted') {
          alert('Разрешение на уведомления не получено');
          return;
        }
        const registration = await navigator.serviceWorker.register('/sw.js');
        const subscription = await registration.pushManager.subscribe({
          userVisibleOnly: true,
          applicationServerKey: this.urlBase64ToUint8Array('BCi1YGq0qGxM3qVxQqFpQqHxM3qVxQqFpQqHxM3qVxQqFpQqHxM3qVxQqFpQqHxM3qVxQqFpQqHxM3qVxQqFpQqHxM3qVxQ'),
        });
        await notificationAPI.subscribe({
          endpoint: subscription.endpoint,
          p256dh: this.arrayBufferToBase64(subscription.getKey('p256dh')),
          auth: this.arrayBufferToBase64(subscription.getKey('auth')),
          deviceName: navigator.userAgent,
        });
        this.pushSubscribed = true;
        alert('Push-уведомления включены');
      } catch (e) {
        console.error('Push subscription error', e);
        alert('Ошибка при подписке на push-уведомления');
      }
    },
    async removeSubscription(id) {
      try {
        await notificationAPI.unsubscribe(id);
        this.subscriptions = this.subscriptions.filter((s) => s.id !== id);
      } catch (e) {
        alert('Ошибка при отписке');
      }
    },
    urlBase64ToUint8Array(base64String) {
      const padding = '='.repeat((4 - (base64String.length % 4)) % 4);
      const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/');
      const rawData = window.atob(base64);
      return Uint8Array.from([...rawData].map((char) => char.charCodeAt(0)));
    },
    arrayBufferToBase64(buffer) {
      const bytes = new Uint8Array(buffer);
      let binary = '';
      for (let i = 0; i < bytes.length; i++) {
        binary += String.fromCharCode(bytes[i]);
      }
      return btoa(binary);
    },
  },
};
</script>

<style scoped>
.notif-settings { padding: 16px; max-width: 600px; width: 100%; }
.section { margin-bottom: 24px; }
.section h3 { margin-bottom: 12px; }
.channel-toggle { margin-bottom: 8px; }
.channel-toggle label { display: flex; align-items: center; gap: 8px; cursor: pointer; }
select { width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 8px; }
.sub-item { display: flex; justify-content: space-between; align-items: center; padding: 8px 0; border-bottom: 1px solid #eee; }
.no-subs { color: #888; font-size: 0.9em; }
.btn-primary { padding: 10px 20px; background: #1976d2; color: white; border: none; border-radius: 8px; cursor: pointer; }
.btn-danger { padding: 4px 12px; background: #c62828; color: white; border: none; border-radius: 6px; cursor: pointer; }
.loading { text-align: center; padding: 40px; }
.push-active { color: #2e7d32; font-weight: 600; }
.push-unsupported { color: #888; }
</style>
