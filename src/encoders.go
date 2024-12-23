package src

import (
	"io"
	"os/exec"
	"strconv"
	"strings"
)

func X264Cli(r io.Reader, w io.WriteCloser, args EncoderArgs) error {
	return runEncoder(x264Cli, r, w, &args)
}

func X265Cli(r io.Reader, w io.WriteCloser, args EncoderArgs) error {
	return runEncoder(x265Cli, r, w, &args)
}

func VP9Cli(r io.Reader, w io.WriteCloser, args EncoderArgs) error {
	return runEncoder(vpxenc, r, w, &args)
}

func runEncoder(enc internalEncoder, r io.Reader, w io.WriteCloser, args *EncoderArgs) error {
	encoder := map[internalEncoder]string{
		x265Cli: "x265", x264Cli: "x264", vpxenc: "vpxenc",
	}[enc]
	baseArgs := getBaseArgs(enc, args)
	for args.pass = 1; args.pass <= args.Passes; args.pass++ {
		passArgs := getPassArgs(enc, args, baseArgs)
		err := runCommand(encoder, passArgs, r, w)
		if err != nil {
			return err
		}
	}
	return nil
}

func runCommand(name string, args []string, r io.Reader, w io.WriteCloser) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = r
	cmd.Stdout = w
	//cmd.Stderr = os.Stderr // for debugging
	//fmt.Println(cmd.Args)
	return cmd.Run()
}

func getBaseArgs(enc internalEncoder, args *EncoderArgs) []string {
	var base []string
	switch enc {
	case x265Cli, x264Cli:
		base = append([]string{"--input", "-", "-o", "-"}, args.ExtendedArgs...)
	case vpxenc:
		base = append([]string{"-", "-o", "-"}, args.ExtendedArgs...)
	default:
		panic(1)
	}

	getInputRes(enc, args, &base)
	getBitDepth(enc, args, &base)
	getPixFmt(enc, args, &base)
	getFpsArgs(enc, args, &base)

	return base
}

func getInputRes(enc internalEncoder, args *EncoderArgs, parsedArgs *[]string) {
	switch enc {
	case x265Cli, x264Cli:
		var builder strings.Builder
		builder.WriteString(strconv.Itoa(args.Width))
		builder.WriteString("x")
		builder.WriteString(strconv.Itoa(args.Height))
		*parsedArgs = append(*parsedArgs, "--input-res", builder.String())
	case vpxenc:
		*parsedArgs = append(*parsedArgs, "--width="+strconv.Itoa(args.Width))
		*parsedArgs = append(*parsedArgs, "--height="+strconv.Itoa(args.Height))
	}
}

func getBitDepth(enc internalEncoder, args *EncoderArgs, parsedArgs *[]string) {
	switch enc {
	case x265Cli, x264Cli:
		*parsedArgs = append(*parsedArgs, "--input-depth", strconv.Itoa(args.Bitdepth))
	case vpxenc:
		*parsedArgs = append(*parsedArgs, "--bit-depth="+strconv.Itoa(args.Bitdepth))
		*parsedArgs = append(*parsedArgs, "--input-bit-depth="+strconv.Itoa(args.Bitdepth))
		var profile = 0
		if args.Bitdepth > 8 {
			profile = 2
		}
		if args.PixFmt != YUV420 {
			profile++
		}
		*parsedArgs = append(*parsedArgs, "--profile="+strconv.Itoa(profile))
	}
}

func getPixFmt(enc internalEncoder, args *EncoderArgs, parsedArgs *[]string) {
	vpxLookup := map[PixelFormat]string{
		YUV420: "--i420", YUV422: "--i422", YUV444: "--i444",
	}

	switch enc {
	case x265Cli, x264Cli:
		if !(args.PixFmt >= YUV400 && args.PixFmt <= YUV444) {
			return
		}
		*parsedArgs = append(*parsedArgs, "--input-csp", strconv.Itoa(int(args.PixFmt)))
	case vpxenc:
		if _, ok := vpxLookup[args.PixFmt]; !ok {
			return
		}
		*parsedArgs = append(*parsedArgs, vpxLookup[args.PixFmt])
	}
}

func getFpsArgs(enc internalEncoder, args *EncoderArgs, parsedArgs *[]string) {
	switch enc {
	case x265Cli, x264Cli:
		*parsedArgs = append(*parsedArgs, "--fps", strconv.Itoa(args.Fps))
	case vpxenc:
		*parsedArgs = append(*parsedArgs, "--fps="+strconv.Itoa(args.Fps)+"/1")
	}
}

func getPassArgs(enc internalEncoder, args *EncoderArgs, parsedArgs []string) []string {
	var passArgs = append([]string{}, parsedArgs...)

	switch enc {
	case x265Cli, x264Cli:
		if !(args.Passes > 1 && args.pass >= 1 && args.pass <= args.Passes) {
			break
		}
		passArgs = append(passArgs, "--pass", strconv.Itoa(args.pass))
		passArgs = append(passArgs, "--passes", strconv.Itoa(args.Passes))
	case vpxenc:
		passArgs = append(passArgs, "--passes="+strconv.Itoa(args.Passes))
		if !(args.Passes > 1 && args.pass >= 1 && args.pass <= args.Passes) {
			break
		}
		passArgs = append(passArgs, "--pass="+strconv.Itoa(args.pass))
	}

	return passArgs
}
