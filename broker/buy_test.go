package broker

import (
	"testing"

	"github.com/imega/stock-miner/domain"
	"github.com/stretchr/testify/assert"
)

func Test_fillBuyTransaction(t *testing.T) {
	type args struct {
		dst domain.Transaction
		src domain.Transaction
		s   domain.Settings
	}
	tests := []struct {
		name string
		args args
		want domain.Transaction
	}{
		{
			name: "optimistic",
			args: args{
				dst: domain.Transaction{},
				src: domain.Transaction{
					Slot: domain.Slot{
						BuyingPrice: 100,
						AmountSpent: 200,
					},
				},
				s: domain.Settings{
					MarketCommission: 10,
					GrossMargin:      0,
				},
			},
			want: domain.Transaction{
				Slot: domain.Slot{
					BuyingPrice: 100,
					AmountSpent: 200,
					Profit:      21,
					TargetPrice: 121,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fillBuyTransaction(tt.args.dst, tt.args.src, tt.args.s)
			assert.Equal(t, got, tt.want)
		})
	}
}
