package videocompressor

import (
	"os/exec"
)

type ToolCheckResult struct {
	FfmpegInstalled   bool
	ExiftoolInstalled bool
}

func CheckIfToolsAreInstalled() *ToolCheckResult {
	ffmpegInstalled := checkIfFFmpegIsInstalled()
	exiftoolInstalled := checkIfExiftoolInstalled()

	return &ToolCheckResult{
		FfmpegInstalled:   ffmpegInstalled,
		ExiftoolInstalled: exiftoolInstalled,
	}
}

func checkIfExiftoolInstalled() bool {
	cmd := exec.Command("exiftool", "-ver")
	err := cmd.Run()

	if err != nil {
		return false
	}

	return true
}

func checkIfFFmpegIsInstalled() bool {
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()

	if err != nil {
		return false
	}

	return true
}
