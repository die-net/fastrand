package fastrand

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkRead16(b *testing.B) {
	b.Run("FastRand", func(b *testing.B) {
		b.SetBytes(16)
		b.RunParallel(func(pb *testing.PB) {
			b := make([]byte, 16)

			for pb.Next() {
				_, _ = Read(b)
			}
		})
	})
	b.Run("MathRandLocal", func(b *testing.B) {
		b.SetBytes(16)
		b.RunParallel(func(pb *testing.PB) {
			b := make([]byte, 16)
			r := rand.New(rand.NewSource(time.Now().UnixNano()))

			for pb.Next() {
				_, _ = r.Read(b)
			}
		})
	})
	b.Run("MathRandGlobal", func(b *testing.B) {
		b.SetBytes(16)
		b.RunParallel(func(pb *testing.PB) {
			b := make([]byte, 16)
			for pb.Next() {
				_, _ = rand.Read(b)
			}
		})
	})
}

func BenchmarkUint32(b *testing.B) {
	b.Run("FastRand", func(b *testing.B) {
		b.SetBytes(4)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = Uint32()
			}
		})
	})
	b.Run("MathRandLocal", func(b *testing.B) {
		b.SetBytes(4)
		b.RunParallel(func(pb *testing.PB) {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			for pb.Next() {
				_ = r.Uint32()
			}
		})
	})
	b.Run("MathRandGlobal", func(b *testing.B) {
		b.SetBytes(4)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = rand.Uint32()
			}
		})
	})
}

func BenchmarkIntn(b *testing.B) {
	b.Run("FastRand", func(b *testing.B) {
		b.SetBytes(8)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = Intn(87654321)
			}
		})
	})
	b.Run("MathRandLocal", func(b *testing.B) {
		b.SetBytes(8)
		b.RunParallel(func(pb *testing.PB) {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))

			for pb.Next() {
				_ = r.Intn(87654321)
			}
		})
	})
	b.Run("MathRandGlobal", func(b *testing.B) {
		b.SetBytes(8)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = rand.Intn(87654321)
			}
		})
	})
}

func BenchmarkInt63(b *testing.B) {
	b.Run("FastRand", func(b *testing.B) {
		b.SetBytes(8)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = Int63()
			}
		})
	})
	b.Run("MathRandLocal", func(b *testing.B) {
		b.SetBytes(8)
		b.RunParallel(func(pb *testing.PB) {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))

			for pb.Next() {
				_ = r.Int63()
			}
		})
	})
	b.Run("MathRandGlobal", func(b *testing.B) {
		b.SetBytes(8)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = rand.Int63()
			}
		})
	})
}
