package main

import (
	"flag"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

func main() {
	saveConfig := flag.Bool("save-config", false, fmt.Sprintf("Save default config to %s", configFileName()))
	showHelp := flag.Bool("h", false, "Show help")
	upgrade := flag.Bool("u", false, "Install updates")

	flag.Parse()

	if *saveConfig {
		saveDefaultConfig()
		fmt.Printf("Config saved to %s\n", configFileName())
		return
	} else if *showHelp {
		flag.PrintDefaults()
		return
	} else if *upgrade {

	}

	config := new(Config)

	readConfig(config)

	fmt.Printf("Use config: %s\n", actualConfigFile())

	runConcurrently(config.PkgConfigs)
}

func runConcurrently(configs []PkgConfig) {
	wait := sync.WaitGroup{}
	wait.Add(len(configs))

	progress := createProgress(&wait)

	for _, pkg := range configs {
		si := &spinnerItem{name: pkg.Name}

		bar := createBar(progress, si)

		go func(pkg PkgConfig) {
			processPackageManager(pkg, si)
			bar.IncrBy(1)
			wait.Done()
		}(pkg)
	}

	progress.Wait()
}

func processPackageManager(cfg PkgConfig, si *spinnerItem) {
	result := make(OutdatedRecords)

	for _, step := range cfg.Flow {
		output := captureOutput(step.Command, cfg.Shell)

		for _, item := range extractVersionInfo(output, step.RegExp) {
			result.Update(cfg.Name, item)
		}
	}

	si.result = result.Filter().List()
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
