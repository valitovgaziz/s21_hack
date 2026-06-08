package repository

import (
	"api_hr/internal/models"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) SaveSubscription(sub *models.PushSubscription) error {
	return r.db.Create(sub).Error
}

func (r *NotificationRepository) FindSubscriptionsByUser(userID uint) ([]models.PushSubscription, error) {
	var subs []models.PushSubscription
	err := r.db.Where("user_id = ? AND active = ?", userID, true).Find(&subs).Error
	return subs, err
}

func (r *NotificationRepository) FindAllActiveSubscriptions() ([]models.PushSubscription, error) {
	var subs []models.PushSubscription
	err := r.db.Where("active = ?", true).Find(&subs).Error
	return subs, err
}

func (r *NotificationRepository) DeactivateSubscription(id uint) error {
	return r.db.Model(&models.PushSubscription{}).Where("id = ?", id).Update("active", false).Error
}

func (r *NotificationRepository) DeactivateSubscriptionByEndpoint(endpoint string) error {
	return r.db.Model(&models.PushSubscription{}).Where("endpoint = ?", endpoint).Update("active", false).Error
}

func (r *NotificationRepository) GetOrCreatePreferences(userID uint) (*models.NotificationPreference, error) {
	var prefs models.NotificationPreference
	err := r.db.Where("user_id = ?", userID).First(&prefs).Error
	if err == gorm.ErrRecordNotFound {
		prefs = models.NotificationPreference{
			UserID:         userID,
			WebPushEnabled: true,
			SMSEnabled:     true,
			EmailEnabled:   true,
		}
		err = r.db.Create(&prefs).Error
	}
	return &prefs, err
}

func (r *NotificationRepository) UpdatePreferences(prefs *models.NotificationPreference) error {
	return r.db.Save(prefs).Error
}

func (r *NotificationRepository) LogNotification(log *models.NotificationLog) error {
	return r.db.Create(log).Error
}

func (r *NotificationRepository) GetLogsBySurvey(surveyID string) ([]models.NotificationLog, error) {
	var logs []models.NotificationLog
	err := r.db.Where("survey_id = ?", surveyID).Order("created_at desc").Find(&logs).Error
	return logs, err
}

func (r *NotificationRepository) GetLogsByUser(userID uint) ([]models.NotificationLog, error) {
	var logs []models.NotificationLog
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&logs).Error
	return logs, err
}

func (r *NotificationRepository) GetChannelStats(surveyID string) (map[string]models.ChannelStat, error) {
	type logStat struct {
		Channel string
		Status  string
		Count   int64
	}
	var stats []logStat
	err := r.db.Model(&models.NotificationLog{}).
		Select("channel, status, COUNT(*) as count").
		Where("survey_id = ?", surveyID).
		Group("channel, status").
		Find(&stats).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]models.ChannelStat)
	for _, s := range stats {
		cs := result[s.Channel]
		switch s.Status {
		case "sent":
			cs.Sent = s.Count
		case "delivered":
			cs.Delivered = s.Count
		case "opened":
			cs.Opened = s.Count
		case "clicked":
			cs.Clicked = s.Count
		}
		result[s.Channel] = cs
	}
	for ch, cs := range result {
		if cs.Sent > 0 {
			cs.CTR = float64(cs.Clicked) / float64(cs.Sent) * 100
		}
		result[ch] = cs
	}
	return result, nil
}

func (r *NotificationRepository) GetRecentLogsByUserWithin(userID uint, days int) ([]models.NotificationLog, error) {
	var logs []models.NotificationLog
	err := r.db.Where("user_id = ? AND created_at > NOW() - INTERVAL '?' DAY", userID, days).Find(&logs).Error
	return logs, err
}
