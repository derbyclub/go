// errchk -0 $G -m $D/$F.go

// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package foo

import "unsafe"

var gxx *int

func foo1(x int) { // ERROR "moved to heap: NAME-x"
	gxx = &x  // ERROR "&x escapes to heap"
}

func foo2(yy *int) { // ERROR "leaking param: NAME-yy"
	gxx = yy
}

func foo3(x int) *int { // ERROR "moved to heap: NAME-x"
	return &x  // ERROR "&x escapes to heap"
}

type T *T

func foo3b(t T) { // ERROR "leaking param: NAME-t"
	*t = t
}

// xx isn't going anywhere, so use of yy is ok
func foo4(xx, yy *int) { // ERROR "xx does not escape" "yy does not escape"
	xx = yy
}

// xx isn't going anywhere, so taking address of yy is ok
func foo5(xx **int, yy *int) { // ERROR "xx does not escape" "yy does not escape"
	xx = &yy  // ERROR "&yy does not escape"
}

func foo6(xx **int, yy *int) { // ERROR "xx does not escape" "leaking param: NAME-yy"
	*xx = yy
}

func foo7(xx **int, yy *int) { // ERROR "xx does not escape" "yy does not escape"
	**xx = *yy
}

func foo8(xx, yy *int) int { // ERROR "xx does not escape" "yy does not escape"
	xx = yy
	return *xx
}

func foo9(xx, yy *int) *int { // ERROR "leaking param: NAME-xx" "leaking param: NAME-yy"
	xx = yy
	return xx
}

func foo10(xx, yy *int) { // ERROR "xx does not escape" "yy does not escape"
	*xx = *yy
}

func foo11() int {
	x, y := 0, 42
	xx := &x  // ERROR "&x does not escape"
	yy := &y  // ERROR "&y does not escape"
	*xx = *yy
	return x
}

var xxx **int

func foo12(yyy **int) { // ERROR "leaking param: NAME-yyy"
	xxx = yyy
}

func foo13(yyy **int) { // ERROR "yyy does not escape"
	*xxx = *yyy
}

func foo14(yyy **int) { // ERROR "yyy does not escape"
	**xxx = **yyy
}

func foo15(yy *int) { // ERROR "moved to heap: NAME-yy"
	xxx = &yy  // ERROR "&yy escapes to heap"
}

func foo16(yy *int) { // ERROR "leaking param: NAME-yy"
	*xxx = yy
}

func foo17(yy *int) { // ERROR "yy does not escape"
	**xxx = *yy
}

func foo18(y int) { // ERROR "moved to heap: "NAME-y"
	*xxx = &y  // ERROR "&y escapes to heap"
}

func foo19(y int) {
	**xxx = y
}

type Bar struct {
	i  int
	ii *int
}

func NewBar() *Bar {
	return &Bar{42, nil} // ERROR "&struct literal escapes to heap"
}

func NewBarp(x *int) *Bar { // ERROR "leaking param: NAME-x"
	return &Bar{42, x} // ERROR "&struct literal escapes to heap"
}

func NewBarp2(x *int) *Bar { // ERROR "x does not escape"
	return &Bar{*x, nil} // ERROR "&struct literal escapes to heap"
}

func (b *Bar) NoLeak() int { // ERROR "b does not escape"
	return *(b.ii)
}

func (b *Bar) AlsoNoLeak() *int { // ERROR "b does not escape"
	return b.ii
}

type Bar2 struct {
	i  [12]int
	ii []int
}

func NewBar2() *Bar2 {
	return &Bar2{[12]int{42}, nil} // ERROR "&struct literal escapes to heap"
}

func (b *Bar2) NoLeak() int { // ERROR "b does not escape"
	return b.i[0]
}

func (b *Bar2) Leak() []int { // ERROR "leaking param: NAME-b"
	return b.i[:]  // ERROR "&b.i escapes to heap"
}

func (b *Bar2) AlsoNoLeak() []int { // ERROR "b does not escape"
	return b.ii[0:1]
}

func (b *Bar2) LeakSelf() { // ERROR "leaking param: NAME-b"
	b.ii = b.i[0:4]  // ERROR "&b.i escapes to heap"
}

func (b *Bar2) LeakSelf2() { // ERROR "leaking param: NAME-b"
	var buf []int
	buf = b.i[0:]  // ERROR "&b.i escapes to heap"
	b.ii = buf
}

