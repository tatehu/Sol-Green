package controller

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadImage 上传图片（返回可访问的 URL）
// POST /api/v1/upload/image (multipart/form-data, field: file)
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少文件", "detail": err.Error()})
		return
	}

	// 仅允许常见图片后缀（避免用户误传可执行文件）
	ext := strings.ToLower(filepath.Ext(file.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
		// ok
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型", "detail": "仅支持 jpg/jpeg/png/webp/gif"})
		return
	}

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败", "detail": err.Error()})
		return
	}

	filename := uuid.New().String() + "_" + time.Now().Format("20060102150405") + ext
	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败", "detail": err.Error()})
		return
	}

	// 通过静态路由 /uploads 提供访问
	c.JSON(http.StatusOK, gin.H{
		"url": "/uploads/" + filename,
	})
}

