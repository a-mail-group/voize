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
)


func EncodeRice(dst bitio.Writer, k, val uint) (err error) {
	k &= 255
	u := uint64(val)&((uint64(1)<<k)-1)
	n := val>>k
	for i := uint(0); i<n; i++ {
		err = dst.WriteBool(true)
		if err!=nil { return }
	}
	err = dst.WriteBool(false)
	if err!=nil { return }
	err = dst.WriteBits(u,byte(k))
	return
}
func decode0(src bitio.Reader) (val uint,err error) {
	var b bool
	for {
		b,err = src.ReadBool()
		if err!=nil { return }
		if !b { return }
		val++
	}
}
func DecodeRice(src bitio.Reader, k uint) (val uint,err error) {
	k &= 255
	val,err = decode0(src)
	if err!=nil { return }
	if k==0 { return }
	val <<= k
	var r uint64
	r,err = src.ReadBits(byte(k))
	val|=uint(r)
	return
}



