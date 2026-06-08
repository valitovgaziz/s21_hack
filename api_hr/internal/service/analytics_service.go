package service

import (
	"fmt"

	"api_hr/internal/models"
	"api_hr/internal/repository"
)

type AnalyticsService struct {
	surveyRepo    *repository.SurveyRepository
	responseRepo  *repository.ResponseRepository
	notifRepo     *repository.NotificationRepository
}

func NewAnalyticsService(surveyRepo *repository.SurveyRepository, responseRepo *repository.ResponseRepository, notifRepo *repository.NotificationRepository) *AnalyticsService {
	return &AnalyticsService{
		surveyRepo:   surveyRepo,
		responseRepo: responseRepo,
		notifRepo:    notifRepo,
	}
}

func (s *AnalyticsService) GetSurveyAnalytics(surveyID string) (*models.SurveyAnalytics, error) {
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		return nil, err
	}

	totalSent, _ := s.responseRepo.CountBySurvey(surveyID)
	totalCompleted, _ := s.responseRepo.CountCompletedBySurvey(surveyID)
	totalStarted, _ := s.responseRepo.GetStartedCount(surveyID)

	analytics := &models.SurveyAnalytics{
		SurveyID:       surveyID,
		Title:          survey.Title,
		TotalSent:      totalSent,
		TotalStarted:   totalStarted,
		TotalCompleted: totalCompleted,
	}

	if totalSent > 0 {
		analytics.CompletionRate = float64(totalCompleted) / float64(totalSent) * 100
	}

	responses, _ := s.responseRepo.FindBySurvey(surveyID)
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

	dailyStats, _ := s.responseRepo.GetDailyStats(surveyID)
	for _, ds := range dailyStats {
		analytics.DailyStats = append(analytics.DailyStats, models.DailyStat{Date: ds.Date, Count: ds.Count})
	}

	channelStats, _ := s.notifRepo.GetChannelStats(surveyID)
	analytics.ChannelStats = channelStats

	return analytics, nil
}
