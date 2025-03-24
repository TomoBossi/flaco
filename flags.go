package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
)

type Flags struct {
	flac    string
	bitrate int
	mp3     string
}

var availableBitrates = []int{8, 16, 24, 32, 40, 48, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320}

func ParseFlags() *Flags {
	flac := flag.String("flac", "", "REQUIRED - Path to the input FLAC file")
	bitrate := flag.Int("bitrate", 128, "OPTIONAL - Constant bitrate of the temporary MP3 file, measured in kbps. It can be 8, 16, 24, 32, 40, 48, 64, 80, 96, 112, 128, 160, 192, 224, 256 or 320") // https://trac.ffmpeg.org/wiki/Encode/MP3#CBREncoding
	mp3 := flag.String("mp3", "", "OPTIONAL - Path to an input MP3 file (default \"\", an MP3 file will be temporarily generated by compressing the FLAC file)")

	flag.Parse()

	if *flac == "" {
		fmt.Print("Error: Required flag missing.\n")
		flag.Usage()
		os.Exit(1)
	}

	if !Exists(*flac) {
		fmt.Print("Error: FLAC file not found.\n")
		os.Exit(1)
	}

	if *mp3 != "" && !Exists(*mp3) {
		fmt.Print("Error: MP3 file not found.\n")
		os.Exit(1)
	}

	if !slices.Contains(availableBitrates, *bitrate) {
		*bitrate = GetNearest(*bitrate, availableBitrates)
		fmt.Printf("Bitrate unavailable. Using nearest available bitrate (%dkbps).\n", *bitrate)
	}

	return &Flags{
		flac:    *flac,
		bitrate: *bitrate,
		mp3:     *mp3,
	}
}
