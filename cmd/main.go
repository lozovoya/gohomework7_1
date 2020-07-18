package main

import (
	"fmt"
	card "github.com/lozovoya/gohomework7_1/cmd/pkg"
	"os"
	"runtime/trace"
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
	const numberOfTransactions = 100
	const parts = 2

	mccList := card.Mcc{
		"5010": "Финансы",
		"6020": "Супермаркеты",
		"7030": "Наличные",
		"8040": "Госуслуги",
		"9050": "Мобильная связь",
	}

	userList := card.User{
		0: "Ivan Ivanov",
		1: "Petr Petrov",
		2: "Dart Vaider",
		3: "Luk I'mYouFarther",
		4: "Vla Pu",
	}

	transactions := make([]card.Transaction, 0, numberOfTransactions)
	transactions = nil

	err = card.GenerateTransactions(&transactions, maxAmount, numberOfTransactions, mccList, userList, parts)
	if err != nil {
		fmt.Println(card.ErrorTransactionFulfill)
		os.Exit(2)
	}

	ts, err := card.SumByCategories(&transactions, 0)
	if err != nil {
		fmt.Println(card.ErrorSummary)
		os.Exit(2)
	}
	fmt.Println("for user:", userList[0])
	fmt.Println("Transactions summary:", ts)

	ts, err = card.SumByCategoriesWithMutex(&transactions, 0, parts)
	if err != nil {
		fmt.Println(card.ErrorSummary)
		os.Exit(2)
	}
	fmt.Println("for user:", userList[0])
	fmt.Println("Transactions summary:", ts)

	ts, err = card.SumByCategoriesWithChannels(&transactions, 0, parts)
	if err != nil {
		fmt.Println(card.ErrorSummary)
		os.Exit(2)
	}
	fmt.Println("for user:", userList[0])
	fmt.Println("Transactions summary:", ts)
}
