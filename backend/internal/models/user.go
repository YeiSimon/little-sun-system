package models

import (
	"time"
	"gorm.io/gorm"
)

/// User 代表系統使用者，使用 GORM 標籤
type User struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Picture   string         `gorm:"type:text" json:"picture"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	LastLogin time.Time      `gorm:"autoUpdateTime:false" json:"last_login"`
	Sessions  []Session      `gorm:"foreignKey:UserID" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 軟刪除
}

// Session 代表使用者會話
type Session struct {
	ID        string         `gorm:"primaryKey;type:uuid" json:"id"`
	UserID    string         `gorm:"size:255;not null;index" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
	IP        string         `gorm:"size:45" json:"ip"`
	UserAgent string         `gorm:"type:text" json:"user_agent"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 軟刪除
}

// AuthResponse 是認證請求的回應
type AuthResponse struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	IsLoggedIn bool   `json:"isLoggedIn"`
}

// TokenRequest 是 Token 驗證的請求
type TokenRequest struct {
	Credential string `json:"credential"`
}