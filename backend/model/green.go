package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// 行为状态常量 / Behavior status constants
const (
	BehaviorStatusPendingReview = "pending_review" // 待审核 / Pending review
	BehaviorStatusApproved      = "approved"       // 认证通过 / Approved
	BehaviorStatusRejected      = "rejected"       // 认证拒绝 / Rejected
)

// StringArray 字符串数组类型，用于 GORM JSON 字段
type StringArray []string

// Value 实现 driver.Valuer 接口
func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan 实现 sql.Scanner 接口
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}

// GreenBehaviorReq 提交请求 / Submission request
type GreenBehaviorReq struct {
	BehaviorType string   `json:"behavior_type" binding:"required,oneof=waste_sorting tree_planting low_carbon_travel"` // 垃圾分类/植树/低碳出行 / Waste sorting/Tree planting/Low-carbon travel
	MediaURLs    []string `json:"media_urls" binding:"required,min=1,max=5"`                                         // 图片/视频链接 / Image/video URLs
	Location     string   `json:"location"`                                                                          // 行为地点（可选）/ Location (optional)
}

// GreenBehavior 环保行为模型 / Environmental behavior model
type GreenBehavior struct {
	ID            string     `gorm:"primaryKey;type:varchar(36)" json:"id"`
	WalletAddr    string     `gorm:"type:varchar(44);index" json:"wallet_addr"` // Solana 钱包地址 / Solana wallet address
	BehaviorType  string     `gorm:"type:varchar(20);index" json:"behavior_type"`
	MediaURLs     StringArray `gorm:"type:json" json:"media_urls"`
	Location      string     `gorm:"type:varchar(100)" json:"location"`
	Status        string     `gorm:"type:varchar(20);index" json:"status"`
	FraudScore    float64    `gorm:"type:decimal(3,2)" json:"fraud_score"` // 欺诈分数 0-1 / Fraud score 0-1
	SubmitTime    time.Time  `gorm:"type:timestamp" json:"submit_time"`
	ApproveTime   *time.Time `gorm:"type:timestamp" json:"approve_time,omitempty"`
	RejectReason  string     `gorm:"type:text" json:"reject_reason,omitempty"`
	TxHash        string     `gorm:"type:varchar(90)" json:"tx_hash,omitempty"` // 奖励交易哈希 / Reward transaction hash
	ProofHash     string     `gorm:"type:varchar(64)" json:"proof_hash,omitempty"` // 链上存证哈希 / On-chain proof hash
	PartnerVerify string     `gorm:"type:varchar(50)" json:"partner_verify,omitempty"` // 第三方认证机构 / Third-party verification organization
}

// User 用户模型
type User struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	WalletAddr string    `gorm:"type:varchar(44);uniqueIndex" json:"wallet_addr"`
	CreatedAt  time.Time `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:timestamp" json:"updated_at"`
}

// RewardRecord 奖励记录
type RewardRecord struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	WalletAddr string    `gorm:"type:varchar(44);index" json:"wallet_addr"`
	BehaviorID string    `gorm:"type:varchar(36);index" json:"behavior_id"`
	Amount     uint64    `gorm:"type:bigint" json:"amount"`
	TxHash     string    `gorm:"type:varchar(90);uniqueIndex" json:"tx_hash"`
	CreatedAt  time.Time `gorm:"type:timestamp" json:"created_at"`
}
