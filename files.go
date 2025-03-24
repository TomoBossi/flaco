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

type Result struct {
	unixTime    int64
	fileName    string
	bitrate     int
	result      bool
	numSwaps    int
	elapsedTime float64
	timestamp   string
	fileHash    string
}

func GetFileHash(path string) string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file %s: %s\n", path, err)
		os.Exit(1)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err = io.Copy(hash, file); err != nil {
		fmt.Printf("Error copying file %s to hash: %s\n", path, err)
		os.Exit(1)
	}

	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func GetFileName(path string) string {
	return filepath.Base(path)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func GetPath(directory, name, extension string) string {
	i := 0
	var path string
	for {
		if i == 0 {
			path = filepath.Join(directory, fmt.Sprintf("%s.%s", name, extension))
		} else {
			path = filepath.Join(directory, fmt.Sprintf("%s_%d.%s", name, i, extension))
		}
		if !Exists(path) {
			return path
		}
		i++
	}
}

func CreateMP3(flac string, bitrate int) string {
	path := GetPath(os.TempDir(), "temp", "mp3")
	cmd := exec.Command("ffmpeg", "-i", flac, "-ab", fmt.Sprintf("%dk", bitrate), path)
	if _, err := cmd.Output(); err != nil {
		fmt.Printf("Error when creating audio file %s: %s\n", path, err)
		os.Exit(1)
	}
	return path
}

func RemoveFile(path string) {
	if err := os.Remove(path); err != nil {
		fmt.Printf("Error when removing file %s: %s\n", path, err)
		os.Exit(1)
	}
}
