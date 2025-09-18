package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".dictextractconfig.json"

type Config struct {
	DBURL           *string `json:"db_url"`
	RawDictDirPath  *string `json:"raw_dict_dir_path"`
	DAWGSaveDirPath *string `json:"dawg_save_dir_path"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting UserHomeDir: %v", err)
	}
	cfgFilePath := filepath.Join(homeDir, configFileName)
	return cfgFilePath, nil
}
