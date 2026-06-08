package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type SurveyStatus string

const (
	SurveyDraft     SurveyStatus = "draft"
	SurveyActive    SurveyStatus = "active"
	SurveyCompleted SurveyStatus = "completed"
	SurveyArchived  SurveyStatus = "archived"
)

type PrivacyMode string

const (
	PrivacyAnonymous  PrivacyMode = "anonymous"
	PrivacyIdentified PrivacyMode = "identified"
)

type QuestionType string

const (
	QTSingleChoice   QuestionType = "single_choice"
	QTMultipleChoice QuestionType = "multiple_choice"
	QTScale          QuestionType = "scale"
	QTText           QuestionType = "text"
	QTMatrix         QuestionType = "matrix"
)

type BranchCondition struct {
	QuestionID  string      `json:"question_id"`
	Operator    string      `json:"operator"`
	Value       interface{} `json:"value"`
	TargetBlock string      `json:"target_block"`
}

type Question struct {
	ID          string            `json:"id" gorm:"size:36;primarykey"`
	SurveyID    string            `json:"survey_id" gorm:"size:36;index"`
	Title       string            `json:"title" gorm:"size:500;not null"`
	Description string            `json:"description,omitempty" gorm:"size:1000"`
	Type        QuestionType      `json:"type" gorm:"size:50"`
	Order       int               `json:"order"`
	Required    bool              `json:"required"`
	Options     JSONMap           `json:"options" gorm:"type:jsonb"`
	BranchLogic JSONBranchRule    `json:"branch_logic" gorm:"type:jsonb"`
	BlockID     string            `json:"block_id,omitempty" gorm:"size:36"`
}

type JSONBranchRule []BranchCondition

func (j JSONBranchRule) Value() (driver.Value, error) {
	return json.Marshal(j)
}
func (j *JSONBranchRule) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

type Survey struct {
	ID             string         `json:"id" gorm:"size:36;primarykey"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	Title          string         `json:"title" gorm:"size:255;not null"`
	Description    string         `json:"description" gorm:"size:2000"`
	Status         SurveyStatus   `json:"status" gorm:"size:50;default:draft"`
	Privacy        PrivacyMode    `json:"privacy" gorm:"size:50"`
	StartDate      *time.Time     `json:"start_date"`
	EndDate        *time.Time     `json:"end_date"`
	TargetRoles    JSONStringSlice `json:"target_roles" gorm:"type:jsonb"`
	AllowAnonymous bool           `json:"allow_anonymous"`
	CreatedBy      uint           `json:"created_by" gorm:"index"`
	Questions      []Question     `json:"questions" gorm:"foreignKey:SurveyID"`
	BlockStructure JSONMap        `json:"block_structure" gorm:"type:jsonb"`
}

type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}
func (j *JSONMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

type JSONStringSlice []string

func (j JSONStringSlice) Value() (driver.Value, error) {
	return json.Marshal(j)
}
func (j *JSONStringSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

type SurveyAnalytics struct {
	SurveyID      string             `json:"survey_id"`
	Title         string             `json:"title"`
	TotalSent     int64              `json:"total_sent"`
	TotalStarted  int64              `json:"total_started"`
	TotalCompleted int64             `json:"total_completed"`
	CompletionRate float64           `json:"completion_rate"`
	AverageScore  *float64           `json:"average_score"`
	ENPS          *float64           `json:"enps"`
	QuestionStats []QuestionStat     `json:"question_stats"`
	DailyStats    []DailyStat        `json:"daily_stats"`
	ChannelStats  map[string]ChannelStat `json:"channel_stats"`
}

type QuestionStat struct {
	QuestionID string            `json:"question_id"`
	Title      string            `json:"title"`
	Type       string            `json:"type"`
	Answers    map[string]int64  `json:"answers"`
	Total      int64             `json:"total"`
}

type DailyStat struct {
	Date   string `json:"date"`
	Count  int64  `json:"count"`
}

type ChannelStat struct {
	Sent     int64   `json:"sent"`
	Delivered int64  `json:"delivered"`
	Opened   int64   `json:"opened"`
	Clicked  int64   `json:"clicked"`
	CTR      float64 `json:"ctr"`
	Cost     float64 `json:"cost"`
}
