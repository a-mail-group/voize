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



package segmenter

import (
	"github.com/go-audio/audio"
	"math"
)

func cloneAF(f *audio.Format) *audio.Format {
	n := new(audio.Format)
	*n = *f
	return n
}

type StreamDecoder struct {
	BlockSize int
	
	srcpos int
	srcbuf *audio.IntBuffer
	
	dstpos int
	
	audio.Format
	Buffer [][]float64
}
func (s *StreamDecoder) reshape() {
	if s.NumChannels==0 {
		s.Format = *(s.srcbuf.Format)
		s.Buffer = make([][]float64,s.NumChannels)
		for i := range s.Buffer {
			s.Buffer[i] = make([]float64,s.BlockSize)
		}
	}
}
func (s *StreamDecoder) pull() {
	sp := s.srcpos
	dp := s.dstpos
	src := s.srcbuf.Data
	se := len(src)
	de := s.BlockSize
	if (sp>=se) || (dp>=de) { return }
	flt := math.Pow(2,float64(s.srcbuf.SourceBitDepth))-1.0
	for (sp<se) && (dp<de) {
		for j,b := range s.Buffer {
			b[dp] = float64(src[sp+j])/flt
		}
		dp++
		sp+=s.NumChannels
	}
	s.srcpos = sp
	s.dstpos = dp
}
func (s *StreamDecoder) Fill() { s.pull() }
func (s *StreamDecoder) Next() bool {
	if s.dstpos<s.BlockSize { return false }
	s.dstpos = 0
	return true
}

func (s *StreamDecoder) Decode(src *audio.IntBuffer) {
	s.srcbuf = src
	s.srcpos = 0
	s.reshape()
}

type StreamEncoder struct {
	BlockSize int
	audio.Format
	Buffer [][]float64
	
	SourceBitDepth int
	
	ibuf *audio.IntBuffer
}
func (s *StreamEncoder) Allocate() {
	s.Buffer = make([][]float64,s.NumChannels)
	for i := range s.Buffer {
		s.Buffer[i] = make([]float64,s.BlockSize)
	}
}
func (s *StreamEncoder) create() {
	if s.ibuf==nil {
		if s.NumChannels!=len(s.Buffer) { panic("assertion failed: s.NumChannels==len(s.Buffer)") }
		s.ibuf = &audio.IntBuffer{
			Format: cloneAF(&s.Format),
			Data: make([]int,s.BlockSize*s.NumChannels),
			SourceBitDepth: s.SourceBitDepth,
		}
	}
}
func (s *StreamEncoder) Encode() *audio.IntBuffer {
	s.create()
	dst := s.ibuf.Data
	flt := math.Pow(2,float64(s.SourceBitDepth))-1.0
	mflt := -flt
	for j,b := range s.Buffer {
		for i := 0;i<s.BlockSize; i++ {
			k := i*s.NumChannels
			dst[k+j] = int(math.Min(flt,math.Max(mflt,b[i]*flt)))
		}
	}
	return s.ibuf
}
//func (s *StreamEncoder) 
//func (s *StreamEncoder) 

