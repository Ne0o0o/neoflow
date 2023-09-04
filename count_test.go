package neoflow

import (
	"math/rand"
	"testing"
	"time"
)

func TestNewCountEngine(t *testing.T) {
	if _, err := NewCountEngine(1*time.Second, 1*time.Second); err != nil {
		if err != nil {
			t.Error(err)
		}
	}
	if _, err := NewCountEngine(1*time.Second, 2*time.Second); err != nil {
		if err == nil {
			t.Error("excepted error but not")
		}
	}
	if _, err := NewCountEngine(0, 1*time.Second); err != nil {
		if err == nil {
			t.Error("excepted error but not")
		}
	}
	if _, err := NewCountEngine(1*time.Second, 0); err != nil {
		if err == nil {
			t.Error("excepted error but not")
		}
	}
}

func TestCountEngineAdd(t *testing.T) {
	e, err := NewCountEngine(5*time.Second, 1*time.Second)
	if err != nil {
		t.Error(err)
	}
	num := rand.Intn(100)
	e.Add(int64(num))
	if ret := e.Total(); ret != int64(num) {
		t.Errorf("add %d but receive %d", num, ret)
	}
	num1 := rand.Intn(1000)
	e.Add(int64(num1))
	if ret := e.Total(); ret != int64(num+num1) {
		t.Errorf("add %d but receive %d", num, ret)
	}
	num2 := rand.Intn(10000)
	e.Add(int64(num2))
	time.Sleep(2*time.Second + 500*time.Millisecond)
	if ret := e.Average(); ret != int64(num+num1+num2)/3 {
		t.Errorf("average %d but receive %d", int64(num+num1+num2)/3, ret)
	}
}
