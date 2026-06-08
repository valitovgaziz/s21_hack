package service

import (
	"api_hr/internal/models"
	"api_hr/internal/repository"
)

type NotificationService struct {
	notifRepo *repository.NotificationRepository
}

func NewNotificationService(notifRepo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{notifRepo: notifRepo}
}

func (s *NotificationService) Subscribe(sub *models.PushSubscription) error {
	return s.notifRepo.SaveSubscription(sub)
}

func (s *NotificationService) GetSubscriptions(userID uint) ([]models.PushSubscription, error) {
	return s.notifRepo.FindSubscriptionsByUser(userID)
}

func (s *NotificationService) Unsubscribe(id uint) error {
	return s.notifRepo.DeactivateSubscription(id)
}

func (s *NotificationService) UnsubscribeByEndpoint(endpoint string) error {
	return s.notifRepo.DeactivateSubscriptionByEndpoint(endpoint)
}

func (s *NotificationService) GetPreferences(userID uint) (*models.NotificationPreference, error) {
	return s.notifRepo.GetOrCreatePreferences(userID)
}

func (s *NotificationService) UpdatePreferences(userID uint, prefs *models.NotificationPreference) (*models.NotificationPreference, error) {
	existing, err := s.notifRepo.GetOrCreatePreferences(userID)
	if err != nil {
		return nil, err
	}
	prefs.ID = existing.ID
	prefs.UserID = userID
	if err := s.notifRepo.UpdatePreferences(prefs); err != nil {
		return nil, err
	}
	return prefs, nil
}

func (s *NotificationService) LogNotification(log *models.NotificationLog) error {
	return s.notifRepo.LogNotification(log)
}

func (s *NotificationService) GetSurveyNotificationLogs(surveyID string) ([]models.NotificationLog, error) {
	return s.notifRepo.GetLogsBySurvey(surveyID)
}

func (s *NotificationService) GetChannelStats(surveyID string) (map[string]models.ChannelStat, error) {
	return s.notifRepo.GetChannelStats(surveyID)
}
