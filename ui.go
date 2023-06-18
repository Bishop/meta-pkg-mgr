package main

import (
	"io"
	"sync"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func createProgress(wg *sync.WaitGroup) *mpb.Progress {
	return mpb.New(mpb.WithWaitGroup(wg))
}

func createBar(progress *mpb.Progress, si *spinnerItem) *mpb.Bar {
	return progress.New(1,
		mpb.SpinnerStyle().PositionLeft(),
		mpb.AppendDecorators(
			decor.Name(si.name, decor.WCSyncSpaceR),
			decor.Elapsed(decor.ET_STYLE_GO, decor.WCSyncWidth),
		),
		mpb.BarWidth(1),
		mpb.BarFillerOnComplete("+"),
		mpb.BarExtender(si, false),
	)
}

type spinnerItem struct {
	name   string
	result []string
}

func (si *spinnerItem) Fill(writer io.Writer, stat decor.Statistics) error {
	if si.result == nil {
		return nil
	}

	for _, line := range si.result {
		_, err := writer.Write([]byte(line + "\n"))
		if err != nil {
			return err
		}
	}

	return nil
}
