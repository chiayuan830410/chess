package main

import (
	"fmt"
)

var ChessLibrary ChessDB
var maxCount = 0

// func Steps(moment Moment, count int) {

// }

var library map[uint64]Moment

func main() {
	library = make(map[uint64]Moment)

	fmt.Println(display(initBoard))

	ChessLibrary.New()

	display(initBoard)
	moment := Moment{
		Board:  initBoard,
		Action: red,
	}
	fmt.Println(moment.Walk(moment.Board[9][0], UP, 1))
	fmt.Println(moment.Action, "\n", display(moment.Board))

	fmt.Println(moment.Walk(moment.Board[0][0], DOWN, 1))
	fmt.Println(moment.Action, "\n", display(moment.Board))

	fmt.Println(moment.Walk(moment.Board[8][0], RIGHT, 6))
	fmt.Println(moment.Action, "\n", display(moment.Board))

	fmt.Println(moment.Walk(moment.Board[1][0], DOWN, 1))
	fmt.Println(moment.Action, "\n", display(moment.Board))

	fmt.Println(moment.Walk(moment.Board[9][7], LEFTUPUP, 0))
	fmt.Println(moment.Action, "\n", display(moment.Board))

	fmt.Println(moment.Walk(moment.Board[2][7], DOWN, 4))
	fmt.Println(moment.Action, "\n", display(moment.Board))

	fmt.Println(moment.Walk(moment.Board[7][7], JUMPUP, 0))
	fmt.Println(moment.Action, "\n", display(moment.Board))

	fmt.Println(moment.Walk(moment.Board[6][7], JUMPLEFT, 4))
	fmt.Println(moment.Action, "\n", display(moment.Board))

}
