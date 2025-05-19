package configs

import (
	"fmt"
	"log"
	"time"

	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"backend/pkg/utils"
	"backend/internal/models"
)

// DBConfig 資料庫連接配置
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DefaultDBConfig 返回預設資料庫配置
func DefaultDBConfig() *DBConfig {
	config := &DBConfig{
		Host:     utils.GetEnv("DB_HOST", "localhost"),
		User:     utils.GetEnv("DB_USER", "postgres"),
		Password: utils.GetEnv("DB_PASSWORD", "password"),
		DBName:   utils.GetEnv("DB_NAME", "userauth"),
		SSLMode:  utils.GetEnv("DB_SSL_MODE", "disable"),
	}

	portStr := utils.GetEnv("DB_PORT", "5432")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("警告: 無法解析 DB_PORT 環境變數，使用預設值 5432: %v", err)
		port = 5432
	}
	config.Port = port

	return config
}

// SetupDB 初始化 GORM 資料庫連接
func SetupDB(config *DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	// 設定自訂 GORM 記錄器
	gormLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,     // 慢查詢閾值
			LogLevel:                  logger.Info,     // 記錄等級，生產環境可改為 Warn
			IgnoreRecordNotFoundError: true,            // 忽略 ErrRecordNotFound 錯誤
			Colorful:                  true,            // 啟用彩色輸出
		},
	)

	// 連接到資料庫
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().Local() // 使用本地時間
		},
	})

	if err != nil {
		return nil, fmt.Errorf("資料庫連接失敗: %w", err)
	}

	// 自動遷移結構到資料庫
	if err := db.AutoMigrate(&models.User{}, &models.Session{}); err != nil {
		return nil, fmt.Errorf("資料庫遷移失敗: %w", err)
	}

	// 取得通用資料庫物件以設定連接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("取得 DB 連接失敗: %w", err)
	}

	// 設定連接池參數
	sqlDB.SetMaxIdleConns(10)     // 最大閒置連接數
	sqlDB.SetMaxOpenConns(100)    // 最大開啟連接數
	sqlDB.SetConnMaxLifetime(time.Hour) // 連接最大生命週期

	log.Println("資料庫連接成功")
	return db, nil
}