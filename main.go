package main

import (
	"encoding/json"
	"context"
	"fmt"
	"time"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"github.com/google/uuid"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"google.golang.org/api/idtoken"
	"github.com/gorilla/sessions"
)

type User struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Picture   string `json:"picture"`
	SessionID string `json:"session_id"`
}

var (
	googleOauthConfig *oauth2.Config
)

func init() {
	// 讀取 client_secret.json
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("無法讀取 client secret 文件: %v", err)
	}

	// 設定 OAuth2 配置
	config, err := google.ConfigFromJSON(b,
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile")
	if err != nil {
		log.Fatalf("無法將 client secret 解析為配置: %v", err)
	}

	googleOauthConfig = config
}

func main() {
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.Use(cors.New(cors.Config{
        // Specific origin instead of wildcard *
        AllowOrigins:     []string{"http://localhost:4200"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        // This is crucial for credentials
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

	// 只保留處理 Google Sign-In 的 POST 路由
	r.POST("/api/login/google", handleGoogleSignInPost)
	r.GET("/api/logout/google", handleLogout)
	r.GET("/api/sheets", SearchCustomerHandler)

	protected := r.Group("/api")
	protected.Use(AuthRequired())
	fmt.Println("服務器運行在 http://localhost:8080")
	r.Run(":8080")
}

var store = sessions.NewCookieStore([]byte("123456"))
// 處理 POST /api/login/google - 由 Google Sign-In 按鈕直接調用
func handleGoogleSignInPost(c *gin.Context) {
        var req struct {
            Credential string `json:"credential"`
        }
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "無法解析 JSON 數據", "details": err.Error()})
            return
        }
    
        if req.Credential == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "缺少憑證"})
            return
        }
    
        log.Printf("收到 Google ID token，長度: %d 字符", len(req.Credential))
    
        payload, err := validateGoogleIDToken(req.Credential, googleOauthConfig.ClientID)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的憑證", "details": err.Error()})
            return
        }
    
		sub, _ := payload["sub"].(string)        // 使用 sub 作為唯一識別符
		email, _ := payload["email"].(string)
		name, _ := payload["name"].(string)
		picture, _ := payload["picture"].(string)
		emailVerified, _ := payload["email_verified"].(bool)
		expTime, _ := payload["exp"].(float64)   // 獲取 token 過期時間

		if !emailVerified{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "電子郵件未驗證"})
			return
		}

		log.Printf("payload詳細資訊: %s", payload)
        log.Printf("用戶登錄成功: %s (%s)", name, email)
		
		session, _ := store.Get(c.Request, "user-session")

		// 設置 session 值
		session.Values["user_id"] = sub
		session.Values["email"] = email
		session.Values["name"] = name
		session.Values["picture"] = picture
		session.Values["auth"] = true

		maxAge := 3600 * 24 * 30

		if expTime > 0 {
			expiresIn := int(expTime - float64(time.Now().Unix()))
			if expiresIn > 0 {
				maxAge = expiresIn
			}
		}
		session.Options.MaxAge = maxAge
		session.Options.Path = "/"
		session.Options.HttpOnly = true
		session.Options.Secure = false

		if err:=session.Save(c.Request, c.Writer);err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存 session 失敗", "details": err.Error()})
        	return
		}

		fileName := "users.txt"
		users := make(map[string]User)

		if data, err := json.MarshalIndent(users, "", "  "); err == nil {
        if err := ioutil.WriteFile(fileName, data, 0644); err != nil {
            log.Printf("警告: 寫入使用者資料失敗: %v", err)
        }
    }

		sessionId := uuid.New().String()
		users[email] = User{
			Email:     email,
			Name:      name,
			Picture:   picture,
			SessionID: sessionId,
		}
	
		// 將更新後的資料存回檔案
		if data, err := json.MarshalIndent(users, "", "  "); err == nil {
			if err := ioutil.WriteFile(fileName, data, 0644); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "寫入使用者資料失敗"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "序列化使用者資料失敗"})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{
			"email": email,
			"name": name,
			"picture": picture,
			"isLoggedIn": true,
		})
    }
// 驗證 Google ID token
func validateGoogleIDToken(idToken string, clientID string) (map[string]interface{}, error) {
	ctx := context.Background()

	// 驗證 ID token
	payload, err := idtoken.Validate(ctx, idToken, clientID)
	if err != nil {
		return nil, err
	}

	return payload.Claims, nil
}

func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        session, _ := store.Get(c.Request, "user-session")
        
        // 檢查用戶是否已認證
        auth, ok := session.Values["auth"].(bool)
        if !ok || !auth {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "需要登入"})
            c.Abort()
            return
        }
        
        // 用戶已認證，繼續下一個處理器
        c.Next()
    }
}

func handleGetProfile(c *gin.Context) {
    session, _ := store.Get(c.Request, "user-session")
    
    email := session.Values["email"].(string)
    name := session.Values["name"].(string)
    picture := session.Values["picture"].(string)
    
    c.JSON(http.StatusOK, gin.H{
        "email": email,
        "name": name,
        "picture": picture,
    })
}

func SearchCustomerHandler(c *gin.Context) {
	// 從 query 參數取得客戶名
	customerName := c.Query("customer")
	if customerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請提供 customer 查詢參數"})
		return
	}

	// 設定服務帳戶 JSON 金鑰、試算表 ID 與範圍
	serviceAccountFile := "little-sun-system-d5e3eda49d9f.json"
	spreadsheetId := "10IIJuGiur0HGpvjAippllfg1XhYq_wIHwR4_xWn-z_c"
	// 假設你的表單有 15 欄（A 至 O），範圍可調整
	readRange := "客戶細項!A1:O"

	// 建立 Sheets API 服務物件
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(serviceAccountFile))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("建立 Sheets 服務失敗: %v", err)})
		return
	}

	// 取得指定範圍內的資料
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("讀取試算表失敗: %v", err)})
		return
	}

	// 確認至少有一列標題與一筆資料
	if len(resp.Values) < 2 {
		c.JSON(http.StatusOK, gin.H{"message": "未找到資料"})
		return
	}

	// 取得第一列作為 header，找出「客戶名」所在的欄位索引
	header := resp.Values[0]
	customerNameIdx := -1
	for i, col := range header {
		if colStr, ok := col.(string); ok && colStr == "客戶名" {
			customerNameIdx = i
			break
		}
	}
	if customerNameIdx == -1 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "找不到 [客戶名] 欄位"})
		return
	}

	// 遍歷資料列，找出符合客戶名的資料
	var results []interface{}
	for _, row := range resp.Values[1:] {
		if len(row) > customerNameIdx {
			if val, ok := row[customerNameIdx].(string); ok && val == customerName {
				results = append(results, row)
			}
		}
	}

	if len(results) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "查無此客戶資料"})
		return
	}

	// 回傳搜尋到的資料列
	c.JSON(http.StatusOK, gin.H{"data": results})
}

// 登出處理函數
func handleLogout(c *gin.Context) {
    session, _ := store.Get(c.Request, "user-session")
    
    // 將認證狀態設為 false
    session.Values["auth"] = false
    
    // 刪除 session（設置過期時間為 -1）
    session.Options.MaxAge = -1
    
    // 保存 session（實際上是刪除）
    if err := session.Save(c.Request, c.Writer); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "登出失敗"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"logout": true})
}