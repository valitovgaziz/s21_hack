package repository

import (
	"api_hr/internal/models"
	"time"

	"gorm.io/gorm"
)

type ResponseRepository struct {
	db *gorm.DB
}

func NewResponseRepository(db *gorm.DB) *ResponseRepository {
	return &ResponseRepository{db: db}
}

func (r *ResponseRepository) Create(response *models.SurveyResponse) error {
	return r.db.Create(response).Error
}

func (r *ResponseRepository) FindByID(id string) (*models.SurveyResponse, error) {
	var resp models.SurveyResponse
	err := r.db.Preload("Answers").First(&resp, "id = ?", id).Error
	return &resp, err
}

func (r *ResponseRepository) FindBySurvey(surveyID string) ([]models.SurveyResponse, error) {
	var responses []models.SurveyResponse
	err := r.db.Where("survey_id = ?", surveyID).Preload("Answers").Find(&responses).Error
	return responses, err
}

func (r *ResponseRepository) FindByUserAndSurvey(userID uint, surveyID string) (*models.SurveyResponse, error) {
	var resp models.SurveyResponse
	err := r.db.Where("user_id = ? AND survey_id = ?", userID, surveyID).First(&resp).Error
	return &resp, err
}

func (r *ResponseRepository) CountBySurvey(surveyID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.SurveyResponse{}).Where("survey_id = ?", surveyID).Count(&count).Error
	return count, err
}

func (r *ResponseRepository) CountCompletedBySurvey(surveyID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.SurveyResponse{}).Where("survey_id = ? AND status = 'completed'", surveyID).Count(&count).Error
	return count, err
}

func (r *ResponseRepository) GetDailyStats(surveyID string) ([]struct {
	Date  string
	Count int64
}, error) {
	var results []struct {
		Date  string
		Count int64
	}
	err := r.db.Model(&models.SurveyResponse{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("survey_id = ? AND status = 'completed'", surveyID).
		Group("DATE(created_at)").
		Order("date asc").
		Find(&results).Error
	return results, err
}

func (r *ResponseRepository) Update(response *models.SurveyResponse) error {
	return r.db.Save(response).Error
}

func (r *ResponseRepository) HasUserResponded(userID uint, surveyID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.SurveyResponse{}).
		Where("user_id = ? AND survey_id = ? AND status = 'completed'", userID, surveyID).
		Count(&count).Error
	return count > 0, err
}

func (r *ResponseRepository) GetStartedCount(surveyID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.SurveyResponse{}).
		Where("survey_id = ? AND status = 'in_progress'", surveyID).
		Count(&count).Error
	return count, err
}

func (r *ResponseRepository) FindBySurveyAndSession(surveyID, sessionID string) (*models.SurveyResponse, error) {
	var resp models.SurveyResponse
	err := r.db.Where("survey_id = ? AND session_id = ?", surveyID, sessionID).First(&resp).Error
	return &resp, err
}

func (r *ResponseRepository) CountUsersWhoResponded(surveyID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.SurveyResponse{}).
		Where("survey_id = ? AND status = 'completed' AND user_id IS NOT NULL", surveyID).
		Distinct("user_id").Count(&count).Error
	return count, err
}

func (r *ResponseRepository) CompletedBetween(surveyID string, start, end time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.SurveyResponse{}).
		Where("survey_id = ? AND status = 'completed' AND completed_at BETWEEN ? AND ?", surveyID, start, end).
		Count(&count).Error
	return count, err
}
