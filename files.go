package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func getFileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("Error opening file %s:\n\t%s", path, err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err = io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("Error copying file %s to get its hash:\n\t%s", path, err)
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
		if pathExists, err := exists(path); !pathExists {
			return path, nil
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("Error getting new file path:\n\t%s", err)
		}
		i++
	}
}

func createMP3(flac string, bitrate int) (string, error) {
	path, err := getPath(os.TempDir(), "temp", "mp3")
	if err != nil {
		return "", err
	}
	cmd := exec.Command("ffmpeg", "-i", flac, "-ab", fmt.Sprintf("%dk", bitrate), path)
	if _, err := cmd.Output(); err != nil {
		return "", fmt.Errorf("Error when creating audio file %s:\n\t%s", path, err)
	}
	return path, nil
}

func removeFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("Error when removing file %s:\n\t%s", path, err)
	}
	return nil
}
