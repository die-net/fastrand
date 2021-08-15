package fastrand

import (
	"bytes"
	"compress/flate"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	t.Parallel()

	n := 1000000
	if testing.Short() {
		n = 10000
	}

	b := make([]byte, n)
	n, err := Read(b)
	if assert.NoError(t, err, "Read() shouldn't return an error") {
		assert.Equal(t, len(b), n, "Read() wrong number of bytes")
	}

	var z bytes.Buffer
	f, _ := flate.NewWriter(&z, 5)
	_, _ = f.Write(b)
	_ = f.Close()
	assert.Greater(t, z.Len(), len(b)*99/100, "shouldn't be able to compress random stream")
}

func TestReadEmpty(t *testing.T) {
	t.Parallel()

	n, err := Read(make([]byte, 0))
	if assert.NoError(t, err, "empty Read() shouldn't error") {
		assert.Empty(t, n, "shouldn't read any bytes")
	}
	n, err = Read(nil)
	if assert.NoError(t, err, "Read(nil) shouldn't error") {
		assert.Empty(t, n, "shouldn't read any bytes")
	}
}

const (
	maxUint   = ^uint(0)
	maxInt    = int(maxUint >> 1)
	maxUint64 = ^uint64(0)
	maxInt64  = int64(maxUint64 >> 1)
	maxUint32 = ^uint32(0)
	maxInt32  = int32(maxUint32 >> 1)
)

func TestUint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		max  uint64
		fn   func() uint64
	}{
		{"Uint64", maxUint64, Uint64},
		{"Uint32", uint64(maxUint32), func() uint64 { return uint64(Uint32()) }},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var min, max, set, unset uint64
			unset = maxUint64
			for i := 0; i < 100; i++ {
				v := test.fn()
				if v < min {
					min = v
				}
				if v > max {
					max = v
				}
				set |= v
				unset &= v
			}
			assert.Greater(t, max, test.max-(test.max>>2), "no output near expected max")
			assert.LessOrEqual(t, max, test.max, "shouldn't exceed test.max")
			assert.Less(t, min, test.max>>2, "no output near expected min")
			assert.Equal(t, test.max, set, "all bits should've been set at least once")
			assert.Empty(t, unset, "all bit should've been unset at least once")
		})
	}
}

func TestInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		max  int64
		fn   func() int64
	}{
		{"Int63", maxInt64, Int63},
		{"Int31", int64(maxInt32), func() int64 { return int64(Int31()) }},
		{"Int", int64(maxInt), func() int64 { return int64(Int()) }},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var min, max, set, unset int64
			unset = maxInt64
			for i := 0; i < 100; i++ {
				v := test.fn()
				if v < min {
					min = v
				}
				if v > max {
					max = v
				}
				set |= v
				unset &= v
			}
			assert.Greater(t, max, test.max-(test.max>>2), "no output near expected max")
			assert.LessOrEqual(t, max, test.max, "shouldn't exceed test.max")
			assert.Less(t, min, test.max>>2, "no output near expected min")
			assert.GreaterOrEqual(t, int64(0), min, "shouldn't be less than 0")
			assert.Equal(t, test.max, set, "all bits should've been set at least once")
			assert.Empty(t, unset, "all bit should've been unset at least once")
		})
	}
}

func TestIntn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		max  []int64
		fn   func(int64) int64
	}{
		{"Intn", []int64{1, 2, 3, 4, 128, 1000000, int64(maxInt32) - 1, int64(maxInt32), int64(maxInt32) + 1, int64(maxInt)}, func(n int64) int64 { return int64(Intn(int(n))) }},
		{"Int63n", []int64{1, 2, 3, 4, 128, 1000000, 10000000000, maxInt64}, Int63n},
		{"Int31n", []int64{1, 2, 3, 4, 128, 1000000, int64(maxInt32)}, func(n int64) int64 { return int64(Int31n(int32(n))) }},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			for _, testMax := range test.max {
				testMax := testMax

				t.Run(strconv.Itoa(int(testMax)), func(t *testing.T) {
					t.Parallel()

					var min, max int64
					for i := 0; i < 100; i++ {
						v := test.fn(testMax)
						if v < min {
							min = v
						}
						if v > max {
							max = v
						}
					}
					assert.GreaterOrEqual(t, max, testMax-(testMax>>2)-1, "no output near expected max")
					assert.Less(t, max, testMax, "should be less than testMax")
					assert.LessOrEqual(t, min, testMax>>2, "no output near expected min")
					assert.GreaterOrEqual(t, int64(0), min, "shouldn't be less than 0")
				})
			}
		})
	}
}

func TestRace(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			buf := make([]byte, 19)
			for j := 0; j < 10; j++ {
				_ = ExpFloat64()
				_ = Float32()
				_ = Float64()
				_ = Intn(Int())
				_ = Int31n(Int31())
				_ = Int63n(Int63())
				_ = NormFloat64()
				_ = Uint32()
				_ = Uint64()
				_ = Perm(10)
				_, _ = Read(buf)
				Shuffle(len(buf), func(i, j int) {
					buf[i], buf[j] = buf[j], buf[i]
				})
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
