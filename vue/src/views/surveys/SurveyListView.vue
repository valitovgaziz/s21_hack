<template>
  <div class="survey-list">
    <div class="page-header">
      <h2>Опросы</h2>
      <router-link v-if="isHR" to="/surveys/new" class="btn-primary">+ Создать опрос</router-link>
    </div>

    <div v-if="loading" class="loading">Загрузка...</div>

    <div v-else-if="error" class="error">{{ error }}</div>

    <div v-else class="surveys">
      <div v-for="survey in surveys" :key="survey.id" class="survey-card" @click="viewSurvey(survey)">
        <div class="survey-header">
          <h3>{{ survey.title }}</h3>
          <span :class="'status-badge status-' + survey.status">{{ statusLabel(survey.status) }}</span>
        </div>
        <p class="survey-desc">{{ survey.description || 'Нет описания' }}</p>
        <div class="survey-meta">
          <span>Статус: {{ statusLabel(survey.status) }}</span>
          <span>{{ survey.questions?.length || 0 }} вопросов</span>
          <span v-if="survey.endDate">До: {{ formatDate(survey.endDate) }}</span>
        </div>
        <div v-if="survey.privacy" class="privacy-badge">
          {{ survey.privacy === 'anonymous' ? 'Анонимный' : 'Идентифицированный' }}
        </div>
      </div>
      <div v-if="surveys.length === 0" class="empty">
        Нет доступных опросов
      </div>
    </div>
  </div>
</template>

<script>
import { surveyAPI } from '../../api/surveys';
import { useAuthStore } from '../../auth/stores/auth.store';

export default {
  name: 'SurveyListView',
  data() {
    return {
      surveys: [],
      loading: true,
      error: null,
    };
  },
  computed: {
    isHR() {
      const store = useAuthStore();
      return store.user?.role === 'hr' || store.user?.role === 'admin';
    },
  },
  async mounted() {
    await this.loadSurveys();
  },
  methods: {
    async loadSurveys() {
      try {
        this.loading = true;
        const res = await surveyAPI.getAll();
        this.surveys = res.data;
      } catch (e) {
        this.error = 'Ошибка загрузки опросов';
      } finally {
        this.loading = false;
      }
    },
    viewSurvey(survey) {
      if (survey.status === 'draft') {
        this.$router.push(`/surveys/${survey.id}/edit`);
      } else {
        this.$router.push(`/survey/${survey.id}`);
      }
    },
    statusLabel(status) {
      const labels = { draft: 'Черновик', active: 'Активный', completed: 'Завершён', archived: 'Архив' };
      return labels[status] || status;
    },
    formatDate(date) {
      return new Date(date).toLocaleDateString();
    },
  },
};
</script>

<style scoped>
.survey-list {
  padding: 16px;
  max-width: 800px;
  width: 100%;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.btn-primary {
  background: #4CAF50;
  color: white;
  padding: 8px 16px;
  border-radius: 8px;
  text-decoration: none;
}
.survey-card {
  background: var(--light-dark-background-color);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
  cursor: pointer;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}
.survey-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.survey-header h3 { margin: 0; }
.survey-desc { color: #666; margin: 8px 0; }
.survey-meta { display: flex; gap: 16px; font-size: 0.85em; color: #888; }
.status-badge {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 0.8em;
}
.status-draft { background: #e0e0e0; color: #666; }
.status-active { background: #e8f5e9; color: #2e7d32; }
.status-completed { background: #e3f2fd; color: #1565c0; }
.status-archived { background: #fce4ec; color: #c62828; }
.privacy-badge {
  margin-top: 8px;
  font-size: 0.8em;
  color: #666;
  background: #f5f5f5;
  padding: 2px 8px;
  border-radius: 8px;
  display: inline-block;
}
.loading, .error, .empty { text-align: center; padding: 40px; color: #888; }
</style>
