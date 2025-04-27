package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	"backend/internal/services"
)

// AuthMiddleware 處理身份驗證中間件
type AuthMiddleware struct {
	store       *sessions.CookieStore
	authService *service.AuthService
}

// NewAuthMiddleware 創建一個新的身份驗證中間件
func NewAuthMiddleware(key []byte, authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		store:       sessions.NewCookieStore(key),
		authService: authService,
	}
}

// AuthRequired 是檢查用戶是否已通過身份驗證的中間件
func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 獲取會話
		session, err := m.store.Get(c.Request, "user-session")
		if err != nil {
			log.Printf("獲取會話失敗: %v", err)
			m.clearSession(c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的會話"})
			c.Abort()
			return
		}
		
		// 檢查用戶是否已通過身份驗證
		auth, ok := session.Values["auth"].(bool)
		if !ok || !auth {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "需要身份驗證"})
			c.Abort()
			return
		}
		
		// 驗證數據庫中的會話
		sessionID, ok := session.Values["session_id"].(string)
		if !ok || sessionID == "" {
			m.clearSession(c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的會話"})
			c.Abort()
			return
		}
		
		// 使用服務層驗證會話
		valid, err := m.authService.ValidateSession(sessionID)
		if err != nil {
			log.Printf("驗證會話失敗: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服務器錯誤"})
			c.Abort()
			return
		}
		
		// 會話無效
		if !valid {
			m.clearSession(c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "會話已過期"})
			c.Abort()
			return
		}
		
		// 刷新會話超時 (可選，根據需求)
		go func() {
			if err := m.authService.RefreshSession(sessionID); err != nil {
				log.Printf("刷新會話失敗: %v", err)
			}
		}()
		
		// 會話有效，繼續處理請求
		c.Next()
	}
}

// GetSessionStore 返回會話存儲
func (m *AuthMiddleware) GetSessionStore() *sessions.CookieStore {
	return m.store
}

// clearSession 清除會話 Cookie
func (m *AuthMiddleware) clearSession(c *gin.Context) {
	session, _ := m.store.Get(c.Request, "user-session")
	session.Options.MaxAge = -1
	session.Save(c.Request, c.Writer)
}