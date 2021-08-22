package broker

import (
	"errors"
	"sort"
	"sync"

	"github.com/imega/stock-miner/domain"
	"github.com/shopspring/decimal"
)

type smaStack struct {
	Stack      map[string]*smaFrame
	stackMutex sync.RWMutex
}

const (
	capacity   = 5 // only odd
	lastItem   = 4
	secondItem = 2
)

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

var errKeyNotExist = errors.New("key does not exist")

func (s *smaStack) IsTrendUp(key string) (bool, error) {
	s.stackMutex.RLock()
	defer s.stackMutex.RUnlock()

	if f, ok := s.Stack[key]; ok {
		return f.IsTrendUp(), nil
	}

	return false, errKeyNotExist
}

func (s *smaStack) Get(key string) (domain.SMAFrame, error) {
	s.stackMutex.RLock()
	defer s.stackMutex.RUnlock()

	if f, ok := s.Stack[key]; ok {
		return f, nil
	}

	return nil, errKeyNotExist
}

func (s *smaStack) Avg(key string) (float64, error) {
	s.stackMutex.RLock()
	defer s.stackMutex.RUnlock()

	if f, ok := s.Stack[key]; ok {
		return f.Avg[1], nil
	}

	return 0, errKeyNotExist
}

func (s *smaStack) Reset() {
	s.stackMutex.Lock()
	defer s.stackMutex.Unlock()

	for _, v := range s.Stack {
		for i := 0; i < capacity; i++ {
			v.Add(0)
		}
	}
}

type smaFrame struct {
	Avg       [2]float64
	Lastt     float64
	Fifo      [capacity]float64
	Cur       int
	RangeHigh float64
	RangeLow  float64
}

func (s *smaFrame) Add(v float64) {
	f, _ := decimal.NewFromFloat(v).Truncate(precision).Float64()
	s.Fifo[s.Cur] = f

	if s.RangeHigh > 0 && s.RangeHigh < f {
		s.RangeHigh = f
	}

	if s.RangeLow > 0 && s.RangeLow > f {
		s.RangeLow = f
	}

	s.Lastt = f
	s.CalcAvg()
	s.NextCur()
}

func (s *smaFrame) NextCur() {
	if s.Cur == lastItem {
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
	m := s.Last() - s.Median()
	a := s.Last() - s.Avg[1]

	return s.Avg[0] <= s.Avg[1] || m >= 0 || a >= 0
}

func (s *smaFrame) Prev() float64 {
	prev := s.Cur - secondItem
	if s.Cur <= 1 {
		prev = capacity - secondItem + s.Cur
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

func (s *smaFrame) RangeHL() (float64, float64) {
	return s.RangeHigh, s.RangeLow
}

func (s *smaFrame) SetRangeHL(h, l float64) {
	s.RangeHigh = h
	s.RangeLow = l
}

func (s *smaFrame) Median() float64 {
	tmp := [capacity]float64{}
	copy(tmp[:], s.Fifo[:])
	sort.Float64s(tmp[:])

	return tmp[(capacity-1)/2]
}

func (s *smaFrame) Distance() float64 {
	tmp := [capacity]float64{}
	copy(tmp[:], s.Fifo[:])
	sort.Float64s(tmp[:])

	return tmp[(capacity-1)] - tmp[0]
}
