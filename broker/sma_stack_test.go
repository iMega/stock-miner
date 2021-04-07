package broker

import (
	"testing"
)

func Test_smaStack(t *testing.T) {
	st := make(smaStack)

	st.Add("AAPL", 1)
	st.Add("AAPL", 2)
	st.Add("AAPL", 3)
	st.Add("AAPL", 4)
	st.Add("AAPL", 5)
	st.Add("AAPL", 6)

	if st["AAPL"].Avg[0] != 3 {
		t.Fatalf("failed to calc avg frame")
	}

	if st["AAPL"].Avg[1] != 4 {
		t.Fatalf("failed to calc avg frame")
	}

	if !st["AAPL"].IsTrendUp() {
		t.Fatalf("failed to calc trend frame")
	}

	if v, err := st.IsTrendUp("AAPL"); err != nil || !v {
		t.Fatalf("failed to calc outer trend frame")
	}
}

func BenchmarkRingAdd(b *testing.B) {
	b.ReportAllocs()

	st := make(smaStack)
	for i := 0; i < b.N; i++ {
		st.Add("AAPL", float64(i))
	}
}

func Test_Regression_1(t *testing.T) {
	st := make(smaStack)

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

	st.Add("AAPL", 125.31)
	if v, err := st.IsTrendUp("AAPL"); err != nil || v == true {
		t.Fatalf("failed to calc outer trend frame")
	}

	st.Add("AAPL", 125.38)
	if v, err := st.IsTrendUp("AAPL"); err != nil || v == false {
		t.Fatalf("failed to calc outer trend frame")
	}

	st.Add("AAPL", 125.31)
	if v, err := st.IsTrendUp("AAPL"); err != nil || v == false {
		t.Fatalf("failed to calc outer trend frame")
	}
}
