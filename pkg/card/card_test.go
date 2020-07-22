package card

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSumByCategories(t *testing.T) {
	type args struct {
		transactions []Transaction
		owner        int
	}

	transactions := makeTestTransactions()
	tests := []struct {
		name       string
		args       args
		wantCatSum map[string]int64
	}{
		{
			name: "По соточке",
			args: args{
				transactions: transactions,
				owner:        0,
			},
			wantCatSum: map[string]int64{
				"5010": 10000000,
				"5020": 20000000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCatSum := SumByCategories(tt.args.transactions, tt.args.owner)
			if !reflect.DeepEqual(gotCatSum, tt.wantCatSum) {
				t.Errorf("SumByCategories() gotCatSum = %v, want %v", gotCatSum, tt.wantCatSum)
			}
		})
	}
}

func TestSumByCategoriesWithMutex(t *testing.T) {
	type args struct {
		transactions []Transaction
		owner        int
		parts        int32
	}

	transactions := makeTestTransactions()
	tests := []struct {
		name       string
		args       args
		wantCatSum map[string]int64
		wantErr    error
	}{
		{
			name: "По соточке",
			args: args{
				transactions: transactions,
				owner:        0,
				parts:        100,
			},
			wantCatSum: map[string]int64{
				"5010": 10000000,
				"5020": 20000000,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCatSum := SumByCategoriesWithMutex(tt.args.transactions, tt.args.owner, tt.args.parts)
			if !reflect.DeepEqual(gotCatSum, tt.wantCatSum) {
				t.Errorf("SumByCategories() gotCatSum = %v, want %v", gotCatSum, tt.wantCatSum)
			}
		})
	}
}

func TestSumByCategoriesWithChannels(t *testing.T) {
	type args struct {
		transactions []Transaction
		owner        int
		parts        int32
	}

	transactions := makeTestTransactions()
	tests := []struct {
		name       string
		args       args
		wantCatSum map[string]int64
		wantErr    error
	}{
		{
			name: "По соточке",
			args: args{
				transactions: transactions,
				owner:        0,
				parts:        100,
			},
			wantCatSum: map[string]int64{
				"5010": 10000000,
				"5020": 20000000,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCatSum := SumByCategoriesWithChannels(tt.args.transactions, tt.args.owner, tt.args.parts)
			if !reflect.DeepEqual(gotCatSum, tt.wantCatSum) {
				t.Errorf("SumByCategories() gotCatSum = %v, want %v", gotCatSum, tt.wantCatSum)
			}
		})
	}
}

func TestSumByCategoriesWithMutex2(t *testing.T) {
	type args struct {
		transactions []Transaction
		owner        int
		parts        int32
	}

	transactions := makeTestTransactions()
	tests := []struct {
		name       string
		args       args
		wantCatSum map[string]int64
		wantErr    error
	}{
		{
			name: "По соточке",
			args: args{
				transactions: transactions,
				owner:        0,
				parts:        100,
			},
			wantCatSum: map[string]int64{
				"5010": 10000000,
				"5020": 20000000,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCatSum := SumByCategoriesWithMutex2(tt.args.transactions, tt.args.owner, tt.args.parts)
			if !reflect.DeepEqual(gotCatSum, tt.wantCatSum) {
				t.Errorf("SumByCategories() gotCatSum = %v, want %v", gotCatSum, tt.wantCatSum)
			}
		})
	}
}

func makeTestTransactions() []Transaction {
	const users = 10_000
	const transactionsPerUser = 1_000
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
	fmt.Println("transaction generating finished")
	return transactions
}

func BenchmarkCategorization(b *testing.B) {
	transactions := makeTestTransactions()
	want := map[string]int64{
		"5010": 10000000,
		"5020": 20000000,
	}
	b.ResetTimer() // сбрасываем таймер, т.к. сама генерация транзакций достаточно ресурсоёмка
	for i := 0; i < b.N; i++ {
		result := SumByCategories(transactions, 0)
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
		"5010": 10000000,
		"5020": 20000000,
	}
	b.ResetTimer() // сбрасываем таймер, т.к. сама генерация транзакций достаточно ресурсоёмка
	for i := 0; i < b.N; i++ {
		result := SumByCategoriesWithMutex(transactions, 0, 100)
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
		"5010": 10000000,
		"5020": 20000000,
	}
	b.ResetTimer() // сбрасываем таймер, т.к. сама генерация транзакций достаточно ресурсоёмка
	for i := 0; i < b.N; i++ {
		result := SumByCategoriesWithChannels(transactions, 0, 100)
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
		"5010": 10000000,
		"5020": 20000000,
	}
	b.ResetTimer() // сбрасываем таймер, т.к. сама генерация транзакций достаточно ресурсоёмка
	for i := 0; i < b.N; i++ {
		result := SumByCategoriesWithMutex2(transactions, 0, 100)
		b.StopTimer() // останавливаем таймер, чтобы время сравнения не учитывалось
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer() // продолжаем работу таймера
	}
}
