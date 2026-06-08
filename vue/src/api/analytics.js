import http from './http';

export const analyticsAPI = {
  getSurveyAnalytics: (surveyId) => http.get(`/analytics/surveys/${surveyId}`),
};
