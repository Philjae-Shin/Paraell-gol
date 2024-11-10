package gol

import (
	"uk.ac.bris.cs/gameoflife/util"
)

// CPU profiling (Before)

//func mod(a, b int) int {
//	return (a%b + b) % b
//}
//
//func calculateNeighbours(height, width int, world [][]byte, y int, x int) int {
//
//	h := height
//	w := width
//	noOfNeighbours := 0
//
//	neighbour := []byte{
//		world[mod(y+1, h)][mod(x, w)],
//		world[mod(y+1, h)][mod(x+1, w)],
//		world[mod(y, h)][mod(x+1, w)],
//		world[mod(y-1, h)][mod(x+1, w)],
//		world[mod(y-1, h)][mod(x, w)],
//		world[mod(y-1, h)][mod(x-1, w)],
//		world[mod(y, h)][mod(x-1, w)],
//		world[mod(y+1, h)][mod(x-1, w)],
//	}
//
//	for i := 0; i < 8; i++ {
//		if neighbour[i] == 255 {
//			noOfNeighbours++
//		}
//	}
//
//	return noOfNeighbours
//}

// CPU profiling (After)
func calculateNeighbours(height, width int, world [][]byte, y int, x int) int {

	h := height
	w := width
	noOfNeighbours := 0

	// Pre-calculate and store the coordinates of adjacent cells
	neighbors := [8][2]int{
		{y + 1, x},     // Below
		{y + 1, x + 1}, // Bottom-right
		{y, x + 1},     // Right
		{y - 1, x + 1}, // Top-right
		{y - 1, x},     // Above
		{y - 1, x - 1}, // Top-left
		{y, x - 1},     // Left
		{y + 1, x - 1}, // Bottom-left
	}

	for _, coord := range neighbors {
		ny, nx := coord[0], coord[1]

		// Handle Y-coordinate wrap-around
		if ny < 0 {
			ny = h - 1
		} else if ny >= h {
			ny = 0
		}

		// Handle X-coordinate wrap-around
		if nx < 0 {
			nx = w - 1
		} else if nx >= w {
			nx = 0
		}

		if world[ny][nx] == 255 {
			noOfNeighbours++
		}
	}

	return noOfNeighbours
}

func calculateNextState(height, width, startY, endY int, world [][]byte) ([][]byte, []util.Cell) {

	newWorld := make([][]byte, endY-startY)
	flipCell := make([]util.Cell, height, width)
	for i := 0; i < endY-startY; i++ {
		newWorld[i] = make([]byte, len(world[0]))
		// copy(newWorld[i], world[startY+i])
	}

	for y := 0; y < endY-startY; y++ {
		for x := 0; x < width; x++ {
			noOfNeighbours := calculateNeighbours(height, width, world, startY+y, x)
			if world[startY+y][x] == 255 {
				if noOfNeighbours < 2 {
					newWorld[y][x] = 0
					flipCell = append(flipCell, util.Cell{X: x, Y: startY + y})
				} else if noOfNeighbours == 2 || noOfNeighbours == 3 {
					newWorld[y][x] = 255
				} else if noOfNeighbours > 3 {
					newWorld[y][x] = 0
					flipCell = append(flipCell, util.Cell{X: x, Y: startY + y})
				}
			} else if world[startY+y][x] == 0 && noOfNeighbours == 3 {
				newWorld[y][x] = 255
				flipCell = append(flipCell, util.Cell{X: x, Y: startY + y})
			}
		}
	}

	return newWorld, flipCell
}

func calculateAliveCells(p Params, world [][]byte) (int, []util.Cell) {

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
