package controller

// i18n.go - 国际化翻译文件
// i18n.go - Internationalization translation file

var (
	// 错误消息 / Error messages
	ErrInvalidParams     = map[string]string{"zh": "参数错误", "en": "Invalid parameters"}
	ErrUnauthorized      = map[string]string{"zh": "未登录", "en": "Unauthorized"}
	ErrNotFound          = map[string]string{"zh": "记录不存在", "en": "Record not found"}
	ErrDuplicateSubmit  = map[string]string{"zh": "请勿重复提交", "en": "Please do not submit duplicate"}
	ErrRewardFailed      = map[string]string{"zh": "奖励发放失败", "en": "Reward distribution failed"}
	ErrProofFailed       = map[string]string{"zh": "链上存证失败", "en": "On-chain proof storage failed"}
	ErrPendingReview     = map[string]string{"zh": "提交成功，进入人工审核", "en": "Submission successful, pending manual review"}
	ErrApproved          = map[string]string{"zh": "认证通过，奖励已发放", "en": "Verification passed, reward distributed"}

	// 成功消息 / Success messages
	MsgSubmitSuccess     = map[string]string{"zh": "提交成功", "en": "Submission successful"}
	MsgApproved          = map[string]string{"zh": "认证通过", "en": "Verification passed"}
	MsgRewardClaimed     = map[string]string{"zh": "奖励领取成功", "en": "Reward claimed successfully"}
	MsgPartnerVerified   = map[string]string{"zh": "第三方认证通过，奖励已发放", "en": "Third-party verification passed, reward distributed"}
)

// GetMessage 根据语言获取消息
// GetMessage gets message by language
func GetMessage(key map[string]string, lang string) string {
	if lang == "" {
		lang = "zh" // 默认中文 / Default Chinese
	}
	if msg, ok := key[lang]; ok {
		return msg
	}
	return key["zh"] // 回退到中文 / Fallback to Chinese
}
