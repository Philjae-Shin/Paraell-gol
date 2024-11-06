package gol

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"uk.ac.bris.cs/gameoflife/util"
)

// Remove duplicate definitions that are already in event.go and io.go
// We can use the types and constants defined in event.go and io.go directly

// Constants for action commands
const (
	Save = iota
	Quit
	Pause
	UnPause
)

// DistributorChannels struct with synchronization primitives
type DistributorChannels struct {
	// Event handling
	eventMu   sync.Mutex
	events    []Event
	eventCond *sync.Cond

	// IO handling
	ioMu         sync.Mutex
	ioCond       *sync.Cond
	ioCommand    ioCommand
	ioCommandSet bool
	ioFilename   string
	ioOutput     []uint8
	ioInput      []uint8
	ioIdle       bool
	ioIdleCond   *sync.Cond
}

// NewDistributorChannels initializes a new DistributorChannels
func NewDistributorChannels() *DistributorChannels {
	c := &DistributorChannels{
		ioIdle: true,
	}
	c.eventCond = sync.NewCond(&c.eventMu)
	c.ioCond = sync.NewCond(&c.ioMu)
	c.ioIdleCond = sync.NewCond(&c.ioMu)
	return c
}

// Methods for DistributorChannels

func (c *DistributorChannels) AddEvent(event Event) {
	c.eventMu.Lock()
	c.events = append(c.events, event)
	c.eventCond.Broadcast()
	c.eventMu.Unlock()
}

func (c *DistributorChannels) GetEvents() []Event {
	c.eventMu.Lock()
	defer c.eventMu.Unlock()
	eventsCopy := make([]Event, len(c.events))
	copy(eventsCopy, c.events)
	return eventsCopy
}

func (c *DistributorChannels) SetIoCommand(cmd ioCommand) {
	c.ioMu.Lock()
	c.ioCommand = cmd
	c.ioCommandSet = true
	c.ioCond.Signal()
	c.ioMu.Unlock()
}

func (c *DistributorChannels) GetIoCommand() ioCommand {
	c.ioMu.Lock()
	for !c.ioCommandSet {
		c.ioCond.Wait()
	}
	cmd := c.ioCommand
	c.ioCommandSet = false
	c.ioMu.Unlock()
	return cmd
}

func (c *DistributorChannels) SetIoFilename(filename string) {
	c.ioMu.Lock()
	c.ioFilename = filename
	c.ioMu.Unlock()
}

func (c *DistributorChannels) GetIoFilename() string {
	c.ioMu.Lock()
	filename := c.ioFilename
	c.ioMu.Unlock()
	return filename
}

func (c *DistributorChannels) AddIoOutput(value uint8) {
	c.ioMu.Lock()
	c.ioOutput = append(c.ioOutput, value)
	c.ioMu.Unlock()
}

func (c *DistributorChannels) GetIoOutput() []uint8 {
	c.ioMu.Lock()
	output := c.ioOutput
	c.ioOutput = nil
	c.ioMu.Unlock()
	return output
}

func (c *DistributorChannels) AddIoInput(data []uint8) {
	c.ioMu.Lock()
	c.ioInput = data
	c.ioMu.Unlock()
}

func (c *DistributorChannels) GetIoInput() uint8 {
	c.ioMu.Lock()
	if len(c.ioInput) == 0 {
		c.ioMu.Unlock()
		return 0
	}
	value := c.ioInput[0]
	c.ioInput = c.ioInput[1:]
	c.ioMu.Unlock()
	return value
}

func (c *DistributorChannels) SetIoIdle(idle bool) {
	c.ioMu.Lock()
	c.ioIdle = idle
	c.ioIdleCond.Broadcast()
	c.ioMu.Unlock()
}

func (c *DistributorChannels) GetIoIdle() bool {
	c.ioMu.Lock()
	idle := c.ioIdle
	c.ioMu.Unlock()
	return idle
}

func (c *DistributorChannels) WaitIoIdle() {
	c.ioMu.Lock()
	for !c.ioIdle {
		c.ioIdleCond.Wait()
	}
	c.ioMu.Unlock()
}

// IOProcessor function
func IOProcessor(c *DistributorChannels, p Params) {
	for {
		cmd := c.GetIoCommand()
		switch cmd {
		case ioInput:
			filename := c.GetIoFilename()
			data, err := readPgmImage(filename, p.ImageHeight, p.ImageWidth)
			if err != nil {
				fmt.Println("Error reading file:", err)
				c.SetIoIdle(true)
				continue
			}
			c.AddIoInput(data)
			c.SetIoIdle(true)
		case ioOutput:
			filename := c.GetIoFilename()
			data := c.GetIoOutput()
			err := writePgmImage(filename, data, p.ImageHeight, p.ImageWidth)
			if err != nil {
				fmt.Println("Error writing file:", err)
			}
			c.SetIoIdle(true)
		case ioCheckIdle:
			c.SetIoIdle(true)
		default:
			// Handle other commands if needed
		}
	}
}

