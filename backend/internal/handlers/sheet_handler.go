package handlers

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// SearchCustomerHandler 處理客戶搜尋請求
func SearchCustomerHandler(c *gin.Context) {
	// 從查詢參數獲取客戶名
	customerName := c.Query("customer")
	if customerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請提供 customer 查詢參數"})
		return
	}

	// 設定服務帳號 JSON 金鑰、試算表 ID 與範圍
	serviceAccountFile := filepath.Join("pkg", "configs", "little-sun-system-d5e3eda49d9f.json")
	spreadsheetId := "10IIJuGiur0HGpvjAippllfg1XhYq_wIHwR4_xWn-z_c"
	readRange := "客戶細項!A1:Q"

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