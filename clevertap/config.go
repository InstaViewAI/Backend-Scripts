package clevertap

import (
	"encoding/json"
	"fmt"
	"os"
)

// AccountConfig stores one CleverTap account's API configuration.
type AccountConfig struct {
	BaseURL   string `json:"base_url"`
	AccountID string `json:"account_id"`
	Passcode  string `json:"passcode"`
}

// AccountsConfig groups multiple named CleverTap accounts.
type AccountsConfig struct {
	Accounts map[string]AccountConfig `json:"accounts"`
}

func LoadAccountsConfig(path string) (*AccountsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read-config-failed: %w", err)
	}

	var cfg AccountsConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("decode-config-failed: %w", err)
	}

	if len(cfg.Accounts) == 0 {
		return nil, fmt.Errorf("no accounts found in config")
	}

	return &cfg, nil
}

func (c *AccountsConfig) Client(name string) (*Client, error) {
	account, ok := c.Accounts[name]
	if !ok {
		return nil, fmt.Errorf("account %q not found in config", name)
	}

	if account.BaseURL == "" || account.AccountID == "" || account.Passcode == "" {
		return nil, fmt.Errorf("account %q has missing base_url/account_id/passcode", name)
	}

	return NewClient(account.BaseURL, account.AccountID, account.Passcode), nil
}
