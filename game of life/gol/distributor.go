package gol

import (
	"fmt"
	"sync"
	"uk.ac.bris.cs/gameoflife/util"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	// TODO: Create a 2D slice to store the world.
	world := make([][]uint8, p.ImageHeight)
	for i := range world {
		world[i] = make([]uint8, p.ImageWidth)
	}
	c.ioCommand <- ioInput
	c.ioFilename <- fmt.Sprintf("%dx%d", p.ImageHeight, p.ImageWidth)
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			world[y][x] = <-c.ioInput
		}
	}

	// TODO: Execute all turns of the Game of Life.
	turn := 0
	c.events <- StateChange{turn, Executing}
	for turn := 0; turn < p.Turns; turn++ {
		// Calculate next step
		nextWorld := make([][]uint8, p.ImageHeight)
		for i := range nextWorld {
			nextWorld[i] = make([]uint8, p.ImageWidth)
		}

		if p.Threads == 1 {
			// Single-thread
			nextWorld, _ = computeNextState(p.ImageHeight, p.ImageWidth, 0, p.ImageHeight, world)
		} else {
			// Multi-thread
			var wg sync.WaitGroup
			rowsPerWorker := p.ImageHeight / p.Threads
			for i := 0; i < p.Threads; i++ {
				wg.Add(1)
				startY := i * rowsPerWorker
				endY := (i + 1) * rowsPerWorker
				if i == p.Threads-1 {
					endY = p.ImageHeight
				}
				go func(startY, endY int) {
					defer wg.Done()
					part, _ := computeNextState(p.ImageHeight, p.ImageWidth, startY, endY, world)
					for y := startY; y < endY; y++ {
						nextWorld[y] = part[y-startY]
					}
				}(startY, endY)
			}
			wg.Wait()
		}

		// Update World
		for y := range world {
			copy(world[y], nextWorld[y])
		}

		// Cells Turn Complete
		c.events <- TurnComplete{
			CompletedTurns: turn + 1,
		}
	}

	// TODO: Report the final state using FinalTurnCompleteEvent.

	// Make sure that the Io has finished any output before exiting.
	finalAliveCells := make([]util.Cell, p.ImageHeight*p.ImageWidth)
	_, finalAliveCells = getAliveCells(p, world)

	finalState := FinalTurnComplete{
		CompletedTurns: p.Turns,
		Alive:          finalAliveCells,
	}

	c.events <- finalState

	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
