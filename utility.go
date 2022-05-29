package main

import (
	"fmt"

	"github.com/shomali11/util/xhashes"
)

func display(board Board) (displayBoard string) {
	displayBoard = ""
	for y, xbroard := range board {
		displayBoard = displayBoard + fmt.Sprint(y, "\t")
		for _, xyboard := range xbroard {
			if xyboard == nil {
				displayBoard = displayBoard + fmt.Sprint("---\t")
			} else if xyboard.Color == black {
				switch xyboard.Piece {
				case 將:
					displayBoard = displayBoard + fmt.Sprint("黑將\t")
				case 士:
					displayBoard = displayBoard + fmt.Sprint("黑士\t")
				case 象:
					displayBoard = displayBoard + fmt.Sprint("黑象\t")
				case 車:
					displayBoard = displayBoard + fmt.Sprint("黑車\t")
				case 馬:
					displayBoard = displayBoard + fmt.Sprint("黑馬\t")
				case 炮:
					displayBoard = displayBoard + fmt.Sprint("黑炮\t")
				case 卒:
					displayBoard = displayBoard + fmt.Sprint("黑卒\t")
				}
			} else if xyboard.Color == red {
				switch xyboard.Piece {
				case 將:
					displayBoard = displayBoard + fmt.Sprint("紅帥\t")
				case 士:
					displayBoard = displayBoard + fmt.Sprint("紅士\t")
				case 象:
					displayBoard = displayBoard + fmt.Sprint("紅象\t")
				case 車:
					displayBoard = displayBoard + fmt.Sprint("紅車\t")
				case 馬:
					displayBoard = displayBoard + fmt.Sprint("紅馬\t")
				case 炮:
					displayBoard = displayBoard + fmt.Sprint("紅炮\t")
				case 卒:
					displayBoard = displayBoard + fmt.Sprint("紅兵\t")
				}
			}

		}
		displayBoard = displayBoard + "\n\n"
	}
	displayBoard = displayBoard + fmt.Sprint("---\t")
	for x := range board[0] {
		displayBoard = displayBoard + fmt.Sprint(x, "\t")
	}
	return displayBoard
}

func (moment *Moment) DisplayStringMoment() string {
	s := display(moment.Board)
	return fmt.Sprint(moment.Action, s)
}

func (moment *Moment) GetAllPiece() (pieces []*Piece) {

	for _, bb := range moment.Board {
		for _, b := range bb {
			if b != nil {
				pieces = append(pieces, b)
			}
		}
	}
	return pieces
}
func (moment *Moment) Hash() string {
	return fmt.Sprint(xhashes.FNV64(moment.DisplayStringMoment()))
}
