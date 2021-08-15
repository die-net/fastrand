FastRand [![Build Status](https://github.com/die-net/fastrand/actions/workflows/go-test.yml/badge.svg)](https://github.com/die-net/fastrand/actions/workflows/go-test.yml) [![Coverage Status](https://coveralls.io/repos/github/die-net/fastrand/badge.svg?branch=master)](https://coveralls.io/github/die-net/fastrand?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/die-net/fastrand)](https://goreportcard.com/report/github.com/die-net/fastrand)
========

FastRand exposes the Go runtime's fastrand(), a fast thread-local pseudorandom
number generator (PRNG).  On x86-64, this is based on hardware-accelerated
AES seeded by /dev/urandom.

This is a partial replacement for math/rand with the following features:

- All of the global functions of math/rand are available, except Seed.
- No seeding is required or allowed.
- You don't need to allocate or maintain goroutine-local copies to get adequate performance.
- The performance is generally somewhat better than the local math/rand functioons, and can be substantially faster than the global functions under load.


License
=======
-------
Copyright 2021 Aaron Hopkins and contributors

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at: http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
