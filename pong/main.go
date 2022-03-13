package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const paddleSymbol = 0x2588
const paddleHeight = 4

var directionMap = map[string]string{
	"Up":      "up",
	"Down":    "down",
	"Rune[w]": "up",
	"Rune[s]": "down",
}

type Paddle struct {
	row, col, width, height int
}

var screen tcell.Screen
var paddle1 *Paddle
var paddle2 *Paddle

func draw(row, col, width, height int) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, paddleSymbol, nil, tcell.StyleDefault)
		}
	}
}

func initGameState() {
	width, height := screen.Size()
	paddleStart := height/2 - paddleHeight/2

	paddle1 = &Paddle{
		row:    paddleStart,
		col:    0,
		width:  1,
		height: paddleHeight,
	}

	paddle2 = &Paddle{
		row:    paddleStart,
		col:    width - 1,
		width:  1,
		height: paddleHeight,
	}
}

func drawState( /* col, row int, key string */ ) {
	screen.Clear()
	draw(paddle1.row, paddle1.col, paddle1.width, paddle1.height)
	draw(paddle2.row, paddle2.col, paddle2.width, paddle2.height)

	/* for _, c := range key {
		screen.SetContent(col, row, c, nil, tcell.StyleDefault)
		col++
	} */

	screen.Show()
}

func main() {
	initScreen()
	initGameState()
	inputChanel := initUserINput()

	for {
		drawState( /* 3, 3, key */ )
		time.Sleep(10 * time.Millisecond)

		key := readUserInput(inputChanel)
		switch key {
		case "Rune[q]":
			screen.Fini()
			os.Exit(0)
		case "Up":
			movePaddle(paddle2, directionMap[key])
		case "Down":
			movePaddle(paddle2, directionMap[key])
		case "Rune[w]":
			movePaddle(paddle1, directionMap[key])
		case "Rune[s]":
			movePaddle(paddle1, directionMap[key])
		}
	}
}

func initScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
}

func initUserINput() chan string {
	inputChannel := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventResize:
				screen.Sync()
				/* initGameState() */
			case *tcell.EventKey:
				inputChannel <- ev.Name()
			}
		}
	}()

	return inputChannel
}

func readUserInput(inputChannel chan string) string {
	var key string
	select {
	case key = <-inputChannel:
	default:
		key = ""

	}

	return key
}

func checkIfAtBoundary(paddle *Paddle) (bool, bool) {
	var isAtTopBoundary, isAtBottomBoundary bool
	_, height := screen.Size()
	if paddle.row == 0 {
		isAtTopBoundary = true
	} else if paddle.row == height-paddleHeight {
		isAtBottomBoundary = true
	}

	return isAtTopBoundary, isAtBottomBoundary
}

func movePaddle(paddle *Paddle, direction string) {
	isAtTopBoundary, isAtBottomBoundary := checkIfAtBoundary(paddle)

	if !isAtTopBoundary && direction == "up" {
		paddle.row--
	} else if !isAtBottomBoundary && direction == "down" {
		paddle.row++
	}
}
