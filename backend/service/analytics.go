package service

import (
	"fmt"
	"sol-green/config"
	"sol-green/model"
	"time"
)

// AnalyticsData 数据分析结果
type AnalyticsData struct {
	DateRange      string                 `json:"date_range"`
	TotalUsers     int64                  `json:"total_users"`
	TotalBehaviors int64                  `json:"total_behaviors"`
	TotalRewards   uint64                 `json:"total_rewards"`
	ActiveUsers    int64                  `json:"active_users"`
	BehaviorStats  map[string]int64      `json:"behavior_stats"`
	DailyStats     []DailyStat            `json:"daily_stats"`
	TopUsers       []UserStat             `json:"top_users"`
	ChallengeStats ChallengeAnalytics      `json:"challenge_stats"`
	MarketingStats MarketingAnalytics      `json:"marketing_stats"`
}

// DailyStat 每日统计
type DailyStat struct {
	Date           string `json:"date"`
	Behaviors      int64  `json:"behaviors"`
	Users          int64  `json:"users"`
	Rewards        uint64 `json:"rewards"`
	Challenges     int64  `json:"challenges"`
	MarketingActs  int64  `json:"marketing_acts"`
}

// UserStat 用户统计
type UserStat struct {
	WalletAddr  string `json:"wallet_addr"`
	Behaviors   int64  `json:"behaviors"`
	Rewards     uint64 `json:"rewards"`
	JoinedAt    string `json:"joined_at"`
}

// ChallengeAnalytics 挑战统计
type ChallengeAnalytics struct {
	TotalChallenges    int64   `json:"total_challenges"`
	ActiveChallenges   int64   `json:"active_challenges"`
	TotalParticipants  int64   `json:"total_participants"`
	CompletionRate     float64 `json:"completion_rate"`
}

// MarketingAnalytics 营销活动统计
type MarketingAnalytics struct {
	TotalActivities    int64   `json:"total_activities"`
	ActiveActivities   int64   `json:"active_activities"`
	TotalParticipants  int64   `json:"total_participants"`
	SignInCount        int64   `json:"sign_in_count"`
	InviteCount        int64   `json:"invite_count"`
}

