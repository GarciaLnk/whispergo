package main

import (
	"fmt"
	"io"
	"os"
	"time"

	// Package imports
	whisper "github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	wav "github.com/go-audio/wav"
)

func Process(model whisper.Model, path string, cb *string) (string, error) {
	var data []float32

	// Create processing context
	context, err := model.NewContext()
	if err != nil {
		return "app error", err
	}
	context.SetTranslate(true)

	fmt.Printf("\n%s\n", context.SystemInfo())

	// Open the file
	fmt.Fprintf(os.Stdout, "Loading %q\n", path)
	fh, err := os.Open(path)
	if err != nil {
		return "error opening file", err
	}
	defer fh.Close()

	// Decode the WAV file - load the full buffer
	dec := wav.NewDecoder(fh)
	if buf, err := dec.FullPCMBuffer(); err != nil {
		return "error decoding wav", err
	} else if dec.SampleRate != whisper.SampleRate {
		return "unsupported sample rate", fmt.Errorf("unsupported sample rate: %d", dec.SampleRate)
	} else if dec.NumChans != 1 {
		return "unsupported number of channels", fmt.Errorf("unsupported number of channels: %d", dec.NumChans)
	} else {
		data = buf.AsFloat32Buffer().Data
	}

	// Process the data
	fmt.Fprintf(os.Stdout, "  ...processing %q\n", path)
	context.ResetTimings()
	if err := context.Process(data, nil, nil); err != nil {
		return "error processing", err
	}

	context.PrintTimings()

	// Print out the results
	return Output(os.Stdout, context, cb)
}

// Output text to terminal
func Output(w io.Writer, context whisper.Context, cb *string) (string, error) {
	for {
		segment, err := context.NextSegment()
		if err == io.EOF {
			return "", nil
		} else if err != nil {
			return "error", err
		}
		fmt.Fprintf(w, "[%6s->%6s]", segment.Start.Truncate(time.Millisecond), segment.End.Truncate(time.Millisecond))
		fmt.Fprintln(w, " ", segment.Text)
		*cb += segment.Text + "\n"
	}
}
