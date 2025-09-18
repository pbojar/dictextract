package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func write(cfg Config) error {
	cfgFilePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error - getConfigFilePath: %v", err)
	}

	jsonData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	outFile, err := os.Create(cfgFilePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer func() {
		err := outFile.Close()
		if err != nil {
			fmt.Printf("error closing file: %v\n", err)
		}
	}()

	_, err = outFile.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}
