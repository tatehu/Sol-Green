package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sol-green/config"
	"sol-green/model"
	"time"
)

// PartnerProvider 第三方认证机构提供商
type PartnerProvider struct {
	ID          string
	Name        string
	APIEndpoint string
	APIKey      string
	APISecret   string
	Enabled     bool
}

// VerifyWithPartner 通过第三方机构验证
func VerifyWithPartner(partnerID, verifyCode string, behavior *model.GreenBehavior) (bool, error) {
	// 获取合作伙伴配置
	partner := getPartnerConfig(partnerID)
	if partner == nil {
		return false, fmt.Errorf("未知的认证机构: %s", partnerID)
	}

	if !partner.Enabled {
		return false, fmt.Errorf("认证机构 %s 未启用", partnerID)
	}

	// 根据不同的合作伙伴调用不同的验证方法
	switch partnerID {
	case "alashan_see":
		return verifyAlashanSEE(partner, verifyCode, behavior)
	case "gov_environment":
		return verifyGovEnvironment(partner, verifyCode, behavior)
	case "unep":
		return verifyUNEP(partner, verifyCode, behavior)
	case "wwf":
		return verifyWWF(partner, verifyCode, behavior)
	case "greenpeace":
		return verifyGreenpeace(partner, verifyCode, behavior)
	default:
		return false, fmt.Errorf("不支持的认证机构: %s", partnerID)
	}
}

// getPartnerConfig 获取合作伙伴配置
func getPartnerConfig(partnerID string) *PartnerProvider {
	// 从环境变量或配置文件读取
	config.Log.Infof("获取合作伙伴配置: %s", partnerID)

	partners := map[string]*PartnerProvider{
		"alashan_see": {
			ID:          "alashan_see",
			Name:        "阿拉善 SEE 生态协会",
			APIEndpoint: config.GetEnv("PARTNER_ALASHAN_API", "https://api.alashansee.org/v1/verify"),
			APIKey:      config.GetEnv("PARTNER_ALASHAN_API_KEY", ""),
			APISecret:   config.GetEnv("PARTNER_ALASHAN_API_SECRET", ""),
			Enabled:     config.GetEnv("PARTNER_ALASHAN_ENABLED", "false") == "true",
		},
		"gov_environment": {
			ID:          "gov_environment",
			Name:        "政府环保部门",
			APIEndpoint: config.GetEnv("PARTNER_GOV_API", ""),
			APIKey:      config.GetEnv("PARTNER_GOV_API_KEY", ""),
			APISecret:   config.GetEnv("PARTNER_GOV_API_SECRET", ""),
			Enabled:     config.GetEnv("PARTNER_GOV_ENABLED", "false") == "true",
		},
		"unep": {
			ID:          "unep",
			Name:        "联合国环境规划署 (UNEP)",
			APIEndpoint: config.GetEnv("PARTNER_UNEP_API", "https://api.unep.org/v1/verify"),
			APIKey:      config.GetEnv("PARTNER_UNEP_API_KEY", ""),
			APISecret:   config.GetEnv("PARTNER_UNEP_API_SECRET", ""),
			Enabled:     config.GetEnv("PARTNER_UNEP_ENABLED", "false") == "true",
		},
		"wwf": {
			ID:          "wwf",
			Name:        "世界自然基金会 (WWF)",
			APIEndpoint: config.GetEnv("PARTNER_WWF_API", "https://api.wwf.org/v1/verify"),
			APIKey:      config.GetEnv("PARTNER_WWF_API_KEY", ""),
			APISecret:   config.GetEnv("PARTNER_WWF_API_SECRET", ""),
			Enabled:     config.GetEnv("PARTNER_WWF_ENABLED", "false") == "true",
		},
		"greenpeace": {
			ID:          "greenpeace",
			Name:        "绿色和平 (Greenpeace)",
			APIEndpoint: config.GetEnv("PARTNER_GREENPEACE_API", "https://api.greenpeace.org/v1/verify"),
			APIKey:      config.GetEnv("PARTNER_GREENPEACE_API_KEY", ""),
			APISecret:   config.GetEnv("PARTNER_GREENPEACE_API_SECRET", ""),
			Enabled:     config.GetEnv("PARTNER_GREENPEACE_ENABLED", "false") == "true",
		},
	}

	return partners[partnerID]
}

