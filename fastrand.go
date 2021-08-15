package fastrand

import (
	"encoding/binary"
	_ "unsafe" // For go:linkname
)

// Uint64 returns a pseudo-random 64-bit value as a uint64.
func Uint64() uint64 {
	return uint64(fastrand())<<32 | uint64(fastrand())
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func Int63() int64 {
	return int64(Uint64() & ((1 << 63) - 1))
}

// Uint32 returns a pseudo-random 32-bit value as a uint32.
func Uint32() uint32 {
	return fastrand()
}

// Int31 returns a non-negative pseudo-random 31-bit integer as an int32.
func Int31() int32 {
	return int32(Uint32() & ((1 << 31) - 1))
}

const uintSize = 32 << (^uint(0) >> 32 & 1) // 32 or 64

// Int returns a non-negative pseudo-random int from the default Source.
func Int() int {
	if uintSize == 32 {
		return int(Int31())
	}
	return int(Int63())
}

// Intn returns, as an int, a non-negative pseudo-random number in [0,n).
// It panics if n <= 0.
func Intn(n int) int {
	if uintSize == 32 {
		return int(Int31n(int32(n)))
	}
	if n <= 1<<31 {
		return int(Int31n(int32(n)))
	}
	return int(Int63n(int64(n)))
}

// Int63n returns, as an int64, a non-negative pseudo-random number in [0,n).
// It panics if n <= 0.
func Int63n(n int64) int64 {
	if n <= 0 {
		panic("invalid argument to Int31n")
	}
	if n&(n-1) == 0 { // n is power of two, can mask
		return Int63() & (n - 1)
	}
	max := int64((1 << 63) - 1 - (1<<63)%uint64(n))
	v := Int63()
	for v > max {
		v = Int63()
	}
	return v % n
}

// Int31n returns, as an int32, a non-negative pseudo-random number in [0,n).
// It panics if n <= 0.
func Int31n(n int32) int32 {
	if n <= 0 {
		panic("invalid argument to Int31n")
	}
	if n&(n-1) == 0 { // n is power of two, can mask
		return Int31() & (n - 1)
	}
	max := int32((1 << 31) - 1 - (1<<31)%uint32(n))
	v := Int31()
	for v > max {
		v = Int31()
	}
	return v % n
}

// Perm returns, as a slice of n ints, a pseudo-random permutation of the integers [0,n).
func Perm(n int) []int {
	p := make([]int, n)

	// We start i at 1 because otherwise m[0] always swaps with m[0] on the first iteration.
	// This is basically a Fisher-Yates shuffle of consecutive integers.
	for i := 1; i < n; i++ {
		j := Intn(i + 1)
		p[i] = p[j]
		p[j] = i
	}
	return p
}

// Shuffle pseudo-randomizes the order of elements.
// n is the number of elements. Shuffle panics if n < 0.
// swap swaps the elements with indexes i and j.
func Shuffle(n int, swap func(i, j int)) {
	if n < 0 {
		panic("invalid argument to Shuffle")
	}

	// Fisher-Yates shuffle: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
	for i := 1; i < n; i++ {
		j := Intn(i + 1)
		swap(i, j)
	}
}

// Read generates len(p) random bytes and writes them into p.  It
// always returns len(p) and a nil error.
func Read(p []byte) (int, error) {
	n := 0

	// Fast path of copying 4 random bytes at a time.
	for len(p) >= 4 {
		binary.LittleEndian.PutUint32(p, Uint32())
		n += 4
		p = p[4:]
	}

	// We have between 0 and 3 bytes remaining to copy.
	if l := len(p); l > 0 {
		// Get 4 random bytes.
		var b [4]byte
		binary.LittleEndian.PutUint32(b[:], Uint32())

		// Copy up to 4 bytes to p.
		copy(p, b[:l])
		n += l
	}

	return n, nil
}

// fastrand is a fast thread local random function built into the Go runtime
// but not normally exposed.  On Linux x86_64, this is aesrand seeded by
// /dev/urandom.
//go:linkname fastrand runtime.fastrand
func fastrand() uint32
