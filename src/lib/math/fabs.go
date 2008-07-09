// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

export	fabs

func
fabs(arg float64) float64
{

	if arg < 0 {
		return -arg;
	}
	return arg;
}
