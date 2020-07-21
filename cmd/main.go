// Package main - домашнее задание.
package main

import (
	"fmt"
	"os"
	"runtime/trace"

	"netology/21_07_2020_Review/gohomework7_1/pkg/card"
)

func main() {

	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	const maxAmount = 1_000_000
	const numberOfTransactions = 1_000_000
	const parts = 100

	MCCList := card.MCC{
		"5010": "Финансы",
		"6020": "Супермаркеты",
		"7030": "Наличные",
		"8040": "Госуслуги",
		"9050": "Мобильная связь",
	}

	users := card.Users{
		0: "Ivan Ivanov",
		1: "Petr Petrov",
		2: "Dart Veider",
		3: "Luke I'm Your Father",
		4: "Vla Pu",
	}

	var transactions []card.Transaction

	err = card.GenerateTransactions(&transactions, maxAmount, numberOfTransactions, MCCList, users, parts)
	if err != nil {
		fmt.Println(card.ErrTransactionFulfill)
		os.Exit(2)
	}

	ts, err := card.SumByCategories(&transactions, 0)
	if err != nil {
		fmt.Println(card.ErrSummary)
		os.Exit(2)
	}
	fmt.Println("for user:", users[0])
	fmt.Println("Transactions summary:", ts)

	ts, err = card.SumByCategoriesWithMutex(&transactions, 0, parts)
	if err != nil {
		fmt.Println(card.ErrSummary)
		os.Exit(2)
	}
	fmt.Println("for user:", users[0])
	fmt.Println("Transactions summary with mutex:", ts)

	ts, err = card.SumByCategoriesWithChannels(&transactions, 0, parts)
	if err != nil {
		fmt.Println(card.ErrSummary)
		os.Exit(2)
	}
	fmt.Println("for user:", users[0])
	fmt.Println("Transactions summary with mutex:", ts)

	ts, err = card.SumByCategoriesWithMutex2(&transactions, 0, parts)
	if err != nil {
		fmt.Println(card.ErrSummary)
		os.Exit(2)
	}
	fmt.Println("for user:", users[0])
	fmt.Println("Transactions summary with mutex2:", ts)
}
