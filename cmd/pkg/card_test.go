package card

import (
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
