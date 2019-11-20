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
import "math/cmplx"
import "fmt"

type DCT4 interface{
	Process(in, out []float64)
}
type DCT4Constructor interface{
	Init(bz int) DCT4
}

type NaiveDCT4 struct {
	twiddles []float64
}
func (n *NaiveDCT4) Init(bz int) DCT4 {
	cf := (.5*math.Pi)/float64(bz)
	twiddles := make([]float64,bz*4)
	for i := range twiddles {
		twiddles[i] = cf * math.Cos(float64(i)+.5)
	}
	n.twiddles = twiddles
	return n
}
func (n *NaiveDCT4) Process(in, out []float64) {
	//if len(in) != (len(out)*2) { panic("DCT4 len(in) != (len(out)*2)") }
	twiddles := n.twiddles
	tl := len(twiddles)
	if (len(in)*4) != tl { panic("len(in) != bz") }
	for k := range out {
		var o float64 = 0
		ti := k
		ts := (k*2)+1
		for _,e := range in {
			o += e * twiddles[ti]
			ti = (ti+ts)%tl
		}
		out[k] = o
	}
}

func single_twiddle(i, fft_len int, inverse bool) complex128 {
	constant := 2*math.Pi
	if !inverse { constant = -constant }
	return cmplx.Rect(1,constant*float64(i)/float64(fft_len))
}

type DCT4ViaFFT struct{
	fft_input  []complex128
	fft_output []complex128
	twiddles   []complex128
	blocksize  int
}
func (n *DCT4ViaFFT) Init(bz int) DCT4 {
	if (bz%4)!=0 { panic("(len(bz)%4)!=0") }
	inner_bz := bz*4
	
	twiddle_scale := 0.25 / math.Cos(math.Pi/float64(inner_bz))
	
	twiddles := make([]complex128,bz)
	for i := range twiddles {
		twiddles[i] = single_twiddle(i,inner_bz,false) * complex(twiddle_scale,0)
	}
	
	n.twiddles   = twiddles
	n.fft_input  = make([]complex128,inner_bz)
	n.fft_output = make([]complex128,inner_bz)
	n.blocksize  = bz
	return n
}
func (n *DCT4ViaFFT) Process(in, out []float64) {
	//if len(in) != (len(out)*2) { panic("DCT4 len(in) != (len(out)*2)") }
	
	if len(in) != n.blocksize {
		//panic("len(in) != bz 1")
		panic(fmt.Sprintf("len(in) (%d) != bz (%d)",len(in),n.blocksize))
	}
	
	{
		O2 := n.blocksize*2
		O1 := O2-1
		O3 := len(n.fft_input)-1
		for i,e := range in {
			n.fft_input[   i] = complex(e,0)
			n.fft_input[O1-i] = complex(e,0) // fft_input.skip(signal.len()  ) & signal.iter().rev()
			n.fft_input[O2+i] = complex(e,0) // fft_input.skip(signal.len()*2)
			n.fft_input[O3-i] = complex(e,0) // fft_input.skip(signal.len()*3) & signal.iter().rev()
		}
	}
	ditfft2c(n.fft_input,n.fft_output,len(n.fft_input),1)
	twiddles := n.twiddles[:len(out)]
	for index,twiddle := range twiddles {
		// e := real(n.fft_output[(index*2)+1] * twiddle)
		e := real(n.fft_output[(index*2)] * twiddle)
		out[index] = e
	}
}

