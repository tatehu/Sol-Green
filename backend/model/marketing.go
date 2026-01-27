package model

import "time"

// MarketingActivityType 营销活动类型
const (
	MarketingTypeSignIn      = "sign_in"       // 签到活动
	MarketingTypeInvite      = "invite"        // 邀请活动
	MarketingTypeDailyTask    = "daily_task"    // 每日任务
	MarketingTypeLuckyDraw   = "lucky_draw"    // 抽奖活动
	MarketingTypeFlashSale    = "flash_sale"    // 限时抢购
	MarketingTypeFestival     = "festival"      // 节日活动
)

// MarketingActivityStatus 活动状态
const (
	MarketingStatusDraft     = "draft"      // 草稿
	MarketingStatusScheduled = "scheduled"  // 已排期
	MarketingStatusActive    = "active"     // 进行中
	MarketingStatusPaused    = "paused"     // 已暂停
	MarketingStatusEnded     = "ended"      // 已结束
)

// MarketingActivity 营销活动模型
type MarketingActivity struct {
	ID              string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Title           string    `gorm:"type:varchar(200);not null" json:"title"`
	Description     string    `gorm:"type:text" json:"description"`
	ActivityType    string    `gorm:"type:varchar(20);index" json:"activity_type"`
	Status          string    `gorm:"type:varchar(20);index" json:"status"`
	StartTime       time.Time `gorm:"type:timestamp" json:"start_time"`
	EndTime         time.Time `gorm:"type:timestamp" json:"end_time"`
	RewardAmount    uint64    `gorm:"type:bigint" json:"reward_amount"`
	MaxParticipants int       `gorm:"type:int;default:0" json:"max_participants"` // 0 表示无限制
	CurrentParticipants int   `gorm:"type:int;default:0" json:"current_participants"`
	Rules           string    `gorm:"type:text" json:"rules,omitempty"`
	ImageURL        string    `gorm:"type:varchar(500)" json:"image_url,omitempty"`
	Config          string    `gorm:"type:json" json:"config,omitempty"` // 活动配置（JSON）
	CreatedBy       string    `gorm:"type:varchar(44)" json:"created_by"` // 创建者钱包地址
	CreatedAt       time.Time `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt       time.Time `gorm:"type:timestamp" json:"updated_at"`
}

// MarketingParticipant 营销活动参与者
type MarketingParticipant struct {
	ID           string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	ActivityID   string    `gorm:"type:varchar(36);index" json:"activity_id"`
	WalletAddr   string    `gorm:"type:varchar(44);index" json:"wallet_addr"`
	JoinedAt     time.Time `gorm:"type:timestamp" json:"joined_at"`
	RewardClaimed bool    `gorm:"type:boolean;default:false" json:"reward_claimed"`
	TxHash       string    `gorm:"type:varchar(90)" json:"tx_hash,omitempty"`
	ExtraData    string    `gorm:"type:json" json:"extra_data,omitempty"` // 额外数据（如邀请人数、签到天数等）
}

// SignInRecord 签到记录
type SignInRecord struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	WalletAddr string    `gorm:"type:varchar(44);index" json:"wallet_addr"`
	ActivityID string    `gorm:"type:varchar(36);index" json:"activity_id"`
	SignInDate time.Time `gorm:"type:date;index" json:"sign_in_date"`
	ConsecutiveDays int  `gorm:"type:int;default:1" json:"consecutive_days"` // 连续签到天数
	RewardAmount uint64  `gorm:"type:bigint" json:"reward_amount"`
	TxHash     string    `gorm:"type:varchar(90)" json:"tx_hash,omitempty"`
}

// InviteRecord 邀请记录
type InviteRecord struct {
	ID           string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	InviterAddr  string    `gorm:"type:varchar(44);index" json:"inviter_addr"` // 邀请人
	InviteeAddr  string    `gorm:"type:varchar(44);index" json:"invitee_addr"` // 被邀请人
	ActivityID   string    `gorm:"type:varchar(36);index" json:"activity_id"`
	InvitedAt    time.Time `gorm:"type:timestamp" json:"invited_at"`
	RewardClaimed bool    `gorm:"type:boolean;default:false" json:"reward_claimed"`
	TxHash       string    `gorm:"type:varchar(90)" json:"tx_hash,omitempty"`
}
