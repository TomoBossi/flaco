package main

import (
	"fmt"
	"os"
	"os/exec"
)

func PlayAudio(path string, timestamp string, volume int) *exec.Cmd {
	cmd := exec.Command("mpv", fmt.Sprintf("--start=%s", timestamp), fmt.Sprintf("--volume=%d", volume), path)
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error when playing audio file %s: %s\n", path, err)
		os.Exit(1)
	}
	return cmd
}

func StopAudio(process *exec.Cmd) {
	err := process.Process.Kill()
	if err != nil {
		fmt.Printf("Error when stopping audio: %s\n", err)
		os.Exit(1)
	}
}

func PlayOneOf(flac, mp3, timestamp string, volume int, playFlac bool) *exec.Cmd {
	if playFlac {
		return PlayAudio(flac, timestamp, volume)
	}
	return PlayAudio(mp3, timestamp, volume)
}
