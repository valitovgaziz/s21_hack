package service

import (
	"errors"
	"time"

	"api_hr/internal/models"
	"api_hr/internal/repository"

	"github.com/google/uuid"
)

type ResponseService struct {
	responseRepo  *repository.ResponseRepository
	surveyRepo    *repository.SurveyRepository
}

func NewResponseService(responseRepo *repository.ResponseRepository, surveyRepo *repository.SurveyRepository) *ResponseService {
	return &ResponseService{responseRepo: responseRepo, surveyRepo: surveyRepo}
}

func (s *ResponseService) StartSurvey(surveyID string, userID *uint, sessionID string) (*models.SurveyResponse, error) {
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		return nil, errors.New("survey not found")
	}
	if survey.Status != models.SurveyActive {
		return nil, errors.New("survey is not active")
	}

	if userID != nil {
		existing, err := s.responseRepo.HasUserResponded(*userID, surveyID)
		if err == nil && existing {
			return nil, errors.New("you have already completed this survey")
		}
	}

	now := time.Now()
	resp := &models.SurveyResponse{
		ID:        uuid.New().String(),
		SurveyID:  surveyID,
		UserID:    userID,
		SessionID: sessionID,
		Anonymous: survey.Privacy == models.PrivacyAnonymous,
		StartedAt: &now,
		Status:    "in_progress",
	}

	if err := s.responseRepo.Create(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *ResponseService) SubmitAnswer(responseID, questionID string, value interface{}) error {
	resp, err := s.responseRepo.FindByID(responseID)
	if err != nil {
		return errors.New("response not found")
	}
	if resp.Status == "completed" {
		return errors.New("survey already completed")
	}

	answer := models.Answer{
		ID:         uuid.New().String(),
		ResponseID: responseID,
		QuestionID: questionID,
		Value:      value.(map[string]interface{}),
	}

	resp.Answers = append(resp.Answers, answer)
	return s.responseRepo.Update(resp)
}

func (s *ResponseService) CompleteSurvey(responseID string) error {
	resp, err := s.responseRepo.FindByID(responseID)
	if err != nil {
		return errors.New("response not found")
	}
	now := time.Now()
	resp.CompletedAt = &now
	resp.Status = "completed"
	return s.responseRepo.Update(resp)
}

func (s *ResponseService) GetResponse(id string) (*models.SurveyResponse, error) {
	return s.responseRepo.FindByID(id)
}

func (s *ResponseService) GetSurveyResponses(surveyID string) ([]models.SurveyResponse, error) {
	return s.responseRepo.FindBySurvey(surveyID)
}

func (s *ResponseService) HasUserResponded(userID uint, surveyID string) (bool, error) {
	return s.responseRepo.HasUserResponded(userID, surveyID)
}
