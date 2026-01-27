package model

import "time"

// ChallengeStatus 挑战状态
const (
	ChallengeStatusDraft      = "draft"       // 草稿
	ChallengeStatusActive     = "active"      // 进行中
	ChallengeStatusCompleted  = "completed"   // 已完成
	ChallengeStatusCancelled  = "cancelled"   // 已取消
)

// ChallengeType 挑战类型
const (
	ChallengeTypeIndividual = "individual" // 个人挑战
	ChallengeTypeTeam       = "team"        // 团队挑战
	ChallengeTypeCommunity  = "community"  // 社区挑战
)

// Challenge 挑战活动模型
type Challenge struct {
	ID              string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Title           string    `gorm:"type:varchar(200);not null" json:"title"`
	Description     string    `gorm:"type:text" json:"description"`
	ChallengeType   string    `gorm:"type:varchar(20);index" json:"challenge_type"`
	BehaviorType    string    `gorm:"type:varchar(20);index" json:"behavior_type"` // 关联的行为类型
	RewardAmount    uint64    `gorm:"type:bigint" json:"reward_amount"`            // 奖励数量
	TargetCount     int       `gorm:"type:int" json:"target_count"`                 // 目标参与人数
	CurrentCount    int       `gorm:"type:int;default:0" json:"current_count"`     // 当前参与人数
	StartTime       time.Time `gorm:"type:timestamp" json:"start_time"`
	EndTime         time.Time `gorm:"type:timestamp" json:"end_time"`
	Status          string    `gorm:"type:varchar(20);index" json:"status"`
	CreatorWallet   string    `gorm:"type:varchar(44);index" json:"creator_wallet"` // 创建者钱包地址
	ImageURL        string    `gorm:"type:varchar(500)" json:"image_url,omitempty"`
	Rules           string    `gorm:"type:text" json:"rules,omitempty"` // 挑战规则
	CreatedAt       time.Time `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt       time.Time `gorm:"type:timestamp" json:"updated_at"`
}

// ChallengeParticipant 挑战参与者
type ChallengeParticipant struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	ChallengeID string    `gorm:"type:varchar(36);index" json:"challenge_id"`
	WalletAddr  string    `gorm:"type:varchar(44);index" json:"wallet_addr"`
	BehaviorID  string    `gorm:"type:varchar(36);index" json:"behavior_id"` // 关联的行为记录
	JoinedAt    time.Time `gorm:"type:timestamp" json:"joined_at"`
	RewardClaimed bool    `gorm:"type:boolean;default:false" json:"reward_claimed"`
	TxHash      string    `gorm:"type:varchar(90)" json:"tx_hash,omitempty"`
}

// ChallengeReq 创建挑战请求
type ChallengeReq struct {
	Title         string    `json:"title" binding:"required"`
	Description   string    `json:"description" binding:"required"`
	ChallengeType string    `json:"challenge_type" binding:"required,oneof=individual team community"`
	BehaviorType  string    `json:"behavior_type" binding:"required,oneof=waste_sorting tree_planting low_carbon_travel"`
	RewardAmount  uint64    `json:"reward_amount" binding:"required,min=1"`
	TargetCount   int       `json:"target_count" binding:"required,min=1"`
	StartTime     time.Time `json:"start_time" binding:"required"`
	EndTime       time.Time `json:"end_time" binding:"required"`
	ImageURL      string    `json:"image_url,omitempty"`
	Rules         string    `json:"rules,omitempty"`
}

// JoinChallengeReq 参与挑战请求
type JoinChallengeReq struct {
	ChallengeID string `json:"challenge_id" binding:"required"`
	BehaviorID  string `json:"behavior_id" binding:"required"` // 关联的环保行为ID
}
