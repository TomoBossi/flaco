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

	flags := NewFlags()

	mp3 := flags.Mp3()
	bitrate := 0
	if mp3 == "" {
		fmt.Print("Converting FLAC to MP3...\n\n")
		bitrate = flags.Bitrate()
		mp3 = createMP3(flags.Flac(), bitrate)
		defer removeFile(mp3)
	}

	fmt.Print("Options:\n- Swap tracks (s, default)\n- Change start timestamp (t)\n- Increase/lower volume (+/-)\n- Make your decision! (d)\n\n")

	NewResult(flags.Flac(), mp3, bitrate, "00:00", flags.Volume())
}
