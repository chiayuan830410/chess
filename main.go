package main

import (
	"fmt"
)

var ChessLibrary ChessDB
var maxCount = 0
var stack = 100

func Steps(moment Moment, count int) (noResult, redWin, blackWin int) {
	var nextMomentResults []NextMomentResult
	if count == 0 { // no result
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

		return 1, 0, 0
	}

	if ChessLibrary.CheckBoard2Redis(moment.Hash()) { // dynamic
		result := ChessLibrary.GetBoard2Redis(moment.Hash())
		for _, r := range result.Next {
			noResult = noResult + r.NoResult
			redWin = redWin + r.NoResult
			blackWin = blackWin + r.BlackWin
		}
		return noResult, redWin, blackWin
	}

	pieces := moment.GetActionPiece()

	for indexPieces, p := range pieces {
		if count > stack/2 {
			for t := 0; t < stack-count; t++ {
				fmt.Print("\t")
			}
			fmt.Println((indexPieces+1)*100/len(pieces), "% ", indexPieces+1, len(pieces), p.Piece, p.Color)
		}

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
							BlackWin: 1,
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

				nextMomentResults = append(nextMomentResults, NextMomentResult{
					Hash:     nextMoment.Hash(),
					NoResult: nr,
					RedWin:   rw,
					BlackWin: bw,
				})
				noResult = noResult + nr
				redWin = redWin + rw
				blackWin = blackWin + bw
			}
		}
	}
	ChessLibrary.SetBoard2Redis(moment.Hash(), MomentResult{
		Moment: moment,
		Next:   nextMomentResults,
	})

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
	fmt.Println(Steps(moment, stack))

	// for {
	// 	h := ""
	// 	fmt.Scanf("%s", &h)
	// 	fmt.Println(h)
	// 	r := ChessLibrary.GetBoard2Redis(h)
	// 	fmt.Println(r.Moment.DisplayStringMoment())
	// 	maxblack := 0
	// 	bi := 0
	// 	maxred := 0
	// 	ri := 0
	// 	maxno := 0
	// 	ni := 0
	// 	for i, n := range r.Next {
	// 		fmt.Println(n)
	// 		if n.BlackWin > maxblack {
	// 			maxblack = n.BlackWin
	// 			bi = i
	// 		}
	// 		if n.RedWin > maxred {
	// 			maxred = n.RedWin
	// 			ri = i
	// 		}
	// 		if n.NoResult > maxno {
	// 			maxno = n.NoResult
	// 			ni = i
	// 		}
	// 	}
	// 	fmt.Println(maxblack, maxred)

	// 	if maxno != 0 {
	// 		fmt.Println("no win")
	// 		fmt.Println(r.Next[ni].Hash)
	// 		bwin := ChessLibrary.GetBoard2Redis(r.Next[ni].Hash)
	// 		fmt.Println(bwin.DisplayStringMoment())
	// 	}
	// 	if maxblack != 0 {
	// 		fmt.Println("black win")
	// 		bwin := ChessLibrary.GetBoard2Redis(r.Next[bi].Hash)
	// 		fmt.Println(bwin.DisplayStringMoment())
	// 	}

	// 	if maxred != 0 {
	// 		fmt.Println("red win")
	// 		rwin := ChessLibrary.GetBoard2Redis(r.Next[ri].Hash)
	// 		fmt.Println(rwin.DisplayStringMoment())
	// 	}
	// }
}
