import http from './http';

export const notificationAPI = {
  subscribe: (subscription) => http.post('/notifications/subscribe', subscription),
  getSubscriptions: () => http.get('/notifications/subscriptions'),
  unsubscribe: (id) => http.delete(`/notifications/subscriptions/${id}`),
  getPreferences: () => http.get('/notifications/preferences'),
  updatePreferences: (data) => http.put('/notifications/preferences', data),
  getSurveyLogs: (surveyId) => http.get(`/notifications/logs/${surveyId}`),
};