// Use the readPgmImage and writePgmImage functions from io.go
// Adjusted to work with DistributorChannels
func readPgmImage(filename string, imageHeight, imageWidth int) ([]uint8, error) {
	// Implement the PGM file reading here
	// For this example, we'll simulate reading by returning a slice of zeros
	data := make([]uint8, imageHeight*imageWidth)
	return data, nil
}

func writePgmImage(filename string, data []uint8, imageHeight, imageWidth int) error {
	// Implement the PGM file writing here
	// For this example, we'll simulate writing by doing nothing
	return nil
}

// Handle output
func handleOutput(p Params, c *DistributorChannels, world [][]uint8, t int) {
	c.SetIoCommand(ioOutput)
	outFilename := strconv.Itoa(p.ImageHeight) + "x" + strconv.Itoa(p.ImageWidth) + "x" + strconv.Itoa(t)
	c.SetIoFilename(outFilename)
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			c.AddIoOutput(world[y][x])
		}
	}
	c.SetIoIdle(false)
	// Wait for IO to finish
	c.WaitIoIdle()

	c.AddEvent(ImageOutputComplete{
		CompletedTurns: t,
		Filename:       outFilename,
	})
}

// Handle input
func handleInput(p Params, c *DistributorChannels, world [][]uint8) [][]uint8 {
	filename := strconv.Itoa(p.ImageHeight) + "x" + strconv.Itoa(p.ImageWidth)
	c.SetIoCommand(ioInput)
	c.SetIoFilename(filename)
	c.SetIoIdle(false)
	// Wait for IO to finish
	c.WaitIoIdle()
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			num := c.GetIoInput()
			world[y][x] = num
			if num == 255 {
				c.AddEvent(CellFlipped{
					CompletedTurns: 0,
					Cell:           util.Cell{X: x, Y: y},
				})
			}
		}
	}
	return world
}

// Handle key presses
func handleKeyPress(p Params, c *DistributorChannels, keyPresses <-chan rune, action *int, actionMu *sync.Mutex, actionCond *sync.Cond) {
	paused := false
	for {
		input := <-keyPresses
		switch input {
		case 's':
			actionMu.Lock()
			*action = Save
			actionMu.Unlock()
			actionCond.Signal()
		case 'q':
			actionMu.Lock()
			*action = Quit
			actionMu.Unlock()
			actionCond.Signal()
			return
		case 'p':
			actionMu.Lock()
			if paused {
				*action = UnPause
				paused = false
			} else {
				*action = Pause
				paused = true
			}
			actionMu.Unlock()
			actionCond.Signal()
		}
	}
}

