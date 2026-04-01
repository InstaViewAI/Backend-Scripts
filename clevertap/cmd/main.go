package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"clevertap"
)

type emailInput struct {
	Emails []string `json:"emails"`
}

func main() {
	ctx := context.Background()

	configPath, err := resolveAccountsConfigPath()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	accountName := resolveAccountName()

	cfg, err := clevertap.LoadAccountsConfig(configPath)
	if err != nil {
		fmt.Printf("failed to load config from %s: %v\n", configPath, err)
		os.Exit(1)
	}

	ct, err := cfg.Client(accountName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	inputPath, err := resolveInputPath(accountName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	input, err := loadEmails(inputPath)
	if err != nil {
		fmt.Printf("failed to read emails from %s: %v\n", inputPath, err)
		os.Exit(1)
	}

	seen := make(map[string]struct{}, len(input.Emails))

	for _, email := range input.Emails {
		email = strings.TrimSpace(email)
		if email == "" {
			continue
		}

		if _, exists := seen[email]; exists {
			continue
		}
		seen[email] = struct{}{}

		res, err := ct.UpdateUser(ctx, email, map[string]interface{}{
			"MSG-sms": true, // Subscribe the SMS communication channel on clevertap
		})
		if err != nil {
			fmt.Printf("failed to update profile for %s: %v\n", email, err)
			continue
		}

		fmt.Printf("[%s] updated profile for %s status:%s processed:%d\n", accountName, email, res.Status, res.Processed)
	}
}

func resolveAccountName() string {
	accountName := strings.TrimSpace(os.Getenv("CLEVERTAP_ACCOUNT"))
	if accountName == "" {
		return "instaview"
	}

	return accountName
}

func resolveAccountsConfigPath() (string, error) {
	candidates := []string{
		"accounts.json",
		filepath.Join("..", "accounts.json"),
		filepath.Join("clevertap", "accounts.json"),
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("could not find accounts.json")
}

func resolveInputPath(accountName string) (string, error) {
	candidates := []string{
		fmt.Sprintf("%s_user.json", accountName),
		filepath.Join("..", fmt.Sprintf("%s_user.json", accountName)),
		filepath.Join("clevertap", fmt.Sprintf("%s_user.json", accountName)),
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("could not find %s_user.json", accountName)
}

func loadEmails(path string) (emailInput, error) {
	var input emailInput

	data, err := os.ReadFile(path)
	if err != nil {
		return input, err
	}

	if err := json.Unmarshal(data, &input); err != nil {
		return input, err
	}

	return input, nil
}
