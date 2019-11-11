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


package fft

import "math"

/*
A Discrete Cosine transform.
*/
type DCT4 struct {
	buffer []complex128
}

func (d *DCT4) Encode(x []float64, y []float64) {
	if len(x)!=len(y) { panic("assertion failed: len(x)==len(y)") }
	if len(d.buffer)<len(x) { d.buffer = make([]complex128,len(x)) }
	ditfft2(x,d.buffer,len(x),1)
	for i := range y { y[i] = real(d.buffer[i]) } // real = cosine
}
func (d *DCT4) Decode(x []float64, y []float64) {
	if len(x)!=len(y) { panic("assertion failed: len(x)==len(y)") }
	if len(d.buffer)<len(x) { d.buffer = make([]complex128,len(x)) }
	ditfft2ic(x,d.buffer,len(x),1,float64(len(x)))
	for i := range y { y[i] = real(d.buffer[i]) } // orig. input = real
}

func f64zero(f []float64) {
	for i := range f { f[i]=0 }
}

/*
A Modified Discrete Cosine Transform.

An instance must not be used for encoding and decoding Simultaneously!
*/
type MDCT struct {
	DCT DCT4
	Half, Full int
	Prev []float64
	Inpt []float64
	Temp []float64
	Outp []float64
}
func (m *MDCT) Init(bz int) {
	m.Full = bz
	m.Half = bz/2
	m.Prev = make([]float64,m.Half)
	m.Inpt = make([]float64,bz)
	m.Temp = make([]float64,bz)
	m.Outp = make([]float64,bz)
}
func (m *MDCT) EncodeHalfBlock(src []float64) []float64 {
	if len(src)!=m.Half { panic("assertion failed: len(src)==m.Half") }
	copy(m.Inpt[:m.Half],m.Inpt[m.Half:])
	copy(m.Inpt[m.Half:],src)
	for i := 0; i<m.Full; i++ {
		f := math.Cos(math.Pi*float64(i)/float64(m.Half))
		m.Temp[i] = (m.Inpt[i]*(1-f))
	}
	m.DCT.Encode(m.Temp,m.Outp)
	for i,coeff := range m.Outp {
		m.Outp[i] = math.Abs(coeff)
	}
	return m.Outp
}
func (m *MDCT) DecodeHalfBlock(src, dst []float64) {
	if len(src)!=m.Full { panic("assertion failed: len(src)==m.Full") }
	if len(dst)!=m.Half { panic("assertion failed: len(dst)==m.Half") }
	m.DCT.Decode(src,m.Inpt)
	
	copy(m.Temp[:m.Half],m.Temp[m.Half:])
	f64zero(m.Temp[m.Half:])
	
	const F = 0.5
	for i := 0; i<m.Full; i++ {
		f := math.Cos(math.Pi*float64(i)/float64(m.Half))
		m.Temp[i] += (m.Inpt[i]*(1-f)*F)
	}
	copy(dst,m.Temp)
}


