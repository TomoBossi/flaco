package main

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func getFileName(path string, removeExtension bool) string {
	fileName := filepath.Base(path)
	if !removeExtension {
		return fileName
	}
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
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
	fmt.Printf("\nConverting FLAC to MP3@%dkbps...\n", bitrate)
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

func matchCSVHeader(path string, columns []string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("Error opening file %s:\n\t%s", path, err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		return false, fmt.Errorf("Error when getting header of CSV file %s:\n\t%s", path, err.Error())
	}

	if len(header) != len(columns) {
		return false, nil
	}

	for i, field := range header {
		if strings.TrimSpace(field) != strings.TrimSpace(columns[i]) {
			return false, nil
		}
	}

	return true, nil
}

func newResultsFromCSV(path string) ([]*result, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Error opening file %s:\n\t%s", path, err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("Error when getting header of CSV file %s:\n\t%s", path, err.Error())
	}

	results := []*result{}
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("Error when reading a record of CSV file %s:\n\t%s", path, err.Error())
		}
		res, err := NewResultFromValues(record)
		if err != nil {
			return nil, fmt.Errorf("Error when parsing a record of CSV file %s as a result struct:\n\t%s", path, err.Error())
		}
		results = append(results, res)
	}

	return results, nil
}
