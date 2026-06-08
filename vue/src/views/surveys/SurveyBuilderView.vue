<template>
  <div class="survey-builder">
    <h2>{{ isEdit ? 'Редактировать опрос' : 'Создать опрос' }}</h2>

    <div class="builder-form">
      <div class="form-group">
        <label>Название</label>
        <input v-model="survey.title" placeholder="Название опроса" />
      </div>

      <div class="form-group">
        <label>Описание</label>
        <textarea v-model="survey.description" placeholder="Описание опроса" rows="3"></textarea>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label>Режим конфиденциальности</label>
          <select v-model="survey.privacy">
            <option value="anonymous">Анонимный</option>
            <option value="identified">Идентифицированный</option>
          </select>
        </div>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label>Дата окончания</label>
          <input type="datetime-local" v-model="endDate" />
        </div>
      </div>

      <div class="form-group">
        <label>Целевая аудитория (роли)</label>
        <div class="roles-select">
          <label v-for="role in availableRoles" :key="role" class="role-checkbox">
            <input type="checkbox" :value="role" v-model="survey.targetRoles" />
            {{ roleLabel(role) }}
          </label>
        </div>
      </div>
    </div>

    <div class="questions-section">
      <h3>Вопросы</h3>

      <div v-for="(q, idx) in survey.questions" :key="q.id" class="question-card">
        <div class="q-header">
          <span class="q-number">Вопрос {{ idx + 1 }}</span>
          <button class="btn-icon" @click="removeQuestion(idx)" title="Удалить">&times;</button>
        </div>

        <div class="form-group">
          <input v-model="q.title" placeholder="Текст вопроса" />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>Тип</label>
            <select v-model="q.type">
              <option value="single_choice">Одиночный выбор</option>
              <option value="multiple_choice">Множественный выбор</option>
              <option value="scale">Шкала (NPS / 1-5)</option>
              <option value="text">Текстовый ответ</option>
              <option value="matrix">Матричный</option>
            </select>
          </div>
          <div class="form-group">
            <label>
              <input type="checkbox" v-model="q.required" />
              Обязательный
            </label>
          </div>
        </div>

        <div v-if="q.type === 'single_choice' || q.type === 'multiple_choice'" class="options-group">
          <label>Варианты ответов</label>
          <div v-for="(opt, oi) in getOptions(q)" :key="oi" class="option-row">
            <input :value="opt" @input="updateOption(q, oi, $event.target.value)" :placeholder="'Вариант ' + (oi+1)" />
            <button class="btn-icon" @click="removeOption(q, oi)">&times;</button>
          </div>
          <button class="btn-sm" @click="addOption(q)">+ Добавить вариант</button>
        </div>
      </div>

      <button class="btn-secondary" @click="addQuestion">+ Добавить вопрос</button>
    </div>

    <div class="builder-actions">
      <button class="btn-primary" @click="saveDraft" :disabled="saving">Сохранить черновик</button>
      <button class="btn-success" @click="publishSurvey" :disabled="saving">Опубликовать</button>
    </div>
  </div>
</template>

<script>
import { surveyAPI } from '../../api/surveys';

export default {
  name: 'SurveyBuilderView',
  props: { surveyId: String },
  data() {
    return {
      isEdit: false,
      saving: false,
      endDate: '',
      survey: this.emptySurvey(),
      availableRoles: ['employee', 'hr', 'manager', 'admin'],
    };
  },
  async mounted() {
    if (this.$route.params.id) {
      this.isEdit = true;
      try {
        const res = await surveyAPI.getByID(this.$route.params.id);
        this.survey = res.data;
        if (this.survey.endDate) {
          this.endDate = this.survey.endDate.slice(0, 16);
        }
      } catch (e) {
        alert('Ошибка загрузки опроса');
        this.$router.push('/surveys');
      }
    }
  },
  methods: {
    emptySurvey() {
      return {
        title: '',
        description: '',
        privacy: 'anonymous',
        targetRoles: ['employee'],
        questions: [],
        blockStructure: {},
      };
    },
    addQuestion() {
      this.survey.questions.push({
        id: 'q_' + Date.now() + '_' + Math.random().toString(36).slice(2, 6),
        title: '',
        type: 'single_choice',
        required: false,
        options: [],
        order: this.survey.questions.length,
        blockId: null,
        branchLogic: [],
      });
    },
    removeQuestion(idx) {
      this.survey.questions.splice(idx, 1);
    },
    getOptions(q) {
      if (!Array.isArray(q.options)) {
        this.$set(q, 'options', []);
      }
      return q.options;
    },
    addOption(q) {
      if (!Array.isArray(q.options)) {
        this.$set(q, 'options', []);
      }
      q.options.push('');
    },
    updateOption(q, idx, value) {
      if (!Array.isArray(q.options)) {
        this.$set(q, 'options', []);
      }
      q.options[idx] = value;
    },
    removeOption(q, idx) {
      q.options.splice(idx, 1);
    },
    async saveDraft() {
      this.saving = true;
      try {
        if (this.isEdit) {
          await surveyAPI.update(this.$route.params.id, this.survey);
        } else {
          await surveyAPI.create(this.survey);
        }
        this.$router.push('/surveys');
      } catch (e) {
        alert('Ошибка сохранения');
      } finally {
        this.saving = false;
      }
    },
    async publishSurvey() {
      const id = this.isEdit ? this.$route.params.id : await this.saveAndGetID();
      if (!id) return;
      try {
        await surveyAPI.publish(id);
        this.$router.push('/surveys');
      } catch (e) {
        alert('Ошибка публикации');
      }
    },
    async saveAndGetID() {
      try {
        const res = await surveyAPI.create(this.survey);
        return res.data.id;
      } catch (e) {
        alert('Ошибка сохранения');
        return null;
      }
    },
    roleLabel(role) {
      const labels = { employee: 'Сотрудник', hr: 'HR', manager: 'Руководитель', admin: 'Администратор' };
      return labels[role] || role;
    },
  },
};
</script>

<style scoped>
.survey-builder {
  padding: 16px;
  max-width: 800px;
  width: 100%;
}
.builder-form { margin-bottom: 24px; }
.form-group { margin-bottom: 12px; }
.form-group label { display: block; margin-bottom: 4px; font-weight: 600; }
.form-group input, .form-group textarea, .form-group select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 14px;
  box-sizing: border-box;
}
.form-row { display: flex; gap: 16px; }
.form-row .form-group { flex: 1; }
.roles-select { display: flex; gap: 16px; flex-wrap: wrap; }
.role-checkbox { display: flex; align-items: center; gap: 4px; }
.question-card {
  background: var(--light-dark-background-color);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
  border: 1px solid #eee;
}
.q-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
.q-number { font-weight: 600; color: #666; }
.options-group { margin-top: 8px; }
.option-row { display: flex; gap: 8px; margin-bottom: 4px; }
.option-row input { flex: 1; }
.btn-icon { background: none; border: none; font-size: 20px; cursor: pointer; color: #c62828; }
.btn-sm { padding: 4px 12px; background: #e0e0e0; border: none; border-radius: 6px; cursor: pointer; }
.btn-secondary { padding: 8px 16px; background: #e0e0e0; border: none; border-radius: 8px; cursor: pointer; }
.btn-primary { padding: 10px 20px; background: #1976d2; color: white; border: none; border-radius: 8px; cursor: pointer; }
.btn-success { padding: 10px 20px; background: #4CAF50; color: white; border: none; border-radius: 8px; cursor: pointer; }
.builder-actions { display: flex; gap: 12px; margin-top: 24px; }
button:disabled { opacity: 0.6; cursor: not-allowed; }
</style>
