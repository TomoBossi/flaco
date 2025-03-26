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
Compile and run the program. Requires `ffmpeg` (to convert FLAC to MP3) and `mpv` (to play audio). I included an example FLAC in the assets folder. It also happens to be the goated opening track of the goat's magnum opus.

## TODOs:
- Add score-keeping using an auto-generated/auto-updated .csv file:
    - By default, generate or update file flaco.csv in the current working directory
    - Add an optional flag to disable scorekeeping during the current execution
    - Add an optional flag to provide a custom path for the .csv (replacing the default)
    - Check that the .csv (default or provided) has only the expected columns as its first line
- Add statistical analysis of previous results .csv
    - By default, use file flaco.csv in the current working directory
    - Reuse the optional flag mentioned above to provide a custom path for the .csv
    - If the .csv file already exists, show the summary after the title when running the program
    - Add an optional flag to only show the summary and exit (ignoring all other flags except the one that optionally provides the path of the .csv)
    - The summary must at the very least include a table showing, for each bitrate present in the data:
        - User precision (% of successful attempts)
        - Total number of attempts
        - Median number of swaps
        - Median total time elapsed
        - Estimation of the degree to which the results are explained by pure randomness