/*
This code is from
https://rosettacode.org/wiki/Fast_Fourier_transform#Go

There is no license for this code.
*/


package dct4

import (
    "math"
    "math/cmplx"
)

/*
The reference implementation from https://rosettacode.org/wiki/Fast_Fourier_transform#Go
Modification: we use complex input!
*/
func ditfft2c(x []complex128, y []complex128, n, s int) {
	if n == 1 {
		y[0] = x[0]
		return
	}
	ditfft2c(x, y, n/2, 2*s)
	ditfft2c(x[s:], y[n/2:], n/2, 2*s)
	for k := 0; k < n/2; k++ {
		tf := cmplx.Rect(1, -2*math.Pi*float64(k)/float64(n)) * y[k+n/2]
		y[k], y[k+n/2] = y[k]+tf, y[k]-tf
	}
}

