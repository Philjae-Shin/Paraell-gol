package main

import (
	"fmt"
	"os"
	"testing"
	"uk.ac.bris.cs/gameoflife/gol"
)

// go test -run ^$ -bench . -benchtime 1x -count 20 | tee result/results.out
// go run golang.org/x/perf/cmd/benchstat -csv result/results.out | tee result/results.csv

//go test -run a$ -bench BenchmarkGol/512x512x1000-1 -timeout 100s -cpuprofile cpu.prof

// go test -bench /8_ -benchtime 1x -count 20 -cpuprofile result/cpu.prof | tee result/benchmark_output.txt
// go tool pprof -pdf -nodefraction=0 -unit=ms cpu.prof

const benchLength = 1000

func BenchmarkGol(b *testing.B) {
	for threads := 1; threads <= 16; threads++ {
		os.Stdout = nil // Disable all program output apart from benchmark results
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
