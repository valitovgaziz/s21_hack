package handlers

import (
	"encoding/json"
	"net/http"

	"api_hr/internal/service"
	"api_hr/internal/utils"

	"github.com/go-chi/chi/v5"
)

type ResponseHandler struct {
	responseService *service.ResponseService
}

func NewResponseHandler(responseService *service.ResponseService) *ResponseHandler {
	return &ResponseHandler{responseService: responseService}
}

func (h *ResponseHandler) StartSurvey(w http.ResponseWriter, r *http.Request) {
	surveyID := chi.URLParam(r, "surveyId")

	var req struct {
		SessionID string `json:"session_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.SessionID = ""
	}

	userID := utils.GetUserIDFromContext(r.Context())

	resp, err := h.responseService.StartSurvey(surveyID, userID, req.SessionID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *ResponseHandler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	responseID := chi.URLParam(r, "responseId")

	var req struct {
		QuestionID string      `json:"question_id"`
		Value      interface{} `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.responseService.SubmitAnswer(responseID, req.QuestionID, req.Value); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *ResponseHandler) CompleteSurvey(w http.ResponseWriter, r *http.Request) {
	responseID := chi.URLParam(r, "responseId")
	if err := h.responseService.CompleteSurvey(responseID); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "completed"})
}

func (h *ResponseHandler) GetResponse(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	resp, err := h.responseService.GetResponse(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Response not found")
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *ResponseHandler) GetSurveyResponses(w http.ResponseWriter, r *http.Request) {
	surveyID := chi.URLParam(r, "surveyId")
	responses, err := h.responseService.GetSurveyResponses(surveyID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	utils.WriteJSON(w, http.StatusOK, responses)
}

func (h *ResponseHandler) CheckResponded(w http.ResponseWriter, r *http.Request) {
	surveyID := chi.URLParam(r, "surveyId")
	userID := utils.GetUserIDFromContext(r.Context())
	if userID == nil {
		utils.WriteJSON(w, http.StatusOK, map[string]bool{"responded": false})
		return
	}
	responded, err := h.responseService.HasUserResponded(*userID, surveyID)
	if err != nil {
		responded = false
	}
	utils.WriteJSON(w, http.StatusOK, map[string]bool{"responded": responded})
}


