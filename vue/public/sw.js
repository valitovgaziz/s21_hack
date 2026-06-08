self.addEventListener('push', function (event) {
  let data = {};
  if (event.data) {
    try {
      data = event.data.json();
    } catch (e) {
      data = { title: 'PulseHR', body: event.data.text() };
    }
  }

  const title = data.title || 'Новый опрос PulseHR';
  const options = {
    body: data.body || 'У вас новый опрос для прохождения',
    icon: '/icons/icon-192.png',
    badge: '/icons/icon-192.png',
    data: {
      surveyId: data.surveyId,
      url: data.url || '/',
    },
    actions: data.surveyId
      ? [{ action: 'open', title: 'Пройти опрос' }]
      : [],
  };

  event.waitUntil(self.registration.showNotification(title, options));
});

self.addEventListener('notificationclick', function (event) {
  event.notification.close();

  const urlToOpen = event.notification.data?.url || '/';

  if (event.action === 'open' && event.notification.data?.surveyId) {
    event.waitUntil(clients.openWindow('/survey/' + event.notification.data.surveyId));
  } else {
    event.waitUntil(clients.openWindow(urlToOpen));
  }
});

self.addEventListener('install', function (event) {
  self.skipWaiting();
});

self.addEventListener('activate', function (event) {
  event.waitUntil(clients.claim());
});
