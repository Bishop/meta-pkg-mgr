package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

func main() {
	config := new(Config)

	readConfig(config)

	result := make(chan *PkgItem)

	go runConcurrently(config.PkgConfigs, result)

	for item := range result {
		fmt.Println(item)
	}
}

func runConcurrently(configs []PkgConfig, result chan *PkgItem) {
	wait := sync.WaitGroup{}

	for _, pkg := range configs {
		wait.Add(1)

		go func(pkg PkgConfig) {
			for _, item := range processPackageManager(pkg) {
				result <- item
			}
			wait.Done()
		}(pkg)
	}

	wait.Wait()

	close(result)
}

func processPackageManager(cfg PkgConfig) OutdatedRecords {
	result := make(OutdatedRecords)

	for _, step := range cfg.Flow {
		output := captureOutput(step.Command, cfg.Shell)

		for _, item := range extractVersionInfo(output, step.RegExp) {
			result.Update(cfg.Name, item)
		}
	}

	result.Filter()

	return result
}

func captureOutput(command string, shell string) string {
	var args []string

	if shell == "" {
		args = strings.Split(command, " ")
	} else {
		args = strings.Split(shell, " ")
		args = append(args, command)
	}

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

func fatalOnError(err error) {
	if err != nil {
		panic(err)
	}
}
