package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Transaction struct {
	From      string
	To        string
	Amount    int
	Timestamp time.Time
}

func loadTransactions() ([]Transaction, error) {
	balanceDir, err := GetBalanceDir()
	if err != nil {
		return nil, err
	}

	transactionsFile := filepath.Join(balanceDir, "transactions.csv")
	file, err := os.Open(transactionsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Transaction{}, nil
		}
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var transactions []Transaction
	for _, record := range records {
		if len(record) < 3 {
			continue
		}
		amount, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			amount = 0
		}
		transactions = append(transactions, Transaction{
			From:   record[0],
			To:     record[1],
			Amount: int(amount),
		})
	}
	return transactions, nil
}

func saveTransactions(transactions []Transaction) error {
	balanceDir, err := GetBalanceDir()
	if err != nil {
		return err
	}

	transactionsFile := filepath.Join(balanceDir, "transactions.csv")
	file, err := os.Create(transactionsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, transaction := range transactions {
		err := writer.Write([]string{transaction.From, transaction.To, strconv.Itoa(transaction.Amount)})
		if err != nil {
			return err
		}
	}

	return nil
}
