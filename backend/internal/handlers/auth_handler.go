package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	"backend/internal/services"
	"backend/internal/models"
)

// AuthHandler 處理身份驗證相關的 HTTP 請求
type AuthHandler struct {
	authService *service.AuthService
	store       *sessions.CookieStore
}

// NewAuthHandler 創建一個新的身份驗證處理器
func NewAuthHandler(authService *service.AuthService, store *sessions.CookieStore) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		store:       store,
	}
}

// HandleGoogleSignIn 處理 Google 登入請求
func (h *AuthHandler) HandleGoogleSignIn(c *gin.Context) {
	var req models.TokenRequest
	
	// 解析請求體
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "無法解析請求",
			"details": err.Error(),
		})
		return
	}
	
	// 驗證請求參數
	if req.Credential == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少憑證"})
		return
	}
	
	log.Printf("收到 Google ID Token，長度: %d 字符", len(req.Credential))
	
	// 驗證 Google Token
	ctx := context.Background()
	payload, err := h.authService.ValidateGoogleToken(ctx, req.Credential)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "無效的憑證",
			"details": err.Error(),
		})
		return
	}
	
	// 提取用戶信息
	sub, _ := payload["sub"].(string)          // Google 提供的唯一 ID
	email, _ := payload["email"].(string)
	name, _ := payload["name"].(string)
	picture, _ := payload["picture"].(string)
	expTime, _ := payload["exp"].(float64)     // Token 過期時間
	
	// 處理用戶身份驗證 (查找或創建用戶)
	userID, err := h.authService.AuthenticateUser(email, name, picture, sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "身份驗證失敗",
			"details": err.Error(),
		})
		return
	}
	
	// 創建用戶會話
	sessionID, expiresAt, err := h.authService.CreateUserSession(
		userID, c.ClientIP(), c.Request.UserAgent(), expTime,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "創建會話失敗",
			"details": err.Error(),
		})
		return
	}
	
	// 設置會話 Cookie
	webSession, _ := h.store.Get(c.Request, "user-session")
	webSession.Values["user_id"] = userID
	webSession.Values["email"] = email
	webSession.Values["name"] = name
	webSession.Values["picture"] = picture
	webSession.Values["auth"] = true
	webSession.Values["session_id"] = sessionID
	
	// 計算 Cookie 過期時間
	maxAge := int(expiresAt.Sub(time.Now()).Seconds())
	if maxAge <= 0 {
		maxAge = 3600 * 24 * 30// 預設 24 * 30 小時
	}
	
	// 配置 Cookie 選項
	webSession.Options.MaxAge = maxAge
	webSession.Options.Path = "/"
	webSession.Options.HttpOnly = true
	webSession.Options.Secure = false // 生產環境設為 true
	webSession.Options.SameSite = http.SameSiteLaxMode
	
	// 保存會話
	if err := webSession.Save(c.Request, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "儲存會話失敗",
			"details": err.Error(),
		})
		return
	}
	
	// 獲取活躍會話數 (可用於限制並發登入)
	activeSessions, _ := h.authService.GetUserActiveSessions(userID)
	
	// 返回成功響應
	c.JSON(http.StatusOK, gin.H{
		"email":          email,
		"name":           name,
		"picture":        picture,
		"isLoggedIn":     true,
		"expire_session": expiresAt,
		"activeSessions": activeSessions,
	})
}

// HandleLogout 處理登出請求
func (h *AuthHandler) HandleLogout(c *gin.Context) {
	// 獲取會話
	session, _ := h.store.Get(c.Request, "user-session")
	
	// 檢查會話 ID 是否存在
	sessionID, ok := session.Values["session_id"].(string)
	if ok && sessionID != "" {
		// 調用服務層登出用戶
		if err := h.authService.LogoutUser(sessionID); err != nil {
			log.Printf("從資料庫刪除會話失敗: %v", err)
		}
	}
	
	// 清除 Cookie
	session.Options.MaxAge = -1
	if err := session.Save(c.Request, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登出失敗"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"logout": true})
}

// HandleGetProfile 處理獲取用戶資料請求
func (h *AuthHandler) HandleGetProfile(c *gin.Context) {
	session, _ := h.store.Get(c.Request, "user-session")
	
	// 獲取會話中的用戶資訊
	email, _ := session.Values["email"].(string)
	name, _ := session.Values["name"].(string)
	picture, _ := session.Values["picture"].(string)
	userID, _ := session.Values["user_id"].(string)
	
	// 獲取活躍會話數 (可選)
	var activeSessions int64
	if userID != "" {
		activeSessions, _ = h.authService.GetUserActiveSessions(userID)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"email":          email,
		"name":           name,
		"picture":        picture,
		"activeSessions": activeSessions,
	})
}