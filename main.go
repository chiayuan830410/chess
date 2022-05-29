package main

import (
	"fmt"
)

var ChessLibrary ChessDB
var maxCount = 0

func Steps(moment Moment, count int) (noResult, redWin, blackWin int) {
	if count == 0 { // no result
		return 1, 0, 0
	}

	pieces := moment.GetAllPiece()

	for _, p := range pieces {
		if p.Color == moment.Action {

			var directions, distances []int
			switch p.Piece {
			case 將:
				directions = []int{UP, DOWN, LEFT, RIGHT, KILL}
				distances = []int{1}
			case 士:
				directions = []int{LEFTUP, LEFTDOWN, RIGHTUP, RIGHTDOWN}
				distances = []int{1}
			case 象:
				directions = []int{LEFTUP, LEFTDOWN, RIGHTUP, RIGHTDOWN}
				distances = []int{1}
			case 車:
				directions = []int{UP, DOWN, LEFT, RIGHT}
				distances = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			case 馬:
				directions = []int{LEFTLEFTUP, LEFTLEFTDOWN, RIGHTRIGHTUP, RIGHTRIGHTDOWN, LEFTUPUP, LEFTDOWNDOWN, RIGHTUPUP, RIGHTDOWNDOWN}
				distances = []int{1}
			case 炮:
				directions = []int{UP, DOWN, LEFT, RIGHT, JUMPUP, JUMPDOWN, JUMPLEFT, JUMPRIGHT}
				distances = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			case 卒:
				directions = []int{UP, DOWN, LEFT, RIGHT}
				distances = []int{1}
			}
			for _, direction := range directions {
				for _, distance := range distances {
					nextMoment := moment
					result, err := nextMoment.Walk(p, direction, distance) // walk
					if err != nil {
						continue
					}
					if result == 2 { // red win
						ChessLibrary.SetBoard2Redis(moment.Hash(), MomentResult{
							Moment: moment,
							Next: []NextMomentResult{{
								Hash:     nextMoment.Hash(),
								NoResult: 0,
								RedWin:   1,
								BlackWin: 0,
							},
							},
						})
						ChessLibrary.SetBoard2Redis(nextMoment.Hash(), MomentResult{
							Moment: nextMoment,
							Next: []NextMomentResult{{
								Hash:     "",
								NoResult: 0,
								RedWin:   0,
								BlackWin: 0,
							},
							},
						})
						return 0, 1, 0
					} else if result == 3 { //black win
						ChessLibrary.SetBoard2Redis(moment.Hash(), MomentResult{
							Moment: moment,
							Next: []NextMomentResult{{
								Hash:     "",
								NoResult: 0,
								RedWin:   0,
								BlackWin: 0,
							},
							},
						})
						ChessLibrary.SetBoard2Redis(nextMoment.Hash(), MomentResult{
							Moment: nextMoment,
							Next: []NextMomentResult{{
								Hash:     "",
								NoResult: 0,
								RedWin:   0,
								BlackWin: 0,
							},
							},
						})
						return 0, 0, 1
					}
					nr, rw, bw := Steps(nextMoment, count-1) // to next
					ChessLibrary.SetBoard2Redis(moment.Hash(), MomentResult{
						Moment: moment,
						Next: []NextMomentResult{{
							Hash:     "",
							NoResult: nr,
							RedWin:   rw,
							BlackWin: bw,
						},
						},
					})
					noResult = noResult + nr
					redWin = redWin + rw
					blackWin = blackWin + bw
				}
			}
		}
	}

	return noResult, redWin, blackWin
}

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
	fmt.Println(Steps(moment, 10))
}
