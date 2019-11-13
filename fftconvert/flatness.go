/*
This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <https://unlicense.org>
*/


package fftconvert

import (
	"math"
)

type Falloff struct{
	falloff []float64
}
func (f *Falloff) Init(bz int) {
	const fall = 0.2
	f.falloff = make([]float64,bz)
	for i := range f.falloff {
		x := (math.Log(float64(i+1))-math.Log(float64(i)))/math.Log(2)
		f.falloff[i] = math.Pow(fall,x)
	}
	f.falloff[0] = fall*fall
	
}
func (f *Falloff) Clear(x []float64) {
	n := len(x)
	last := x[0]
	for i:=1 ; i<n; i++ {
		cur := x[i]
		if x[i]<(last*f.falloff[i]) { x[i] = 0 }
		last = cur
	}
	last = x[n-1]
	for i:=n-2; i>=0 ; i-- {
		cur := x[i]
		if x[i]<(last*f.falloff[i]) { x[i] = 0 }
		last = cur
	}
}


func Flatten (x []float64) {
	const f = 0.5
	n := len(x)
	last := x[0]
	for i:=1 ; i<n; i++ {
		cur := x[i]
		if x[i]<last*f { x[i] = 0 }
		last = cur
	}
	last = x[n-1]
	for i:=n-2; i>=0 ; i-- {
		cur := x[i]
		if x[i]<last*f { x[i] = 0 }
		last = cur
		//x[i] = math.Max(x[i],x[i+1]*f)
	}
}
