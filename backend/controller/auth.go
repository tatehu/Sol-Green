package controller

import (
	"net/http"
	"sol-green/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// WalletAuthReq 钱包认证请求 / Wallet authentication request
type WalletAuthReq struct {
	WalletAddr string `json:"wallet_addr" binding:"required"`
	Signature  string `json:"signature" binding:"required"` // Solana 钱包签名 / Solana wallet signature
	Message    string `json:"message" binding:"required"`   // 签名消息 / Signature message
}

// WalletAuth 钱包登录认证
// WalletAuth wallet login authentication
func WalletAuth(c *gin.Context) {
	var req WalletAuthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "参数错误", // "Invalid parameters"
			"error_en": "Invalid parameters",
		})
		return
	}

	// TODO: 验证 Solana 钱包签名
	// TODO: Verify Solana wallet signature
	// 实际应调用 Solana 验证签名接口
	// Should call Solana signature verification API
	// if !verifySolanaSignature(req.WalletAddr, req.Message, req.Signature) {
	//     c.JSON(http.StatusUnauthorized, gin.H{"error": "签名验证失败"}) // "Signature verification failed"
	//     return
	// }

	// 生成 JWT Token
	// Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"wallet_addr": req.WalletAddr,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.GetEnv("JWT_SECRET", "sol-green-secret-key")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Token 生成失败", // "Token generation failed"
			"error_en": "Token generation failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":       tokenString,
		"wallet_addr": req.WalletAddr,
	})
}
