package repository

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"backend/internal/models"
)

// UserRepository 提供使用者相關的資料存取方法
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 創建一個新的使用者資料庫存取層
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetUserByEmail 透過電子郵件查找使用者
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		// 如果記錄不存在，返回 nil 而非錯誤
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查詢使用者失敗: %w", result.Error)
	}
	
	return &user, nil
}

// GetUserByID 透過 ID 查找使用者
func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	
	result := r.db.First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查詢使用者失敗: %w", result.Error)
	}
	
	return &user, nil
}

// CreateUser 創建新使用者
func (r *UserRepository) CreateUser(user *models.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		return fmt.Errorf("創建使用者失敗: %w", result.Error)
	}
	
	return nil
}

// UpdateUser 更新現有使用者
func (r *UserRepository) UpdateUser(user *models.User) error {
	// 只更新特定欄位，避免覆蓋其他欄位
	result := r.db.Model(user).
		Select("name", "picture", "last_login").
		Updates(map[string]interface{}{
			"name":       user.Name,
			"picture":    user.Picture,
			"last_login": time.Now(),
		})
	
	if result.Error != nil {
		return fmt.Errorf("更新使用者失敗: %w", result.Error)
	}
	
	return nil
}

// CreateSession 創建新會話
func (r *UserRepository) CreateSession(session *models.Session) error {
	result := r.db.Create(session)
	if result.Error != nil {
		return fmt.Errorf("創建會話失敗: %w", result.Error)
	}
	
	return nil
}

// GetSessionByID 透過 ID 查找會話
func (r *UserRepository) GetSessionByID(id string) (*models.Session, error) {
	var session models.Session
	
	result := r.db.First(&session, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查詢會話失敗: %w", result.Error)
	}
	
	return &session, nil
}

// UpdateSessionExpiry 更新會話過期時間
func (r *UserRepository) UpdateSessionExpiry(sessionID string, expiresAt time.Time) error {
	result := r.db.Model(&models.Session{}).
		Where("id = ?", sessionID).
		Update("expires_at", expiresAt)
		
	if result.Error != nil {
		return fmt.Errorf("更新會話過期時間失敗: %w", result.Error)
	}
	
	return nil
}

// DeleteSession 刪除會話
func (r *UserRepository) DeleteSession(sessionID string) error {
	result := r.db.Where("id = ?", sessionID).Delete(&models.Session{})
	if result.Error != nil {
		return fmt.Errorf("刪除會話失敗: %w", result.Error)
	}
	
	return nil
}

// DeleteExpiredSessions 刪除所有過期會話
func (r *UserRepository) DeleteExpiredSessions() error {
	result := r.db.Where("expires_at < ?", time.Now()).Delete(&models.Session{})
	if result.Error != nil {
		return fmt.Errorf("刪除過期會話失敗: %w", result.Error)
	}
	
	return nil
}

// GetUserWithSessions 獲取用戶及其所有會話
func (r *UserRepository) GetUserWithSessions(userID string) (*models.User, error) {
	var user models.User
	
	result := r.db.Preload("Sessions").First(&user, "id = ?", userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查詢使用者及會話失敗: %w", result.Error)
	}
	
	return &user, nil
}

// CountUserSessions 計算使用者的活躍會話數
func (r *UserRepository) CountUserSessions(userID string) (int64, error) {
	var count int64
	
	result := r.db.Model(&models.Session{}).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Count(&count)
		
	if result.Error != nil {
		return 0, fmt.Errorf("計算使用者會話失敗: %w", result.Error)
	}
	
	return count, nil
}