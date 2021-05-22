package broker

import (
	"fmt"

	"github.com/imega/stock-miner/domain"
	"github.com/shopspring/decimal"
)

type smaStack map[string]*smaFrame

const capacity = 5

func NewSMAStack() domain.SMAStack {
	return make(smaStack)
}

func (s smaStack) Add(stack string, v float64) bool {
	if _, ok := s[stack]; !ok {
		s[stack] = &smaFrame{}
	}

	if v == s[stack].Lastt {
		return false
	}

	s[stack].Add(v)

	return true
}

func (s smaStack) IsTrendUp(stack string) (bool, error) {
	if f, ok := s[stack]; ok {
		return f.IsTrendUp(), nil
	}

	return false, fmt.Errorf("stack does not exist")
}

func (s smaStack) Get(stack string) (domain.SMAFrame, error) {
	if f, ok := s[stack]; ok {
		return f, nil
	}

	return nil, fmt.Errorf("stack does not exist")
}

type smaFrame struct {
	Avg   [2]float64
	Lastt float64
	Fifo  [capacity]float64
	Cur   int
}

func (s *smaFrame) Add(v float64) {
	f, _ := decimal.NewFromFloat(v).Truncate(2).Float64()
	s.Fifo[s.Cur] = f
	s.Lastt = f
	s.CalcAvg()
	s.NextCur()
}

func (s *smaFrame) NextCur() {
	if s.Cur == 4 {
		s.Cur = 0

		return
	}

	s.Cur++
}

func (s *smaFrame) CalcAvg() {
	r := decimal.NewFromFloat(0)
	for _, v := range s.Fifo {
		r = r.Add(decimal.NewFromFloat(v))
	}

	s.Avg[0] = s.Avg[1]
	f, _ := r.DivRound(decimal.NewFromInt(capacity), 4).Float64()
	s.Avg[1] = f
}

func (s *smaFrame) IsTrendUp() bool {
	return s.Avg[0] <= s.Avg[1]
}

func (s *smaFrame) Prev() float64 {
	prev := s.Cur - 2
	if s.Cur <= 1 {
		prev = capacity - 2 + s.Cur
	}

	return s.Fifo[prev]
}

func (s *smaFrame) Last() float64 {
	return s.Lastt
}

func (s *smaFrame) IsFull() bool {
	for _, v := range s.Fifo {
		if v == 0 {
			return false
		}
	}

	return true
}
