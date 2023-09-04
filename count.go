package neoflow

import (
	"errors"
	"sync"
	"time"
)

type CountEngine struct {
	Window   time.Duration
	Interval time.Duration
	Lane     []int64
	Pos      int
	Cap      int
	sync.RWMutex
}

func NewCountEngine(windowSize, Interval time.Duration) (*CountEngine, error) {
	if windowSize == 0 || Interval == 0 {
		return nil, errors.New("params can not be zero")
	}

	if windowSize < Interval || windowSize%Interval != 0 {
		return nil, errors.New("illegal params value")
	}
	engine := &CountEngine{
		Window:   windowSize,
		Interval: Interval,
		Lane:     make([]int64, int(windowSize/Interval)),
		Pos:      0,
		Cap:      0,
	}
	go engine.tiktok()
	return engine, nil
}

func (e *CountEngine) tiktok() {
	ticker := time.NewTicker(e.Interval)
	e.Cap++
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			e.Lock()
			if e.Pos = e.Pos + 1; e.Pos >= len(e.Lane) {
				e.Pos = 0
			}
			// reset this lane
			e.Lane[e.Pos] = 0
			if e.Cap < len(e.Lane) {
				e.Cap++
			}
			e.Unlock()
		}
	}
}

func (e *CountEngine) Add(v int64) {
	e.Lock()
	e.Lane[e.Pos] += v
	e.Unlock()
}

func (e *CountEngine) Total() int64 {
	e.RLock()
	defer e.RUnlock()
	var total int64
	pos := e.Pos
	for i := 1; i <= len(e.Lane); i++ {
		total += e.Lane[(pos+i)%len(e.Lane)]
	}
	return total
}

func (e *CountEngine) Average() int64 {
	return e.Total() / int64(e.Cap)
}
