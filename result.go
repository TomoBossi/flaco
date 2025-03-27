package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
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
	flacFileName := getFileName(flac, true)
	flacFileHash, err := getFileHash(flac)
	if err != nil {
		return nil, err
	}
	flacSamplingRate, flacBitDepth, err := getFlacInfo(flac)
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

func NewResultFromValues(values []string) (*result, error) {
	if len(values) != len(ResultFields()) {
		return nil, fmt.Errorf("Error when creating new result from values: The number of values does not match the number of fields in the result struct.")
	}

	unixTime, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Error when converting string %s to int64 (unixTime):\n\t%s", values[0], err.Error())
	}

	flacSamplingRate, err := strconv.Atoi(values[3])
	if err != nil {
		return nil, fmt.Errorf("Error when converting string %s to int (flacSamplingRate):\n\t%s", values[3], err.Error())
	}

	flacBitDepth, err := strconv.Atoi(values[4])
	if err != nil {
		return nil, fmt.Errorf("Error when converting string %s to int (flacBitDepth):\n\t%s", values[4], err.Error())
	}

	mp3Bitrate, err := strconv.Atoi(values[5])
	if err != nil {
		return nil, fmt.Errorf("Error when converting string %s to int (mp3Bitrate):\n\t%s", values[5], err.Error())
	}

	success, err := strconv.ParseBool(values[6])
	if err != nil {
		return nil, fmt.Errorf("Error when converting string %s to bool (result):\n\t%s", values[7], err.Error())
	}

	numSwaps, err := strconv.Atoi(values[7])
	if err != nil {
		return nil, fmt.Errorf("Error when converting string %s to int (numSwaps):\n\t%s", values[7], err.Error())
	}

	elapsedTime, err := strconv.ParseFloat(values[8], 64)
	if err != nil {
		return nil, fmt.Errorf("Error when converting string %s to float64 (elapsedTime):\n\t%s", values[8], err.Error())
	}

	return &result{
		unixTime:         unixTime,
		flacFileName:     values[1],
		flacFileHash:     values[2],
		flacSamplingRate: flacSamplingRate,
		flacBitDepth:     flacBitDepth,
		mp3Bitrate:       mp3Bitrate,
		result:           success,
		numSwaps:         numSwaps,
		elapsedTime:      elapsedTime,
		timestamp:        values[9],
	}, nil
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

func ResultFields() []string {
	return []string{"unixTime(s)", "flacFileName", "flacFileHash", "flacSamplingRate(Hz)", "flacBitDepth", "mp3Bitrate(kbps)", "result", "numSwaps", "elapsedTime(s)", "timestamp(mm:ss)"}
}

func (r result) ResultValuesCSV() string {
	return fmt.Sprintf("%d,%s,%s,%d,%d,%d,%t,%d,%f,%s", r.unixTime, r.flacFileName, r.flacFileHash, r.flacSamplingRate, r.flacBitDepth, r.mp3Bitrate, r.result, r.numSwaps, r.elapsedTime, r.timestamp)
}
