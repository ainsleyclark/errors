<div align="center">
<img height="250" src="res/logo.svg" alt="Errors Logo" style="margin-bottom: 1rem" />

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/ainsleyclark/errors?color=success&label=version&sort=semver)
[![Go Report Card](https://goreportcard.com/badge/github.com/ainsleyclark/errors)](https://goreportcard.com/report/github.com/ainsleyclark/errors)
[![Maintainability](https://api.codeclimate.com/v1/badges/b3afd7bf115341995077/maintainability)](https://codeclimate.com/github/ainsleyclark/errors/maintainability)
[![Test](https://github.com/ainsleyclark/errors/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/ainsleyclark/errors/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/ainsleyclark/errors/branch/master/graph/badge.svg?token=K27L8LS7DA)](https://codecov.io/gh/ainsleyclark/errors)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/ainsleyclark/errors)

</div>


# Errors
A drop-in replacement for Go errors, with some added sugar! Error handling in Go made easy with codes, messages and more. Failure is your domain!


## Why?

## How to use

## Available Error Codes



## Benchmarks
Ran on 19/05/2022

```bash
$ go version
go version go1.18.1 darwin/amd64

$ go test -benchmem -bench .
goos: darwin
goarch: amd64
pkg: github.com/ainsleyclark/errors
cpu: AMD Ryzen 7 5800X 8-Core Processor
BenchmarkNew-16                          1000000              1026 ns/op            1320 B/op          6 allocs/op
BenchmarkNewInternal-16                  1000000              1023 ns/op            1320 B/op          6 allocs/op
BenchmarkError_Error-16                  6582385               180.6 ns/op           464 B/op          4 allocs/op
BenchmarkError_Code-16                  637953824                1.883 ns/op           0 B/op          0 allocs/op
BenchmarkError_Message-16               628838649                1.906 ns/op           0 B/op          0 allocs/op
BenchmarkError_ToError-16               705541435                1.699 ns/op           0 B/op          0 allocs/op
BenchmarkError_HTTPStatusCode-16        1000000000               0.6282 ns/op          0 B/op          0 allocs/op
```

## Contributing


## Credits
Shout out to the incredible [Maria Letta](https://github.com/MariaLetta) for her excellent Gopher illustrations
