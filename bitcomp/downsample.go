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


package bitcomp

import (
	"github.com/icza/bitio"
	"math"
	//"fmt"
)

func Delta2Bits(f float64) int {
	i := int(math.Sqrt(f))
	return i
}
func Bits2Delta(i int) float64 {
	f := float64(i)
	return f*f
}

func iabs(i int) uint {
	if i<0 { return uint(-i) }
	return uint(i)
}

type Sampler struct{
	buf []float64
	rbuf []int
	
	Last float64
}
func (s *Sampler) Init(bz int) {
	s.buf = make([]float64,bz/4)
	s.rbuf = make([]int,bz/4)
}
func (s *Sampler) pack(dst bitio.Writer) (err error) {
	k,sk := uint(0),uint(0)
	sbc := uint(0)
	for ; k<4 ; k++ {
		bc := (k+2)*uint(len(s.rbuf))
		for _,r := range s.rbuf {
			bc += iabs(r)>>k
		}
		if bc<sbc || k==0 {
			sbc,sk = bc,k
		}
	}
	err = dst.WriteBits(uint64(sk),2)
	if err!=nil { return }
	//fmt.Println("bit-count/8 = ",float64(sbc)/8)
	for _,r := range s.rbuf {
		err = dst.WriteBool(r<0)
		if err!=nil { return }
		err = EncodeRice(dst,sk,iabs(r))
		if err!=nil { return }
	}
	return
}
func (s *Sampler) upsample(dec []float64) {
	last := s.Last
	for i := range s.buf {
		j := i*4
		d := s.buf[i]-last
		for k := 0; k<4 ; k++ {
			dec[j+k] = last + d*(float64(k+1)/4)
		}
		last = s.buf[i]
	}
	last2 := s.Last
	for i,v := range dec[:(len(s.buf)*4)-1] {
		interp := (last2+dec[i+1])/2
		last2 = (v+interp)/2
		dec[i] = last2
	}
	s.Last = last
}
func (s *Sampler) Compress(raw, dec []float64,dst bitio.Writer) (err error) {
	for i := range s.buf {
		j := i*4
		s.buf[i] = (raw[j]+raw[j+1]+raw[j+2]+raw[j+3])/4
	}
	for i,f := range s.buf {
		f = math.Min(f,1)
		f = math.Max(f,-1)
		s.rbuf[i] = int(f*127)
		s.buf[i] = float64(s.rbuf[i])/127
	}
	//fmt.Printf("%d\n",s.rbuf)
	err = s.pack(dst)
	if err!=nil { return }
	s.upsample(dec)
	return
}

//func (s *Sampler) 
//func (s *Sampler) 


