package fastrand

import (
	"bytes"
	"compress/flate"
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

func TestRace(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			buf := make([]byte, 16)
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
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
