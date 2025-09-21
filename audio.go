package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"syscall"
)

func getFlacInfo(flac string) (int, int, error) {
	cmd := exec.Command("ffmpeg", "-f", "null", "-", "-i", flac)
	var out bytes.Buffer
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return 0, 0, fmt.Errorf("error when getting info of FLAC file %s:\n\t%s", flac, err.Error())
	}
	info := out.String()

	samplingRate := 0
	matches := regexp.MustCompile(`, (\d+) Hz`).FindStringSubmatch(info)
	if len(matches) > 1 {
		_, err = fmt.Sscanf(matches[1], "%d", &samplingRate)
		if err != nil {
			return 0, 0, fmt.Errorf("error when parsing sampling rate info of FLAC file %s:\n\t%s", flac, err.Error())
		}
	}

	bitDepth := 0
	matches = regexp.MustCompile(`, s(\d+)`).FindStringSubmatch(info)
	if len(matches) > 1 {
		_, err = fmt.Sscanf(matches[1], "%d", &bitDepth)
		if err != nil {
			return 0, 0, fmt.Errorf("error when parsing bit depth info of FLAC file %s:\n\t%s", flac, err.Error())
		}
	}

	return samplingRate, bitDepth, nil
}

func playAudio(path string, timestamp string, volume int) (*exec.Cmd, error) {
	cmd := exec.Command("mpv", fmt.Sprintf("--start=%s", timestamp), fmt.Sprintf("--volume=%d", volume), path)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error when playing audio file %s:\n\t%s", path, err.Error())
	}
	return cmd, nil
}

func stopAudio(process *exec.Cmd) error {
	if process.ProcessState != nil && process.ProcessState.Exited() {
		return nil
	}
	if err := syscall.Kill(-process.Process.Pid, syscall.SIGKILL); err != nil {
		return fmt.Errorf("error when stopping audio:\n\t%s", err.Error())
	}
	if _, err := process.Process.Wait(); err != nil {
		return fmt.Errorf("error when waiting for stopped audio:\n\t%s", err.Error())
	}
	return nil
}

func playOneOf(flac, mp3, timestamp string, volume int, playFlac bool, currentAudioProcess *exec.Cmd) (*exec.Cmd, error) {
	if currentAudioProcess != nil {
		stopAudio(currentAudioProcess)
	}

	if playFlac {
		return playAudio(flac, timestamp, volume)
	}
	return playAudio(mp3, timestamp, volume)
}
