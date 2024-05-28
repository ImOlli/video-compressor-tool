package main

import (
	"flag"
	"fmt"
	"os"
	videocompressor "video-compressor/pkg/video-compressor"
)

var encoder = flag.String("encoder", "libx265", "Encoder to use for compression")
var crf = flag.Uint("crf", 18, "CRF value to use for compression")
var output = flag.String("o", "compressed", "The output location for the file")
var debug = flag.Bool("debug", false, "Enabled debug output")

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("You need to specify a file or directory to compress")
		os.Exit(1)
	}

	file := flag.Args()[0]

	checkTools()

	compressor := videocompressor.NewCompressor(file, *encoder)
	compressor.CRF = *crf
	compressor.OutputLocation = *output
	compressor.Verbose = *debug
	compressor.Compress()
}

func checkTools() {
	checkResult := videocompressor.CheckIfToolsAreInstalled()

	if !checkResult.FfmpegInstalled {
		fmt.Println("ffmpeg is not installed. Please install it and try again")
	}

	if !checkResult.ExiftoolInstalled {
		fmt.Println("exiftool is not installed. Please install it and try again")
	}

	if !checkResult.FfmpegInstalled || !checkResult.ExiftoolInstalled {
		os.Exit(1)
	}
}
