import http from './http';

export const responseAPI = {
  start: (surveyId, sessionId) => http.post(`/responses/${surveyId}/start`, { session_id: sessionId }),
  submitAnswer: (responseId, questionId, value) =>
    http.post(`/responses/${responseId}/answer`, { question_id: questionId, value }),
  complete: (responseId) => http.post(`/responses/${responseId}/complete`),
  getByID: (id) => http.get(`/responses/${id}`),
  getBySurvey: (surveyId) => http.get(`/responses/survey/${surveyId}`),
  checkResponded: (surveyId) => http.get(`/responses/survey/${surveyId}/check`),
};
