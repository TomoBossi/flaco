package main

import (
	"fmt"
)

func main() {
	fmt.Print(`Welcome to flaco! Can you hear the difference?
 _______ ___     _______ _______ _______ 
|   _   |   |   |   _   |   _   |   _   | ♪♬
|.  1___|.  |   |.  1   |.  1___|.  |   |
|.  __) |.  |___|.  _   |.  |___|.  |   |
|:  |   |:  1   |:  |   |:  1   |:  1   |
|::.|   |::.. . |::.|:. |::.. . |::.. . |
'---'   '-------'--- ---'-------'-------'

	`)

	flags, err := NewFlags()
	if err != nil {
		fmt.Println(err)
		return
	}

	mp3 := flags.Mp3()
	bitrate := 0
	if mp3 == "" {
		fmt.Print("Converting FLAC to MP3...\n\n")
		bitrate = flags.Bitrate()
		mp3, err := createMP3(flags.Flac(), bitrate)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer removeFile(mp3)
	}

	fmt.Print("Options:\n- Swap tracks (s, default)\n- Change start timestamp (t)\n- Increase/lower volume (+/-)\n- Make your decision! (d)\n\n")

	_, err = NewResult(flags.Flac(), mp3, bitrate, "00:00", flags.Volume())
	if err != nil {
		fmt.Println(err)
		return
	}
}
