package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
)

type Account struct {
	Name    string
	Balance int
}

func loadAccounts() ([]Account, error) {
	balanceDir, err := GetBalanceDir()
	if err != nil {
		return nil, err
	}

	accountsFile := filepath.Join(balanceDir, "accounts.csv")
	file, err := os.Open(accountsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Account{}, nil
		}
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var accounts []Account
	for _, record := range records {
		if len(record) < 2 {
			continue
		}
		balance, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			balance = 0
		}
		accounts = append(accounts, Account{
			Name:    record[0],
			Balance: int(balance),
		})
	}
	return accounts, nil
}

func saveAccounts(accounts []Account) error {
	balanceDir, err := GetBalanceDir()
	if err != nil {
		return err
	}

	accountsFile := filepath.Join(balanceDir, "accounts.csv")
	file, err := os.Create(accountsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, account := range accounts {
		err := writer.Write([]string{account.Name, strconv.Itoa(account.Balance)})
		if err != nil {
			return err
		}
	}

	return nil
}
