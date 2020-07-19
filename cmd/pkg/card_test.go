package card

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSumByCategories(t *testing.T) {
	type args struct {
		transactions *[]Transaction
		owner        int
	}

	tr := []Transaction{
		{OwnerId: 0, Amount: 100, Mcc: "5010"},
		{OwnerId: 1, Amount: 100, Mcc: "5010"},
		{OwnerId: 0, Amount: 100, Mcc: "5020"},
		{OwnerId: 0, Amount: 100, Mcc: "5010"},
		{OwnerId: 1, Amount: 100, Mcc: "5010"},
	}

	tests := []struct {
		name       string
		args       args
		wantCatSum map[string]int64
		wantErr    error
	}{
		{
			name: "По соточке",
			args: args{
				transactions: &tr,
				owner:        0,
			},
			wantCatSum: map[string]int64{
				"5010": 200,
				"5020": 100,
			},
			wantErr: nil,
		},
	}

	transactions := makeTestTransactions()
	fmt.Println(SumByCategories(&transactions, 0))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCatSum, err := SumByCategories(tt.args.transactions, tt.args.owner)
			if err != tt.wantErr {
				t.Errorf("SumByCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCatSum, tt.wantCatSum) {
				t.Errorf("SumByCategories() gotCatSum = %v, want %v", gotCatSum, tt.wantCatSum)
			}
		})
	}
}

func makeTestTransactions() []Transaction {
	const users = 10_000
	const transactionsPerUser = 10_000
	const transactionAmount1 = 1_00
	const transactionAmount2 = 2_00
	const testcategory1 = "5010"
	const testcategory2 = "5020"

	transactions := make([]Transaction, users*transactionsPerUser)
	for index := range transactions {
		switch index % 100 {
		case 0:
			transactions[index] = Transaction{
				OwnerId: 0,
				Amount:  transactionAmount1,
				Mcc:     testcategory1,
			}
		case 20:
			transactions[index] = Transaction{
				OwnerId: 0,
				Amount:  transactionAmount2,
				Mcc:     testcategory2,
			}
		default:
			transactions[index] = Transaction{
				OwnerId: 1,
				Amount:  transactionAmount1,
				Mcc:     "5030",
			}
		}
	}
	return transactions
}

func BenchmarkCategorization(b *testing.B) {
	transactions := makeTestTransactions()
	want := map[string]int64{
		"5010": 100000000,
		"5020": 200000000,
	}
	b.ResetTimer() // сбрасываем таймер, т.к. сама генерация транзакций достаточно ресурсоёмка
	for i := 0; i < b.N; i++ {
		result, _ := SumByCategories(&transactions, 0)
		b.StopTimer() // останавливаем таймер, чтобы время сравнения не учитывалось
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer() // продолжаем работу таймера
	}
}

func BenchmarkCategorizationWithMutex(b *testing.B) {
	transactions := makeTestTransactions()
	want := map[string]int64{
		"5010": 100000000,
		"5020": 200000000,
	}
	b.ResetTimer() // сбрасываем таймер, т.к. сама генерация транзакций достаточно ресурсоёмка
	for i := 0; i < b.N; i++ {
		result, _ := SumByCategoriesWithMutex(&transactions, 0, 100)
		b.StopTimer() // останавливаем таймер, чтобы время сравнения не учитывалось
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer() // продолжаем работу таймера
	}
}

func BenchmarkCategorizationWithChannels(b *testing.B) {
	transactions := makeTestTransactions()
	want := map[string]int64{
		"5010": 100000000,
		"5020": 200000000,
	}
	b.ResetTimer() // сбрасываем таймер, т.к. сама генерация транзакций достаточно ресурсоёмка
	for i := 0; i < b.N; i++ {
		result, _ := SumByCategoriesWithChannels(&transactions, 0, 100)
		b.StopTimer() // останавливаем таймер, чтобы время сравнения не учитывалось
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer() // продолжаем работу таймера
	}
}

func BenchmarkCategorizationWithMutex2(b *testing.B) {
	transactions := makeTestTransactions()
	want := map[string]int64{
		"5010": 100000000,
		"5020": 200000000,
	}
	b.ResetTimer() // сбрасываем таймер, т.к. сама генерация транзакций достаточно ресурсоёмка
	for i := 0; i < b.N; i++ {
		result, _ := SumByCategoriesWithMutex2(&transactions, 0, 100)
		b.StopTimer() // останавливаем таймер, чтобы время сравнения не учитывалось
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer() // продолжаем работу таймера
	}
}
