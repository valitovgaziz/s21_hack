import http from './http';

export const surveyAPI = {
  getAll: () => http.get('/surveys'),
  getByID: (id) => http.get(`/surveys/${id}`),
  create: (data) => http.post('/surveys', data),
  update: (id, data) => http.put(`/surveys/${id}`, data),
  delete: (id) => http.delete(`/surveys/${id}`),
  publish: (id) => http.post(`/surveys/${id}/publish`),
  complete: (id) => http.post(`/surveys/${id}/complete`),
  archive: (id) => http.post(`/surveys/${id}/archive`),
  getStats: (id) => http.get(`/surveys/${id}/stats`),
};
