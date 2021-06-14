package broker

import (
	"testing"

	"github.com/imega/stock-miner/domain"
	"github.com/stretchr/testify/assert"
)

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
						BuyingPrice: 1,
					},
				},
				price: 100.01,
			},
			want: []domain.Slot{
				{
					ID:          "1",
					TargetPrice: 100.01,
					BuyingPrice: 1,
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
						BuyingPrice: 1,
					},
					{
						ID:          "1",
						TargetPrice: 100.01,
						BuyingPrice: 1,
					},
					{
						ID:          "2",
						TargetPrice: 100.02,
						BuyingPrice: 1,
					},
				},
				price: 100.01,
			},
			want: []domain.Slot{
				{
					ID:          "0",
					TargetPrice: 100,
					BuyingPrice: 1,
				},
				{
					ID:          "1",
					TargetPrice: 100.01,
					BuyingPrice: 1,
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
						BuyingPrice: 1,
					},
					{
						ID:          "1",
						TargetPrice: 101,
						BuyingPrice: 1,
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

func Test_minBuyingPrice(t *testing.T) {
	type args struct {
		slots []domain.Slot
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "optimistic",
			args: args{
				slots: []domain.Slot{
					{BuyingPrice: 3},
					{BuyingPrice: 2},
					{BuyingPrice: 1},
				},
			},
			want: 1,
		},
		{
			name: "empty slots",
			args: args{
				slots: []domain.Slot{},
			},
			want: -1,
		},
		{
			name: "slots with zero price",
			args: args{
				slots: []domain.Slot{
					{BuyingPrice: 0},
					{BuyingPrice: 1},
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minBuyingPrice(tt.args.slots); got != tt.want {
				t.Errorf("minBuyingPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
