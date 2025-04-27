package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"backend/internal/repository"
	"backend/internal/models"
	"backend/internal/auth"
)

// AuthService 提供身份驗證相關的業務邏輯
type AuthService struct {
	googleAuth *auth.GoogleAuth
	userRepo   *repository.UserRepository
}

// NewAuthService 創建一個新的身份驗證服務
func NewAuthService(googleAuth *auth.GoogleAuth, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		googleAuth: googleAuth,
		userRepo:   userRepo,
	}
}

// ValidateGoogleToken 驗證 Google JWT token
func (s *AuthService) ValidateGoogleToken(ctx context.Context, token string) (map[string]interface{}, error) {
	return s.googleAuth.ValidateIDToken(ctx, token)
}

// AuthenticateUser 處理使用者身份驗證，返回用戶 ID
func (s *AuthService) AuthenticateUser(email, name, picture, subID string) (string, error) {
	// 查找現有使用者
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", fmt.Errorf("查詢使用者失敗: %w", err)
	}
	
	if user == nil {
		// 創建新使用者
		newUser := &models.User{
			ID:        subID,
			Email:     email,
			Name:      name,
			Picture:   picture,
			LastLogin: time.Now(),
		}
		
		if err := s.userRepo.CreateUser(newUser); err != nil {
			return "", fmt.Errorf("創建使用者失敗: %w", err)
		}
		
		log.Printf("新使用者註冊: %s (%s)", name, email)
		return subID, nil
	}
	
	// 更新現有使用者
	user.Name = name
	user.Picture = picture
	user.LastLogin = time.Now()
	
	if err := s.userRepo.UpdateUser(user); err != nil {
		return "", fmt.Errorf("更新使用者失敗: %w", err)
	}
	
	log.Printf("使用者登入: %s (%s)", name, email)
	return user.ID, nil
}

// CreateUserSession 創建使用者會話
func (s *AuthService) CreateUserSession(userID, ip, userAgent string, expTime float64) (string, time.Time, error) {
	// 產生唯一會話 ID
	sessionID := uuid.New().String()
	
	// 計算過期時間
	expiresAt := time.Now().Add(30 * 24 * time.Hour) // 預設 1 個月
	
	// if expTime > 0 {
	// 	googleExpireTime := time.Unix(int64(expTime), 0)//確定google身分驗證是否一致?
    //     // 可以決定是否採用 Google 的過期時間
    //     if googleExpireTime.After(time.Now()) {
    //         expiresAt = googleExpireTime
    //     }
	// }
	
	// 建立會話記錄
	session := &models.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
		IP:        ip,
		UserAgent: userAgent,
	}
	
	// 儲存到資料庫
	if err := s.userRepo.CreateSession(session); err != nil {
		return "", time.Time{}, fmt.Errorf("創建會話失敗: %w", err)
	}
	
	// 清理過期會話 (非同步執行)
	go func() {
		if err := s.userRepo.DeleteExpiredSessions(); err != nil {
			log.Printf("清理過期會話失敗: %v", err)
		}
	}()
	
	return sessionID, expiresAt, nil
}

// ValidateSession 驗證會話有效性
func (s *AuthService) ValidateSession(sessionID string) (bool, error) {
	// 獲取會話
	session, err := s.userRepo.GetSessionByID(sessionID)
	if err != nil {
		return false, fmt.Errorf("查詢會話失敗: %w", err)
	}
	
	// 會話不存在
	if session == nil {
		return false, nil
	}
	
	// 檢查會話是否過期
	if time.Now().After(session.ExpiresAt) {
		// 刪除過期會話
		if err := s.userRepo.DeleteSession(sessionID); err != nil {
			log.Printf("刪除過期會話失敗: %v", err)
		}
		return false, nil
	}
	
	// 會話有效
	return true, nil
}

// RefreshSession 延長會話有效期
func (s *AuthService) RefreshSession(sessionID string) error {
	// 新的過期時間 (1 小時後)
	newExpiryTime := time.Now().Add(30 * 24 * time.Hour)
	
	// 更新資料庫中的過期時間
	if err := s.userRepo.UpdateSessionExpiry(sessionID, newExpiryTime); err != nil {
		return fmt.Errorf("更新會話失敗: %w", err)
	}
	
	return nil
}

// LogoutUser 登出使用者
func (s *AuthService) LogoutUser(sessionID string) error {
	// 從資料庫刪除會話
	if err := s.userRepo.DeleteSession(sessionID); err != nil {
		return fmt.Errorf("登出失敗: %w", err)
	}
	
	return nil
}

// GetUserActiveSessions 獲取使用者活躍會話數
func (s *AuthService) GetUserActiveSessions(userID string) (int64, error) {
	count, err := s.userRepo.CountUserSessions(userID)
	if err != nil {
		return 0, fmt.Errorf("查詢活躍會話失敗: %w", err)
	}
	
	return count, nil
}