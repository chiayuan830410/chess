package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/shomali11/util/xhashes"
)

type ChessDB struct {
	RDB *redis.Client
	ctx context.Context
}

func (c *ChessDB) New() {
	c.ctx = context.Background()
	c.RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "a830410a", // no password set
		DB:       0,          // use default DB
	})
}

func (c *ChessDB) CheckBoard2Redis(hash string) bool {
	n, err := c.RDB.Exists(c.ctx, hash).Result()
	if err != nil {
		panic(err)
	}
	if n > 0 {
		return true
	} else {
		return false
	}
}

func (c *ChessDB) SetBoard2Redis(hash string, moment Moment) {
	err := c.RDB.Set(c.ctx, hash, c.MarshalBinary(moment), 0).Err()
	if err != nil {
		panic(err)
	}
}
func (c *ChessDB) GetBoard2Redis(hash string) (moment Moment) {
	b, err := c.RDB.Get(c.ctx, hash).Result()
	if err != nil {
		panic(err)
	}
	moment = c.UnMarshalBinary(b)
	return moment
}
func (c *ChessDB) MarshalBinary(moment Moment) []byte {
	b, err := json.Marshal(moment)
	if err != nil {
		panic(err)
	}
	return b
}

func (c *ChessDB) UnMarshalBinary(data string) (moment Moment) {

	b := []byte(data)
	err := json.Unmarshal(b, &moment)
	if err != nil {
		panic(err)
	}
	return moment
}

const (
	將 = iota
	士
	象
	車
	馬
	炮
	卒
)
const (
	red   = 1
	black = 2
)

type Piece struct {
	Piece int
	Color int // red: 1, black: 2
}

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

type Board [10][9]*Piece

type Moment struct {
	Board  Board
	Action int // red: 1, black: 2
}

var ChessLibrary ChessDB
var maxCount = 0

