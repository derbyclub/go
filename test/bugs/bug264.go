// $G $D/$F.go && $L $F.$A && ./$A.out || echo BUG

// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test case for http://code.google.com/p/go/issues/detail?id=692

package main

import "fmt"

var fooCount = 0
var barCount = 0
var balCount = 0

func foo() (int, int) {
	fooCount++
	fmt.Println("foo")
	return 0, 0
}

func bar() (int, int) {
	barCount++
	fmt.Println("bar")
	return 0, 0
}

func bal() (int, int) {
	balCount++
	fmt.Println("bal")
	return 0, 0
}

var a, b = foo() // foo is called once
var c, _ = bar() // bar is called twice
var _, _ = bal() // bal is called twice

func main() {
	if fooCount != 1 {
		panic("fooCount != 1")
	}
	if barCount != 1 {
		panic("barCount != 1")
	}
	if balCount != 1 {
		panic("balCount != 1")
	}
}
