package card

// Банковская карта.

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

var ErrTransactionFulfill = errors.New("Slice of transactions is empty after generating func")
var ErrSummary = errors.New("Wrong summaring")

type Transaction struct {
	OwnerId int
	Amount  int64
	MCC     string
}

type MCC map[string]string
type Users map[int]string

func GenerateTransactions(transactions *[]Transaction, max int, amount int64, mccList MCC, userList Users, parts int32) error {

	var mccTemp = make(map[int]string)
	i := 0
	for m, _ := range mccList {
		mccTemp[i] = m
		i++
	}

	for i := int64(0); i < amount; i++ {
		rand.Seed(int64(time.Now().Nanosecond()))
		t := Transaction{OwnerId: rand.Intn(len(userList)), Amount: int64(rand.Intn(max)), MCC: mccTemp[rand.Intn(len(mccTemp))]}

		*transactions = append(*transactions, t)

	}
	if *transactions == nil {
		return ErrTransactionFulfill
	}
	return nil
}

func SumByCategories(transactions *[]Transaction, owner int) (catSum map[string]int64, error error) {

	catSum = make(map[string]int64, 10)
	error = nil
	for _, t := range *transactions {
		if t.OwnerId == owner {
			catSum[t.MCC] = catSum[t.MCC] + t.Amount
		}
	}
	if catSum == nil {
		return catSum, ErrSummary
	}
	return catSum, nil
}

func SumByCategoriesWithMutex(transactions *[]Transaction, owner int, parts int32) (catSum map[string]int64, error error) {

	wg := sync.WaitGroup{}
	wg.Add(int(parts))
	mu := sync.Mutex{}

	catSum = make(map[string]int64, 10)
	input := *transactions
	partSize := int32(len(*transactions)) / parts

	for i := int32(0); i < parts; i++ {
		part := input[i*partSize : (i+1)*partSize]
		go func() {
			partSum, _ := SumByCategories(&part, owner)
			mu.Lock()
			for key, value := range partSum {
				catSum[key] += value
			}
			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	if catSum == nil {
		return catSum, ErrSummary
	}
	return catSum, error
}

func SumByCategoriesWithChannels(transactions *[]Transaction, owner int, parts int32) (catSum map[string]int64, error error) {

	error = nil
	catSum = make(map[string]int64, 5)
	input := *transactions
	partSize := int32(len(*transactions)) / parts
	result := make(chan map[string]int64)

	for i := int32(0); i < parts; i++ {
		part := input[i*partSize : (i+1)*partSize]
		sum2(&part, owner, result)
	}
	finished := int32(0)
	for sums := range result {
		finished++
		for key, value := range sums {
			catSum[key] += value
		}
		if finished == parts {
			close(result)
			break
		}
	}
	if catSum == nil {
		return catSum, ErrSummary
	}
	return catSum, error

}

func sum2(part *[]Transaction, owner int, result chan<- map[string]int64) {
	go func() {
		partSum, _ := SumByCategories(part, owner)
		result <- partSum
	}()
}

func SumByCategoriesWithMutex2(transactions *[]Transaction, owner int, parts int32) (catSum map[string]int64, error error) {

	wg := sync.WaitGroup{}
	wg.Add(int(parts))
	mu := sync.Mutex{}

	error = nil
	catSum = make(map[string]int64, 10)
	input := *transactions
	partSize := int32(len(*transactions)) / parts

	for i := int32(0); i < parts; i++ {
		part := input[i*partSize : (i+1)*partSize]
		go func() {
			//partSum, _ := SumByCategories(&part, owner)
			for _, t := range part {
				if t.OwnerId == owner {
					mu.Lock()
					catSum[t.MCC] = catSum[t.MCC] + t.Amount
					mu.Unlock()
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if catSum == nil {
		return catSum, ErrSummary
	}
	return catSum, error
}
