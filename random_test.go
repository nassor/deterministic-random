package random

import "testing"

func BenchmarkRandomFloatPool(b *testing.B) {
	b.Run("sequential", func(b *testing.B) {
		var rnd *Random
		fp := NewFloatPool()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rnd = fp.Borrow()
			fp.Return(rnd)
		}
	})

	b.Run("parallel", func(b *testing.B) {
		fp := NewFloatPool()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			var rnd *Random
			for pb.Next() {
				rnd = fp.Borrow()
				fp.Return(rnd)
			}
		})
	})

	b.Run("no pool", func(b *testing.B) {
		fp := NewFloatPool()
		var rnd *Random
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rnd = &Random{}
			rnd.seed = fp.newFloatArray(&rnd.rnd)
		}
	})

}
