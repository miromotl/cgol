// Implementing Conway's Game Of Life
// ----------------------------------
//
// Using a map for storing the current state of the world.
//
// We are printing the successive populations in a format that can be fed
// to gnuplot and creating in this way an animated view of the population.
//
// This is just an exercise for using maps in go! Do not take this
// too serious...
//
// To see the simulation in gnuplot, call the program like this:
// ./gol.exe | gnuplot --persist

package main

import (
	"flag"
	"fmt"
	"strings"
	"strconv"
	"os"
	"math/rand"
)

// We are storing the cells (alive or dead) in a map. The keys are the Cartesian
// coordinates of the cells and the values are the properties of the cells,
// namely their state and number of alive neighbours.

// A cell has its state, and its number of life neighbours
type Cell struct {
	alive bool
	n int
}

// The coordinates are plain 2-d cartesian coordinates
type Coord struct {
	x int
	y int
}

// The world is a map of Coord and Cell
type World map[Coord]Cell

// Inflate inflates the world with dead cells surrounding
// the live cells
func (world World) Inflate() World {
	var newWorld World
	newWorld = make(World)
	
	for coord, cell := range world {
		newWorld[coord] = cell
		for i := -1; i < 2; i++ {
			for j := -1; j < 2; j++ {
				c := Coord{coord.x + i, coord.y + j}
				if _, found := newWorld[c]; !found {
					newWorld[c] = Cell{false, 0}
				}
			}
		}
	}

	return newWorld
}

// Deflate deflates the world: only the live cells remain
func (world World) Deflate() World {
	var newWorld World
	newWorld = make(World)
	
	for coord, cell := range world {
		if cell.alive {
			newWorld[coord] = cell
		}
	}

	return newWorld
}

// CountLiveNeighbours counts for each cell in the world its neighbouring
// alive cells and updates its counter
func (world World) CountLiveNeighbours() World {
	var newWorld World
	newWorld = make(World)
	
	for coord, cell := range world {
		n := 0
		for i := -1; i < 2; i++ {
			for j := -1; j < 2; j++ {
				c := Coord{coord.x + i, coord.y + j}
				if (i != 0 || j != 0) && world[c].alive {
					n = n+1
				}
			}
		}
		newWorld[coord] = Cell{cell.alive, n}
	}
	
	return newWorld
}

// ApplyRules applies the rules to each cell of the world. This determines
// the fate of the cell for the next tick.
func (world World) ApplyRules() World {
	var newWorld World
	newWorld = make(World)

	// apply the rules of the game to each cell
	for coord, cell := range world {
		if cell.alive {
			if 1 < cell.n && cell.n < 4 {
				newWorld[coord] = Cell{true, 0}
			}
		} else {
			if cell.n == 3 {
				newWorld[coord] = Cell{true, 0}
			}
		}
	}

	return newWorld
}

// Tick computes the next generation of live cells in the world
func (world World) Tick() World {
	return world.Inflate().CountLiveNeighbours().ApplyRules().Deflate()
}

// gnuplotHeader prints the header for gnuplot
func gnuplotHeader(d int) {
	fmt.Printf("unset key; set xrange[-%[1]d:%[1]d]\n", d/2)
	fmt.Printf("set yrange[-%[1]d:%[1]d]\n", d/2)
	fmt.Println("set style line 1 lc rgb '#0060ad' pt 7")
}

// gnuplotWorld prints the coordinates of the cells in the world
func gnuplotWorld(world World) {
	fmt.Println("plot '-' with points ls 1")

	for coord := range world {
		fmt.Printf("%d, %d\n", coord.x, coord.y)
	}
	
	fmt.Println("e")
}

func main() {
	ticks, size, pattern := handleCommandLine()
	// The world
	var world World
	world = make(World)

	for _, coord := range pattern {
		world[coord] = Cell{true, 0}
	}
	
	gnuplotHeader(size)

	gnuplotWorld(world)
	
	for i := 0; i < ticks; i++ {
		world = world.Inflate().Tick().Deflate()
		gnuplotWorld(world)
	}
}

func handleCommandLine() (ticks, size int, pattern []Coord) {
	// Define our own usage message, overwriting the default one
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage: cgol [flags] [pattern] | gnuplot --persist\n")
		flag.PrintDefaults()
	}

	// Define the command line flags
	flag.IntVar(&ticks, "ticks", 10, "number of iterations running the game")
	flag.IntVar(&size, "size", 50, "size of the visible world in x and y direction")
	var random *bool = flag.Bool("random", false, "generate a random pattern to start with")
	var coordinatesOpt *string = flag.String("coordinates", "1,0;0,1;1,1;1,2;2,2", "semi-colon-separated list of coordinates")
	flag.Parse()
	
	// Create a ranodm starting pattern or use the r-pentomino pattern
	if *random {
		// Generate a random pattern
		//pattern = make([]Coord, size*size/5, size*size)
		pattern = []Coord{}
		rand.Seed(42)
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				if rand.Intn(100) < 20 {
					pattern = append(pattern, Coord{i, j})
				}
			}
		}
	} else {
		coordinates := strings.Split(*coordinatesOpt, ";")
		pattern = make([]Coord, len(coordinates))
		for idx := range coordinates {
			xy := strings.Split(coordinates[idx], ",")
			x, err := strconv.Atoi(xy[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			y, err := strconv.Atoi(xy[1])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			pattern[idx] = Coord{x, y}
		}
	}
	
	return ticks, size, pattern
}
