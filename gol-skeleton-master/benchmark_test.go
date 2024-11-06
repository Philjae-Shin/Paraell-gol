package main

import (
	"fmt"
	"os"
	"testing"
	"uk.ac.bris.cs/gameoflife/gol"
)

// go test -run ^$ -bench . -benchtime 1x -count 10 | tee result/resultsNew.out
// go run golang.org/x/perf/cmd/benchstat -csv result/resultsNew.out | tee result/resultsNew.csv

const benchLength = 1000

func BenchmarkGol(b *testing.B) {
	// Disable all program output apart from benchmark results
	os.Stdout = nil

	for threads := 1; threads <= 16; threads++ {
		p := gol.Params{
			Turns:       benchLength,
			Threads:     threads,
			ImageWidth:  512,
			ImageHeight: 512,
		}
		name := fmt.Sprintf("%dx%dx%d-%d", p.ImageWidth, p.ImageHeight, p.Turns, p.Threads)
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				events := make(chan gol.Event)
				go gol.Run(p, events, nil)
				for range events {

				}
			}
		})
	}
}
