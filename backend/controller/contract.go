package controller

import (
	"net/http"
	"sol-green/config"
	"sol-green/service"

	"github.com/gin-gonic/gin"
)

// GetContractInfo 获取合约信息
func GetContractInfo(c *gin.Context) {
	programID := c.DefaultQuery("program_id", config.SolanaProofContract())
	if programID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "程序 ID 未配置"})
		return
	}

	manager := service.NewContractUpgradeManager()
	version, err := manager.GetCurrentProgramVersion(programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取版本失败", "detail": err.Error()})
		return
	}

	available, latestVersion, err := manager.CheckUpgradeAvailable(programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查升级失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"program_id":      programID,
		"current_version": version,
		"latest_version":  latestVersion,
		"upgrade_available": available,
	})
}

// CheckUpgrade 检查升级
func CheckUpgrade(c *gin.Context) {
	walletAddr := c.GetString("wallet_addr")
	if walletAddr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// TODO: 检查管理员权限

	programID := c.DefaultQuery("program_id", config.SolanaProofContract())
	manager := service.NewContractUpgradeManager()

	available, latestVersion, err := manager.CheckUpgradeAvailable(programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查升级失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"upgrade_available": available,
		"current_version":   "",
		"latest_version":    latestVersion,
	})
}

// PrepareUpgrade 准备升级
func PrepareUpgrade(c *gin.Context) {
	walletAddr := c.GetString("wallet_addr")
	if walletAddr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// TODO: 检查管理员权限

	var req struct {
		ProgramID    string `json:"program_id" binding:"required"`
		NewProgramID string `json:"new_program_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	manager := service.NewContractUpgradeManager()

	// 验证新程序
	if err := manager.ValidateUpgrade(req.NewProgramID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "升级验证失败", "detail": err.Error()})
		return
	}

	// 准备升级
	if err := manager.PrepareUpgrade(req.ProgramID, req.NewProgramID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "准备升级失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "升级准备完成",
	})
}

// ExecuteUpgrade 执行升级
func ExecuteUpgrade(c *gin.Context) {
	walletAddr := c.GetString("wallet_addr")
	if walletAddr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// TODO: 检查管理员权限

	var req struct {
		ProgramID    string `json:"program_id" binding:"required"`
		NewProgramID string `json:"new_program_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	manager := service.NewContractUpgradeManager()

	// 执行升级
	txHash, err := manager.ExecuteUpgrade(req.ProgramID, req.NewProgramID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "升级失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":     "升级成功",
		"tx_hash": txHash,
	})
}

// GetUpgradeHistory 获取升级历史
func GetUpgradeHistory(c *gin.Context) {
	programID := c.DefaultQuery("program_id", config.SolanaProofContract())
	manager := service.NewContractUpgradeManager()

	history, err := manager.GetUpgradeHistory(programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取升级历史失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": history,
	})
}
