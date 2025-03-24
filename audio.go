package main

import (
	"fmt"
	"os"
	"os/exec"
)

func playAudio(path string, timestamp string, volume int) *exec.Cmd {
	cmd := exec.Command("mpv", fmt.Sprintf("--start=%s", timestamp), fmt.Sprintf("--volume=%d", volume), path)
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error when playing audio file %s: %s\n", path, err)
		os.Exit(1)
	}
	return cmd
}

func StopAudio(process *exec.Cmd) {
	if err := process.Process.Kill(); err != nil {
		fmt.Printf("Error when stopping audio: %s\n", err)
		os.Exit(1)
	}
}

func PlayOneOf(flac, mp3, timestamp string, volume int, playFlac bool, currentAudioProcess *exec.Cmd) *exec.Cmd {
	if currentAudioProcess != nil {
		StopAudio(currentAudioProcess)
	}

	if playFlac {
		return playAudio(flac, timestamp, volume)
	}
	return playAudio(mp3, timestamp, volume)
}
