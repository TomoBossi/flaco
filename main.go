package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"time"
)

var title string = `Welcome to flaco! Can you hear the difference?
 _______ ___     _______ _______ _______ 
|   _   |   |   |   _   |   _   |   _   |
|.  1___|.  |   |.  1   |.  1___|.  |   |
|.  __) |.  |___|.  _   |.  |___|.  |   |
|:  |   |:  1   |:  |   |:  1   |:  1   |
|::.|   |::.. . |::.|:. |::.. . |::.. . |
'---'   '-------'--- ---'-------'-------' ♬

`

func mainLoop(flac, mp3 string, bitrate int, timestamp string, volume int) *Result {
	numSwaps := 0
	start := time.Now()
	playFlac := rand.IntN(2) == 0
	startedWithFlac := playFlac
	scanner := bufio.NewScanner(os.Stdin)
	audioProcess := PlayOneOf(flac, mp3, timestamp, volume, playFlac)
	for {
		fmt.Printf("PLAYING TRACK %d (started at %s)\n", (numSwaps)%2+1, timestamp)
		fmt.Print("What will you do? (s/t/d) ")
		if scanner.Scan() {
			fmt.Print("\n")
			input := scanner.Text()
			switch input {
			default:
				StopAudio(audioProcess)
				numSwaps++
				playFlac = !playFlac
				audioProcess = PlayOneOf(flac, mp3, timestamp, volume, playFlac)
			case "t":
				fmt.Print("Please provide a timestamp in mm:ss format: ")
				if scanner.Scan() {
					timestamp = scanner.Text()
				}
				StopAudio(audioProcess)
				audioProcess = PlayOneOf(flac, mp3, timestamp, volume, playFlac)
			case "d":
				StopAudio(audioProcess)
				for {
					fmt.Print("Which of the two was the FLAC file? (1/2) ")
					if scanner.Scan() {
						fmt.Print("\n")
						input = scanner.Text()
						if input != "1" && input != "2" {
							fmt.Print("Invalid input.\n\n")
						} else {
							result := (input == "1" && startedWithFlac) || (input == "2" && !startedWithFlac)
							if result {
								fmt.Print("Correct! Bye-bye.\n")
							} else {
								fmt.Print("Wrong! Bye-bye.\n")
							}
							return &Result{
								unixTime:    time.Now().Unix(),
								fileName:    GetFileName(flac),
								bitrate:     bitrate,
								result:      result,
								numSwaps:    numSwaps,
								elapsedTime: time.Now().Sub(start).Seconds(),
								timestamp:   timestamp,
								fileHash:    GetFileHash(flac),
							}
						}
					}
				}
			}
		}
	}
}

func main() {
	fmt.Print(title)

	flags := ParseFlags()

	fmt.Print("Converting FLAC to MP3...\n\n")
	mp3 := CreateMP3(flags.flac, flags.bitrate)
	defer RemoveFile(mp3)

	fmt.Print("Options:\n- Swap tracks (s, default)\n- Change start timestamp (t)\n- Make your decision! (d)\n\n")

	mainLoop(flags.flac, mp3, flags.bitrate, "00:00", 70)
}
