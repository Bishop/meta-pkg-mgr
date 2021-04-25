package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	config := new(Config)

	fileName, err := filepath.Abs("config.json")

	fatalOnError(err)

	readConfig(fileName, config)

	result := make([]*PkgItem, 0, 100)

	for _, pkg := range config.PkgConfigs {
		for _, item := range processPackageManager(pkg) {
			result = append(result, item)
		}
	}

	for _, item := range result {
		fmt.Println(item)
	}
}

func processPackageManager(cfg PkgConfig) OutdatedRecords {
	result := make(OutdatedRecords)

	for _, step := range cfg.Flow {
		output := captureOutput(step.Command)

		for _, item := range extractVersionInfo(output, step.RegExp) {
			result.Update(cfg.Name, item)
		}
	}

	result.Filter()

	return result
}

func captureOutput(command string) string {
	args := strings.Split(command, " ")

	cmd := exec.Command(args[0], args[1:]...)

	output, err := cmd.Output()
	fatalOnError(err)

	return string(output)
}

func extractVersionInfo(text string, re string) []Hash {
	if re == "" {
		return []Hash{}
	} else {
		re := regexp.MustCompile(re)
		matches := re.FindAllStringSubmatch(text, 1000)

		items := make([]Hash, len(matches))

		for _, row := range matches {
			item := make(Hash)

			for i, name := range re.SubexpNames() {
				item[name] = row[i]
			}

			items = append(items, item)
		}

		return items
	}
}

func readConfig(filename string, data interface{}) {
	content, err := os.ReadFile(filename)

	fatalOnError(err)

	err = json.Unmarshal(content, data)

	fatalOnError(err)

	return
}

func fatalOnError(err error) {
	if err != nil {
		panic(err)
	}
}
