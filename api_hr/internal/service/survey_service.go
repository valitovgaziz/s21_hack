package service

import (
	"errors"
	"fmt"
	"time"

	"api_hr/internal/models"
	"api_hr/internal/repository"

	"github.com/google/uuid"
)

type SurveyService struct {
	surveyRepo   *repository.SurveyRepository
	responseRepo *repository.ResponseRepository
}

func NewSurveyService(surveyRepo *repository.SurveyRepository, responseRepo *repository.ResponseRepository) *SurveyService {
	return &SurveyService{surveyRepo: surveyRepo, responseRepo: responseRepo}
}

func (s *SurveyService) CreateSurvey(req *models.Survey) (*models.Survey, error) {
	if req.ID == "" {
		req.ID = uuid.New().String()
	}
	for i := range req.Questions {
		if req.Questions[i].ID == "" {
			req.Questions[i].ID = uuid.New().String()
		}
		req.Questions[i].SurveyID = req.ID
	}
	req.Status = models.SurveyDraft
	if err := s.surveyRepo.Create(req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *SurveyService) GetSurvey(id string) (*models.Survey, error) {
	return s.surveyRepo.FindByID(id)
}

func (s *SurveyService) GetAllSurveys() ([]models.Survey, error) {
	return s.surveyRepo.FindAll()
}

func (s *SurveyService) UpdateSurvey(id string, req *models.Survey) (*models.Survey, error) {
	existing, err := s.surveyRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("survey not found")
	}
	req.ID = existing.ID
	req.CreatedAt = existing.CreatedAt
	for i := range req.Questions {
		if req.Questions[i].ID == "" {
			req.Questions[i].ID = uuid.New().String()
		}
		req.Questions[i].SurveyID = req.ID
	}
	if err := s.surveyRepo.Update(req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *SurveyService) DeleteSurvey(id string) error {
	return s.surveyRepo.Delete(id)
}

func (s *SurveyService) PublishSurvey(id string) error {
	survey, err := s.surveyRepo.FindByID(id)
	if err != nil {
		return errors.New("survey not found")
	}
	if survey.Status != models.SurveyDraft {
		return errors.New("only draft surveys can be published")
	}
	now := time.Now()
	survey.StartDate = &now
	return s.surveyRepo.UpdateStatus(id, models.SurveyActive)
}

func (s *SurveyService) CompleteSurvey(id string) error {
	return s.surveyRepo.UpdateStatus(id, models.SurveyCompleted)
}

func (s *SurveyService) ArchiveSurvey(id string) error {
	return s.surveyRepo.UpdateStatus(id, models.SurveyArchived)
}

func (s *SurveyService) GetActiveSurveysForRole(role string) ([]models.Survey, error) {
	return s.surveyRepo.FindActiveByTargetRole(role)
}

func (s *SurveyService) GetSurveyStats(id string) (*models.SurveyAnalytics, error) {
	survey, err := s.surveyRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("survey not found")
	}

	totalSent, _ := s.responseRepo.CountBySurvey(id)
	totalCompleted, _ := s.responseRepo.CountCompletedBySurvey(id)
	totalStarted, _ := s.responseRepo.GetStartedCount(id)

	analytics := &models.SurveyAnalytics{
		SurveyID:       id,
		Title:          survey.Title,
		TotalSent:      totalSent + totalStarted,
		TotalStarted:   totalStarted,
		TotalCompleted: totalCompleted,
	}

	if totalSent+totalStarted > 0 {
		analytics.CompletionRate = float64(totalCompleted) / float64(totalSent+totalStarted) * 100
	}

	responses, _ := s.responseRepo.FindBySurvey(id)
	questionStats := make(map[string]*models.QuestionStat)
	for _, q := range survey.Questions {
		questionStats[q.ID] = &models.QuestionStat{
			QuestionID: q.ID,
			Title:      q.Title,
			Type:       string(q.Type),
			Answers:    make(map[string]int64),
		}
	}

	var totalScore float64
	var scoreCount int64
	for _, resp := range responses {
		if resp.Status != "completed" {
			continue
		}
		for _, ans := range resp.Answers {
			qs, ok := questionStats[ans.QuestionID]
			if !ok {
				continue
			}
			qs.Total++

			for _, val := range ans.Value {
				valStr := fmt.Sprintf("%v", val)
				qs.Answers[valStr]++

				if qs.Type == "scale" {
					if score, ok := toFloat64(val); ok {
						totalScore += score
						scoreCount++
					}
				}
			}
		}
	}

	for _, qs := range questionStats {
		analytics.QuestionStats = append(analytics.QuestionStats, *qs)
	}

	if scoreCount > 0 {
		avg := totalScore / float64(scoreCount)
		analytics.AverageScore = &avg
		enps := (avg - 1) / 9 * 100
		analytics.ENPS = &enps
	}

	dailyStats, _ := s.responseRepo.GetDailyStats(id)
	for _, ds := range dailyStats {
		analytics.DailyStats = append(analytics.DailyStats, models.DailyStat{
			Date:  ds.Date,
			Count: ds.Count,
		})
	}

	return analytics, nil
}

func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case string:
		var f float64
		_, _ = fmt.Sscanf(val, "%f", &f)
		return f, f != 0 || val == "0"
	}
	return 0, false
}
