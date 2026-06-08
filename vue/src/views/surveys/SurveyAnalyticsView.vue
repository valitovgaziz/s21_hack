<template>
  <div class="analytics">
    <h2>Аналитика: {{ analytics?.title || 'Загрузка...' }}</h2>

    <div v-if="loading" class="loading">Загрузка аналитики...</div>

    <div v-else-if="analytics" class="analytics-dashboard">
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-value">{{ analytics.totalCompleted }}</div>
          <div class="stat-label">Завершено</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{{ analytics.totalStarted }}</div>
          <div class="stat-label">Начато</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{{ analytics.completionRate?.toFixed(1) || 0 }}%</div>
          <div class="stat-label">Прохождение</div>
        </div>
        <div v-if="analytics.enps != null" class="stat-card highlight">
          <div class="stat-value">{{ analytics.enps.toFixed(1) }}</div>
          <div class="stat-label">eNPS</div>
        </div>
      </div>

      <div v-if="analytics.dailyStats?.length" class="section">
        <h3>Прохождения по дням</h3>
        <div class="daily-chart">
          <div v-for="d in analytics.dailyStats" :key="d.date" class="daily-bar-wrapper">
            <div class="daily-bar" :style="{ height: barHeight(d.count) + 'px' }"></div>
            <div class="daily-label">{{ formatShortDate(d.date) }}</div>
          </div>
        </div>
      </div>

      <div v-if="analytics.questionStats?.length" class="section">
        <h3>Результаты по вопросам</h3>
        <div v-for="qs in analytics.questionStats" :key="qs.questionId" class="question-stat">
          <h4>{{ qs.title }}</h4>
          <div v-for="(count, answer) in qs.answers" :key="answer" class="answer-row">
            <span class="answer-text">{{ answer }}</span>
            <div class="answer-bar-bg">
              <div class="answer-bar" :style="{ width: answerPercent(count, qs.total) + '%' }"></div>
            </div>
            <span class="answer-count">{{ count }} ({{ answerPercent(count, qs.total).toFixed(1) }}%)</span>
          </div>
        </div>
      </div>

      <div v-if="analytics.channelStats && Object.keys(analytics.channelStats).length" class="section">
        <h3>Эффективность каналов уведомлений</h3>
        <table class="channel-table">
          <thead>
            <tr>
              <th>Канал</th>
              <th>Отправлено</th>
              <th>Доставлено</th>
              <th>Открыто</th>
              <th>Переходы</th>
              <th>CTR</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(stat, channel) in analytics.channelStats" :key="channel">
              <td>{{ channelLabel(channel) }}</td>
              <td>{{ stat.sent }}</td>
              <td>{{ stat.delivered }}</td>
              <td>{{ stat.opened }}</td>
              <td>{{ stat.clicked }}</td>
              <td>{{ stat.ctr?.toFixed(1) || 0 }}%</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-else class="error">Нет данных для аналитики</div>
  </div>
</template>

<script>
import { analyticsAPI } from '../../api/analytics';

export default {
  name: 'SurveyAnalyticsView',
  data() {
    return {
      analytics: null,
      loading: true,
      maxDailyCount: 0,
    };
  },
  async mounted() {
    await this.loadAnalytics();
  },
  methods: {
    async loadAnalytics() {
      try {
        const res = await analyticsAPI.getSurveyAnalytics(this.$route.params.id);
        this.analytics = res.data;
        if (this.analytics.dailyStats?.length) {
          this.maxDailyCount = Math.max(...this.analytics.dailyStats.map((d) => d.count));
        }
      } catch (e) {
        console.error('Analytics load error', e);
      } finally {
        this.loading = false;
      }
    },
    barHeight(count) {
      return this.maxDailyCount > 0 ? (count / this.maxDailyCount) * 120 : 0;
    },
    answerPercent(count, total) {
      return total > 0 ? (count / total) * 100 : 0;
    },
    formatShortDate(dateStr) {
      const d = new Date(dateStr);
      return d.toLocaleDateString('ru', { day: 'numeric', month: 'short' });
    },
    channelLabel(ch) {
      const labels = { web_push: 'Web Push', sms: 'SMS', telegram: 'Telegram', email: 'Email' };
      return labels[ch] || ch;
    },
  },
};
</script>

<style scoped>
.analytics { padding: 16px; max-width: 800px; width: 100%; }
.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(140px, 1fr)); gap: 12px; margin-bottom: 24px; }
.stat-card {
  background: var(--light-dark-background-color);
  border-radius: 12px;
  padding: 20px;
  text-align: center;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}
.stat-card.highlight { background: #e8f5e9; }
.stat-value { font-size: 2em; font-weight: 700; color: #1976d2; }
.stat-card.highlight .stat-value { color: #2e7d32; }
.stat-label { font-size: 0.85em; color: #666; margin-top: 4px; }
.section { margin-bottom: 24px; }
.section h3 { margin-bottom: 12px; }
.daily-chart { display: flex; align-items: flex-end; gap: 8px; height: 140px; padding: 10px 0; }
.daily-bar-wrapper { flex: 1; display: flex; flex-direction: column; align-items: center; }
.daily-bar { width: 100%; max-width: 30px; background: #1976d2; border-radius: 4px 4px 0 0; min-height: 4px; }
.daily-label { font-size: 0.7em; color: #888; margin-top: 4px; transform: rotate(-45deg); white-space: nowrap; }
.question-stat {
  background: var(--light-dark-background-color);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
}
.question-stat h4 { margin: 0 0 8px; }
.answer-row { display: flex; align-items: center; gap: 12px; margin-bottom: 6px; }
.answer-text { min-width: 120px; font-size: 0.9em; }
.answer-bar-bg { flex: 1; height: 20px; background: #e0e0e0; border-radius: 10px; overflow: hidden; }
.answer-bar { height: 100%; background: #1976d2; border-radius: 10px; transition: width 0.3s; }
.answer-count { min-width: 80px; font-size: 0.85em; color: #666; }
.channel-table { width: 100%; border-collapse: collapse; }
.channel-table th, .channel-table td { padding: 8px 12px; border-bottom: 1px solid #eee; text-align: left; }
.channel-table th { font-weight: 600; color: #666; }
.loading, .error { text-align: center; padding: 40px; color: #888; }
</style>