func Steps(moment Moment, count int) {
	if count > maxCount {
		maxCount = count
		fmt.Println(count)
	}
	if count > 50 {
		return
	}
	hash := xhashes.FNV64(display(moment.Board) + strconv.Itoa(moment.Action))
	// fmt.Println(hash)
	if ChessLibrary.CheckBoard2Redis(fmt.Sprint(hash)) {
		return
	} else {
		ChessLibrary.SetBoard2Redis(fmt.Sprint(hash), moment)
	}

	// finished
	finishFlag := 0
	for _, bb := range moment.Board {
		for _, b := range bb {
			if b == nil {
				continue
			}
			if b.Piece == 將 {
				finishFlag++
			}
		}
	}
	if finishFlag < 2 {
		return
	}
	// finished

	if moment.Action == red {
		moment.Action = black
	} else if moment.Action == black {
		moment.Action = red
	}

	// moment.Board := moment.Board
	// move
	basicBoard := moment.Board
	for y, xbroard := range moment.Board {
		for x, xyboard := range xbroard {
			if xyboard != nil && xyboard.Color == moment.Action {
				switch xyboard.Piece {
				case 將:
					if xyboard.Color == red {
						//left
						basicBoard = moment.Board
						if (x - 1) >= 3 {
							if basicBoard[y][x-1] == nil {
								basicBoard[y][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y][x-1].Color == black {
								basicBoard[y][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)
						//right
						basicBoard = moment.Board
						if (x + 1) <= 5 {
							if basicBoard[y][x+1] == nil {
								basicBoard[y][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y][x+1].Color == black {
								basicBoard[y][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)
						//up
						basicBoard = moment.Board
						if (y + 1) <= 9 {
							if basicBoard[y+1][x] == nil {
								basicBoard[y+1][x] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y+1][x].Color == black {
								basicBoard[y+1][x] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)
						//down
						basicBoard = moment.Board
						if (y - 1) >= 7 {
							if basicBoard[y-1][x] == nil {
								basicBoard[y-1][x] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y-1][x].Color == black {
								basicBoard[y-1][x] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)
					} else if xyboard.Color == black {
						//left
						basicBoard = moment.Board
						if (x - 1) >= 3 {
							if basicBoard[y][x-1] == nil {
								basicBoard[y][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y][x-1].Color == red {
								basicBoard[y][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)
						//right
						basicBoard = moment.Board
						if (x + 1) <= 5 {
							if basicBoard[y][x+1] == nil {
								basicBoard[y][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y][x+1].Color == red {
								basicBoard[y][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)
						//up
						basicBoard = moment.Board
						if (y + 1) <= 2 {
							if basicBoard[y+1][x] == nil {
								basicBoard[y+1][x] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y+1][x].Color == red {
								basicBoard[y+1][x] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)
						//down
						basicBoard = moment.Board
						if (y - 1) >= 0 {
							if basicBoard[y-1][x] == nil {
								basicBoard[y-1][x] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y-1][x].Color == red {
								basicBoard[y-1][x] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)
					}
				case 士:
					if xyboard.Color == red {
						//left up
						basicBoard = moment.Board
						if (x-1) >= 3 && (y+1) <= 9 {
							if basicBoard[y+1][x-1] == nil {
								basicBoard[y+1][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y+1][x-1].Color == black {
								basicBoard[y+1][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)
						//right up
						basicBoard = moment.Board
						if (x+1) <= 5 && (y+1) <= 9 {
							if basicBoard[y+1][x+1] == nil {
								basicBoard[y+1][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y+1][x+1].Color == black {
								basicBoard[y+1][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)

						//left down
						basicBoard = moment.Board
						if (x-1) >= 3 && (y-1) >= 7 {
							if basicBoard[y-1][x-1] == nil {
								basicBoard[y-1][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y-1][x-1].Color == black {
								basicBoard[y-1][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)

						//right down
						basicBoard = moment.Board
						if (x+1) <= 5 && (y-1) >= 7 {
							if basicBoard[y-1][x+1] == nil {
								basicBoard[y-1][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y-1][x+1].Color == black {
								basicBoard[y-1][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)

					} else if xyboard.Color == black {
						//left up
						basicBoard = moment.Board
						if (x-1) >= 3 && (y+1) <= 2 {
							if basicBoard[y+1][x-1] == nil {
								basicBoard[y+1][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y+1][x-1].Color == red {
								basicBoard[y+1][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)
						//right up
						basicBoard = moment.Board
						if (x+1) <= 5 && (y+1) <= 2 {
							if basicBoard[y+1][x+1] == nil {
								basicBoard[y+1][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y+1][x+1].Color == red {
								basicBoard[y+1][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)

						//left down
						basicBoard = moment.Board
						if (x-1) >= 3 && (y-1) >= 0 {
							if basicBoard[y-1][x-1] == nil {
								basicBoard[y-1][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y-1][x-1].Color == red {
								basicBoard[y-1][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)

						//right down
						basicBoard = moment.Board
						if (x+1) <= 5 && (y-1) >= 0 {
							if basicBoard[y-1][x+1] == nil {
								basicBoard[y-1][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y-1][x+1].Color == red {
								basicBoard[y-1][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)
					}
				case 象:
					if xyboard.Color == red {
						//left up
						basicBoard = moment.Board
						if (x-2) >= 0 && (y+2) <= 9 {
							if basicBoard[y+1][x-1] == nil && basicBoard[y+2][x-2] == nil {
								basicBoard[y+2][x-2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y+1][x-1] == nil && basicBoard[y+2][x-2].Color == black {
								basicBoard[y+2][x-2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)
						//right up
						basicBoard = moment.Board
						if (x+2) <= 8 && (y+2) <= 9 {
							if basicBoard[y+1][x+1] == nil && basicBoard[y+2][x+2] == nil {
								basicBoard[y+2][x+2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y+1][x+1] == nil && basicBoard[y+2][x+2].Color == black {
								basicBoard[y+2][x+2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)

						//left down
						basicBoard = moment.Board
						if (x-2) >= 0 && (y-2) >= 5 {
							if basicBoard[y-1][x-1] == nil && basicBoard[y-2][x-2] == nil {
								basicBoard[y-2][x-2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y-1][x-1] == nil && basicBoard[y-2][x-2].Color == black {
								basicBoard[y-2][x-2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)

						//right down
						basicBoard = moment.Board
						if (x+2) <= 8 && (y-2) >= 5 {
							if basicBoard[y-1][x+1] == nil && basicBoard[y-2][x+2] == nil {
								basicBoard[y-2][x+2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y-1][x+1] == nil && basicBoard[y-2][x+2].Color == black {
								basicBoard[y-2][x+2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)

					} else if xyboard.Color == black {
						//left up
						basicBoard = moment.Board
						if (x-2) >= 0 && (y+2) <= 4 {
							if basicBoard[y+1][x-1] == nil && basicBoard[y+2][x-2] == nil {
								basicBoard[y+2][x-2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y+1][x-1] == nil && basicBoard[y+2][x-2].Color == red {
								basicBoard[y+2][x-2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)
						//right up
						basicBoard = moment.Board
						if (x+2) <= 8 && (y+2) <= 4 {
							if basicBoard[y+1][x+1] == nil && basicBoard[y+2][x+2] == nil {
								basicBoard[y+2][x+2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y+1][x+1] == nil && basicBoard[y+2][x+2].Color == red {
								basicBoard[y+2][x+2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)

						//left down
						basicBoard = moment.Board
						if (x-2) >= 0 && (y-2) >= 0 {
							if basicBoard[y-1][x-1] == nil && basicBoard[y-2][x-2] == nil {
								basicBoard[y-2][x-2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y-1][x-1] == nil && basicBoard[y-2][x-2].Color == red {
								basicBoard[y-2][x-2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)

						//right down
						basicBoard = moment.Board
						if (x+2) <= 8 && (y-2) >= 0 {
							if basicBoard[y-1][x+1] == nil && basicBoard[y-2][x+2] == nil {
								basicBoard[y-2][x+2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y-1][x+1] == nil && basicBoard[y-2][x+2].Color == red {
								basicBoard[y-2][x+2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)
					}
				case 車:
					if xyboard.Color == red {
						//left
						for xx := x - 1; xx >= 0; xx-- {
							basicBoard = moment.Board
							if basicBoard[y][xx] == nil {
								basicBoard[y][xx] = basicBoard[y][x]
								basicBoard[y][x] = nil
								Steps(Moment{
									Board:  basicBoard,
									Action: black,
								}, count+1)
							} else if basicBoard[y][xx].Color == black {
								basicBoard[y][xx] = basicBoard[y][x]
								basicBoard[y][x] = nil
								Steps(Moment{
									Board:  basicBoard,
									Action: black,
								}, count+1)
								break
							} else {
								break
							}

						}

						//right
						basicBoard = moment.Board
						for xx := x + 1; xx <= 8; xx++ {
							if basicBoard[y][xx] == nil {
								basicBoard[y][xx] = basicBoard[y][x]
								basicBoard[y][x] = nil
								Steps(Moment{
									Board:  basicBoard,
									Action: black,
								}, count+1)
							} else if basicBoard[y][xx].Color == black {
								basicBoard[y][xx] = basicBoard[y][x]
								basicBoard[y][x] = nil
								Steps(Moment{
									Board:  basicBoard,
									Action: black,
								}, count+1)
								break
							} else {
								break
							}
						}

						//up
						basicBoard = moment.Board
						for yy := y + 1; yy <= 9; yy++ {
							if basicBoard[yy][x] == nil {
								basicBoard[yy][x] = basicBoard[y][x]
								basicBoard[y][x] = nil
								Steps(Moment{
									Board:  basicBoard,
									Action: black,
								}, count+1)
							} else if basicBoard[yy][x].Color == black {
								basicBoard[yy][x] = basicBoard[y][x]
								basicBoard[y][x] = nil
								Steps(Moment{
									Board:  basicBoard,
									Action: black,
								}, count+1)
								break
							} else {
								break
							}
						}
						//down
						basicBoard = moment.Board
						for yy := y - 1; yy >= 0; yy-- {
							if basicBoard[yy][x] == nil {
								basicBoard[yy][x] = basicBoard[y][x]
								basicBoard[y][x] = nil
								Steps(Moment{
									Board:  basicBoard,
									Action: black,
								}, count+1)
							} else if basicBoard[yy][x].Color == black {
								basicBoard[yy][x] = basicBoard[y][x]
								basicBoard[y][x] = nil
								Steps(Moment{
									Board:  basicBoard,
									Action: black,
								}, count+1)
								break
							} else {
								break
							}
						}
					} else if xyboard.Color == black {
						//left
						for xx := x - 1; xx >= 0; xx-- {
							basicBoard = moment.Board
							if basicBoard[y][xx] == nil {
								basicBoard[y][xx] = basicBoard[y][x]
								basicBoard[y][x] = nil
								Steps(Moment{
									Board:  basicBoard,
									Action: red,
								}, count+1)
							} else if basicBoard[y][xx].Color == red {
								basicBoard[y][xx] = basicBoard[y][x]
								basicBoard[y][x] = nil
								Steps(Moment{
									Board:  basicBoard,
									Action: red,
								}, count+1)
								break
							} else {
								break
							}

						}

						//right
						basicBoard = moment.Board
						for xx := x + 1; xx <= 8; xx++ {
							if basicBoard[y][xx] == nil {
								basicBoard[y][xx] = basicBoard[y][x]
								basicBoard[y][x] = nil
								Steps(Moment{
									Board:  basicBoard,
									Action: red,
								}, count+1)
							} else if basicBoard[y][xx].Color == red {
								basicBoard[y][xx] = basicBoard[y][x]
								basicBoard[y][x] = nil
								Steps(Moment{
									Board:  basicBoard,
									Action: red,
								}, count+1)
								break
							} else {
								break
							}
						}

						//up
						basicBoard = moment.Board
						for yy := y + 1; yy <= 9; yy++ {
							if basicBoard[yy][x] == nil {
								basicBoard[yy][x] = basicBoard[y][x]
								basicBoard[y][x] = nil
								Steps(Moment{
									Board:  basicBoard,
									Action: red,
								}, count+1)
							} else if basicBoard[yy][x].Color == red {
								basicBoard[yy][x] = basicBoard[y][x]
								basicBoard[y][x] = nil
								Steps(Moment{
									Board:  basicBoard,
									Action: red,
								}, count+1)
								break
							} else {
								break
							}
						}
						//down
						basicBoard = moment.Board
						for yy := y - 1; yy >= 0; yy-- {
							if basicBoard[yy][x] == nil {
								basicBoard[yy][x] = basicBoard[y][x]
								basicBoard[y][x] = nil
								Steps(Moment{
									Board:  basicBoard,
									Action: red,
								}, count+1)
							} else if basicBoard[yy][x].Color == red {
								basicBoard[yy][x] = basicBoard[y][x]
								basicBoard[y][x] = nil
								Steps(Moment{
									Board:  basicBoard,
									Action: red,
								}, count+1)
								break
							} else {
								break
							}
						}
					}

				case 馬:
					if xyboard.Color == red {
						//left up
						basicBoard = moment.Board
						if (x-2) >= 0 && (y+1) <= 9 {
							if basicBoard[y][x-1] == nil && basicBoard[y+1][x-2] == nil {
								basicBoard[y+1][x-2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y][x-1] == nil && basicBoard[y+1][x-2].Color == black {
								basicBoard[y+1][x-2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)
						//right up
						basicBoard = moment.Board
						if (x+2) <= 8 && (y+1) <= 9 {
							if basicBoard[y][x+1] == nil && basicBoard[y+1][x+2] == nil {
								basicBoard[y+1][x+2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y][x+1] == nil && basicBoard[y+1][x+2].Color == black {
								basicBoard[y+1][x+2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)

						//left down
						basicBoard = moment.Board
						if (x-2) >= 0 && (y-1) >= 0 {
							if basicBoard[y][x-1] == nil && basicBoard[y-1][x-2] == nil {
								basicBoard[y-1][x-2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y][x-1] == nil && basicBoard[y-1][x-2].Color == black {
								basicBoard[y-1][x-2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)

						//right down
						basicBoard = moment.Board
						if (x+2) <= 8 && (y-1) >= 0 {
							if basicBoard[y][x+1] == nil && basicBoard[y-1][x+2] == nil {
								basicBoard[y-1][x+2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y][x+1] == nil && basicBoard[y-1][x+2].Color == black {
								basicBoard[y-1][x+2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)

						// -----
						//left up
						basicBoard = moment.Board
						if (x+1) <= 8 && (y-2) >= 0 {
							if basicBoard[y-1][x] == nil && basicBoard[y-2][x+1] == nil {
								basicBoard[y-2][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y-1][x] == nil && basicBoard[y-2][x+1].Color == black {
								basicBoard[y-2][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)
						//right up
						basicBoard = moment.Board
						if (x+1) <= 8 && (y+2) <= 9 {
							if basicBoard[y+1][x] == nil && basicBoard[y+2][x+1] == nil {
								basicBoard[y+2][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y+1][x] == nil && basicBoard[y+2][x+1].Color == black {
								basicBoard[y+2][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)

						//left down
						basicBoard = moment.Board
						if (x-1) >= 0 && (y-2) >= 0 {
							if basicBoard[y-1][x] == nil && basicBoard[y-2][x-1] == nil {
								basicBoard[y-2][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y-1][x] == nil && basicBoard[y-2][x-1].Color == black {
								basicBoard[y-2][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)

						//right down
						basicBoard = moment.Board
						if (x-1) >= 0 && (y+2) <= 9 {
							if basicBoard[y+1][x] == nil && basicBoard[y+2][x-1] == nil {
								basicBoard[y+2][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y+1][x] == nil && basicBoard[y+2][x-1].Color == black {
								basicBoard[y+2][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: black,
						}, count+1)
						// -----
					} else if xyboard.Color == black {
						//left up
						basicBoard = moment.Board
						if (x-2) >= 0 && (y+1) <= 9 {
							if basicBoard[y][x-1] == nil && basicBoard[y+1][x-2] == nil {
								basicBoard[y+1][x-2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y][x-1] == nil && basicBoard[y+1][x-2].Color == red {
								basicBoard[y+1][x-2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)
						//right up
						basicBoard = moment.Board
						if (x+2) <= 8 && (y+1) <= 9 {
							if basicBoard[y][x+1] == nil && basicBoard[y+1][x+2] == nil {
								basicBoard[y+1][x+2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y][x+1] == nil && basicBoard[y+1][x+2].Color == red {
								basicBoard[y+1][x+2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)

						//left down
						basicBoard = moment.Board
						if (x-2) >= 0 && (y-1) >= 0 {
							if basicBoard[y][x-1] == nil && basicBoard[y-1][x-2] == nil {
								basicBoard[y-1][x-2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y][x-1] == nil && basicBoard[y-1][x-2].Color == red {
								basicBoard[y-1][x-2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)

						//right down
						basicBoard = moment.Board
						if (x+2) <= 8 && (y-1) >= 0 {
							if basicBoard[y][x+1] == nil && basicBoard[y-1][x+2] == nil {
								basicBoard[y-1][x+2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y][x+1] == nil && basicBoard[y-1][x+2].Color == red {
								basicBoard[y-1][x+2] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)
						//-------
						//left up
						basicBoard = moment.Board
						if (x+1) <= 8 && (y-2) >= 0 {
							if basicBoard[y-1][x] == nil && basicBoard[y-2][x+1] == nil {
								basicBoard[y-2][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y-1][x] == nil && basicBoard[y-2][x+1].Color == red {
								basicBoard[y-2][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)
						//right up
						basicBoard = moment.Board
						if (x+1) <= 8 && (y+2) <= 9 {
							if basicBoard[y+1][x] == nil && basicBoard[y+2][x+1] == nil {
								basicBoard[y+2][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y+1][x] == nil && basicBoard[y+2][x+1].Color == red {
								basicBoard[y+2][x+1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)

						//left down
						basicBoard = moment.Board
						if (x-1) >= 0 && (y-2) >= 0 {
							if basicBoard[y-1][x] == nil && basicBoard[y-2][x-1] == nil {
								basicBoard[y-2][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y-1][x] == nil && basicBoard[y-2][x-1].Color == red {
								basicBoard[y-2][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)

						//right down
						basicBoard = moment.Board
						if (x-1) >= 0 && (y+2) <= 9 {
							if basicBoard[y+1][x] == nil && basicBoard[y+2][x-1] == nil {
								basicBoard[y+2][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							} else if basicBoard[y+1][x] == nil && basicBoard[y+2][x-1].Color == red {
								basicBoard[y+2][x-1] = basicBoard[y][x]
								basicBoard[y][x] = nil
							}
						}
						Steps(Moment{
							Board:  basicBoard,
							Action: red,
						}, count+1)

					}
				case 炮:

				case 卒:
				}
			}
		}
	}

}

var library map[uint64]Moment

func main() {
	library = make(map[uint64]Moment)
	var initBoard Board = Board{
		{{車, black}, {馬, black}, {象, black}, {士, black}, {將, black}, {士, black}, {象, black}, {馬, black}, {車, black}},
		{},
		{nil, {炮, black}, nil, nil, nil, nil, nil, {炮, black}, nil},
		{{卒, black}, nil, {卒, black}, nil, {卒, black}, nil, {卒, black}, nil, {卒, black}},
		{},
		{},
		{{卒, red}, nil, {卒, red}, nil, {卒, red}, nil, {卒, red}, nil, {卒, red}},
		{nil, {炮, red}, nil, nil, nil, nil, nil, {炮, red}, nil},
		{},
		{{車, red}, {馬, red}, {象, red}, {士, red}, {將, red}, {士, red}, {象, red}, {馬, red}, {車, red}},
	}
	fmt.Println(display(initBoard))

	ChessLibrary.New()
	// for {
	// 	var strHash string
	// 	fmt.Scanf("%s", &strHash)
	// 	fmt.Println(display((ChessLibrary.GetBoard2Redis(strHash).Board)))
	// }

	var moment Moment
	// library := map[int64]Moment{}
	moment.Action = red
	moment.Board = initBoard
	Steps(moment, 0)

	// initBoard[0][0] = nil
	// initBoard[9][0] = &Piece{車, black}
	// hash := xhashes.FNV64(display(initBoard) + strconv.Itoa(red))
	// fmt.Println(hash)
	// fmt.Println(display(initBoard))

}
