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


package dct4


// The whole code is inspired from https://docs.rs/crate/rustdct/0.1.3

import "math"

func mp3(res []float64) {
	n := len(res)
	ct := math.Pi/float64(n)
	
	for i := range res {
		
		res[i] = math.Sin(ct * (float64(n) + 0.5))
	}
}
func vorbis(res []float64) {
	n := len(res)
	ct := math.Pi/float64(n)
	
	for i := range res {
		inner_sin := math.Sin(ct * (float64(n) + 0.5))
		res[i] = math.Sin(math.Pi * 0.5 * inner_sin * inner_sin)
	}
}

type WindowFn func(res []float64)
var WF_MP3 WindowFn = mp3
var WF_Vorbis WindowFn = vorbis


func mult(I1, I2, O []float64) {
	for i,e := range I1 {
		O[i] = e*I2[i]
	}
}

// TODO: there is no inverse that seems to work


type MDCTViaDCT4 struct {
	dct DCT4
	dctbuf []float64
	window []float64
	tempbf []float64
	bsize  int
}
func (d *MDCTViaDCT4) Init(d4 DCT4Constructor, wf WindowFn, bz int) *MDCTViaDCT4 {
	d.dctbuf = make([]float64,bz)
	d.window = make([]float64,bz*2)
	wf(d.window)
	d.dct = d4.Init(bz)
	d.bsize = bz
	d.tempbf = make([]float64,bz)
	return d
}
func (d *MDCTViaDCT4) ProcessSplit(input_a,input_b,output []float64) {
	halfz  := d.bsize/2
	halfz1 := halfz-1
	bsize1 := d.bsize-1
	
	//we're going to divide input_a into two subgroups, (a,b), and input_b into two subgroups: (c,d)
        //then scale them by the window function, then combine them into two subgroups: (-D-Cr, A-Br) where R means reversed
	
	//the first half of the dct input is -Cr - D
	mult(input_b,d.window[d.bsize:],d.tempbf)
	hb := d.dctbuf[:halfz]
	for i := range hb {
		Cr := d.tempbf[halfz1-i]
		D  := d.tempbf[halfz +i]
		hb[i] = -Cr - D
	}
	
	//the second half of the dct input is A - Br
	mult(input_a,d.window[:d.bsize],d.tempbf)
	hb = d.dctbuf[halfz:]
	for i := range hb {
		A  := d.tempbf[       i]
		Br := d.tempbf[bsize1-i]
		hb[i] = A - Br
	}
	d.dct.Process(d.dctbuf,output)
}

