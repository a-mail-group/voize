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
	"github.com/a-mail-group/voize/fft"
	"math/cmplx"
)

func czero128(c []complex128) {
	for i := range c { c[i] = 0 }
}

type Codec struct {
	buffer []complex128
	mbuffer []complex128
	nbuffer []complex128
	Abs   []float64
	Phase []float64
}
func (e *Codec) Init(bz int) {
	e.buffer = make([]complex128,bz)
	e.mbuffer = make([]complex128,bz)
	e.nbuffer = make([]complex128,bz)
	e.Abs = make([]float64,bz/2)
	e.Phase = make([]float64,bz/2)
}
func (e *Codec) Encode(data []float64) {
	fft.Raw_FFT(data,e.buffer)
	for i := range e.Abs {
		e.Abs[i],e.Phase[i] = cmplx.Polar(e.buffer[i]*2)
	}
}
func (e *Codec) Decode(data []float64) {
	for i := range e.Abs {
		e.nbuffer[i] = cmplx.Rect(e.Abs[i],e.Phase[i])
	}
	fft.Raw_iFFT(e.nbuffer,e.buffer)
	for i,c := range e.buffer { data[i] = real(c) }
}
func (e *Codec) ScanAndDiff(data,diff []float64) {
	fft.Raw_FFT(data,e.buffer)
	fft.Raw_iFFT(e.buffer,e.nbuffer)
	for i,c := range e.nbuffer { diff[i] = real(c)-data[i] }
}
func (e *Codec) Extract(data,diff []float64,min,max int) {
	czero128(e.mbuffer)
	for i := min; i<max ; i++ {
		e.mbuffer[i] = e.buffer[i]
	}
	fft.Raw_iFFT(e.mbuffer,e.nbuffer)
	if len(diff)==len(data) {
		for i,c := range e.nbuffer { data[i] = real(c)+diff[i] }
	} else {
		for i,c := range e.nbuffer { data[i] = real(c) }
	}
}


