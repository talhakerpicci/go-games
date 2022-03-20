package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const ballSymbol = 0x25CF
const paddleSymbol = 0x2588
const paddleHeight = 4

var directionMap = map[string]string{
	"Up":      "up",
	"Down":    "down",
	"Rune[w]": "up",
	"Rune[s]": "down",
}

type GameObject struct {
	row, col, width, height int
	velRow, velCol          int
	symbol                  rune
}

var screen tcell.Screen
var paddle1 *GameObject
var paddle2 *GameObject
var ball *GameObject

var initialBallVelocityRow = -1
var initialBallVelocityCol = -2

var gameObjects []*GameObject

func main() {
	initScreen()
	initGameState()
	inputChanel := initUserINput()

	for {
		handleUserInput(readUserInput(inputChanel))
		updateState()
		drawState( /* 3, 3, key */ )
		time.Sleep(30 * time.Millisecond)
	}
}

func initGameState() {
	width, height := screen.Size()
	paddleStart := height/2 - paddleHeight/2

	paddle1 = &GameObject{
		row:    paddleStart,
		col:    0,
		width:  1,
		height: paddleHeight,
		symbol: paddleSymbol,
		velRow: 0,
		velCol: 0,
	}

	paddle2 = &GameObject{
		row:    paddleStart,
		col:    width - 1,
		width:  1,
		height: paddleHeight,
		symbol: paddleSymbol,
		velRow: 0,
		velCol: 0,
	}

	ball = &GameObject{
		row:    height / 2,
		col:    width / 2,
		width:  1,
		height: 1,
		symbol: ballSymbol,
		velRow: initialBallVelocityRow,
		velCol: initialBallVelocityCol,
	}

	gameObjects = []*GameObject{
		paddle1, paddle2, ball,
	}
}

func draw(gameObject *GameObject) {
	for r := 0; r < gameObject.height; r++ {
		for c := 0; c < gameObject.width; c++ {
			screen.SetContent(gameObject.col+c, gameObject.row+r, gameObject.symbol, nil, tcell.StyleDefault)
		}
	}
}

func drawState( /* col, row int, key string */ ) {
	screen.Clear()
	for _, object := range gameObjects {
		draw(object)
	}

	screen.Show()
}

func updateState() {
	width, _ := screen.Size()
	ball.row += ball.velRow
	ball.col += ball.velCol

	if checkIfBallAtBoundary(ball) {
		ball.velRow *= -1
	} else if checkIfTouchingPaddle(ball, paddle1) || checkIfTouchingPaddle(ball, paddle2) {
		ball.velCol *= -1
	}

	if ball.col < 0 || ball.col > width {
		fmt.Println("FINISHED")
		screen.Fini()
		os.Exit(0)
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

func handleUserInput(key string) {
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

func checkIfAtBoundary(paddle *GameObject) (bool, bool) {
	var isAtTopBoundary, isAtBottomBoundary bool
	_, height := screen.Size()
	if paddle.row == 0 {
		isAtTopBoundary = true
	} else if paddle.row == height-paddleHeight {
		isAtBottomBoundary = true
	}

	return isAtTopBoundary, isAtBottomBoundary
}

func checkIfBallAtBoundary(obj *GameObject) bool {
	_, screenHeight := screen.Size()
	return obj.row+obj.velRow < 0 || obj.row+obj.velRow >= screenHeight
}

func checkIfTouchingPaddle(ball *GameObject, paddle *GameObject) bool {
	var collidesOnColumn bool
	if ball.col < paddle.col {
		collidesOnColumn = ball.col+ball.velCol >= paddle.col
	} else {
		collidesOnColumn = ball.col+ball.velCol <= paddle.col
	}

	return collidesOnColumn && ball.row >= paddle.row && ball.row < paddle.row+paddle.height
}

func movePaddle(paddle *GameObject, direction string) {
	isAtTopBoundary, isAtBottomBoundary := checkIfAtBoundary(paddle)

	if !isAtTopBoundary && direction == "up" {
		paddle.row--
	} else if !isAtBottomBoundary && direction == "down" {
		paddle.row++
	}
}
