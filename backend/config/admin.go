package config

import (
	"strings"
)

// IsAdminWallet checks whether a wallet address is configured as an admin.
// Configure with env var: ADMIN_WALLETS="addr1,addr2,addr3"
func IsAdminWallet(walletAddr string) bool {
	walletAddr = strings.TrimSpace(walletAddr)
	if walletAddr == "" {
		return false
	}

	raw := GetEnv("ADMIN_WALLETS", "")
	if strings.TrimSpace(raw) == "" {
		return false
	}

	for _, part := range strings.Split(raw, ",") {
		if strings.TrimSpace(part) == walletAddr {
			return true
		}
	}
	return false
}

