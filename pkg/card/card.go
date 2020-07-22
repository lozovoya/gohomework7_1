package card

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

var ErrTransactionFulfill = errors.New("Slice of transactions is empty after generating func")

type Transaction struct {
	OwnerId int
	Amount  int64
	Mcc     string
}

type Mcc map[string]string
type User map[int]string

func GenerateTransactions(transactions *[]Transaction, max int, amount int64, mccList Mcc, userList User) error {

	var mccTemp = make(map[int]string)
	i := 0
	for m := range mccList {
		mccTemp[i] = m
		i++
	}

	for i := int64(0); i < amount; i++ {
		rand.Seed(int64(time.Now().Nanosecond()))
		t := Transaction{OwnerId: rand.Intn(len(userList)), Amount: int64(rand.Intn(max)), Mcc: mccTemp[rand.Intn(len(mccTemp))]}
		*transactions = append(*transactions, t)
	}
	if *transactions == nil {
		return ErrTransactionFulfill
	}
	return nil
}

func SumByCategories(transactions []Transaction, owner int) (catSum map[string]int64) {

	catSum = make(map[string]int64, 10)
	for _, t := range transactions {
		if t.OwnerId == owner {
			catSum[t.Mcc] = catSum[t.Mcc] + t.Amount
		}
	}

	return catSum
}

func SumByCategoriesWithMutex(transactions []Transaction, owner int, parts int32) (catSum map[string]int64) {

	wg := sync.WaitGroup{}
	wg.Add(int(parts))
	mu := sync.Mutex{}

	catSum = make(map[string]int64, 10)
	input := transactions
	partSize := int32(len(transactions)) / parts

	for i := int32(0); i < parts; i++ {
		part := input[i*partSize : (i+1)*partSize]
		go func() {
			partSum := SumByCategories(part, owner)
			mu.Lock()
			for key, value := range partSum {
				catSum[key] += value
			}
			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()

	return catSum
}

func SumByCategoriesWithChannels(transactions []Transaction, owner int, parts int32) (catSum map[string]int64) {

	catSum = make(map[string]int64, 5)
	input := transactions
	partSize := int32(len(transactions)) / parts
	result := make(chan map[string]int64)

	for i := int32(0); i < parts; i++ {
		part := input[i*partSize : (i+1)*partSize]
		sum2(part, owner, result)
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

	return catSum

}

func sum2(part []Transaction, owner int, result chan<- map[string]int64) {
	go func() {
		partSum := SumByCategories(part, owner)
		result <- partSum
	}()
}

func SumByCategoriesWithMutex2(transactions []Transaction, owner int, parts int32) (catSum map[string]int64) {

	wg := sync.WaitGroup{}
	wg.Add(int(parts))
	mu := sync.Mutex{}

	catSum = make(map[string]int64, 10)
	input := transactions
	partSize := int32(len(transactions)) / parts

	for i := int32(0); i < parts; i++ {
		part := input[i*partSize : (i+1)*partSize]
		go func() {
			//partSum, _ := SumByCategories(&part, owner)
			for _, t := range part {
				if t.OwnerId == owner {
					mu.Lock()
					catSum[t.Mcc] = catSum[t.Mcc] + t.Amount
					mu.Unlock()
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	return catSum
}
