package main

// IMPORTANT!
// all funcitons uses y, x order ( not like x, y in mathematics)
// all indexation starting from 0

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func makeGrid(height, width int, allocation float64) []byte {
	grid := make([]byte, height*width/8)

	for i := 0; i < len(grid); i++ {
		total := 0
		border := int(allocation * 255)

		for j := 0; j < 8; j++ {
			randNum := rand.Intn(256)
			switch {
			case randNum <= border:
				total += 1 << j
			default:
			}
		}
		grid[i] = byte(total)
	}
	return grid
}

var (
	height        int          = 48
	width         int          = 96
	allocation    float64      = 0.3
	grid          []byte       = makeGrid(height, width, allocation)
	filled, empty string       = "██", "  "
	birth         map[int]bool = map[int]bool{3: true}
	alive         map[int]bool = map[int]bool{2: true, 3: true}
	duration      int          = 10
)

func indexToBitValue(index int) int {
	outB, inB := index/8, index%8
	val := (grid[outB] >> (7 - inB)) & 1
	return int(val)
}

func indexToCoord(index int) (int, int) {
	return index / width, index % width
}

func coordToIndex(coord [2]int) int {
	return coord[0]*width + coord[1]
}

func getSurroundingCellsCount(index int) (count int) {
	y_init, x_init := indexToCoord(index)

	for y := y_init - 1; y < y_init+2; y++ {
		for x := x_init - 1; x < x_init+2; x++ {
			if (x == x_init && y == y_init) || (y < 0 || y > height-1 || x < 0 || x > width-1) {
				continue
			}

			count += indexToBitValue(coordToIndex([2]int{y, x}))

		}
	}
	return count
}

func render() int {
	cur := 0
	for i := 0; i < height; i++ {

		for j := 0; j < width; j++ {
			if indexToBitValue(cur) == 1 {
				fmt.Printf(filled)
			} else {
				fmt.Printf(empty)
			}
			cur++
		}
		fmt.Print("\n")
	}
	return cur
}

func update() {
	buff := make([]byte, len(grid))

	for pos := 0; pos < len(grid); pos++ {
		buffByte := byte(0)
		for j := 0; j < 8; j++ {
			curIndex := pos*8 + j
			curValue := indexToBitValue(curIndex)
			neighbors := getSurroundingCellsCount(curIndex)
			switch {
			case curValue == 1 && alive[neighbors]:
				buffByte += 1 << (7 - j)
			case curValue == 0 && birth[neighbors]:
				buffByte += 1 << (7 - j)
			default:
			}

		}
		buff[pos] = buffByte
	}
	copy(grid, buff)
	clear(buff)

}

func main() {
	fmt.Println(os.Args)
	if len(os.Args) == 3 && (os.Args[1] == "-d" || os.Args[1] == "--duration") {
		durationTemp, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("Invalid arguments")
		}
		duration = durationTemp
	}

	for {
		fmt.Println()
		render()
		fmt.Println()
		time.Sleep(time.Duration(duration) * time.Millisecond)
		fmt.Print("\033[H\033[2J") // clear terminal
		update()
	}
}
