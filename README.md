FastRand [![Build Status](https://github.com/die-net/fastrand/actions/workflows/go-test.yml/badge.svg)](https://github.com/die-net/fastrand/actions/workflows/go-test.yml) [![Coverage Status](https://coveralls.io/repos/github/die-net/fastrand/badge.svg?branch=main)](https://coveralls.io/github/die-net/fastrand?branch=main) [![Go Report Card](https://goreportcard.com/badge/github.com/die-net/fastrand)](https://goreportcard.com/report/github.com/die-net/fastrand)
========

## This project is now archived. It has been obsoleted by Go 1.22's `math/rand/v2`, which uses per-thread random generator results.

---

FastRand exposes the Go internal `runtime.fastrand()`, a fast thread-local pseudorandom
number generator (PRNG).  On x86-64, this is based on hardware-accelerated
AES seeded by `/dev/urandom`.

This is a partial replacement for `math/rand` with the following features:

- Almost all of the global functions of `math/rand` are available, except `Seed()`.
- No seeding is required or allowed, meaning you can't test reproducible random sequences.
- You don't need to allocate or maintain goroutine-local copies to get adequate performance.

Because `runtime.fastrand()` doesn't have any locking overhead, simple 32-bit operations
like `fastrand.Uint32()` are 5.9x as fast as `math/rand.Uint32()`.  Because
`runtime.fastrand()` is 32-bit oriented, 64-bit operations like
`fastrand.Int63()` are only 2.5x as fast as `rand.Int63()`.

```
$ go test -cpu=1 -bench=.
goos: linux
goarch: amd64
pkg: github.com/die-net/fastrand
cpu: Intel(R) Xeon(R) Platinum 8175M CPU @ 2.50GHz
BenchmarkRead16/FastRand              68981697	    18.10 ns/op	       883.77 MB/s
BenchmarkRead16/MathRandLocal         61125805	    20.24 ns/op	       790.47 MB/s
BenchmarkRead16/MathRandGlobal        35250096	    33.95 ns/op	       471.22 MB/s
BenchmarkUint32/FastRand             434866173	     2.665 ns/op      1500.92 MB/s
BenchmarkUint32/MathRandLocal        285497366	     4.283 ns/op       933.86 MB/s
BenchmarkUint32/MathRandGlobal        76039005	    15.75 ns/op	       254.01 MB/s
BenchmarkIntn/FastRand               118486633	     9.983 ns/op       801.39 MB/s
BenchmarkIntn/MathRandLocal           90589077	    13.21 ns/op	       605.54 MB/s
BenchmarkIntn/MathRandGlobal          54485478	    22.15 ns/op	       361.09 MB/s
BenchmarkInt63/FastRand              190791255	     6.256 ns/op      1278.86 MB/s
BenchmarkInt63/MathRandLocal         233707926	     5.161 ns/op      1550.23 MB/s
BenchmarkInt63/MathRandGlobal         75825789	    15.80 ns/op	       506.37 MB/s
PASS
```

Because the `math/rand` global operations have a `sync.Mutex` around global state,
they slow down quite a bit under heavy load, if you somehow are generating a
lot of randomness.  To the contrary, fastrand gets faster almost linearly
with the number of CPUs.  To get similar performance out of `math/rand`, you
have to use call `math/rand.New()` and manage per goroutine state yourself.

```
$ go test -cpu=16 -bench=.
goos: linux
goarch: amd64
pkg: github.com/die-net/fastrand
cpu: Intel(R) Xeon(R) Platinum 8175M CPU @ 2.50GHz
BenchmarkRead16/FastRand-16          871821174	     1.161 ns/op     13786.61 MB/s
BenchmarkRead16/MathRandLocal-16     823987732	     1.250 ns/op     12798.27 MB/s
BenchmarkRead16/MathRandGlobal-16      4095805	   319.0 ns/op	        50.16 MB/s
BenchmarkUint32/FastRand-16         1000000000	     0.1719 ns/op    23275.32 MB/s
BenchmarkUint32/MathRandLocal-16    1000000000	     0.2824 ns/op    14165.72 MB/s
BenchmarkUint32/MathRandGlobal-16      5422761	   247.8 ns/op	        16.14 MB/s
BenchmarkIntn/FastRand-16           1000000000	     0.6563 ns/op    12188.68 MB/s
BenchmarkIntn/MathRandLocal-16      1000000000	     0.8258 ns/op     9687.31 MB/s
BenchmarkIntn/MathRandGlobal-16        5305188	   235.0 ns/op	        34.04 MB/s
BenchmarkInt63/FastRand-16          1000000000	     0.4474 ns/op    17880.81 MB/s
BenchmarkInt63/MathRandLocal-16     1000000000	     0.2998 ns/op    26688.37 MB/s
BenchmarkInt63/MathRandGlobal-16       5921886	   237.3 ns/op	        33.71 MB/s
PASS
```


License
-------
Copyright 2021 Aaron Hopkins and contributors.

Portions copied from Go's math/rand, Copyright 2009 The Go Authors.

All rights reserved.  Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