func foo21() func() int {
	x := 42 // ERROR "moved to heap: NAME-x"
	return func() int {  // ERROR "func literal escapes to heap"
		return x  // ERROR "&x escapes to heap"
	}
}

func foo22() int {
	x := 42
	return func() int {  // ERROR "func literal does not escape"
		return x
	}()
}

func foo23(x int) func() int { // ERROR "moved to heap: NAME-x"
	return func() int {  // ERROR "func literal escapes to heap"
		return x  // ERROR "&x escapes to heap"
	}
}

func foo23a(x int) func() int { // ERROR "moved to heap: NAME-x"
	f := func() int {  // ERROR "func literal escapes to heap"
		return x  // ERROR "&x escapes to heap"
	}
	return f
}

func foo23b(x int) *(func() int) { // ERROR "moved to heap: NAME-x"
	f := func() int { return x } // ERROR "moved to heap: NAME-f" "func literal escapes to heap" "&x escapes to heap"
	return &f  // ERROR "&f escapes to heap"
}

func foo24(x int) int {
	return func() int {  // ERROR "func literal does not escape"
		return x
	}()
}

var x *int

func fooleak(xx *int) int { // ERROR "leaking param: NAME-xx"
	x = xx
	return *x
}

func foonoleak(xx *int) int { // ERROR "xx does not escape"
	return *x + *xx
}

func foo31(x int) int { // ERROR "moved to heap: NAME-x"
	return fooleak(&x)  // ERROR "&x escapes to heap"
}

func foo32(x int) int {
	return foonoleak(&x)  // ERROR "&x does not escape"
}

type Foo struct {
	xx *int
	x  int
}

var F Foo
var pf *Foo

func (f *Foo) fooleak() { // ERROR "leaking param: NAME-f"
	pf = f
}

func (f *Foo) foonoleak() { // ERROR "f does not escape"
	F.x = f.x
}

func (f *Foo) Leak() { // ERROR "leaking param: NAME-f"
	f.fooleak()
}

func (f *Foo) NoLeak() { // ERROR "f does not escape"
	f.foonoleak()
}

func foo41(x int) { // ERROR "moved to heap: NAME-x"
	F.xx = &x  // ERROR "&x escapes to heap"
}

func (f *Foo) foo42(x int) { // ERROR "f does not escape" "moved to heap: NAME-x"
	f.xx = &x  // ERROR "&x escapes to heap"
}

func foo43(f *Foo, x int) { // ERROR "f does not escape" "moved to heap: NAME-x"
	f.xx = &x  // ERROR "&x escapes to heap"
}

func foo44(yy *int) { // ERROR "leaking param: NAME-yy"
	F.xx = yy
}

func (f *Foo) foo45() { // ERROR "f does not escape"
	F.x = f.x
}

func (f *Foo) foo46() { // ERROR "f does not escape"
	F.xx = f.xx
}

func (f *Foo) foo47() { // ERROR "leaking param: NAME-f"
	f.xx = &f.x  // ERROR "&f.x escapes to heap"
}

var ptrSlice []*int

func foo50(i *int) { // ERROR "leaking param: NAME-i"
	ptrSlice[0] = i
}

var ptrMap map[*int]*int

func foo51(i *int) { // ERROR "leaking param: NAME-i"
	ptrMap[i] = i
}

func indaddr1(x int) *int { // ERROR "moved to heap: NAME-x"
	return &x  // ERROR "&x escapes to heap"
}

func indaddr2(x *int) *int { // ERROR "leaking param: NAME-x"
	return *&x  // ERROR "&x does not escape"
}

func indaddr3(x *int32) *int { // ERROR "leaking param: NAME-x"
	return *(**int)(unsafe.Pointer(&x))  // ERROR "&x does not escape"
}

// From package math:

func Float32bits(f float32) uint32 {
	return *(*uint32)(unsafe.Pointer(&f))  // ERROR "&f does not escape"
}

func Float32frombits(b uint32) float32 {
	return *(*float32)(unsafe.Pointer(&b))  // ERROR "&b does not escape"
}

func Float64bits(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))  // ERROR "&f does not escape"
}

