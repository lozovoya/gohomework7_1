package card

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var ErrorTransactionFulfill = errors.New("Slice of transactions is empty after generating func")
var ErrorSummary = errors.New("Wrong suumaring")

type Transaction struct {
	OwnerId int
	Amount  int64
	Mcc     string
	//Moment time.Time
}

type Mcc map[string]string
type User map[int]string

func GenerateTransactions(transactions *[]Transaction, max int, amount int64, mccList Mcc, userList User) error {

	var mccTemp = make(map[int]string)
	i := 0
	for m, _ := range mccList {
		mccTemp[i] = m
		i++
	}
	fmt.Println(mccTemp)

	for i := int64(0); i < amount; i++ {
		rand.Seed(int64(time.Now().Nanosecond()))
		t := Transaction{OwnerId: rand.Intn(len(userList)), Amount: int64(rand.Intn(max)), Mcc: mccTemp[rand.Intn(len(mccTemp))]}
		//t := Transaction{OwnerId: rand.Intn(len(userList)), Amount: 1, Mcc: mccTemp[rand.Intn(len(mccTemp))]}
		//fmt.Println(t)
		*transactions = append(*transactions, t)
		//transactions[index].Amount = 1
	}
	if *transactions == nil {
		return ErrorTransactionFulfill
	}
	return nil
}

func SumByCategories(transactions *[]Transaction, owner int) (catSum map[string]int64, error error) {

	catSum = make(map[string]int64, 10)
	error = nil
	for _, t := range *transactions {
		if t.OwnerId == owner {
			catSum[t.Mcc] = catSum[t.Mcc] + t.Amount
		}
	}
	if catSum == nil {
		return catSum, ErrorSummary
	}
	return catSum, nil
}
