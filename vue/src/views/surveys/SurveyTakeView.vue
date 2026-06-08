<template>
  <div class="survey-take">
    <div v-if="loading" class="loading">Загрузка...</div>

    <div v-else-if="completed" class="completed-message">
      <h2>Спасибо!</h2>
      <p>Ваши ответы успешно сохранены.</p>
      <router-link to="/" class="btn-primary">На главную</router-link>
    </div>

    <div v-else-if="survey" class="survey-content">
      <div v-if="survey.privacy === 'anonymous'" class="privacy-notice">
        Этот опрос анонимный. HR не увидит ваши ответы.
      </div>
      <div v-else class="privacy-notice identified">
        Ваши ответы будут видны HR с указанием вашего имени.
      </div>

      <h2>{{ survey.title }}</h2>
      <p class="survey-desc">{{ survey.description }}</p>
      <p v-if="survey.questions" class="question-count">{{ survey.questions.length }} вопросов</p>

      <div v-for="(q, idx) in visibleQuestions" :key="q.id" class="question-block">
        <div class="question-text">
          <span class="q-num">{{ idx + 1 }}.</span>
          {{ q.title }}
          <span v-if="q.required" class="required">*</span>
        </div>

        <div v-if="q.type === 'single_choice'" class="options">
          <label v-for="opt in getOptions(q)" :key="opt" class="option-label">
            <input type="radio" :name="'q_' + q.id" :value="opt"
              @change="setAnswer(q.id, { value: opt })" />
            {{ opt }}
          </label>
        </div>

        <div v-else-if="q.type === 'multiple_choice'" class="options">
          <label v-for="opt in getOptions(q)" :key="opt" class="option-label">
            <input type="checkbox" :value="opt"
              @change="toggleMulti(q.id, opt)" />
            {{ opt }}
          </label>
        </div>

        <div v-else-if="q.type === 'scale'" class="scale">
          <div class="scale-options">
            <button v-for="n in scaleRange(q)" :key="n"
              :class="['scale-btn', { active: getAnswer(q.id)?.value == n }]"
              @click="setAnswer(q.id, { value: n })">
              {{ n }}
            </button>
          </div>
        </div>

        <div v-else-if="q.type === 'text'">
          <textarea v-model="textAnswers[q.id]" @input="setAnswer(q.id, { value: textAnswers[q.id] })"
            placeholder="Ваш ответ..." rows="3"></textarea>
        </div>
      </div>

      <div class="take-actions">
        <button class="btn-primary" @click="submitSurvey" :disabled="submitting">
          {{ submitting ? 'Отправка...' : 'Отправить ответы' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { surveyAPI } from '../../api/surveys';
import { responseAPI } from '../../api/responses';

export default {
  name: 'SurveyTakeView',
  data() {
    return {
      survey: null,
      loading: true,
      submitting: false,
      completed: false,
      responseId: null,
      answers: {},
      textAnswers: {},
      currentBlock: null,
    };
  },
  computed: {
    visibleQuestions() {
      if (!this.survey?.questions) return [];
      return this.survey.questions.filter((q) => {
        if (!q.branchLogic || q.branchLogic.length === 0) return true;
        for (const rule of q.branchLogic) {
          const parentAnswer = this.answers[rule.questionId]?.value;
          if (!parentAnswer) continue;
          if (rule.operator === 'eq' && parentAnswer === rule.value) return true;
          if (rule.operator === 'in' && Array.isArray(rule.value) && rule.value.includes(parentAnswer)) return true;
        }
        return false;
      });
    },
  },
  async mounted() {
    await this.loadSurvey();
  },
  methods: {
    async loadSurvey() {
      try {
        const res = await surveyAPI.getByID(this.$route.params.id);
        this.survey = res.data;
        await this.startResponse();
      } catch (e) {
        alert('Ошибка загрузки опроса');
        this.$router.push('/');
      } finally {
        this.loading = false;
      }
    },
    async startResponse() {
      try {
        const sessionId = 'sess_' + Date.now();
        const res = await responseAPI.start(this.survey.id, sessionId);
        this.responseId = res.data.id;
      } catch (e) {
        if (e.response?.status === 400) {
          this.completed = true;
        }
      }
    },
    getOptions(q) {
      if (Array.isArray(q.options)) return q.options;
      return [];
    },
    scaleRange(q) {
      if (Array.isArray(q.options) && q.options.length >= 2) {
        const min = parseInt(q.options[0]) || 1;
        const max = parseInt(q.options[1]) || 5;
        const arr = [];
        for (let i = min; i <= max; i++) arr.push(i);
        return arr;
      }
      return [1, 2, 3, 4, 5];
    },
    setAnswer(questionId, value) {
      this.answers[questionId] = value;
    },
    getAnswer(questionId) {
      return this.answers[questionId] || null;
    },
    toggleMulti(questionId, opt) {
      const current = this.answers[questionId]?.value || [];
      if (Array.isArray(current)) {
        const idx = current.indexOf(opt);
        if (idx >= 0) current.splice(idx, 1);
        else current.push(opt);
        this.answers[questionId] = { value: current };
      } else {
        this.answers[questionId] = { value: [opt] };
      }
    },
    async submitSurvey() {
      if (!this.responseId) return;
      this.submitting = true;
      try {
        for (const [qId, ans] of Object.entries(this.answers)) {
          await responseAPI.submitAnswer(this.responseId, qId, ans);
        }
        await responseAPI.complete(this.responseId);
        this.completed = true;
      } catch (e) {
        alert('Ошибка при отправке ответов');
      } finally {
        this.submitting = false;
      }
    },
  },
};
</script>

<style scoped>
.survey-take {
  padding: 16px;
  max-width: 700px;
  width: 100%;
}
.privacy-notice {
  background: #e3f2fd;
  padding: 10px 16px;
  border-radius: 8px;
  margin-bottom: 16px;
  font-size: 0.9em;
  color: #1565c0;
}
.privacy-notice.identified {
  background: #fff3e0;
  color: #e65100;
}
.survey-desc { color: #666; }
.question-count { font-size: 0.9em; color: #888; }
.question-block {
  background: var(--light-dark-background-color);
  border-radius: 12px;
  padding: 16px;
  margin: 12px 0;
  border: 1px solid #eee;
}
.question-text {
  font-weight: 600;
  margin-bottom: 12px;
}
.q-num { color: #1976d2; }
.required { color: #c62828; }
.options { display: flex; flex-direction: column; gap: 8px; }
.option-label {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  cursor: pointer;
}
.option-label:hover { background: #f5f5f5; }
.scale-options { display: flex; gap: 8px; flex-wrap: wrap; }
.scale-btn {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  border: 2px solid #ddd;
  background: white;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.2s;
}
.scale-btn.active { background: #1976d2; color: white; border-color: #1976d2; }
.scale-btn:hover { border-color: #1976d2; }
textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 14px;
  box-sizing: border-box;
  resize: vertical;
}
.take-actions { margin-top: 24px; text-align: center; }
.btn-primary {
  padding: 12px 32px;
  background: #4CAF50;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
}
.loading { text-align: center; padding: 40px; color: #888; }
.completed-message { text-align: center; padding: 60px 20px; }
.completed-message h2 { color: #4CAF50; }
</style>
