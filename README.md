# flaco
*Welcome to flaco! Can **you** hear the difference?* (please read using the ["Gatorade. Is it in you?"](https://www.youtube.com/watch?v=E4B8owXm0Co) cadence)
```
 _______ ___     _______ _______ _______ 
|   _   |   |   |   _   |   _   |   _   | ♪♬
|.  1___|.  |   |.  1   |.  1___|.  |   |
|.  __) |.  |___|.  _   |.  |___|.  |   |
|:  |   |:  1   |:  |   |:  1   |:  1   |
|::.|   |::.. . |::.|:. |::.. . |::.. . |
'---'   '-------'--- ---'-------'-------'
```

## Why?
Because I wanted to know if I can reliably tell the difference between 4.1kHz FLAC and MP3 at bitrates upwards of 128kbps.

## What?
This is an interactive CLI program that tests the user's (and the user's setup) ability to discern high-fidelity uncompressed audio from standard audio. It takes a FLAC file as a parameter and allows the user to alternate between playing the original FLAC and a compressed MP3 version of it, without revealing which is which. When the user is ready, the program will ask and finally reveal which of the two was the FLAC file.

## How?
Compile and run the program. Requires `ffmpeg` (to convert FLAC to MP3 and get file info) and `mpv` (to play audio).