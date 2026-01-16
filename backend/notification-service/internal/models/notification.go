package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"not null;index:idx_notifications_user_read" json:"user_id"`
	Type      string     `gorm:"not null;size:50" json:"type"`
	Title     string     `gorm:"not null;size:255" json:"title"`
	Message   string     `gorm:"not null;type:text" json:"message"`
	Data      string     `gorm:"type:text" json:"data"`
	Channel   string     `gorm:"not null;size:20;default:in_app" json:"channel"`
	IsRead    bool       `gorm:"default:false;index:idx_notifications_user_read" json:"is_read"`
	SentAt    *time.Time `json:"sent_at"`
	ReadAt    *time.Time `json:"read_at"`
	gorm.Model
}

func (n *Notification) ValidateTitle() error {
	if n.Title == "" {
		return errors.New("title cannot be empty")
	}
	if len(n.Title) > 255 {
		return errors.New("title cannot exceed 255 characters")
	}
	return nil
}

func (n *Notification) ValidateMessage() error {
	if n.Message == "" {
		return errors.New("message cannot be empty")
	}
	return nil
}

func (n *Notification) ValidateUserID() error {
	if n.UserID == 0 {
		return errors.New("user ID cannot be empty")
	}
	return nil
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if err := n.ValidateUserID(); err != nil {
		return err
	}
	if err := n.ValidateTitle(); err != nil {
		return err
	}
	if err := n.ValidateMessage(); err != nil {
		return err
	}
	if n.Channel == "" {
		n.Channel = "in_app"
	}
	return nil
}

type NotificationPreference struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	UserID         uint   `gorm:"not null;uniqueIndex:idx_user_channel" json:"user_id"`
	Channel        string `gorm:"not null;size:20;uniqueIndex:idx_user_channel" json:"channel"`
	Enabled        bool   `gorm:"default:true" json:"enabled"`
	TelegramChatID int64  `json:"telegram_chat_id"`
	Email          string `gorm:"size:255" json:"email"`
	gorm.Model
}

func (p *NotificationPreference) ValidateUserID() error {
	if p.UserID == 0 {
		return errors.New("user ID cannot be empty")
	}
	return nil
}

func (p *NotificationPreference) ValidateChannel() error {
	if p.Channel == "" {
		return errors.New("channel cannot be empty")
	}
	return nil
}

func (p *NotificationPreference) BeforeCreate(tx *gorm.DB) error {
	if err := p.ValidateUserID(); err != nil {
		return err
	}
	return p.ValidateChannel()
}
