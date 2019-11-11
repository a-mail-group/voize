/*
This code is from
https://rosettacode.org/wiki/Fast_Fourier_transform#Go

There is no license for this code.
*/


package fft

import (
    "math"
    "math/cmplx"
)

/*
The reference implementation from https://rosettacode.org/wiki/Fast_Fourier_transform#Go
*/
func ditfft2(x []float64, y []complex128, n, s int) {
	if n == 1 {
		y[0] = complex(x[0], 0)
		return
	}
	ditfft2(x, y, n/2, 2*s)
	ditfft2(x[s:], y[n/2:], n/2, 2*s)
	for k := 0; k < n/2; k++ {
		tf := cmplx.Rect(1, -2*math.Pi*float64(k)/float64(n)) * y[k+n/2]
		y[k], y[k+n/2] = y[k]+tf, y[k]-tf
	}
}



func Raw_FFT(x []float64, y []complex128) {
	if len(x)!=len(y) { panic("assertion failed: len(x)==len(y)") }
	ditfft2(x,y,len(x),1)
}

func cdiv(a complex128,b float64) complex128 {
	return complex(real(a)/b,imag(a)/b)
}
/*
The inverse FFT variation from the RosettaCode implementation.
*/
func ditfft2i(x []complex128, y []complex128, n, s int,div float64) {
	if n == 1 {
		y[0] = cdiv(x[0],div)
		return
	}
	ditfft2i(x, y, n/2, 2*s, div)
	ditfft2i(x[s:], y[n/2:], n/2, 2*s, div)
	for k := 0; k < n/2; k++ {
		tf := cmplx.Rect(1, 2*math.Pi*float64(k)/float64(n)) * y[k+n/2]
		y[k], y[k+n/2] = y[k]+tf, y[k]-tf
	}
}
func Raw_iFFT(x []complex128, y []complex128) {
	if len(x)!=len(y) { panic("assertion failed: len(x)==len(y)") }
	ditfft2i(x,y,len(x),1,float64(len(x)))
}

/*
The inverse FFT variation from the RosettaCode implementation for inverse discrete cosine transform (iDCT).
*/
func ditfft2ic(x []float64, y []complex128, n, s int,div float64) {
	if n == 1 {
		y[0] = complex(x[0]/div,0)
		return
	}
	ditfft2ic(x, y, n/2, 2*s, div)
	ditfft2ic(x[s:], y[n/2:], n/2, 2*s, div)
	for k := 0; k < n/2; k++ {
		tf := cmplx.Rect(1, 2*math.Pi*float64(k)/float64(n)) * y[k+n/2]
		y[k], y[k+n/2] = y[k]+tf, y[k]-tf
	}
}

