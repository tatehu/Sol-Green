package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sol-green/config"
	"sol-green/controller"
	"sol-green/middleware"
	"sol-green/model"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := r.Group("/api/v1")
	{
		api.POST("/auth/wallet", controller.WalletAuth)
		api.POST("/green/behavior/submit", middleware.Auth(), controller.SubmitGreenBehavior)
		api.GET("/green/behavior/:id", middleware.Auth(), controller.GetBehaviorStatus)
		api.POST("/challenges", middleware.Auth(), controller.CreateChallenge)
		api.POST("/challenges/:id/activate", middleware.Auth(), controller.ActivateChallenge)
		api.POST("/challenges/:id/join", middleware.Auth(), controller.JoinChallenge)
		api.POST("/challenges/:id/claim", middleware.Auth(), controller.ClaimChallengeReward)
	}

	return r
}

func TestChallengeClaimUsesBonus(t *testing.T) {
	config.InitConfig()
	config.InitDB()
	config.InitRedis()
	config.InitAIFraudDetector()

	router := setupRouter()

	// 登录获取 token
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

	// 先创建一个 approved 行为记录（用于 join 校验）
	behavior := &model.GreenBehavior{
		ID:           "behavior-1",
		WalletAddr:   "TestWallet123",
		BehaviorType: "waste_sorting",
		MediaURLs:    model.StringArray{"https://example.com/image1.jpg"},
		Status:       model.BehaviorStatusApproved,
		FraudScore:   0.1,
		SubmitTime:   config.Now(),
	}
	config.DB.Create(behavior)

	// 创建挑战（开始/结束时间用 RFC3339）
	now := config.Now()
	// 让挑战在当前时间窗口内，便于 join/claim
	start := now.Add(-1 * time.Minute).Format(time.RFC3339)
	end := now.Add(10 * time.Minute).Format(time.RFC3339)
	createReq := map[string]interface{}{
		"title":          "test",
		"description":    "desc",
		"challenge_type": "individual",
		"behavior_type":  "waste_sorting",
		"reward_amount":  1000,
		"target_count":   1,
		"start_time":     start,
		"end_time":       end,
	}
	createJson, _ := json.Marshal(createReq)
	createHttpReq, _ := http.NewRequest("POST", "/api/v1/challenges", bytes.NewBuffer(createJson))
	createHttpReq.Header.Set("Content-Type", "application/json")
	createHttpReq.Header.Set("Authorization", "Bearer "+token)
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createHttpReq)
	assert.Equal(t, http.StatusOK, createW.Code)

	var createResp map[string]interface{}
	json.Unmarshal(createW.Body.Bytes(), &createResp)
	ch := createResp["data"].(map[string]interface{})
	chID := ch["id"].(string)

	// 激活挑战
	actReq, _ := http.NewRequest("POST", "/api/v1/challenges/"+chID+"/activate", nil)
	actReq.Header.Set("Authorization", "Bearer "+token)
	actW := httptest.NewRecorder()
	router.ServeHTTP(actW, actReq)
	assert.Equal(t, http.StatusOK, actW.Code)

	// 参与挑战
	joinReqBody := map[string]string{"challenge_id": chID, "behavior_id": "behavior-1"}
	joinJson, _ := json.Marshal(joinReqBody)
	joinReq, _ := http.NewRequest("POST", "/api/v1/challenges/"+chID+"/join", bytes.NewBuffer(joinJson))
	joinReq.Header.Set("Content-Type", "application/json")
	joinReq.Header.Set("Authorization", "Bearer "+token)
	joinW := httptest.NewRecorder()
	router.ServeHTTP(joinW, joinReq)
	assert.Equal(t, http.StatusOK, joinW.Code)

	// 领取奖励前需要挑战 completed（target_count=1，join 后应自动完成）
	claimReq, _ := http.NewRequest("POST", "/api/v1/challenges/"+chID+"/claim", nil)
	claimReq.Header.Set("Authorization", "Bearer "+token)
	claimW := httptest.NewRecorder()
	router.ServeHTTP(claimW, claimReq)
	assert.Equal(t, http.StatusOK, claimW.Code)
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
