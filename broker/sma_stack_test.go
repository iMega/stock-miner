package broker

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_smaStack(t *testing.T) {
	st := &smaStack{
		Stack:      make(map[string]*smaFrame),
		stackMutex: sync.RWMutex{},
	}

	st.Add("AAPL", 1)
	st.Add("AAPL", 2)
	st.Add("AAPL", 3)
	st.Add("AAPL", 4)
	st.Add("AAPL", 5)
	st.Add("AAPL", 6)

	if st.Stack["AAPL"].Avg[0] != 3 {
		t.Fatalf("failed to calc avg frame")
	}

	if st.Stack["AAPL"].Avg[1] != 4 {
		t.Fatalf("failed to calc avg frame")
	}

	if !st.Stack["AAPL"].IsTrendUp() {
		t.Fatalf("failed to calc trend frame")
	}

	if v, err := st.IsTrendUp("AAPL"); err != nil || !v {
		t.Fatalf("failed to calc outer trend frame")
	}
}

func BenchmarkRingAdd(b *testing.B) {
	b.ReportAllocs()

	st := &smaStack{
		Stack:      make(map[string]*smaFrame),
		stackMutex: sync.RWMutex{},
	}
	for i := 0; i < b.N; i++ {
		st.Add("AAPL", float64(i))
	}
}

func est_Regression_1(t *testing.T) {
	st := &smaStack{
		Stack:      make(map[string]*smaFrame),
		stackMutex: sync.RWMutex{},
	}

	st.Add("AAPL", 125.38)
	st.Add("AAPL", 125.39)
	st.Add("AAPL", 125.38)
	st.Add("AAPL", 125.4)
	st.Add("AAPL", 125.38)

	st.Add("AAPL", 125.3)

	if v, err := st.IsTrendUp("AAPL"); err != nil || v == true {
		t.Fatalf("failed to calc outer trend frame")
	}

	st.Add("AAPL", 125.33)
	if v, err := st.IsTrendUp("AAPL"); err != nil || v == true {
		t.Fatalf("failed to calc outer trend frame")
	}

	st.Add("AAPL", 125.3)
	if v, err := st.IsTrendUp("AAPL"); err != nil || v == true {
		t.Fatalf("failed to calc outer trend frame")
	}

	// st.Add("AAPL", 125.31)
	// if v, err := st.IsTrendUp("AAPL"); err != nil || v == true {
	// 	t.Fatalf("failed to calc outer trend frame")
	// }

	st.Add("AAPL", 125.38)
	if v, err := st.IsTrendUp("AAPL"); err != nil || v == false {
		t.Fatalf("failed to calc outer trend frame")
	}

	st.Add("AAPL", 125.31)
	if v, err := st.IsTrendUp("AAPL"); err != nil || v == false {
		t.Fatalf("failed to calc outer trend frame")
	}
}

func Test_PrevPrice(t *testing.T) {
	st := &smaStack{
		Stack:      make(map[string]*smaFrame),
		stackMutex: sync.RWMutex{},
	}

	st.Add("AAPL", 1)

	for i := 2; i < 11; i++ {
		st.Add("AAPL", float64(i))

		frame, err := st.Get("AAPL")
		if err != nil {
			t.Fatalf("failed getting frame")
		}

		if frame.Prev() != float64(i-1) {
			t.Fatalf("failed getting previous price, %f not equal %f", frame.Prev(), float64(i-1))
		}
	}

}

func Test_smaFrame_Median(t *testing.T) {
	type fields struct {
		Avg       [2]float64
		Lastt     float64
		Fifo      [capacity]float64
		Cur       int
		RangeHigh float64
		RangeLow  float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			fields: fields{
				Fifo: [capacity]float64{5, 3, 2, 4, 1},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &smaFrame{
				Avg:       tt.fields.Avg,
				Lastt:     tt.fields.Lastt,
				Fifo:      tt.fields.Fifo,
				Cur:       tt.fields.Cur,
				RangeHigh: tt.fields.RangeHigh,
				RangeLow:  tt.fields.RangeLow,
			}
			if got := s.Median(); got != tt.want {
				t.Errorf("smaFrame.Median() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_smaStack_Reset(t *testing.T) {
	s := &smaStack{
		Stack: map[string]*smaFrame{
			"TEST": {
				Fifo: [capacity]float64{1, 2, 3, 4, 5},
			},
			"TEST2": {
				Fifo: [capacity]float64{1, 2, 3, 4, 5},
			},
		},
		stackMutex: sync.RWMutex{},
	}
	s.Reset()

	f, err := s.Get("TEST2")
	if err != nil {
		t.Errorf("failed getting stack, %s", err)
	}

	frame, ok := f.(*smaFrame)
	if !ok {
		t.Error("failed to cast smaFrame")
	}

	expected := [capacity]float64{0, 0, 0, 0, 0}

	assert.Equal(t, expected, frame.Fifo)
}
