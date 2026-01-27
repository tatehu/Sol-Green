package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sol-green/config"
	"sol-green/controller"
	"sol-green/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	
	api := r.Group("/api/v1")
	{
		api.POST("/auth/wallet", controller.WalletAuth)
		api.POST("/green/behavior/submit", controller.SubmitGreenBehavior)
		api.GET("/green/behavior/:id", controller.GetBehaviorStatus)
	}
	
	return r
}

func TestWalletAuth(t *testing.T) {
	config.InitConfig()
	
	router := setupRouter()
	
	reqBody := map[string]string{
		"wallet_addr": "TestWallet123",
		"signature":   "test_signature",
		"message":     "test_message",
	}
	jsonValue, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", "/api/v1/auth/wallet", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response["token"])
}

func TestSubmitGreenBehavior(t *testing.T) {
	config.InitConfig()
	config.InitDB()
	config.InitRedis()
	config.InitAIFraudDetector()
	
	router := setupRouter()
	
	// 先登录获取 token
	authReq := map[string]string{
		"wallet_addr": "TestWallet123",
		"signature":   "test_signature",
		"message":     "test_message",
	}
	authJson, _ := json.Marshal(authReq)
	authHttpReq, _ := http.NewRequest("POST", "/api/v1/auth/wallet", bytes.NewBuffer(authJson))
	authHttpReq.Header.Set("Content-Type", "application/json")
	authW := httptest.NewRecorder()
	router.ServeHTTP(authW, authHttpReq)
	
	var authResponse map[string]interface{}
	json.Unmarshal(authW.Body.Bytes(), &authResponse)
	token := authResponse["token"].(string)
	
	// 提交行为
	reqBody := model.GreenBehaviorReq{
		BehaviorType: "waste_sorting",
		MediaURLs:    []string{"https://example.com/image1.jpg"},
		Location:     "测试地点",
	}
	jsonValue, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", "/api/v1/green/behavior/submit", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response["id"])
}
