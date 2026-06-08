package models

import (
	"time"

	"gorm.io/gorm"
)

type PushSubscription struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	UserID       uint   `json:"user_id" gorm:"index;not null"`
	Endpoint     string `json:"endpoint" gorm:"type:text;not null"`
	P256DH       string `json:"p256dh" gorm:"type:text;not null"`
	Auth         string `json:"auth" gorm:"type:text;not null"`
	DeviceName   string `json:"device_name" gorm:"size:255"`
	UserAgent    string `json:"user_agent" gorm:"type:text"`
	Active       bool   `json:"active" gorm:"default:true"`
}

type NotificationChannel string

const (
	ChannelWebPush  NotificationChannel = "web_push"
	ChannelSMS      NotificationChannel = "sms"
	ChannelTelegram NotificationChannel = "telegram"
	ChannelEmail    NotificationChannel = "email"
)

type NotificationPreference struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	UserID           uint               `json:"user_id" gorm:"index;unique;not null"`
	WebPushEnabled   bool               `json:"web_push_enabled" gorm:"default:true"`
	SMSEnabled       bool               `json:"sms_enabled" gorm:"default:true"`
	TelegramEnabled  bool               `json:"telegram_enabled" gorm:"default:false"`
	EmailEnabled     bool               `json:"email_enabled" gorm:"default:true"`
	PreferredTime    string             `json:"preferred_time" gorm:"size:20;default:day"`
	DNDStart         *time.Time         `json:"dnd_start"`
	DNDEnd           *time.Time         `json:"dnd_end"`
}

type NotificationLog struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	UserID     uint               `json:"user_id" gorm:"index;not null"`
	SurveyID   string             `json:"survey_id" gorm:"size:36;index;not null"`
	Channel    NotificationChannel `json:"channel" gorm:"size:50;not null"`
	Status     string             `json:"status" gorm:"size:20;not null"` // sent, delivered, opened, clicked, failed
	OpenedAt   *time.Time         `json:"opened_at"`
	ClickedAt  *time.Time         `json:"clicked_at"`
	CostCents  float64            `json:"cost_cents,omitempty"`
	ErrorMsg   string             `json:"error_msg,omitempty" gorm:"type:text"`
}
