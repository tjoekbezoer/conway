package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
)

type Game struct {
	w, h  int
	Cells World
}
type World [][]bool
type LoopFunc func(x, y int, value bool)

func NewGame(w, h int) *Game {
	return &Game{
		w:     w,
		h:     h,
		Cells: initWorld(w, h),
	}
}

// Tick performs the cell calculations for this generation
func (g *Game) Tick() {
	newWorld := makeWorld(g.w, g.h)

	LoopCells(g.Cells, func(x, y int, v bool) {
		newWorld[y][x] = g.calcCell(x, y)
	})

	g.Cells = newWorld
}

func (g *Game) calcCell(x, y int) bool {
	c := g.cell(x, y)
	n := g.neighbours(x, y)

	switch {
	case c == 1 && (n == 2 || n == 3):
		return true
	case c == 0 && n == 3:
		return true
	default:
		return false
	}
}

func (g *Game) neighbours(x, y int) int {
	return g.cell(x-1, y-1) + g.cell(x, y-1) + g.cell(x+1, y-1) +
		g.cell(x-1, y) + g.cell(x+1, y) +
		g.cell(x-1, y+1) + g.cell(x, y+1) + g.cell(x+1, y+1)
}

func (g *Game) cell(x, y int) int {
	// Wrap around the screen
	if x == -1 {
		x = g.w - 1
	} else if x == g.w {
		x = 0
	}
	if y == -1 {
		y = g.h - 1
	} else if y == g.h {
		y = 0
	}

	if g.Cells[y][x] {
		return 1
	} else {
		return 0
	}
}

func LoopCells(cells World, handler LoopFunc) {
	for y, row := range cells {
		for x, v := range row {
			handler(x, y, v)
		}
	}
}

func readFile(fileName string) (World, int, int) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Panic(err)
		}
	}()

	// Read in the game file, splitting it up into lines
	// for processing.
	data := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err == nil && n > 0 {
			data = append(data, buf[:n]...)
		} else if err != io.EOF {
			log.Panic(err)
		} else {
			break
		}
	}
	data = bytes.Trim(data, "\n")
	lines := strings.Split(string(data), "\n")

	// Build a World data structure from the split string lines.
	h := len(lines)
	w := 0
	res := make(World, h)
	for i := range lines {
		cw := len(lines[i])
		w = max(w, cw)
		res[i] = make([]bool, cw)

		for j := range lines[i] {
			cell := false
			if lines[i][j] != ' ' {
				cell = true
			}
			res[i][j] = cell
		}
	}

	return res, w, h
}

func initWorld(w, h int) [][]bool {
	start, gw, gh := readFile("game.txt")
	cw := (w - gw) / 2
	ch := (h - gh) / 2
	world := makeWorld(w, h)

	// Copy the game from file to the in-memory world, centering
	// the starting game in the process.
	LoopCells(world, func(x, y int, _ bool) {
		if y >= ch && y-ch < len(start) && x >= cw && x-cw < len(start[y-ch]) {
			world[y][x] = start[y-ch][x-cw]
		}
	})

	return world
}

func makeWorld(w, h int) [][]bool {
	cells := make([][]bool, h)

	for i := range cells {
		cells[i] = make([]bool, w)
	}

	return cells
}
