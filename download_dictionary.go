package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func downloadDictionary(url string) (filepath string, err error) {

	// Create tmp file
	tmpFile, err := os.CreateTemp("", "dict-*")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	// Get data
	resp, err := http.Get(url)
	if err != nil {
		defer os.Remove(tmpFile.Name())
		return "", err
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		defer os.Remove(tmpFile.Name())
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write data to tmp file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		defer os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}