func Float64frombits(b uint64) float64 {
	return *(*float64)(unsafe.Pointer(&b))  // ERROR "&b does not escape"
}

// contrast with
func float64bitsptr(f float64) *uint64 { // ERROR "moved to heap: NAME-f"
	return (*uint64)(unsafe.Pointer(&f))  // ERROR "&f escapes to heap"
}

func float64ptrbitsptr(f *float64) *uint64 { // ERROR "leaking param: NAME-f"
	return (*uint64)(unsafe.Pointer(f))
}

func typesw(i interface{}) *int { // ERROR "leaking param: NAME-i"
	switch val := i.(type) {
	case *int:
		return val
	case *int8:
		v := int(*val) // ERROR "moved to heap: NAME-v"
		return &v  // ERROR "&v escapes to heap"
	}
	return nil
}

func exprsw(i *int) *int { // ERROR "leaking param: NAME-i"
	switch j := i; *j + 110 {
	case 12:
		return j
	case 42:
		return nil
	}
	return nil

}

// assigning to an array element is like assigning to the array
func foo60(i *int) *int { // ERROR "leaking param: NAME-i"
	var a [12]*int
	a[0] = i
	return a[1]
}

func foo60a(i *int) *int { // ERROR "i does not escape"
	var a [12]*int
	a[0] = i
	return nil
}

// assigning to a struct field  is like assigning to the struct
func foo61(i *int) *int { // ERROR "leaking param: NAME-i"
	type S struct {
		a, b *int
	}
	var s S
	s.a = i
	return s.b
}

func foo61a(i *int) *int { // ERROR "i does not escape"
	type S struct {
		a, b *int
	}
	var s S
	s.a = i
	return nil
}

// assigning to a struct field is like assigning to the struct but
// here this subtlety is lost, since s.a counts as an assignment to a
// track-losing dereference.
func foo62(i *int) *int { // ERROR "leaking param: NAME-i"
	type S struct {
		a, b *int
	}
	s := new(S) // ERROR "new[(]S[)] does not escape"
	s.a = i
	return nil // s.b
}

type M interface {
	M()
}

func foo63(m M) { // ERROR "m does not escape"
}

func foo64(m M) { // ERROR "leaking param: NAME-m"
	m.M()
}

type MV int

func (MV) M() {}

func foo65() {
	var mv MV
	foo63(&mv)  // ERROR "&mv does not escape"
}

func foo66() {
	var mv MV // ERROR "moved to heap: NAME-mv"
	foo64(&mv)  // ERROR "&mv escapes to heap"
}

func foo67() {
	var mv MV
	foo63(mv)
}

func foo68() {
	var mv MV
	foo64(mv) // escapes but it's an int so irrelevant
}

func foo69(m M) { // ERROR "leaking param: NAME-m"
	foo64(m)
}

func foo70(mv1 *MV, m M) { // ERROR "leaking param: NAME-mv1" "leaking param: NAME-m"
	m = mv1
	foo64(m)
}

func foo71(x *int) []*int { // ERROR "leaking param: NAME-x"
	var y []*int
	y = append(y, x)
	return y
}

func foo71a(x int) []*int { // ERROR "moved to heap: NAME-x"
	var y []*int
	y = append(y, &x)  // ERROR "&x escapes to heap"
	return y
}

func foo72() {
	var x int
	var y [1]*int
	y[0] = &x  // ERROR "&x does not escape"
}

func foo72aa() [10]*int {
	var x int // ERROR "moved to heap: NAME-x"
	var y [10]*int
	y[0] = &x  // ERROR "&x escapes to heap"
	return y
}

func foo72a() {
	var y [10]*int
	for i := 0; i < 10; i++ {
		// escapes its scope
		x := i // ERROR "moved to heap: NAME-x"
		y[i] = &x // ERROR "&x escapes to heap"
	}
	return
}

func foo72b() [10]*int {
	var y [10]*int
	for i := 0; i < 10; i++ {
		x := i // ERROR "moved to heap: NAME-x"
		y[i] = &x  // ERROR "&x escapes to heap"
	}
	return y
}

// issue 2145
func foo73() {
	s := []int{3, 2, 1} // ERROR "slice literal does not escape"
	for _, v := range s {
		vv := v        // ERROR "moved to heap: NAME-vv"
		// actually just escapes its scope
		defer func() { // ERROR "func literal escapes to heap"
			println(vv)  // ERROR "&vv escapes to heap"
		}()
	}
}

