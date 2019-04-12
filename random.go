package random

import (
	"math/rand"
	"sync"
	"time"
)

const size = 100

type Random struct {
	mu       sync.Mutex
	seed     int64
	lastUsed int
	rnd      [size]float64
}

func (r *Random) Float64() float64 {
	r.mu.Lock()
	r.lastUsed++
	if r.lastUsed > len(r.rnd)-1 {
		r.lastUsed = 0
	}
	f := r.rnd[r.lastUsed]
	r.mu.Unlock()
	return f
}

type FloatPool struct {
	pool sync.Pool
}

func NewFloatPool() FloatPool {
	fp := FloatPool{pool: sync.Pool{}}
	fp.pool.New = func() interface{} {
		r := &Random{rnd: [size]float64{}}
		r.seed = fp.newFloatArray(&r.rnd)
		return r
	}
	return fp
}

func (*FloatPool) newFloatArray(fa *[size]float64) int64 {
	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	rnd := rand.New(src)
	r := &Random{}
	for i := 0; i < len(fa); i++ {
		r.rnd[i] = rnd.Float64()
	}
	return seed
}

func (fp *FloatPool) Borrow() *Random {
	r := fp.pool.Get().(*Random)
	r.seed = fp.newFloatArray(&r.rnd)
	r.lastUsed = 0
	return r
}

func (fp *FloatPool) Return(r *Random) {
	fp.pool.Put(r)
}