// Distributor function
func Distributor(p Params, c *DistributorChannels, keyPresses <-chan rune) {
	world := make([][]uint8, p.ImageHeight)
	prevWorld := make([][]uint8, p.ImageHeight)
	for i := range world {
		world[i] = make([]uint8, p.ImageWidth)
		prevWorld[i] = make([]uint8, p.ImageWidth)
	}

	world = handleInput(p, c, world)

	turn := 0
	ticker := time.NewTicker(2 * time.Second)
	done := make(chan bool)
	pause := false
	quit := false
	finished := false

	var mu sync.Mutex

	action := 0
	actionMu := &sync.Mutex{}
	actionCond := sync.NewCond(actionMu)

	// Start IOProcessor goroutine
	go IOProcessor(c, p)

	// Start handling key presses
	go handleKeyPress(p, c, keyPresses, &action, actionMu, actionCond)

	// Send StateChange event indicating Executing state at the start
	c.AddEvent(StateChange{CompletedTurns: turn, NewState: Executing})

	// Start ticker for AliveCellsCount events
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				mu.Lock()
				snapshot := make([][]uint8, p.ImageHeight)
				for i := range prevWorld {
					snapshot[i] = make([]uint8, p.ImageWidth)
					copy(snapshot[i], prevWorld[i])
				}
				currentTurn := turn
				mu.Unlock()
				aliveCount, _ := calculateAliveCells(p, snapshot)
				c.AddEvent(AliveCellsCount{
					CompletedTurns: currentTurn,
					CellsCount:     aliveCount,
				})
			}
		}
	}()

	for !finished && (turn < p.Turns || pause) {
		if pause {
			actionMu.Lock()
			for action == 0 {
				actionCond.Wait()
			}
			command := action
			action = 0
			actionMu.Unlock()

			switch command {
			case UnPause:
				pause = false
				c.AddEvent(StateChange{CompletedTurns: turn, NewState: Executing})
			case Quit:
				quit = true
				finished = true
			case Save:
				mu.Lock()
				snapshot := make([][]uint8, p.ImageHeight)
				for i := range world {
					snapshot[i] = make([]uint8, p.ImageWidth)
					copy(snapshot[i], world[i])
				}
				currentTurn := turn
				mu.Unlock()
				handleOutput(p, c, snapshot, currentTurn)
			}
			continue
		}

		select {
		default:
			// Check for any actions
			actionMu.Lock()
			if action != 0 {
				command := action
				action = 0
				actionMu.Unlock()

				switch command {
				case Pause:
					pause = true
					c.AddEvent(StateChange{CompletedTurns: turn, NewState: Paused})
				case Quit:
					quit = true
					finished = true
				case Save:
					mu.Lock()
					snapshot := make([][]uint8, p.ImageHeight)
					for i := range world {
						snapshot[i] = make([]uint8, p.ImageWidth)
						copy(snapshot[i], world[i])
					}
					currentTurn := turn
					mu.Unlock()
					handleOutput(p, c, snapshot, currentTurn)
				}
			} else {
				actionMu.Unlock()
			}

			// Update game state
			if !quit && turn < p.Turns {
				mu.Lock()
				for i := range world {
					copy(prevWorld[i], world[i])
				}
				mu.Unlock()
				var flipFragment []util.Cell
				world, flipFragment = calculateNextState(p.ImageHeight, p.ImageWidth, 0, p.ImageHeight, prevWorld)
				c.AddEvent(CellsFlipped{
					CompletedTurns: turn,
					Cells:          flipFragment,
				})
				mu.Lock()
				turn++
				mu.Unlock()
				c.AddEvent(TurnComplete{CompletedTurns: turn})
			} else if quit || turn >= p.Turns {
				finished = true
			}
		}
	}

	ticker.Stop()
	done <- true

	_, aliveCells := calculateAliveCells(p, world)

	handleOutput(p, c, world, turn)

	if quit {
		c.AddEvent(StateChange{CompletedTurns: turn, NewState: Quitting})
	}
	c.AddEvent(FinalTurnComplete{
		CompletedTurns: turn,
		Alive:          aliveCells,
	})

	// Wait for IO to finish
	c.SetIoCommand(ioCheckIdle)
	c.WaitIoIdle()
}

// Utility functions for Game of Life calculations

func mod(a, b int) int {
	return (a%b + b) % b
}

func calculateNeighbours(height, width int, world [][]uint8, y int, x int) int {
	h := height
	w := width
	noOfNeighbours := 0

	neighbour := []uint8{
		world[mod(y+1, h)][mod(x, w)],
		world[mod(y+1, h)][mod(x+1, w)],
		world[mod(y, h)][mod(x+1, w)],
		world[mod(y-1, h)][mod(x+1, w)],
		world[mod(y-1, h)][mod(x, w)],
		world[mod(y-1, h)][mod(x-1, w)],
		world[mod(y, h)][mod(x-1, w)],
		world[mod(y+1, h)][mod(x-1, w)],
	}

	for i := 0; i < 8; i++ {
		if neighbour[i] == 255 {
			noOfNeighbours++
		}
	}

	return noOfNeighbours
}

func calculateNextState(height, width, startY, endY int, world [][]uint8) ([][]uint8, []util.Cell) {
	newWorld := make([][]uint8, endY-startY)
	var flipCell []util.Cell
	for i := 0; i < endY-startY; i++ {
		newWorld[i] = make([]uint8, width)
	}

	for y := 0; y < endY-startY; y++ {
		for x := 0; x < width; x++ {
			noOfNeighbours := calculateNeighbours(height, width, world, startY+y, x)
			if world[startY+y][x] == 255 {
				if noOfNeighbours < 2 || noOfNeighbours > 3 {
					newWorld[y][x] = 0
					flipCell = append(flipCell, util.Cell{X: x, Y: startY + y})
				} else {
					newWorld[y][x] = 255
				}
			} else if world[startY+y][x] == 0 && noOfNeighbours == 3 {
				newWorld[y][x] = 255
				flipCell = append(flipCell, util.Cell{X: x, Y: startY + y})
			} else {
				newWorld[y][x] = world[startY+y][x]
			}
		}
	}

	return newWorld, flipCell
}

func calculateAliveCells(p Params, world [][]uint8) (int, []util.Cell) {
	var aliveCells []util.Cell
	count := 0
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			if world[y][x] == 255 {
				count++
				aliveCells = append(aliveCells, util.Cell{X: x, Y: y})
			}
		}
	}
	return count, aliveCells
}
