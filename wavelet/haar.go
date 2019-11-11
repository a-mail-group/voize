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


type HaarWavelet struct {
	buf1 []float64
	buf2 []float64
}
func (d4 *HaarWavelet) step(n int) {
	half := n >> 1
	
	a := d4.buf1
	tmp := d4.buf2
	
	for i := 0; i<half; i++ {
		j := i<<1
		tmp[i     ] = a[j]+a[j+1]
		tmp[i+half] = a[j]-a[j+1]
	}
	
	copy(a[:n],tmp[:n])
}
func (d4 *HaarWavelet) Forward(src, dst []float64) {
	d4.buf1 = append(d4.buf1[:0],src...)
	d4.buf2 = append(d4.buf2[:0],src...)
	for n := len(src); n>=1; n>>=1 {	
		d4.step(n)
	}
	copy(dst,d4.buf1)
}
func (d4 *HaarWavelet) istep(n int) {
	half := n >> 1
	
	a := d4.buf1
	tmp := d4.buf2
	
	for i := 0; i<half; i++ {
		j := i<<1
		tmp[j  ] = (a[i]+a[i+half])/2.0
		tmp[j+1] = (a[i]-a[i+half])/2.0
	}
	copy(a[:n],tmp[:n])
}
func (d4 *HaarWavelet) Reverse(src, dst []float64) {
	d4.buf1 = append(d4.buf1[:0],src...)
	d4.buf2 = append(d4.buf2[:0],src...)
	N := len(src)
	for n := 2; n<=N; n<<=1 {
		d4.istep(n)
	}
	copy(dst,d4.buf1)
}
