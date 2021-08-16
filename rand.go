package fastrand

import (
	"math/rand"
)

// Wrap our Uint64() and Int63() in a math/rand.Source-compatible interface,
// so we can use math/rand's implementation of everything but the integer
// functions.

type fastSource struct{}

func (s *fastSource) Uint64() uint64 {
	return Uint64()
}

func (s *fastSource) Int63() int64 {
	return Int63()
}

func (s *fastSource) Seed(seed int64) {
	panic("Seed() not implemented")
}

var globalRand = rand.New(&fastSource{})

// Float64 returns, as a float64, a pseudo-random number in [0.0,1.0).
func Float64() float64 { return globalRand.Float64() }

// Float32 returns, as a float32, a pseudo-random number in [0.0,1.0).
func Float32() float32 { return globalRand.Float32() }

// NormFloat64 returns a normally distributed float64 in the range
// [-math.MaxFloat64, +math.MaxFloat64] with
// standard normal distribution (mean = 0, stddev = 1).
// To produce a different normal distribution, callers can
// adjust the output using:
//
//  sample = NormFloat64() * desiredStdDev + desiredMean
//
func NormFloat64() float64 { return globalRand.NormFloat64() }

// ExpFloat64 returns an exponentially distributed float64 in the range
// (0, +math.MaxFloat64] with an exponential distribution whose rate parameter
// (lambda) is 1 and whose mean is 1/lambda (1).
// To produce a distribution with a different rate parameter,
// callers can adjust the output using:
//
//  sample = ExpFloat64() / desiredRateParameter
//
func ExpFloat64() float64 { return globalRand.ExpFloat64() }

// NewZipf returns a Zipf variate generator.
// The generator generates values k ∈ [0, imax]
// such that P(k) is proportional to (v + k) ** (-s).
// Requirements: s > 1 and v >= 1.
//
// Call the Uint64 method on the returned object to get a value drawn from
// the Zipf distribution.
func NewZipf(s, v float64, imax uint64) *rand.Zipf {
	return rand.NewZipf(globalRand, s, v, imax)
}
