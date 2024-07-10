package main

import (
	// "fmt"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
)

/*
	1. Any live cell with fewer than two live neighbours dies, as if by underpopulation.
	2. Any live cell with two or three live neighbours lives on to the next generation.
	3. Any live cell with more than three live neighbours dies, as if by overpopulation.
	4. Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
*/

var defStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

func main() {
	quit := make(chan bool)

	s, err := tcell.NewScreen()
	if err != nil {
		log.Panic(err)
	}

	if err := s.Init(); err != nil {
		log.Panic(err)
	}

	s.SetStyle(defStyle)
	s.Clear()

	defer cleanup(s)

	go func() {
		for {
			ev := s.PollEvent()

			switch ev := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventKey:
				if ev.Rune() == 'q' {
					quit <- true
					return
				}
			}
		}
	}()

	w, h := s.Size()
	g := NewGame(w, h)

	// Main loop
	for {
		g.Tick()
		g.LoopCells(func(x, y int, v bool) {
			DrawCell(s, x, y, v)
		})

		s.Show()

		select {
		case <-time.After(25 * time.Millisecond):
		case <-quit:
			return
		}
	}
}

func DrawCell(s tcell.Screen, x, y int, state bool) {
	var ch rune
	if state {
		ch = 0x2588
	}

	s.SetContent(x, y, ch, nil, defStyle)
}

func cleanup(s tcell.Screen) {
	maybePanic := recover()
	s.Fini()
	if maybePanic != nil {
		panic(maybePanic)
	}
}
