package experiment

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
)

type SafeRnd struct {
	mu  sync.Mutex
	rnd *rand.Rand
}

func (sr *SafeRnd) Float64() float64 {
	sr.mu.Lock()
	v := sr.rnd.Float64()
	sr.mu.Unlock()
	return v
}

func TestSeed(t *testing.T) {
	nCalls := 9000

	src := rand.NewSource(1)
	rnd := rand.New(src)
	expect := make([]float64, nCalls, nCalls)
	for i := 0; i < nCalls; i++ {
		expect[i] = rnd.Float64()
	}

	// Simple and easy, but you must keep your processing sequential
	t.Run("sequential generation and access", func(t *testing.T) {
		src := rand.NewSource(1)
		rnd := rand.New(src)
		seq := make([]float64, nCalls, nCalls)
		for i := 0; i < nCalls; i++ {
			seq[i] = rnd.Float64()
		}
		check(t, expect, seq)
	})

	// It will forces you to keep list of random generator and last used during your processing pipeline
	t.Run("sequential generation and control with concurrent access", func(t *testing.T) {
		src := rand.NewSource(1)
		rnd := rand.New(src)

		// pre-generate all calls
		seq := make([]float64, nCalls, nCalls)
		for i := 0; i < nCalls; i++ {
			seq[i] = rnd.Float64()
		}

		use := make([]float64, nCalls, nCalls)
		processed := uint64(0)
		returnCh := make(chan struct{})
		for i := 0; i < nCalls; i++ {
			go func(i int) {
				use[i] = seq[i]
				p := atomic.AddUint64(&processed, 1)
				if int(p) == nCalls {
					returnCh <- struct{}{}
				}
			}(i)
		}
		<-returnCh
		check(t, expect, use)
	})

	// Run multiple times, sometimes will pass, sometimes don't
	t.Run("concurrent generation and access", func(t *testing.T) {
		src := rand.NewSource(1)
		rnd := rand.New(src)
		sr := SafeRnd{mu: sync.Mutex{}, rnd: rnd}

		seq := make([]float64, nCalls, nCalls)
		processed := uint64(0)
		returnCh := make(chan struct{})
		for i := 0; i < nCalls; i++ {
			go func(i int) {
				seq[i] = sr.Float64()
				p := atomic.AddUint64(&processed, 1)
				if int(p) == nCalls {
					returnCh <- struct{}{}
				}
			}(i)
		}
		<-returnCh
		check(t, expect, seq)
	})
}

func check(t *testing.T, expect, seq []float64) {
	for i := 0; i < len(expect); i++ {
		if expect[i] != seq[i] {
			t.Errorf("don't match:\nexpected:%+v\nreceived:%+v\n", expect[i], seq[i])
		}
	}
}
