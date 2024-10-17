package gol

import (
	"uk.ac.bris.cs/gameoflife/util"
)

func modulo(a, b int) int {
	return (a%b + b) % b
}

func countNeighbors(height, width int, world [][]byte, y int, x int) int {
	numNeighbors := 0
	neighbors := []byte{
		world[modulo(y+1, height)][modulo(x, width)],
		world[modulo(y+1, height)][modulo(x+1, width)],
		world[modulo(y+1, height)][modulo(x-1, width)],
		world[modulo(y-1, height)][modulo(x, width)],
		world[modulo(y-1, height)][modulo(x+1, width)],
		world[modulo(y-1, height)][modulo(x-1, width)],
		world[modulo(y, height)][modulo(x+1, width)],
		world[modulo(y, height)][modulo(x-1, width)],
	}
	for i := 0; i < 8; i++ {
		if neighbors[i] == 255 {
			numNeighbors++
		}
	}
	return numNeighbors
}

func computeNextState(height, width, startY, endY int, world [][]byte) ([][]byte, []util.Cell) {
	newWorld := make([][]byte, endY-startY)
	changedCells := make([]util.Cell, height, width)
	for i := 0; i < endY-startY; i++ {
		newWorld[i] = make([]byte, len(world[0]))
	}
	for y := 0; y < endY-startY; y++ {
		for x := 0; x < width; x++ {
			numNeighbors := countNeighbors(height, width, world, startY+y, x)
			if world[startY+y][x] == 255 {

				if numNeighbors < 2 {
					newWorld[y][x] = 0
					changedCells = append(changedCells, util.Cell{X: x, Y: startY + y})

				} else if numNeighbors == 2 || numNeighbors == 3 {
					newWorld[y][x] = 255

				} else if numNeighbors > 3 {
					newWorld[y][x] = 0

					changedCells = append(changedCells, util.Cell{X: x, Y: startY + y})
				}

			} else if world[startY+y][x] == 0 && numNeighbors == 3 {
				newWorld[y][x] = 255

				changedCells = append(changedCells, util.Cell{X: x, Y: startY + y})
			}
		}
	}
	return newWorld, changedCells
}

func getAliveCells(p Params, world [][]byte) (int, []util.Cell) {
	var activeCells []util.Cell
	count := 0
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			if world[y][x] == 255 {
				count++
				activeCells = append(activeCells, util.Cell{X: x, Y: y})
			}
		}
	}
	return count, activeCells
}
