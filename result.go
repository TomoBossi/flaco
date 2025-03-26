package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"time"
)

type result struct {
	unixTime         int64
	flacFileName     string
	flacFileHash     string
	flacSamplingRate int
	flacBitDepth     int
	mp3Bitrate       int
	result           bool
	numSwaps         int
	elapsedTime      float64
	timestamp        string
}

func NewResult(flac, mp3 string, bitrate int, timestamp string, volume int) (*result, error) {
	flacFileName := getFileName(flac)
	flacFileHash, err := getFileHash(flac)
	if err != nil {
		return nil, err
	}
	flacSamplingRate, flacBitDepth, err := getAudioInfo(flac)
	if err != nil {
		return nil, err
	}
	numSwaps := 0
	start := time.Now()
	playFlac := rand.IntN(2) == 0
	startedWithFlac := playFlac
	scanner := bufio.NewScanner(os.Stdin)
	audioProcess, err := playOneOf(flac, mp3, timestamp, volume, playFlac, nil)
	if err != nil {
		return nil, err
	}
	for {
		fmt.Printf("PLAYING TRACK %d (started at %s on %d%% volume)\n", (numSwaps)%2+1, timestamp, volume)
		fmt.Print("What will you do? (S/t/+/-/d) ")
		if scanner.Scan() {
			fmt.Print("\n")
			switch scanner.Text() {
			default:
				numSwaps++
				playFlac = !playFlac
			case "t":
				fmt.Print("Please provide a timestamp in mm:ss format: ")
				if scanner.Scan() {
					timestamp = scanner.Text()
				}
			case "d":
				if err := stopAudio(audioProcess); err != nil {
					return nil, err
				}
				for {
					fmt.Print("Which of the two was the FLAC file? (1/2) ")
					if scanner.Scan() {
						fmt.Print("\n")
						input := scanner.Text()
						if input != "1" && input != "2" {
							fmt.Print("Invalid input.\n\n")
						} else {
							success := (input == "1" && startedWithFlac) || (input == "2" && !startedWithFlac)
							if success {
								fmt.Print("Correct!\n\n")
							} else {
								fmt.Print("Wrong!\n\n")
							}
							return &result{
								unixTime:         time.Now().Unix(),
								flacFileName:     flacFileName,
								flacFileHash:     flacFileHash,
								flacSamplingRate: flacSamplingRate,
								flacBitDepth:     flacBitDepth,
								mp3Bitrate:       bitrate,
								result:           success,
								numSwaps:         numSwaps,
								elapsedTime:      time.Now().Sub(start).Seconds(),
								timestamp:        timestamp,
							}, nil
						}
					}
				}
			case "+":
				volume = clamp(volume+5, 0, 100)
			case "-":
				volume = clamp(volume-5, 0, 100)

			}
			audioProcess, err = playOneOf(flac, mp3, timestamp, volume, playFlac, audioProcess)
			if err != nil {
				return nil, err
			}
		}
	}
}

func (r result) UnixTime() int64 {
	return r.unixTime
}

func (r result) FlacFileName() string {
	return r.flacFileName
}

func (r result) FlacFileHash() string {
	return r.flacFileHash
}

func (r result) FlacSamplingRate() int {
	return r.flacSamplingRate
}

func (r result) FlacBitDepth() int {
	return r.flacBitDepth
}

func (r result) Mp3Bitrate() int {
	return r.mp3Bitrate
}

func (r result) Result() bool {
	return r.result
}

func (r result) NumSwaps() int {
	return r.numSwaps
}

func (r result) ElapsedTime() float64 {
	return r.elapsedTime
}

func (r result) Timestamp() string {
	return r.timestamp
}
