package broker

import "testing"

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
