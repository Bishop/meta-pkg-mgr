package main

import (
	"flag"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

func main() {
	saveConfig := flag.Bool("save-config", false, fmt.Sprintf("Save default config to %s", configFileName()))
	showHelp := flag.Bool("h", false, "Show help")

	flag.Parse()

	if *saveConfig {
		saveDefaultConfig()
		fmt.Printf("Config saved to %s\n", configFileName())
		return
	} else if *showHelp {
		flag.PrintDefaults()
		return
	}

	config := new(Config)

	readConfig(config)

	fmt.Printf("Use config: %s\n", actualConfigFile())

	for pkgItems := range runConcurrently(config.PkgConfigs) {
		for _, item := range *pkgItems {
			fmt.Println(item)
		}
	}
}

func runConcurrently(configs []PkgConfig) chan *OutdatedRecords {
	result := make(chan *OutdatedRecords, len(configs))

	wait := sync.WaitGroup{}
	progress := mpb.New(mpb.WithWaitGroup(&wait))

	for _, pkg := range configs {
		wait.Add(1)
		bar := createBar(progress, pkg.Name)

		go func(pkg PkgConfig) {
			result <- processPackageManager(pkg)
			bar.IncrBy(1)
			wait.Done()
		}(pkg)
	}

	progress.Wait()

	close(result)

	return result
}

func processPackageManager(cfg PkgConfig) *OutdatedRecords {
	result := make(OutdatedRecords)

	for _, step := range cfg.Flow {
		output := captureOutput(step.Command, cfg.Shell)

		for _, item := range extractVersionInfo(output, step.RegExp) {
			result.Update(cfg.Name, item)
		}
	}

	result.Filter()

	return &result
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

func createBar(progress *mpb.Progress, name string) *mpb.Bar {
	return progress.AddSpinner(1, mpb.SpinnerOnLeft,
		mpb.PrependDecorators(
			decor.Name(name, decor.WCSyncSpaceR),
		),
		mpb.AppendDecorators(
			decor.Elapsed(decor.ET_STYLE_GO, decor.WCSyncWidth),
		),
		mpb.BarClearOnComplete(),
		mpb.BarWidth(4),
	)
}
