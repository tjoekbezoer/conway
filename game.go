package main

var data = `
 x
xx
 xx
`

type Game struct {
	w, h  int
	Cells [][]bool
}

type LoopFunc func(x, y int, value bool)

func NewGame(w, h int) *Game {
	// state := convertFile(data)

	return &Game{
		w:     w,
		h:     h,
		Cells: makeWorld(w, h),
	}
}

func (g *Game) Tick() {
	newWorld := makeWorld(g.w, g.h)

	g.LoopCells(func(x, y int, v bool) {
		newWorld[y][x] = g.calcCell(x, y)
	})

	g.Cells = newWorld
}

func (g *Game) LoopCells(handler LoopFunc) {
	for y, row := range g.Cells {
		for x, v := range row {
			handler(x, y, v)
		}
	}
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

// func convertFile(data string) ([][]bool, int, int) {
// 	lines := strings.Split(data, "\n")
// 	h := len(lines)
//
// 	res := make([][]byte)
// 	for i := range lines {
// 		res.
// 	}
// }

func makeWorld(w, h int) [][]bool {
	cells := make([][]bool, h)

	for i := range cells {
		cells[i] = make([]bool, w)
	}

	// temp creation
	cells[25+0][25+1] = true
	cells[25+1][25+0] = true
	cells[25+1][25+1] = true
	cells[25+2][25+1] = true
	cells[25+2][25+2] = true

	return cells
}
