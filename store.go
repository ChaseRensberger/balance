package main

import (
	"os"
	"path/filepath"
)

func GetBalanceDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	balanceDir := filepath.Join(homeDir, ".balance")
	if err := os.MkdirAll(balanceDir, 0755); err != nil {
		return "", err
	}
	return balanceDir, nil
}
