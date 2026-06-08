package repository

import (
	"api_hr/internal/models"

	"gorm.io/gorm"
)

type SurveyRepository struct {
	db *gorm.DB
}

func NewSurveyRepository(db *gorm.DB) *SurveyRepository {
	return &SurveyRepository{db: db}
}

func (r *SurveyRepository) Create(survey *models.Survey) error {
	return r.db.Create(survey).Error
}

func (r *SurveyRepository) FindByID(id string) (*models.Survey, error) {
	var survey models.Survey
	err := r.db.Preload("Questions").First(&survey, "id = ?", id).Error
	return &survey, err
}

func (r *SurveyRepository) FindAll() ([]models.Survey, error) {
	var surveys []models.Survey
	err := r.db.Order("created_at desc").Find(&surveys).Error
	return surveys, err
}

func (r *SurveyRepository) FindByStatus(status models.SurveyStatus) ([]models.Survey, error) {
	var surveys []models.Survey
	err := r.db.Where("status = ?", status).Order("created_at desc").Find(&surveys).Error
	return surveys, err
}

func (r *SurveyRepository) FindActiveByTargetRole(role string) ([]models.Survey, error) {
	var surveys []models.Survey
	err := r.db.Where("status = ? AND (target_roles IS NULL OR target_roles = '[]'::jsonb OR target_roles @> ?::jsonb)", models.SurveyActive, `["`+role+`"]`).
		Order("created_at desc").Find(&surveys).Error
	return surveys, err
}

func (r *SurveyRepository) Update(survey *models.Survey) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(survey).Error
}

func (r *SurveyRepository) Delete(id string) error {
	return r.db.Delete(&models.Survey{}, "id = ?", id).Error
}

func (r *SurveyRepository) UpdateStatus(id string, status models.SurveyStatus) error {
	return r.db.Model(&models.Survey{}).Where("id = ?", id).Update("status", status).Error
}
