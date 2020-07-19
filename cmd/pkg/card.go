package card

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

var ErrorTransactionFulfill = errors.New("Slice of transactions is empty after generating func")
var ErrorSummary = errors.New("Wrong summaring")

type Transaction struct {
	OwnerId int
	Amount  int64
	Mcc     string
	//Moment time.Time
}

type Mcc map[string]string
type User map[int]string

func GenerateTransactions(transactions *[]Transaction, max int, amount int64, mccList Mcc, userList User, parts int32) error {

	var mccTemp = make(map[int]string)
	i := 0
	for m, _ := range mccList {
		mccTemp[i] = m
		i++
	}

	//wg := sync.WaitGroup{}
	//wg.Add(int(parts))
	//mu := sync.Mutex{}
	//
	//partSize := (int32(amount)/parts)
	//
	//for i := 0; i < parts; i++ {
	//	rand.Seed(int64(time.Now().Nanosecond()))
	//}
	for i := int64(0); i < amount; i++ {
		rand.Seed(int64(time.Now().Nanosecond()))
		t := Transaction{OwnerId: rand.Intn(len(userList)), Amount: int64(rand.Intn(max)), Mcc: mccTemp[rand.Intn(len(mccTemp))]}
		//t := Transaction{OwnerId: 0, Amount: 1, Mcc: mccTemp[rand.Intn(len(mccTemp))]}
		//fmt.Println(t)
		//go func() {
		//	mu.Lock()
		*transactions = append(*transactions, t)
		//	mu.Unlock()
		//	wg.Done()
		//}()

		//transactions[index].Amount = 1
	}
	//wg.Wait()
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
		return catSum, ErrorSummary
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
		return catSum, ErrorSummary
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
					catSum[t.Mcc] = catSum[t.Mcc] + t.Amount
					mu.Unlock()
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if catSum == nil {
		return catSum, ErrorSummary
	}
	return catSum, error
}
