package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"slices"
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

type flags struct {
	flac    string
	bitrate int
}

func parseFlags() *flags {
	flac := flag.String("flac", "", "Path to the input FLAC file (required)")
	bitrate := flag.Int("bitrate", 128, "Constant bitrate of the temporary MP3 file, measured in kbps. It can be 8, 16, 24, 32, 40, 48, 64, 80, 96, 112, 128, 160, 192, 224, 256 or 320") // https://trac.ffmpeg.org/wiki/Encode/MP3#CBREncoding

	flag.Parse()

	if *flac == "" {
		fmt.Print("Error: Required flag missing.\n")
		flag.Usage()
		os.Exit(1)
	}

	if !slices.Contains([]int{8, 16, 24, 32, 40, 48, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320}, *bitrate) {
		fmt.Print("Error: Invalid bitrate.\n")
		flag.Usage()
		os.Exit(1)
	}

	return &flags{
		flac:    *flac,
		bitrate: *bitrate,
	}
}

func createMP3path(directory, fileName string) string {
	i := 0
	var path string
	for {
		if i == 0 {
			path = fmt.Sprintf("%s/%s.mp3", directory, fileName)
		} else {
			path = fmt.Sprintf("%s/%s_%d.mp3", directory, fileName, i)
		}
		if _, err := os.Stat(path); err != nil {
			return path
		}
		i++
	}
}

func createMP3(flags *flags, fileName string) string {
	path := createMP3path(".", fileName)
	cmd := exec.Command("ffmpeg", "-i", flags.flac, "-ab", fmt.Sprintf("%dk", flags.bitrate), path)
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error when creating audio file %s: %s\n", path, err)
		os.Exit(1)
	}
	return path
}

func removeFile(path string) {
	err := os.Remove(path)
	if err != nil {
		fmt.Printf("Error when removing file %s: %s\n", path, err)
		os.Exit(1)
	}
}

func playAudio(path string, timestamp string, volume int) *exec.Cmd {
	cmd := exec.Command("mpv", fmt.Sprintf("--start=%s", timestamp), fmt.Sprintf("--volume=%d", volume), path)
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error when playing audio file %s: %s\n", path, err)
		os.Exit(1)
	}
	return cmd
}

func stopAudio(process *exec.Cmd) {
	err := process.Process.Kill()
	if err != nil {
		fmt.Printf("Error when stopping audio: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	fmt.Print(title)

	flags := parseFlags()

	fmt.Print("Converting FLAC to MP3...\n\n")
	mp3 := createMP3(flags, "temp")

	fmt.Print("Options:\n- Swap tracks (s, default)\n- Change start timestamp (t)\n- Make your decision! (d)\n\n")
	i := rand.IntN(2)
	start := i
	timestamp := "00:00"
	volume := 50
	flac := []string{"1", "2"}[start]
	scanner := bufio.NewScanner(os.Stdin)
	audioFilePaths := []string{flags.flac, mp3}
	audio := playAudio(audioFilePaths[i], timestamp, volume)
	for {
		fmt.Printf("PLAYING TRACK %d (started at %s)\n", (i-start)%2+1, timestamp)
		fmt.Print("What will you do? (s/t/d) ")
		if scanner.Scan() {
			fmt.Print("\n")
			input := scanner.Text()
			switch input {
			default:
				stopAudio(audio)
				i++
				audio = playAudio(audioFilePaths[i%2], timestamp, volume)
			case "t":
				fmt.Print("Please provide a timestamp in mm:ss format: ")
				if scanner.Scan() {
					timestamp = scanner.Text()
				}
				stopAudio(audio)
				audio = playAudio(audioFilePaths[i%2], timestamp, volume)
			case "d":
				stopAudio(audio)
				removeFile(mp3)
				for {
					fmt.Print("Which of the two was the FLAC file? (1/2) ")
					if scanner.Scan() {
						fmt.Print("\n")
						input = scanner.Text()
						if input != "1" && input != "2" {
							fmt.Print("Invalid input.\n\n")
						} else {
							if input == flac {
								fmt.Print("Correct! Bye-bye.\n")
							} else {
								fmt.Print("Wrong! Bye-bye.\n")
							}
							return
						}
					}
				}
			}
		}
	}
}