func foo74() {
	s := []int{3, 2, 1} // ERROR "slice literal does not escape"
	for _, v := range s {
		vv := v        // ERROR "moved to heap: NAME-vv"
		// actually just escapes its scope
		fn := func() { // ERROR "func literal escapes to heap"
			println(vv)  // ERROR "&vv escapes to heap"
		}
		defer fn()
	}
}

func myprint(y *int, x ...interface{}) *int { // ERROR "x does not escape" "leaking param: NAME-y"
	return y
}

func myprint1(y *int, x ...interface{}) *interface{} { // ERROR "y does not escape" "leaking param: NAME-x"
	return &x[0]  // ERROR "&x.0. escapes to heap"
}

func foo75(z *int) { // ERROR "leaking param: NAME-z"
	myprint(z, 1, 2, 3) // ERROR "[.][.][.] argument does not escape"
}

func foo75a(z *int) { // ERROR "z does not escape"
	myprint1(z, 1, 2, 3) // ERROR "[.][.][.] argument escapes to heap"
}

func foo76(z *int) { // ERROR "leaking param: NAME-z"
	myprint(nil, z) // ERROR "[.][.][.] argument does not escape"
}

func foo76a(z *int) { // ERROR "leaking param: NAME-z"
	myprint1(nil, z) // ERROR "[.][.][.] argument escapes to heap"
}

func foo76b() {
	myprint(nil, 1, 2, 3) // ERROR "[.][.][.] argument does not escape"
}

func foo76c() {
	myprint1(nil, 1, 2, 3) // ERROR "[.][.][.] argument escapes to heap"
}

func foo76d() {
	defer myprint(nil, 1, 2, 3) // ERROR "[.][.][.] argument does not escape"
}

func foo76e() {
	defer myprint1(nil, 1, 2, 3) // ERROR "[.][.][.] argument escapes to heap"
}

func foo76f() {
	for {
		// TODO: This one really only escapes its scope, but we don't distinguish yet.
		defer myprint(nil, 1, 2, 3) // ERROR "[.][.][.] argument escapes to heap"
	}
}

func foo76g() {
	for {
		defer myprint1(nil, 1, 2, 3) // ERROR "[.][.][.] argument escapes to heap"
	}
}

func foo77(z []interface{}) { // ERROR "z does not escape"
	myprint(nil, z...) // z does not escape
}

func foo77a(z []interface{}) { // ERROR "leaking param: NAME-z"
	myprint1(nil, z...)
}

func foo78(z int) *int { // ERROR "moved to heap: NAME-z"
	return &z  // ERROR "&z escapes to heap"
}

func foo78a(z int) *int { // ERROR "moved to heap: NAME-z"
	y := &z  // ERROR "&z escapes to heap"
	x := &y  // ERROR "&y does not escape"
	return *x // really return y
}

func foo79() *int {
	return new(int) // ERROR "new[(]int[)] escapes to heap"
}

func foo80() *int {
	var z *int
	for {
		// Really just escapes its scope but we don't distinguish
		z = new(int) // ERROR "new[(]int[)] escapes to heap"
	}
	_ = z
	return nil
}

func foo81() *int {
	for {
		z := new(int) // ERROR "new[(]int[)] does not escape"
		_ = z
	}
	return nil
}

type Fooer interface {
	Foo()
}

type LimitedFooer struct {
	Fooer
	N int64
}

func LimitFooer(r Fooer, n int64) Fooer { // ERROR "leaking param: NAME-r"
	return &LimitedFooer{r, n} // ERROR "&struct literal escapes to heap"
}

func foo90(x *int) map[*int]*int { // ERROR "leaking param: NAME-x"
	return map[*int]*int{nil: x} // ERROR "map literal escapes to heap"
}

func foo91(x *int) map[*int]*int { // ERROR "leaking param: NAME-x"
	return map[*int]*int{x: nil} // ERROR "map literal escapes to heap"
}

func foo92(x *int) [2]*int { // ERROR "leaking param: NAME-x"
	return [2]*int{x, nil}
}

// does not leak c
func foo93(c chan *int) *int { // ERROR "c does not escape"
	for v := range c {
		return v
	}
	return nil
}

