package videocompressor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Compressor struct {
	File           string
	Encoder        string
	CRF            uint
	OutputLocation string
	Verbose        bool
}

func NewCompressor(file string, encoder string) *Compressor {
	return &Compressor{
		File:    file,
		Encoder: encoder,
		CRF:     28,
		Verbose: false,
	}
}

type compressibleFile struct {
	path    string
	file    string
	size    int64
	newSize int64
}

func (c *Compressor) Compress() {
	stat, err := os.Stat(c.File)
	if err != nil {
		fmt.Println("Failed to read file")
		os.Exit(1)
	}

	outputStat, err := os.Stat(c.OutputLocation)
	if err != nil {
		fmt.Println("Invalid output location:", c.OutputLocation)
		os.Exit(1)
	}

	if !outputStat.IsDir() {
		fmt.Println("Output location is not a directory")
		os.Exit(1)
	}

	if stat.IsDir() {
		var filesToCompress []*compressibleFile

		// get all files in directory
		files, err := os.ReadDir(c.File)
		if err != nil {
			fmt.Println("Failed to read directory")
			os.Exit(1)
		}

		for _, entry := range files {
			// Wenn der Eintrag kein Ordner ist und der Name mit mp4 endet und nicht mit compress_ anfängt
			name := entry.Name()

			if c.Verbose {
				fmt.Println("Checking file", name)
			}

			if entry.IsDir() {

				continue
			}

			ext := filepath.Ext(name)
			fileName := strings.TrimSuffix(name, ext)
			ext = strings.ToLower(ext)

			if !strings.HasSuffix(fileName, "_compressed") && (ext == ".mp4" || ext == ".mov" || ext == ".mkv" || ext == ".avi" || ext == ".mpg") {
				// Get size of file
				s, err := os.Stat(filepath.Join(c.File, name))

				var size int64
				if err != nil {
					size = 0
				} else {
					size = s.Size()
				}

				filesToCompress = append(filesToCompress, &compressibleFile{
					path: c.File,
					file: name,
					size: size,
				})
			}
		}

		if len(filesToCompress) == 0 {
			fmt.Println("No files to compress found!")
			os.Exit(0)
		}

		fmt.Println("Found", len(filesToCompress), "files to compress")

		for _, file := range filesToCompress {
			c.compressFile(file)
		}

		fmt.Println("Finished compressing all files")
	} else {
		dir, fileName := filepath.Split(c.File)
		c.compressFile(&compressibleFile{
			path: dir,
			file: fileName,
			size: stat.Size(),
		})
	}
}

func (c *Compressor) compressFile(file *compressibleFile) {
	fmt.Println("Compressing file:", file.file)

	ext := filepath.Ext(file.file)
	onlyName := strings.TrimSuffix(file.file, ext)

	newFileName := filepath.Join(c.OutputLocation, onlyName+"_compressed.mp4")
	fullFilePath := filepath.Join(file.path, file.file)

	cmd := exec.Command("ffmpeg", "-i", fullFilePath, "-c:v", c.Encoder, "-crf", strconv.Itoa(int(c.CRF)), "-movflags", "use_metadata_tags", "-map_metadata", "0", "-tag:v", "hvc1", newFileName)
	err := cmd.Run()

	if err != nil {
		return
	}

	cmd = exec.Command("exiftool", "-overwrite_original", "-tagsFromFile", fullFilePath, "-MediaCreateDate", "-FileModifyDate", "-MediaModifyDate", "-TrackModifyDate", "-TrackCreateDate", "-ModifyDate", "-CreateDate", newFileName)
	err = cmd.Run()

	if err != nil {
		return
	}

	s, err := os.Stat(newFileName)

	if err != nil {
		fmt.Println("Compressed file", file.file)
		return
	}

	file.newSize = s.Size()

	fmt.Println("Compressed file", file.file, "Old Size:", ByteCountSI(file.size), "New Size:", ByteCountSI(s.Size()), "Saved: ", ByteCountSI(file.size-s.Size()))
}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