// GetAnalytics 获取数据分析
func GetAnalytics(startDate, endDate time.Time) (*AnalyticsData, error) {
	data := &AnalyticsData{
		DateRange:     fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		BehaviorStats: make(map[string]int64),
	}

	// 总用户数
	config.DB.Model(&model.User{}).Count(&data.TotalUsers)

	// 总行为数
	config.DB.Model(&model.GreenBehavior{}).
		Where("submit_time BETWEEN ? AND ?", startDate, endDate).
		Count(&data.TotalBehaviors)

	// 总奖励
	var totalRewards struct {
		Total uint64
	}
	config.DB.Model(&model.RewardRecord{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&totalRewards)
	data.TotalRewards = totalRewards.Total

	// 活跃用户（有行为的用户）
	config.DB.Model(&model.GreenBehavior{}).
		Where("submit_time BETWEEN ? AND ?", startDate, endDate).
		Distinct("wallet_addr").
		Count(&data.ActiveUsers)

	// 行为类型统计
	behaviorTypes := []string{"waste_sorting", "tree_planting", "low_carbon_travel"}
	for _, bt := range behaviorTypes {
		var count int64
		config.DB.Model(&model.GreenBehavior{}).
			Where("behavior_type = ? AND submit_time BETWEEN ? AND ?", bt, startDate, endDate).
			Count(&count)
		data.BehaviorStats[bt] = count
	}

	// 每日统计
	data.DailyStats = getDailyStats(startDate, endDate)

	// 排行榜用户
	data.TopUsers = getTopUsers(startDate, endDate, 10)

	// 挑战统计
	data.ChallengeStats = getChallengeStats(startDate, endDate)

	// 营销活动统计
	data.MarketingStats = getMarketingStats(startDate, endDate)

	return data, nil
}

// getDailyStats 获取每日统计
func getDailyStats(startDate, endDate time.Time) []DailyStat {
	var stats []DailyStat
	
	current := startDate
	for current.Before(endDate) || current.Equal(endDate) {
		nextDay := current.AddDate(0, 0, 1)
		
		var behaviors, users int64
		var rewards uint64
		
		config.DB.Model(&model.GreenBehavior{}).
			Where("submit_time >= ? AND submit_time < ?", current, nextDay).
			Count(&behaviors)
		
		config.DB.Model(&model.GreenBehavior{}).
			Where("submit_time >= ? AND submit_time < ?", current, nextDay).
			Distinct("wallet_addr").
			Count(&users)
		
		var totalRewards struct {
			Total uint64
		}
		config.DB.Model(&model.RewardRecord{}).
			Select("COALESCE(SUM(amount), 0) as total").
			Where("created_at >= ? AND created_at < ?", current, nextDay).
			Scan(&totalRewards)
		rewards = totalRewards.Total
		
		var challenges, marketingActs int64
		config.DB.Model(&model.Challenge{}).
			Where("created_at >= ? AND created_at < ?", current, nextDay).
			Count(&challenges)
		
		config.DB.Model(&model.MarketingActivity{}).
			Where("created_at >= ? AND created_at < ?", current, nextDay).
			Count(&marketingActs)
		
		stats = append(stats, DailyStat{
			Date:          current.Format("2006-01-02"),
			Behaviors:     behaviors,
			Users:         users,
			Rewards:       rewards,
			Challenges:    challenges,
			MarketingActs: marketingActs,
		})
		
		current = nextDay
	}
	
	return stats
}

// getTopUsers 获取排行榜用户
func getTopUsers(startDate, endDate time.Time, limit int) []UserStat {
	var users []UserStat
	
	config.DB.Model(&model.GreenBehavior{}).
		Select("wallet_addr, COUNT(*) as behaviors").
		Where("submit_time BETWEEN ? AND ?", startDate, endDate).
		Group("wallet_addr").
		Order("behaviors DESC").
		Limit(limit).
		Scan(&users)
	
	// 补充奖励信息
	for i := range users {
		var totalRewards struct {
			Total uint64
		}
		config.DB.Model(&model.RewardRecord{}).
			Select("COALESCE(SUM(amount), 0) as total").
			Where("wallet_addr = ? AND created_at BETWEEN ? AND ?", users[i].WalletAddr, startDate, endDate).
			Scan(&totalRewards)
		users[i].Rewards = totalRewards.Total
		
		var user model.User
		if err := config.DB.Where("wallet_addr = ?", users[i].WalletAddr).First(&user).Error; err == nil {
			users[i].JoinedAt = user.CreatedAt.Format("2006-01-02")
		}
	}
	
	return users
}

// getChallengeStats 获取挑战统计
func getChallengeStats(startDate, endDate time.Time) ChallengeAnalytics {
	var stats ChallengeAnalytics
	
	config.DB.Model(&model.Challenge{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&stats.TotalChallenges)
	
	config.DB.Model(&model.Challenge{}).
		Where("status = ? AND created_at BETWEEN ? AND ?", model.ChallengeStatusActive, startDate, endDate).
		Count(&stats.ActiveChallenges)
	
	config.DB.Model(&model.ChallengeParticipant{}).
		Where("joined_at BETWEEN ? AND ?", startDate, endDate).
		Count(&stats.TotalParticipants)
	
	// 计算完成率
	if stats.TotalChallenges > 0 {
		var completed int64
		config.DB.Model(&model.Challenge{}).
			Where("status = ? AND created_at BETWEEN ? AND ?", model.ChallengeStatusCompleted, startDate, endDate).
			Count(&completed)
		stats.CompletionRate = float64(completed) / float64(stats.TotalChallenges) * 100
	}
	
	return stats
}

// getMarketingStats 获取营销活动统计
func getMarketingStats(startDate, endDate time.Time) MarketingAnalytics {
	var stats MarketingAnalytics
	
	config.DB.Model(&model.MarketingActivity{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&stats.TotalActivities)
	
	config.DB.Model(&model.MarketingActivity{}).
		Where("status = ? AND created_at BETWEEN ? AND ?", model.MarketingStatusActive, startDate, endDate).
		Count(&stats.ActiveActivities)
	
	config.DB.Model(&model.MarketingParticipant{}).
		Where("joined_at BETWEEN ? AND ?", startDate, endDate).
		Count(&stats.TotalParticipants)
	
	config.DB.Model(&model.SignInRecord{}).
		Where("sign_in_date BETWEEN ? AND ?", startDate, endDate).
		Count(&stats.SignInCount)
	
	config.DB.Model(&model.InviteRecord{}).
		Where("invited_at BETWEEN ? AND ?", startDate, endDate).
		Count(&stats.InviteCount)
	
	return stats
}
