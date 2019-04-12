package experiment

import (
	"crypto/rand"
	"math/big"
	"testing"
)

// I don't know I'm testing this, I know that will be different. 😅
func TestCryptoSeed(t *testing.T) {
	nCalls := 5

	expect := make([]*big.Int, nCalls, nCalls)
	for i := 0; i < nCalls; i++ {
		expect[i], _ = rand.Int(rand.Reader, big.NewInt(100000))
	}

	t.Run("sequential", func(t *testing.T) {
		seq := make([]*big.Int, nCalls, nCalls)
		for i := 0; i < nCalls; i++ {
			seq[i], _ = rand.Int(rand.Reader, big.NewInt(100000))
		}
		checkBigInt(t, expect, seq)
	})

	// t.Run("concurrent", func(t *testing.T) {
	// 	src := rand.NewSource(1)
	// 	rnd := rand.New(src)
	// 	sr := SafeRnd{mu: sync.Mutex{}, rnd: rnd}

	// 	seq := make([]float64, nCalls, nCalls)
	// 	processed := uint64(0)
	// 	returnCh := make(chan struct{})
	// 	for i := uint64(0); i < nCalls; i++ {
	// 		go func(i uint64) {
	// 			seq[i] = sr.Float64()
	// 			p := atomic.AddUint64(&processed, 1)
	// 			if p == nCalls {
	// 				returnCh <- struct{}{}
	// 			}
	// 		}(i)
	// 	}
	// 	<-returnCh
	// 	check(t, expect, seq)
	// })
}

func checkBigInt(t *testing.T, expect, seq []*big.Int) {
	for i := 0; i < len(expect); i++ {
		if expect[i].Cmp(seq[i]) != 0 {
			t.Errorf("don't match:\nexpected:%+v\nreceived:%+v\n", expect[i], seq[i])
		}
	}
}
