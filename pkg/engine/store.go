package newflow

import (
	"errors"
	"sync"
	"time"
)

type StoreEngine struct {
	Window   time.Duration
	Interval time.Duration
	Lane     []Operator
	Pos      int
	Cap      int
	sync.RWMutex
}

type Operator struct {
	Elem []any
	Pos  int
}

func NewStoreEngine(windowSize, Interval time.Duration, limit int) (*StoreEngine, error) {
	if windowSize == 0 || Interval == 0 {
		return nil, errors.New("params can not be zero")
	}

	if windowSize <= Interval || windowSize%Interval != 0 {
		return nil, errors.New("illegal params value")
	}
	engine := &StoreEngine{
		Window:   windowSize,
		Interval: Interval,
		Lane:     make([]Operator, int(windowSize/Interval)),
		Pos:      0,
		Cap:      0,
	}
	// set each lane
	for i := 0; i < int(windowSize/Interval); i++ {
		engine.Lane[i] = Operator{
			Elem: make([]any, 0),
			Pos:  0,
		}
	}
	go engine.tiktok()
	return engine, nil
}

func (e *StoreEngine) tiktok() {
	ticker := time.NewTicker(e.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			e.Lock()
			if e.Pos = e.Pos + 1; e.Pos >= len(e.Lane) {
				e.Pos = 0
			}
			// reset this pipeline
			e.Lane[e.Pos].Elem = make([]any, 0)
			e.Lane[e.Pos].Pos = 0
			if e.Cap < len(e.Lane) {
				e.Cap++
			}
			e.Unlock()
		}
	}
}

func (e *StoreEngine) Push(elem any) {
	e.Lock()
	defer e.Unlock()
	lane := e.Lane[e.Pos]
	e.Lane[e.Pos].Elem = append(lane.Elem, elem)
	e.Lane[e.Pos].Pos++
}

func (e *StoreEngine) Pull() ([]any, int) {
	var ret []any
	e.RLock()
	defer e.RUnlock()
	pos := e.Pos
	for i := 1; i <= len(e.Lane); i++ {
		ret = append(ret, e.Lane[(pos+i)%len(e.Lane)].Elem...)
	}
	return ret, len(ret)
}
