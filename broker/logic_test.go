package broker

import (
	"testing"

	"github.com/imega/stock-miner/domain"
	"github.com/stretchr/testify/assert"
)

func Test_calcTargetPrice(t *testing.T) {
	type args struct {
		commission  float64
		buyingPrice float64
		margin      float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "",
			args: args{
				commission:  0.3,
				buyingPrice: 100,
				margin:      0.2,
			},
			want: 100.8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcTargetPrice(tt.args.commission, tt.args.buyingPrice, tt.args.margin); got != tt.want {
				t.Errorf("calcTargetPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getItemsForSale(t *testing.T) {
	type args struct {
		slots []domain.Slot
		price float64
	}
	tests := []struct {
		name string
		args args
		want []domain.Slot
	}{
		{
			name: "one item equal from one",
			args: args{
				slots: []domain.Slot{
					{
						ID:          "1",
						TargetPrice: 100.01,
					},
				},
				price: 100.01,
			},
			want: []domain.Slot{
				{
					ID:          "1",
					TargetPrice: 100.01,
				},
			},
		},
		{
			name: "returns two slots",
			args: args{
				slots: []domain.Slot{
					{
						ID:          "0",
						TargetPrice: 100,
					},
					{
						ID:          "1",
						TargetPrice: 100.01,
					},
					{
						ID:          "2",
						TargetPrice: 100.02,
					},
				},
				price: 100.01,
			},
			want: []domain.Slot{
				{
					ID:          "0",
					TargetPrice: 100,
				},
				{
					ID:          "1",
					TargetPrice: 100.01,
				},
			},
		},
		{
			name: "returns empty",
			args: args{
				slots: []domain.Slot{
					{
						ID:          "0",
						TargetPrice: 100,
					},
					{
						ID:          "1",
						TargetPrice: 101,
					},
				},
				price: 99,
			},
			want: []domain.Slot{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getItemsForSale(tt.args.slots, tt.args.price)
			assert.Equal(t, got, tt.want)
		})
	}
}
