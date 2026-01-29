package controller

import (
	"net/http"
	"os"
	"sol-green/config"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminMe 用于排查“无权限”问题：返回当前登录钱包与是否管理员
// GET /api/v1/admin/me
func AdminMe(c *gin.Context) {
	walletAddr := c.GetString("wallet_addr")
	if walletAddr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	raw := os.Getenv("ADMIN_WALLETS")
	parts := []string{}
	for _, p := range strings.Split(raw, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			parts = append(parts, p)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"wallet_addr":        walletAddr,
		"is_admin":           config.IsAdminWallet(walletAddr),
		"admin_wallets_raw":  raw,
		"admin_wallets_list": parts,
	})
}

