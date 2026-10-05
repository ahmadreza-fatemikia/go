// run

// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Comparing values through nil pointers must panic, even when
// both pointers are nil and the comparison is implemented by a
// call to runtime.memequal, which returns early without loading
// anything when its two pointer arguments are equal.

package main

import (
	"fmt"
	"runtime"
	"strings"
)

type S struct {
	a [40]byte
	b int
}

//go:noinline
func eq17(p, q *[17]byte) bool { return *p == *q }

//go:noinline
func eq24(p, q *[24]byte) bool { return *p == *q }

//go:noinline
func eq32(p, q *[32]byte) bool { return *p == *q }

//go:noinline
func eq33(p, q *[33]byte) bool { return *p == *q }

//go:noinline
func eq64(p, q *[64]byte) bool { return *p == *q }

//go:noinline
func eq1024(p, q *[1024]byte) bool { return *p == *q }

//go:noinline
func eqS(p, q *S) bool { return *p == *q }

//go:noinline
func eqSelf(p *[64]byte) bool { return *p == *p }

func main() {
	var x [64]byte
	var s S
	tests := []struct {
		name string
		f    func() bool
	}{
		{"eq17(nil, nil)", func() bool { return eq17(nil, nil) }},
		{"eq24(nil, nil)", func() bool { return eq24(nil, nil) }},
		{"eq32(nil, nil)", func() bool { return eq32(nil, nil) }},
		{"eq33(nil, nil)", func() bool { return eq33(nil, nil) }},
		{"eq64(nil, nil)", func() bool { return eq64(nil, nil) }},
		{"eq1024(nil, nil)", func() bool { return eq1024(nil, nil) }},
		{"eqS(nil, nil)", func() bool { return eqS(nil, nil) }},
		{"eqSelf(nil)", func() bool { return eqSelf(nil) }},
		{"eq64(nil, &x)", func() bool { return eq64(nil, &x) }},
		{"eq64(&x, nil)", func() bool { return eq64(&x, nil) }},
		{"eqS(nil, &s)", func() bool { return eqS(nil, &s) }},
		{"eqS(&s, nil)", func() bool { return eqS(&s, nil) }},
	}
	failed := false
	for _, tt := range tests {
		if err := shouldPanic(tt.f); err != nil {
			fmt.Printf("%s: %v\n", tt.name, err)
			failed = true
		}
	}
	if failed {
		panic("FAIL")
	}
}

func shouldPanic(f func() bool) (err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		re, ok := r.(runtime.Error)
		if !ok || !strings.Contains(re.Error(), "nil pointer dereference") {
			err = fmt.Errorf("want nil pointer dereference panic, got %v", r)
		}
	}()
	res := f()
	return fmt.Errorf("did not panic, returned %v", res)
}
