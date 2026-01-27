package tests

import (
	"sol-green/model"
	"sol-green/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateBehaviorHash(t *testing.T) {
	behavior := &model.GreenBehavior{
		WalletAddr:   "TestWallet123",
		BehaviorType: "waste_sorting",
		SubmitTime:   time.Now(),
	}
	
	hash1 := service.GenerateBehaviorHash(behavior)
	hash2 := service.GenerateBehaviorHash(behavior)
	
	// 相同输入应产生相同哈希
	assert.Equal(t, hash1, hash2)
	assert.NotEmpty(t, hash1)
	assert.Len(t, hash1, 64) // SHA256 哈希长度为 64
}
