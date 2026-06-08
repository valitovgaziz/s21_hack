package handlers

import (
	"net/http"

	"api_hr/internal/service"
	"api_hr/internal/utils"

	"github.com/go-chi/chi/v5"
)

type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

func (h *AnalyticsHandler) GetSurveyAnalytics(w http.ResponseWriter, r *http.Request) {
	surveyID := chi.URLParam(r, "surveyId")

	analytics, err := h.analyticsService.GetSurveyAnalytics(surveyID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Survey not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, analytics)
}
