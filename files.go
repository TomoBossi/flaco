package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func getFileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("Error opening file %s:\n\t%s", path, err.Error())
	}
	defer file.Close()

	hash := sha256.New()
	if _, err = io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("Error copying file %s to get its hash:\n\t%s", path, err.Error())
	}

	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes), nil
}

func getFileName(path string) string {
	return filepath.Base(path)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	return err == nil, err
}

func getPath(directory, name, extension string) (string, error) {
	i := 0
	var path string
	for {
		if i == 0 {
			path = filepath.Join(directory, fmt.Sprintf("%s.%s", name, extension))
		} else {
			path = filepath.Join(directory, fmt.Sprintf("%s_%d.%s", name, i, extension))
		}
		if pathExists, err := exists(path); !pathExists && errors.Is(err, fs.ErrNotExist) {
			return path, nil
		} else if err != nil {
			return "", fmt.Errorf("Error getting new file path:\n\t%s", err.Error())
		}
		i++
	}
}

func flac2Mp3(flac string, bitrate int) (string, error) {
	fmt.Print("Converting FLAC to MP3...\n")
	path, err := getPath(os.TempDir(), "temp", "mp3")
	if err != nil {
		return "", err
	}
	cmd := exec.Command("ffmpeg", "-i", flac, "-ab", fmt.Sprintf("%dk", bitrate), path)
	if _, err := cmd.Output(); err != nil {
		return "", fmt.Errorf("Error when creating audio file %s:\n\t%s", path, err.Error())
	}
	return path, nil
}

func removeFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("Error when removing file %s:\n\t%s", path, err.Error())
	}
	return nil
}
