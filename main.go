package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
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
		fmt.Printf("%s\n\n", err.Error())
		flag.Usage()
		return
	}

	var csvFile *os.File
	if flags.History() != "" {
		csvFile, err = os.OpenFile(flags.History(), os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			fmt.Printf("Error when opening CSV file %s:\n\t%s", flags.History(), err.Error())
			return
		}
		defer csvFile.Close()

		// TODO: Read and process CSV, output and print statistical summary
		// results, _ := newResultsFromCSV(flags.History())
		// for _, r := range results {
		// 	fmt.Println(*r)
		// }
		// fmt.Println()
	}

	if !flags.Summary() {
		mp3 := flags.Mp3()
		bitrate := 0
		timestamp := "00:00"
		if mp3 == "" {
			bitrate = flags.Bitrate()
			mp3, err = flac2Mp3(flags.Flac(), bitrate)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			defer removeFile(mp3)
		}

		fmt.Print("\nOptions:\n- Swap tracks (S)\n- Change start timestamp (t)\n- Increase/lower volume (+/-)\n- Make your decision! (d)\n\n")

		for {
			res, err := NewResult(flags.Flac(), mp3, bitrate, timestamp, flags.Volume())
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if flags.History() != "" {
				_, err = csvFile.WriteString(fmt.Sprintf("%s\n", res.ResultValuesCSV()))
				if err != nil {
					fmt.Printf("Error when appending row to CSV file %s:\n\t%s", flags.History(), err.Error())
					return
				}
			}

			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Again? (y/N) ")
			if scanner.Scan(); scanner.Text() == "y" {
				fmt.Print("\n")
				timestamp = res.Timestamp()
				if flags.Mp3() == "" {
					fmt.Printf("Continue using bitrate of %dkbps? (Y/n) ", bitrate)
					if scanner.Scan(); scanner.Text() == "n" {
						for {
							fmt.Print("Please provide a bitrate: ")
							if scanner.Scan() {
								if newBitrate, err := strconv.Atoi(scanner.Text()); err != nil {
									fmt.Print("\nInvalid input.\n\n")
								} else {
									newBitrate = getNearestBitrate(newBitrate)
									if newBitrate != bitrate {
										bitrate = newBitrate
										mp3, err = flac2Mp3(flags.Flac(), bitrate)
										if err != nil {
											fmt.Println(err.Error())
											return
										}
										defer removeFile(mp3)
									}
									break
								}
							}
						}
					}
					fmt.Print("\n")
				}
			} else {
				fmt.Print("\nSee you space flaco...\n")
				return
			}
		}
	}
}
