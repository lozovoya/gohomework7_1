package main

import (
	"fmt"
	card2 "github.com/lozovoya/gohomework7_1/pkg/card"
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
	const numberOfTransactions = 1_000
	const parts = 100

	mccList := card2.Mcc{
		"5010": "Финансы",
		"6020": "Супермаркеты",
		"7030": "Наличные",
		"8040": "Госуслуги",
		"9050": "Мобильная связь",
	}

	userList := card2.User{
		0: "Ivan Ivanov",
		1: "Petr Petrov",
		2: "Dart Vaider",
		3: "Luk I'mYouFarther",
		4: "Vla Pu",
	}

	var transactions []card2.Transaction

	err = card2.GenerateTransactions(transactions, maxAmount, numberOfTransactions, mccList, userList)
	if err != nil {
		fmt.Println(card2.ErrTransactionFulfill)
		os.Exit(2)
	}

	ts := card2.SumByCategories(transactions, 0)
	fmt.Println("for user:", userList[0])
	fmt.Println("Transactions summary:", ts)

	ts = card2.SumByCategoriesWithMutex(transactions, 0, parts)
	fmt.Println("for user:", userList[0])
	fmt.Println("Transactions summary with mutex:", ts)

	ts = card2.SumByCategoriesWithChannels(transactions, 0, parts)
	fmt.Println("for user:", userList[0])
	fmt.Println("Transactions summary with mutex:", ts)

	ts = card2.SumByCategoriesWithMutex2(transactions, 0, parts)
	fmt.Println("for user:", userList[0])
	fmt.Println("Transactions summary with mutex2:", ts)
}
