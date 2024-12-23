package src

import (
	"io"
)

type EncoderArgs struct {
	Width, Height int
	Passes, pass  int
	Bitdepth      int
	PixFmt        PixelFormat
	Fps           int
	ExtendedArgs  []string
}

type PixelFormat int

const (
	YUV400 = iota
	YUV420
	YUV422
	YUV444
)

type Encoder func(io.Reader, io.WriteCloser, EncoderArgs) error

type internalEncoder int

const (
	x265Cli = iota
	x264Cli
	vpxenc
)
