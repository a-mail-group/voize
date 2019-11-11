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


package wavelet

import "math"

/*
Transcribed from: http://bearcave.com/misl/misl_tech/wavelets/daubechies/daub.java
*/
var d4_sqrt3 = math.Sqrt(3)
var d4_denom = 4*math.Sqrt(3)

var d4_h0 = (1+d4_sqrt3)/d4_denom
var d4_h1 = (3+d4_sqrt3)/d4_denom
var d4_h2 = (3-d4_sqrt3)/d4_denom
var d4_h3 = (1-d4_sqrt3)/d4_denom

var d4_g0 =  d4_h3
var d4_g1 = -d4_h2
var d4_g2 =  d4_h1
var d4_g3 = -d4_h0

var d4_ih0 = d4_h2
var d4_ih1 = d4_g2 // h1
var d4_ih2 = d4_h0
var d4_ih3 = d4_g2 // h3

var d4_ig0 = d4_h3
var d4_ig1 = d4_g3 // -h0
var d4_ig2 = d4_h1
var d4_ig3 = d4_g1 // -h2

/*

Daubechies wavelet. Warning: I thing it has a bug. this Impl. is unusable!

Special thanks to: http://bearcave.com/software/java/wavelets/daubechies.html
*/
type D4Wavelet struct {
	buf1 []float64
	buf2 []float64
}
func (d4 *D4Wavelet) step(n int) {
	half := n >> 1
	
	a := d4.buf1
	tmp := d4.buf2
	
	i := 0
	for j := 0 ; j<n-3 ; j += 2 {
		tmp[i     ] = a[j]*d4_h0 + a[j+1]*d4_h1 + a[j+2]*d4_h2 + a[j+3]*d4_h3
		tmp[i+half] = a[j]*d4_g0 + a[j+1]*d4_g1 + a[j+2]*d4_g2 + a[j+3]*d4_g3
		i++
	}
	tmp[i     ] = a[n-2]*d4_h0 + a[n-1]*d4_h1 + a[0]*d4_h2 + a[1]*d4_h3
	tmp[i+half] = a[n-2]*d4_g0 + a[n-1]*d4_g1 + a[0]*d4_g2 + a[1]*d4_g3
	
	copy(a[:n],tmp[:n])
}
func (d4 *D4Wavelet) Forward(src, dst []float64) {
	d4.buf1 = append(d4.buf1[:0],src...)
	d4.buf2 = append(d4.buf2[:0],src...)
	for n := len(src); n>=4; n>>=1 {	
		d4.step(n)
	}
	copy(dst,d4.buf1)
}
func (d4 *D4Wavelet) istep(n int) {
	half := n >> 1
	hpl1 := half+1
	hmn1 := half-1
	
	a := d4.buf1
	tmp := d4.buf2
	
	tmp[0] = a[hmn1]*d4_ih0 + a[n-1]*d4_ih1 + a[0]*d4_ih2 + a[half]*d4_ih3
	tmp[1] = a[hmn1]*d4_ig0 + a[n-1]*d4_ig1 + a[0]*d4_ig2 + a[half]*d4_ig3
	
	j := 2
	for i := 0 ; i < hmn1 ; i++ {
		tmp[j] = a[i]*d4_ih0 + a[i+half]*d4_ih1 + a[i+1]*d4_ih2 + a[i+hpl1]*d4_ih3 ; j++
		tmp[j] = a[i]*d4_ig0 + a[i+half]*d4_ig1 + a[i+1]*d4_ig2 + a[i+hpl1]*d4_ig3 ; j++
	}
	
	copy(a[:n],tmp[:n])
}
func (d4 *D4Wavelet) Reverse(src, dst []float64) {
	d4.buf1 = append(d4.buf1[:0],src...)
	d4.buf2 = append(d4.buf2[:0],src...)
	N := len(src)
	for n := 4; n<=N; n<<=1 {
		d4.istep(n)
	}
	copy(dst,d4.buf1)
}
