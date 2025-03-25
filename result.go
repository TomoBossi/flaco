package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"time"
)

type result struct {
	unixTime    int64
	fileName    string
	bitrate     int
	result      bool
	numSwaps    int
	elapsedTime float64
	timestamp   string
	fileHash    string
}

func NewResult(flac, mp3 string, bitrate int, timestamp string, volume int) *result {
	numSwaps := 0
	start := time.Now()
	playFlac := rand.IntN(2) == 0
	startedWithFlac := playFlac
	scanner := bufio.NewScanner(os.Stdin)
	audioProcess := playOneOf(flac, mp3, timestamp, volume, playFlac, nil)
	for {
		fmt.Printf("PLAYING TRACK %d (started at %s on %d%% volume)\n", (numSwaps)%2+1, timestamp, volume)
		fmt.Print("What will you do? (s/t/+/-/d) ")
		if scanner.Scan() {
			fmt.Print("\n")
			input := scanner.Text()
			switch input {
			default:
				numSwaps++
				playFlac = !playFlac
			case "t":
				fmt.Print("Please provide a timestamp in mm:ss format: ")
				if scanner.Scan() {
					timestamp = scanner.Text()
				}
			case "d":
				stopAudio(audioProcess)
				for {
					fmt.Print("Which of the two was the FLAC file? (1/2) ")
					if scanner.Scan() {
						fmt.Print("\n")
						input = scanner.Text()
						if input != "1" && input != "2" {
							fmt.Print("Invalid input.\n\n")
						} else {
							success := (input == "1" && startedWithFlac) || (input == "2" && !startedWithFlac)
							if success {
								fmt.Print("Correct! Bye-bye.\n")
							} else {
								fmt.Print("Wrong! Bye-bye.\n")
							}
							return &result{
								unixTime:    time.Now().Unix(),
								fileName:    getFileName(flac),
								bitrate:     bitrate,
								result:      success,
								numSwaps:    numSwaps,
								elapsedTime: time.Now().Sub(start).Seconds(),
								timestamp:   timestamp,
								fileHash:    getFileHash(flac),
							}
						}
					}
				}
			case "+":
				volume = clamp(volume+5, 0, 100)
			case "-":
				volume = clamp(volume-5, 0, 100)

			}
			audioProcess = playOneOf(flac, mp3, timestamp, volume, playFlac, audioProcess)
		}
	}
}

func (r result) UnixTime() int64 {
	return r.unixTime
}

func (r result) FileName() string {
	return r.fileName
}

func (r result) Bitrate() int {
	return r.bitrate
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

func (r result) FileHash() string {
	return r.fileHash
}
