package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Params struct {
	input        string
	format       string
	output       string
	inwildcard   string
	debug        bool
	ffmpegbindir string
}

func main() {
	fmt.Println("convertallaudio 0.1 - written in golang")
	params := Params{}
	flag.Usage = func() {
		fmt.Println("Scans a directory recursively for wav files and converts them to a specific output format using FFMPEG.")
		flag.PrintDefaults()
	}
	flag.StringVar(&params.input, "input", "", "Input directory.")
	flag.StringVar(&params.format, "format", "ogg", "Output format.")
	flag.StringVar(&params.output, "output", "", "Output directory.")
	flag.StringVar(&params.inwildcard, "inwildcard", ".wav", "Defines the type of files to look for.")
	flag.BoolVar(&params.debug, "debug", false, "Turns on more output.")
	flag.Parse()

	// check if FFMPEG is available as environment variable
	params.ffmpegbindir = strings.ReplaceAll(os.Getenv("FFMPEG"), "\"", "")
	if params.ffmpegbindir == "" {
		log.Fatal("Please set the environment variable FFMPEG to the ffmpeg bin directory!")
	}

	// check parameter values
	if params.input == "" {
		log.Fatal("Please provide a input directory via -input=...")
	}
	VALID_FORMATS := []string{"flac", "ogg", "mp3"}
	if !contains(VALID_FORMATS, params.format) {
		log.Fatal("-format must be one of " + strings.Join(VALID_FORMATS, ","))
	}

	// browse the files
	var err = scanForInputFiles(params)
	if err != nil {
		log.Fatal(err)
	}
}

// search for wav files and convert them
func scanForInputFiles(params Params) error {
	var inputfiles []string
	err := filepath.Walk(params.input,
		func(fullpath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(fullpath, params.inwildcard) {
				inputfiles = append(inputfiles, fullpath)
			}
			return nil
		})
	if err != nil {
		return err
	}
	fmt.Printf("Processing %d files\n", len(inputfiles))

	for _, fullpath := range inputfiles {
		var outdir = filepath.Join(params.output, strings.TrimPrefix(filepath.Dir(fullpath), filepath.VolumeName(fullpath)))
		os.MkdirAll(outdir, os.ModePerm)
		var filenameNoExt = strings.TrimRight(filepath.Base(fullpath), params.inwildcard)
		var target = filepath.Join(outdir, filenameNoExt+"."+params.format)
		fmt.Println("Create: " + target)
		cmd := exec.Command(filepath.Clean(filepath.Join(params.ffmpegbindir, "ffmpeg.exe")), "-y", "-i", fullpath, target)
		var output bytes.Buffer
		cmd.Stdout = &output
		cmd.Stderr = &output
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		if params.debug {
			fmt.Println(output.String())
		}
	}
	return nil
}

func contains(arr []string, search string) bool {
	for _, value := range arr {
		if value == search {
			return true
		}
	}
	return false
}
