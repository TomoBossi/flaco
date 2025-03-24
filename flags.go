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
}

var availableBitrates = []int{8, 16, 24, 32, 40, 48, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320}

func ParseFlags() *Flags {
	flac := flag.String("flac", "", "Path to the input FLAC file (required)")
	bitrate := flag.Int("bitrate", 128, "Constant bitrate of the temporary MP3 file, measured in kbps. It can be 8, 16, 24, 32, 40, 48, 64, 80, 96, 112, 128, 160, 192, 224, 256 or 320") // https://trac.ffmpeg.org/wiki/Encode/MP3#CBREncoding

	flag.Parse()

	if *flac == "" {
		fmt.Print("Error: Required flag missing.\n")
		flag.Usage()
		os.Exit(1)
	}

	if !slices.Contains(availableBitrates, *bitrate) {
		*bitrate = GetNearest(*bitrate, availableBitrates)
		fmt.Printf("Bitrate unavailable. Using nearest available bitrate (%dkbps).\n", *bitrate)
	}

	return &Flags{
		flac:    *flac,
		bitrate: *bitrate,
	}
}
