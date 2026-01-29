package controller

import (
	"fmt"
	"strings"
	"time"
)

// parseFlexibleTime parses time values commonly produced by browsers/forms.
// Supported examples:
// - 2026-01-25T02:07           (HTML datetime-local)
// - 2026-01-25T02:07:00
// - 2026-01-25T02:07:00Z
// - 2026-01-25T02:07:00+08:00 (RFC3339)
func parseFlexibleTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, fmt.Errorf("时间不能为空")
	}

	// 对于不带时区的输入（例如 datetime-local），优先按本地时区解析，避免被当成 UTC 导致“活动未开始/已结束”的错觉
	noZoneLayouts := []string{
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	for _, layout := range noZoneLayouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, nil
		}
	}

	// Try the layouts from most specific/likely to generic.
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("时间格式不正确：%q（建议使用 RFC3339，如 2026-01-25T02:07:00+08:00）", s)
}

