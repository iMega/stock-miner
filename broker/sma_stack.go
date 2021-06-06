package broker

import (
	"fmt"
	"sync"

	"github.com/imega/stock-miner/domain"
	"github.com/shopspring/decimal"
)

type smaStack struct {
	Stack      map[string]*smaFrame
	stackMutex sync.RWMutex
}

const capacity = 5

func NewSMAStack() domain.SMAStack {
	return &smaStack{
		Stack:      make(map[string]*smaFrame),
		stackMutex: sync.RWMutex{},
	}
}

func (s *smaStack) Add(key string, v float64) bool {
	s.stackMutex.Lock()
	defer s.stackMutex.Unlock()

	if _, ok := s.Stack[key]; !ok {
		s.Stack[key] = &smaFrame{}
	}

	if v == s.Stack[key].Lastt {
		return false
	}

	s.Stack[key].Add(v)

	return true
}

func (s *smaStack) IsTrendUp(key string) (bool, error) {
	s.stackMutex.RLock()
	defer s.stackMutex.RUnlock()

	if f, ok := s.Stack[key]; ok {
		return f.IsTrendUp(), nil
	}

	return false, fmt.Errorf("key does not exist")
}

func (s *smaStack) Get(key string) (domain.SMAFrame, error) {
	s.stackMutex.RLock()
	defer s.stackMutex.RUnlock()

	if f, ok := s.Stack[key]; ok {
		return f, nil
	}

	return nil, fmt.Errorf("key does not exist")
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
