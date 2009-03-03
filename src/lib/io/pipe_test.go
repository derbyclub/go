// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"io";
	"os";
	"testing";
	"time";
)

func checkWrite(t *testing.T, w io.Write, data []byte, c chan int) {
	n, err := w.Write(data);
	if err != nil {
		t.Errorf("write: %v", err);
	}
	if n != len(data) {
		t.Errorf("short write: %d != %d", n, len(data));
	}
	c <- 0;
}

// Test a single read/write pair.
func TestPipe1(t *testing.T) {
	c := make(chan int);
	r, w := io.Pipe();
	var buf [64]byte;
	go checkWrite(t, w, io.StringBytes("hello, world"), c);
	n, err := r.Read(buf);
	if err != nil {
		t.Errorf("read: %v", err);
	}
	else if n != 12 || string(buf[0:12]) != "hello, world" {
		t.Errorf("bad read: got %q", buf[0:n]);
	}
	<-c;
	r.Close();
	w.Close();
}

func reader(t *testing.T, r io.Read, c chan int) {
	var buf [64]byte;
	for {
		n, err := r.Read(buf);
		if err != nil {
			t.Errorf("read: %v", err);
		}
		c <- n;
		if n == 0 {
			break;
		}
	}
}

// Test a sequence of read/write pairs.
func TestPipe2(t *testing.T) {
	c := make(chan int);
	r, w := io.Pipe();
	go reader(t, r, c);
	var buf [64]byte;
	for i := 0; i < 5; i++ {
		p := buf[0:5+i*10];
		n, err := w.Write(p);
		if n != len(p) {
			t.Errorf("wrote %d, got %d", len(p), n);
		}
		if err != nil {
			t.Errorf("write: %v", err);
		}
		nn := <-c;
		if nn != n {
			t.Errorf("wrote %d, read got %d", n, nn);
		}
	}
	w.Close();
	nn := <-c;
	if nn != 0 {
		t.Errorf("final read got %d", nn);
	}
}

// Test a large write that requires multiple reads to satisfy.
func writer(w io.WriteClose, buf []byte, c chan pipeReturn) {
	n, err := w.Write(buf);
	w.Close();
	c <- pipeReturn{n, err};
}

func TestPipe3(t *testing.T) {
	c := make(chan pipeReturn);
	r, w := io.Pipe();
	var wdat [128]byte;
	for i := 0; i < len(wdat); i++ {
		wdat[i] = byte(i);
	}
	go writer(w, wdat, c);
	var rdat [1024]byte;
	tot := 0;
	for n := 1; n <= 256; n *= 2 {
		nn, err := r.Read(rdat[tot:tot+n]);
		if err != nil {
			t.Fatalf("read: %v", err);
		}

		// only final two reads should be short - 1 byte, then 0
		expect := n;
		if n == 128 {
			expect = 1;
		} else if n == 256 {
			expect = 0;
		}
		if nn != expect {
			t.Fatalf("read %d, expected %d, got %d", n, expect, nn);
		}
		tot += nn;
	}
	pr := <-c;
	if pr.n != 128 || pr.err != nil {
		t.Fatalf("write 128: %d, %v", pr.n, pr.err);
	}
	if tot != 128 {
		t.Fatalf("total read %d != 128", tot);
	}
	for i := 0; i < 128; i++ {
		if rdat[i] != byte(i) {
			t.Fatalf("rdat[%d] = %d", i, rdat[i]);
		}
	}
}

// Test read after/before writer close.

func delayClose(t *testing.T, cl io.Close, ch chan int) {
	time.Sleep(1000*1000);	// 1 ms
	if err := cl.Close(); err != nil {
		t.Errorf("delayClose: %v", err);
	}
	ch <- 0;
}

func testPipeReadClose(t *testing.T, async bool) {
	c := make(chan int, 1);
	r, w := io.Pipe();
	if async {
		go delayClose(t, w, c);
	} else {
		delayClose(t, w, c);
	}
	var buf [64]int;
	n, err := r.Read(buf);
	<-c;
	if err != nil {
		t.Errorf("read from closed pipe: %v", err);
	}
	if n != 0 {
		t.Errorf("read on closed pipe returned %d", n);
	}
	if err = r.Close(); err != nil {
		t.Errorf("r.Close: %v", err);
	}
}

// Test write after/before reader close.

func testPipeWriteClose(t *testing.T, async bool) {
	c := make(chan int, 1);
	r, w := io.Pipe();
	if async {
		go delayClose(t, r, c);
	} else {
		delayClose(t, r, c);
	}
	n, err := io.WriteString(w, "hello, world");
	<-c;
	if err != os.EPIPE {
		t.Errorf("write on closed pipe: %v", err);
	}
	if n != 0 {
		t.Errorf("write on closed pipe returned %d", n);
	}
	if err = w.Close(); err != nil {
		t.Errorf("w.Close: %v", err);
	}
}

func TestPipeReadCloseAsync(t *testing.T) {
	testPipeReadClose(t, true);
}

func TestPipeReadCloseSync(t *testing.T) {
	testPipeReadClose(t, false);
}

func TestPipeWriteCloseAsync(t *testing.T) {
	testPipeWriteClose(t, true);
}

func TestPipeWriteCloseSync(t *testing.T) {
	testPipeWriteClose(t, false);
}

