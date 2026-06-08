package models

import (
	"time"

	"gorm.io/gorm"
)

type SurveyResponse struct {
	ID        string         `json:"id" gorm:"size:36;primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	SurveyID   string `json:"survey_id" gorm:"size:36;index;not null"`
	UserID     *uint  `json:"user_id,omitempty" gorm:"index"`
	SessionID  string `json:"session_id,omitempty" gorm:"size:64;index"`
	Anonymous  bool   `json:"anonymous"`
	StartedAt  *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Status     string `json:"status" gorm:"size:20;default:in_progress"`

	Answers []Answer `json:"answers" gorm:"foreignKey:ResponseID"`
}

type Answer struct {
	ID         string `json:"id" gorm:"size:36;primarykey"`
	ResponseID string `json:"response_id" gorm:"size:36;index;not null"`
	QuestionID string `json:"question_id" gorm:"size:36;not null"`
	Value      JSONMap `json:"value" gorm:"type:jsonb"`
}

type ResponseSummary struct {
	ID        string    `json:"id"`
	SurveyID  string    `json:"survey_id"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
}
