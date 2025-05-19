package main

import (
	"fmt"
	"net/http" 
	"log"
	"path/filepath"
	"time"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/services"
	"backend/pkg/configs"
	"backend/internal/auth"
	"backend/pkg/utils"
)

func main() {
	// 設定日誌
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("啟動服務...")

	// 初始化資料庫連接
	log.Println("連接資料庫...")
	dbConfig := configs.DefaultDBConfig()
	db, err := configs.SetupDB(dbConfig)
	if err != nil {
		log.Fatalf("資料庫初始化失敗: %v", err)
	}

	// 確保資料庫連接在程式結束時關閉
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("獲取資料庫連接失敗: %v", err)
	}
	defer sqlDB.Close()

	// 初始化儲存庫
	userRepo := repository.NewUserRepository(db)

	// 設定 Google Auth
	clientSecretPath := filepath.Join("pkg", "configs", "client_secret.json")
	googleAuth, err := auth.NewGoogleAuth(clientSecretPath)
	if err != nil {
		log.Fatalf("Google Auth 設定失敗: %v", err)
	}

	// 設定服務層
	authService := service.NewAuthService(googleAuth, userRepo)

	// 設定 session 存儲
	// 注意: 在生產環境中，應使用加密的金鑰，並考慮使用 Redis 等外部存儲
	sessionKey := []byte("verySecureSessionKey123456") // 生產環境應使用環境變數或密鑰管理服務
	sessionStore := sessions.NewCookieStore(sessionKey)
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 7 * 30, // 預設 1 小時
		HttpOnly: true,
		Secure:   false, // 生產環境設為 true
		SameSite: http.SameSiteLaxMode,
	}

	// 設定中間件
	authMiddleware := middleware.NewAuthMiddleware(sessionKey, authService)

	// 設定處理器
	authHandler := handlers.NewAuthHandler(authService, sessionStore)

	// 創建 Gin 引擎
	r := gin.Default()
	
	// 設定可信代理
	r.SetTrustedProxies([]string{"127.0.0.1", "172.16.0.0/12", "192.168.0.0/16"})

	// 設定 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{ "http://localhost:4200", 
        "http://frontend:4200",
        utils.GetEnv("FRONTEND_URL", "http://localhost:4200")}, // 前端網址
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // 允許攜帶憑證
		MaxAge:           12 * time.Hour,
	}))

	// 公開路由
	r.POST("/api/login/google", authHandler.HandleGoogleSignIn)
	r.GET("/api/logout", authHandler.HandleLogout)
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	// 受保護路由
	api := r.Group("/api")
	api.Use(authMiddleware.AuthRequired())
	{
		api.GET("/profile", authHandler.HandleGetProfile)
		api.GET("/sheets", handlers.SearchCustomerHandler)
		
		// 可以添加更多受保護的路由
	}

	// 啟動服務器
	port := utils.GetEnv("PORT", "8080")
	portInt, _ := strconv.Atoi(port)
	log.Printf("服務器啟動於 :%d", portInt)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("服務器啟動失敗: %v", err)
	}
}