package gol

import (
	"fmt"
	"strconv"
	"sync"
	"time"

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

const Save int = 0
const Quit int = 1
const Pause int = 2
const unPause int = 3

var wg sync.WaitGroup
var rWg sync.WaitGroup
var rMutex = new(sync.Mutex)
var mutex = new(sync.Mutex)
var cond = sync.NewCond(mutex)
var rCond = sync.NewCond(rMutex)
var readInProgress bool = false
var writeInProgress bool = false

// Modify params in calculateNextState
func worker(p Params, startY, endY, startX, endX int, world [][]uint8, c distributorChannels, turn int) {
	cellUpdates := make([]util.Cell, (endY-startY)*endX/2)
	nextWorldPart := make([][]uint8, endY-startY)
	previousWorld := make([][]uint8, p.ImageHeight)
	for h := range world {
		previousWorld[h] = make([]uint8, endX)
	}
	for i := range nextWorldPart {
		nextWorldPart[i] = make([]uint8, endX)
	}

	rCond.L.Lock()
	for readInProgress == false {
		rCond.Wait()
	}
	for j := range world {
		copy(previousWorld[j], world[j])
	}
	rCond.L.Unlock()
	rWg.Done()
	nextWorldPart, cellUpdates = calculateNextState(p.ImageHeight, p.ImageWidth, startY, endY, previousWorld)

	// Waits other goroutines to copy the previous world

	cond.L.Lock()
	for writeInProgress == false {
		cond.Wait()
	}
	for j := range nextWorldPart {
		copy(world[startY+j], nextWorldPart[j])
	}
	for _, cell := range cellUpdates {
		c.events <- CellFlipped{
			CompletedTurns: turn,
			Cell:           cell,
		}
	}
	cond.L.Unlock()
	wg.Done()
}

func handleOutput(p Params, c distributorChannels, world [][]uint8, t int) {
	c.ioCommand <- 0
	outFilename := strconv.Itoa(p.ImageHeight) + "x" + strconv.Itoa(p.ImageWidth) + "x" + strconv.Itoa(t)
	c.ioFilename <- outFilename
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			c.ioOutput <- world[y][x]
		}
	}
	c.events <- ImageOutputComplete{
		CompletedTurns: t,
		Filename:       outFilename,
	}
}

// Gets input from IO and initialises cellflip
func handleInput(p Params, c distributorChannels, world [][]uint8) [][]uint8 {
	filename := strconv.Itoa(p.ImageHeight) + "x" + strconv.Itoa(p.ImageWidth)
	c.ioCommand <- 1
	c.ioFilename <- filename
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			num := <-c.ioInput
			world[y][x] = num
			if num == 255 {
				c.events <- CellFlipped{
					CompletedTurns: 0,
					Cell:           util.Cell{X: x, Y: y},
				}
			}
		}
	}
	return world
}

