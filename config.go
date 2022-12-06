package main

import (
	_ "embed"
	"encoding/json"
	"os"
	"os/user"
	"path"
)

const configFsPath = ".upt/config.json"

//go:embed "config.json"
var configFileContent []byte

func readConfig(data interface{}) {
	var err error

	if configFileExist() {
		configFileContent, err = os.ReadFile(configFileName())
		fatalOnError(err)
	}

	err = json.Unmarshal(configFileContent, data)

	fatalOnError(err)

	return
}

func configFileExist() bool {
	_, err := os.Stat(configFileName())

	return !os.IsNotExist(err)
}

func configFileName() string {
	u, _ := user.Current()

	return path.Join(u.HomeDir, configFsPath)
}

func saveDefaultConfig() {
	var err error

	if !configFileExist() {
		err = os.MkdirAll(path.Dir(configFileName()), 0740)
		fatalOnError(err)
	}

	err = os.WriteFile(configFileName(), configFileContent, 0640)
	fatalOnError(err)
}