// does not leak m
func foo94(m map[*int]*int, b bool) *int { // ERROR "m does not escape"
	for k, v := range m {
		if b {
			return k
		}
		return v
	}
	return nil
}

// does leak x
func foo95(m map[*int]*int, x *int) { // ERROR "m does not escape" "leaking param: NAME-x"
	m[x] = x
}

// does not leak m
func foo96(m []*int) *int { // ERROR "m does not escape"
	return m[0]
}

// does leak m
func foo97(m [1]*int) *int { // ERROR "leaking param: NAME-m"
	return m[0]
}

// does not leak m
func foo98(m map[int]*int) *int { // ERROR "m does not escape"
	return m[0]
}

// does leak m
func foo99(m *[1]*int) []*int { // ERROR "leaking param: NAME-m"
	return m[:]
}

// does not leak m
func foo100(m []*int) *int { // ERROR "m does not escape"
	for _, v := range m {
		return v
	}
	return nil
}

// does leak m
func foo101(m [1]*int) *int { // ERROR "leaking param: NAME-m"
	for _, v := range m {
		return v
	}
	return nil
}

// does not leak m
func foo101a(m [1]*int) *int { // ERROR "m does not escape"
	for i := range m { // ERROR "moved to heap: NAME-i"
		return &i  // ERROR "&i escapes to heap"
	}
	return nil
}

// does leak x
func foo102(m []*int, x *int) { // ERROR "m does not escape" "leaking param: NAME-x"
	m[0] = x
}

// does not leak x
func foo103(m [1]*int, x *int) { // ERROR "m does not escape" "x does not escape"
	m[0] = x
}

var y []*int

// does not leak x
func foo104(x []*int) {  // ERROR "x does not escape"
	copy(y, x)
}

// does not leak x
func foo105(x []*int) {  // ERROR "x does not escape"
	_ = append(y, x...)
}

// does leak x
func foo106(x *int) { // ERROR "leaking param: NAME-x"
	_ = append(y, x)
}

func foo107(x *int) map[*int]*int { // ERROR "leaking param: NAME-x"
	return map[*int]*int{x: nil} // ERROR "map literal escapes to heap"
}

func foo108(x *int) map[*int]*int { // ERROR "leaking param: NAME-x"
	return map[*int]*int{nil: x} // ERROR "map literal escapes to heap"
}

func foo109(x *int) *int { // ERROR "leaking param: NAME-x"
	m := map[*int]*int{x: nil}  // ERROR "map literal does not escape"
	for k, _ := range m {
		return k
	}
	return nil
}

func foo110(x *int) *int { // ERROR "leaking param: NAME-x"
	m := map[*int]*int{nil: x}  // ERROR "map literal does not escape"
	return m[nil]
}

func foo111(x *int) *int { // ERROR "leaking param: NAME-x"
	m := []*int{x}  // ERROR "slice literal does not escape"
	return m[0]
}

func foo112(x *int) *int { // ERROR "leaking param: NAME-x"
	m := [1]*int{x}
	return m[0]
}

func foo113(x *int) *int { // ERROR "leaking param: NAME-x"
	m := Bar{ii: x}
	return m.ii
}

func foo114(x *int) *int { // ERROR "leaking param: NAME-x"
	m := &Bar{ii: x}  // ERROR "&struct literal does not escape"
	return m.ii
}

func foo115(x *int) *int { // ERROR "leaking param: NAME-x"
	return (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(x)) + 1))
}

func foo116(b bool) *int {
	if b {
		x := 1  // ERROR "moved to heap: NAME-x"
		return &x  // ERROR "&x escapes to heap"
	} else {
		y := 1  // ERROR "moved to heap: NAME-y"
		return &y  // ERROR "&y escapes to heap"
	}
	return nil
}

func foo117(unknown func(interface{})) {  // ERROR "unknown does not escape"
	x := 1 // ERROR "moved to heap: NAME-x"
	unknown(&x) // ERROR "&x escapes to heap"
}

func foo118(unknown func(*int)) {  // ERROR "unknown does not escape"
	x := 1 // ERROR "moved to heap: NAME-x"
	unknown(&x) // ERROR "&x escapes to heap"
}

func external(*int)

func foo119(x *int) {  // ERROR "leaking param: NAME-x"
	external(x)
}