func handleKeyPress(p Params, c distributorChannels, keyPresses <-chan rune, worldChannel <-chan [][]uint8, turnChannel <-chan int, actions chan int) {
	paused := false
	for {
		input := <-keyPresses

		switch input {
		case 's':
			actions <- Save
			w := <-worldChannel
			turn := <-turnChannel
			go handleOutput(p, c, w, turn)

		case 'q':
			actions <- Quit
			w := <-worldChannel
			turn := <-turnChannel
			go handleOutput(p, c, w, turn)

			newState := StateChange{CompletedTurns: turn, NewState: State(Quitting)}
			fmt.Println(newState.String())

			c.events <- newState
			c.events <- FinalTurnComplete{CompletedTurns: turn}
		case 'p':
			if paused {
				actions <- unPause
				turn := <-turnChannel
				paused = false
				newState := StateChange{CompletedTurns: turn, NewState: State(Executing)}
				fmt.Println(newState.String())
				c.events <- newState
			} else {
				actions <- Pause
				turn := <-turnChannel
				paused = true
				newState := StateChange{CompletedTurns: turn, NewState: State(Paused)}
				fmt.Println(newState.String())
				c.events <- newState
			}

		case 'k':
		}

	}

}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels, keyPresses <-chan rune) {
	world := make([][]uint8, p.ImageHeight)
	previousWorld := make([][]uint8, p.ImageHeight)
	sharedWorld := make([][]uint8, p.ImageHeight)
	for i := range world {
		world[i] = make([]uint8, p.ImageWidth)
		previousWorld[i] = make([]uint8, p.ImageWidth)
		sharedWorld[i] = make([]uint8, p.ImageWidth)
	}

	world = handleInput(p, c, world)

	currentTurn := 0
	ticker := time.NewTicker(2 * time.Second)
	terminateSignal := make(chan bool)
	isPaused := false
	shouldQuit := false
	resumeSignal := make(chan bool)
	go func() {
		for {
			if !shouldQuit {
				select {
				case <-terminateSignal:
					return
				case <-ticker.C:
					aliveCount, _ := calculateAliveCells(p, previousWorld)
					aliveReport := AliveCellsCount{
						CompletedTurns: currentTurn,
						CellsCount:     aliveCount,
					}
					c.events <- aliveReport
				}
			} else {
				return
			}
		}
	}()

	turnChannel := make(chan int)
	worldChannel := make(chan [][]uint8)
	userAction := make(chan int)
	go handleKeyPress(p, c, keyPresses, worldChannel, turnChannel, userAction)
	go func() {
		for {

			select {
			case command := <-userAction:
				switch command {
				case Pause:
					isPaused = true
					turnChannel <- currentTurn
				case unPause:
					isPaused = false
					turnChannel <- currentTurn
					resumeSignal <- true
				case Quit:
					worldChannel <- world
					turnChannel <- currentTurn
					shouldQuit = true
				case Save:
					worldChannel <- world
					turnChannel <- currentTurn
				}
			}
		}
	}()

	channels := make([]chan [][]uint8, p.Threads)
	cellUpdates := make([]util.Cell, p.ImageHeight*p.ImageWidth)
	unit := int(p.ImageHeight / p.Threads)

	for t := 0; t < p.Turns; t++ {
		cellUpdates = make([]util.Cell, p.ImageHeight*p.ImageWidth)
		if isPaused {
			<-resumeSignal
		}
		if !isPaused && !shouldQuit {
			currentTurn = t
			for j := range world {
				copy(previousWorld[j], world[j])
			}
			if p.Threads == 1 {
				world, cellUpdates = calculateNextState(p.ImageHeight, p.ImageWidth, 0, p.ImageHeight, world)
				for _, cell := range cellUpdates {
					c.events <- CellFlipped{
						CompletedTurns: currentTurn,
						Cell:           cell,
					}
				}
			} else {
				rWg.Add(p.Threads)
				wg.Add(p.Threads)
				for i := 0; i < p.Threads; i++ {
					channels[i] = make(chan [][]uint8)
					if i == p.Threads-1 {
						go worker(p, i*unit, p.ImageHeight, 0, p.ImageWidth, world, c, currentTurn)
					} else {
						go worker(p, i*unit, (i+1)*unit, 0, p.ImageWidth, world, c, currentTurn)
					}
				}
				rCond.L.Lock()
				readInProgress = true
				rCond.Broadcast()
				rCond.L.Unlock()
				rWg.Wait()

				cond.L.Lock()
				readInProgress = false
				writeInProgress = true
				cond.Broadcast()
				cond.L.Unlock()
				wg.Wait()
				writeInProgress = false
			}

			c.events <- TurnComplete{
				CompletedTurns: currentTurn,
			}
		} else {
			if shouldQuit {
				break
			} else {
				continue
			}
		}
	}
	ticker.Stop()
	terminateSignal <- true

	handleOutput(p, c, world, p.Turns)

	aliveCells := make([]util.Cell, p.ImageHeight*p.ImageWidth)
	_, aliveCells = calculateAliveCells(p, world)
	report := FinalTurnComplete{
		CompletedTurns: p.Turns,
		Alive:          aliveCells,
	}
	c.events <- report
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{currentTurn, Quitting}
	close(c.events)
}
