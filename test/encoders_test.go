package test

import (
	"os"
	"testing"

	"github.com/GreatValueCreamSoda/video-encoders/src"
)

// 854x4080
var testingYUVData = "./samples/sample.yuv"

var testingEncoderArgs = src.EncoderArgs{
	Width:        854,
	Height:       480,
	Passes:       1,
	Bitdepth:     10,
	PixFmt:       src.YUV420,
	Fps:          30,
	ExtendedArgs: make([]string, 0),
}

func init() {

	if _, err := os.Stat("./encodes"); os.IsNotExist(err) {
		err := os.Mkdir("./encodes", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func TestX265Cli(t *testing.T) {
	r, err := os.Open(testingYUVData)
	if err != nil {
		t.Error(err)
	}

	w, err := os.Create("./encodes/encode_x265.ivf")
	if err != nil {
		t.Error(err)
	}

	err = src.X265Cli(r, w, testingEncoderArgs)
	if err != nil {
		t.Error(err)
	}
}

func TestX264Cli(t *testing.T) {
	r, err := os.Open(testingYUVData)
	if err != nil {
		t.Error(err)
	}

	w, err := os.Create("./encodes/encode_x264.ivf")
	if err != nil {
		t.Error(err)
	}

	err = src.X264Cli(r, w, testingEncoderArgs)
	if err != nil {
		t.Error(err)
	}
}

func TestVP9CLI(t *testing.T) {
	r, err := os.Open(testingYUVData)
	if err != nil {
		t.Error(err)
	}

	w, err := os.Create("./encodes/encode_vp9.ivf")
	if err != nil {
		t.Error(err)
	}

	err = src.VP9Cli(r, w, testingEncoderArgs)
	if err != nil {
		t.Error(err)
	}
}
