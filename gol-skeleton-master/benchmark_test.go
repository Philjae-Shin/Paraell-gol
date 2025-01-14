package main

import (
	"fmt"
	"os"
	"testing"
	"uk.ac.bris.cs/gameoflife/gol"
)

// Benchmarking (Can change -count)
// go test -run ^$ -bench . -benchtime 1x -count 8 | tee result/resultsNew.out
// go run golang.org/x/perf/cmd/benchstat -format csv result/resultsNew.out | tee result/resultsNew.csv

// CPU profiling (Can change -count)
// go test -bench /8_ -benchtime 1x -count 20 -cpuprofile cpu.prof

//go test -run ^$ -bench BenchmarkGol/512x512x1000-1 -timeout 100s -cpuprofile cpu.prof

// Convert to PDF
// go tool pprof -pdf -nodefraction=0 -unit=ms cpu.prof

const benchLength = 1000

func BenchmarkGol(b *testing.B) {

	for threads := 1; threads <= 16; threads++ {
		os.Stdout = nil
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
				//for range events {
				//
				//}
			}
		})
	}
}

//func BenchmarkGol(b *testing.B) {
//	// Disable all program output apart from benchmark results
//	os.Stdout = nil
//
//	for threads := 1; threads <= 16; threads++ {
//		b.Run(fmt.Sprintf("%d_workers", threads), func(b *testing.B) {
//			traceParams := gol.Params{
//				Turns:       100,
//				Threads:     threads,
//				ImageWidth:  512,
//				ImageHeight: 512,
//			}
//			events := make(chan gol.Event)
//			for i := 0; i < b.N; i++ {
//				go gol.Run(traceParams, events, nil)
//				complete := false
//				for !complete {
//					event := <-events
//					switch event.(type) {
//					case gol.FinalTurnComplete:
//						complete = true
//					}
//				}
//			}
//		})
//	}
//}