// verifyAlashanSEE 阿拉善 SEE 生态协会认证
// 阿拉善 SEE 是中国知名的环保 NGO，专注于荒漠化防治和生态保护
func verifyAlashanSEE(partner *PartnerProvider, verifyCode string, behavior *model.GreenBehavior) (bool, error) {
	config.Log.Infof("调用阿拉善 SEE 认证: behavior_id=%s, verify_code=%s", behavior.ID, verifyCode)

	if !partner.Enabled || partner.APIKey == "" {
		// 模拟验证（开发/测试环境）
		if len(verifyCode) >= 3 && verifyCode[:3] == "SEE" {
			return true, nil
		}
		return false, nil
	}

	// 实际 API 调用
	reqBody := map[string]interface{}{
		"verify_code":  verifyCode,
		"behavior_id":   behavior.ID,
		"wallet_addr":  behavior.WalletAddr,
		"behavior_type": behavior.BehaviorType,
		"timestamp":    time.Now().Unix(),
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", partner.APIEndpoint, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+partner.APIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("API 调用失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if resp.StatusCode == 200 {
		if verified, ok := result["verified"].(bool); ok {
			return verified, nil
		}
	}

	return false, fmt.Errorf("认证失败: %v", result)
}

// verifyGovEnvironment 政府环保部门认证
func verifyGovEnvironment(partner *PartnerProvider, verifyCode string, behavior *model.GreenBehavior) (bool, error) {
	config.Log.Infof("调用政府环保部门认证: behavior_id=%s, verify_code=%s", behavior.ID, verifyCode)

	if !partner.Enabled || partner.APIKey == "" {
		// 模拟验证
		if len(verifyCode) >= 3 && verifyCode[:3] == "GOV" {
			return true, nil
		}
		return false, nil
	}

	// 实际 API 调用（类似上面的实现）
	return false, fmt.Errorf("政府环保部门 API 未配置")
}

// verifyUNEP 联合国环境规划署认证
// UNEP 是联合国系统内负责环境事务的牵头机构
func verifyUNEP(partner *PartnerProvider, verifyCode string, behavior *model.GreenBehavior) (bool, error) {
	config.Log.Infof("调用 UNEP 认证: behavior_id=%s", behavior.ID)

	if !partner.Enabled || partner.APIKey == "" {
		// 模拟验证
		if len(verifyCode) >= 4 && verifyCode[:4] == "UNEP" {
			return true, nil
		}
		return false, nil
	}

	// 实际 API 调用
	return false, fmt.Errorf("UNEP API 未配置")
}

// verifyWWF 世界自然基金会认证
// WWF 是全球最大的独立性非政府环境保护组织之一
func verifyWWF(partner *PartnerProvider, verifyCode string, behavior *model.GreenBehavior) (bool, error) {
	config.Log.Infof("调用 WWF 认证: behavior_id=%s", behavior.ID)

	if !partner.Enabled || partner.APIKey == "" {
		// 模拟验证
		if len(verifyCode) >= 3 && verifyCode[:3] == "WWF" {
			return true, nil
		}
		return false, nil
	}

	// 实际 API 调用
	return false, fmt.Errorf("WWF API 未配置")
}

// verifyGreenpeace 绿色和平认证
// Greenpeace 是国际知名的环保组织
func verifyGreenpeace(partner *PartnerProvider, verifyCode string, behavior *model.GreenBehavior) (bool, error) {
	config.Log.Infof("调用 Greenpeace 认证: behavior_id=%s", behavior.ID)

	if !partner.Enabled || partner.APIKey == "" {
		// 模拟验证
		if len(verifyCode) >= 2 && verifyCode[:2] == "GP" {
			return true, nil
		}
		return false, nil
	}

	// 实际 API 调用
	return false, fmt.Errorf("Greenpeace API 未配置")
}
