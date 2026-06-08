package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api_hr/internal/models"
	"api_hr/internal/service"
	"api_hr/internal/utils"

	"github.com/go-chi/chi/v5"
)

type NotificationHandler struct {
	notifService *service.NotificationService
}

func NewNotificationHandler(notifService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifService: notifService}
}

func (h *NotificationHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())
	if userID == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var sub models.PushSubscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	sub.UserID = *userID
	sub.Active = true

	if err := h.notifService.Subscribe(&sub); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to save subscription")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"status": "subscribed"})
}

func (h *NotificationHandler) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())
	if userID == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	subs, err := h.notifService.GetSubscriptions(*userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get subscriptions")
		return
	}

	utils.WriteJSON(w, http.StatusOK, subs)
}

func (h *NotificationHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	if err := h.notifService.Unsubscribe(uint(id)); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to unsubscribe")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "unsubscribed"})
}

func (h *NotificationHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())
	if userID == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	prefs, err := h.notifService.GetPreferences(*userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get preferences")
		return
	}

	utils.WriteJSON(w, http.StatusOK, prefs)
}

func (h *NotificationHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())
	if userID == nil {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var prefs models.NotificationPreference
	if err := json.NewDecoder(r.Body).Decode(&prefs); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.notifService.UpdatePreferences(*userID, &prefs)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update preferences")
		return
	}

	utils.WriteJSON(w, http.StatusOK, result)
}

func (h *NotificationHandler) GetSurveyLogs(w http.ResponseWriter, r *http.Request) {
	surveyID := chi.URLParam(r, "surveyId")

	logs, err := h.notifService.GetSurveyNotificationLogs(surveyID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get logs")
		return
	}

	utils.WriteJSON(w, http.StatusOK, logs)
}
