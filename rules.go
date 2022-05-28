package main

import (
	"errors"
)

const (
	KILL = iota
	UP
	DOWN
	LEFT
	RIGHT

	LEFTUP
	LEFTDOWN
	RIGHTUP
	RIGHTDOWN

	LEFTLEFTUP
	LEFTLEFTDOWN
	LEFTUPUP
	LEFTDOWNDOWN
	RIGHTRIGHTUP
	RIGHTRIGHTDOWN
	RIGHTUPUP
	RIGHTDOWNDOWN

	JUMPUP
	JUMPDOWN
	JUMPLEFT
	JUMPRIGHT
)

func (m *Moment) Walk(p *Piece, direction int, distance int) error {
	if p == nil {
		return errors.New("empty piece")
	}

	// find piece
	x, y, err := func() (int, int, error) {
		for y, b := range m.Board {
			for x := range b {
				if m.Board[y][x] == p {
					return x, y, nil
				}
			}
		}
		return -1, -1, errors.New("piece doesn't on board")
	}()
	if err != nil {
		panic(err)
	}

	if m.Action != p.Color {
		return errors.New("wrong move")
	}

	var xx, yy int
	type Stock struct {
		XX int
		YY int
	}
	var stock []Stock

	switch p.Piece {
	case 將:
		switch direction {
		case KILL:
			if m.Action == red {
				for yyy := y - 1; yyy >= 0; yyy-- {
					if m.Board[yyy][x].Piece == 將 {
						xx = x
						yy = yyy
						stock = nil
						break
					}
				}
				return errors.New("wrong move")
			} else if m.Action == black {
				for yyy := y + 1; yyy <= 9; yyy++ {
					if m.Board[yyy][x].Piece == 將 {
						xx = x
						yy = yyy
						stock = nil
						break
					}
				}
				return errors.New("wrong move")
			}
		case UP:
			xx = x
			yy = y - 1
			stock = nil
		case DOWN:
			xx = x
			yy = y + 1
			stock = nil
		case LEFT:
			xx = x - 1
			yy = y
			stock = nil
		case RIGHT:
			xx = x + 1
			yy = y
			stock = nil
		default:
			return errors.New("wrong move")
		}

		if direction != KILL {
			if xx < 3 || xx > 5 {
				return errors.New("wrong move")
			}
			if yy > 2 && yy < 7 {
				return errors.New("wrong move")
			}
		}
	case 士:
		switch direction {
		case LEFTUP:
			xx = x - 1
			yy = y - 1
			stock = nil
		case LEFTDOWN:
			xx = x - 1
			yy = y + 1
			stock = nil
		case RIGHTUP:
			xx = x + 1
			yy = y - 1
			stock = nil
		case RIGHTDOWN:
			xx = x - 1
			yy = y + 1
			stock = nil
		default:
			return errors.New("wrong move")
		}
	case 象:

		switch direction {
		case LEFTUP:
			xx = x - 2
			yy = y - 2
			stock = append(stock, Stock{
				XX: x - 1,
				YY: y - 1,
			})
		case LEFTDOWN:
			xx = x - 2
			yy = y + 2
			stock = append(stock, Stock{
				XX: x - 1,
				YY: y + 1,
			})
		case RIGHTUP:
			xx = x + 2
			yy = y - 2
			stock = append(stock, Stock{
				XX: x + 1,
				YY: y - 1,
			})
		case RIGHTDOWN:
			xx = x - 2
			yy = y + 2
			stock = append(stock, Stock{
				XX: x - 1,
				YY: y + 1,
			})
		default:
			return errors.New("wrong move")
		}
	case 車:
		if distance < 1 {
			return errors.New("wrong move")
		}
		switch direction {
		case UP:
			xx = x
			yy = y - distance
			for yyy := y - 1; yyy >= y-distance; yyy-- {
				stock = append(stock, Stock{
					XX: x,
					YY: yyy,
				})
			}
		case DOWN:
			xx = x
			yy = y + distance
			for yyy := y + 1; yyy <= y+distance; yyy++ {
				stock = append(stock, Stock{
					XX: x,
					YY: yyy,
				})
			}
		case LEFT:
			xx = x - distance
			yy = y
			for xxx := x - 1; xxx >= y-distance; xxx-- {
				stock = append(stock, Stock{
					XX: xxx,
					YY: y,
				})
			}
		case RIGHT:
			xx = x + distance
			yy = y
			for xxx := x + 1; xxx <= y+distance; xxx++ {
				stock = append(stock, Stock{
					XX: xxx,
					YY: y,
				})
			}
		default:
			return errors.New("wrong move")
		}

	case 馬:
		switch direction {
		case LEFTLEFTUP:
			xx = x - 2
			yy = y - 1
			stock = append(stock, Stock{
				XX: x - 1,
				YY: y,
			})
		case RIGHTDOWNDOWN:
			xx = x + 1
			yy = y + 2
			stock = append(stock, Stock{
				XX: x,
				YY: y + 1,
			})
		case LEFTLEFTDOWN:
			xx = x - 2
			yy = y + 1
			stock = append(stock, Stock{
				XX: x - 1,
				YY: y,
			})
		case LEFTUPUP:
			xx = x - 1
			yy = y - 2
			stock = append(stock, Stock{
				XX: x,
				YY: y - 1,
			})
		case LEFTDOWNDOWN:
			xx = x - 1
			yy = y + 2
			stock = append(stock, Stock{
				XX: x,
				YY: y + 1,
			})
		case RIGHTRIGHTUP:
			xx = x + 2
			yy = y - 1
			stock = append(stock, Stock{
				XX: x + 1,
				YY: y,
			})
		case RIGHTRIGHTDOWN:
			xx = x + 2
			yy = y + 1
			stock = append(stock, Stock{
				XX: x + 1,
				YY: y,
			})
		case RIGHTUPUP:
			xx = x + 1
			yy = y - 2
			stock = append(stock, Stock{
				XX: x,
				YY: y - 1,
			})
		default:
			return errors.New("wrong move")
		}
	case 炮:
		if (direction == UP ||
			direction == DOWN ||
			direction == LEFT ||
			direction == RIGHT) &&
			distance < 1 {
			return errors.New("wrong move")
		}

		switch direction {
		case UP:
			xx = x
			yy = y - distance
			for yyy := y - 1; yyy >= y-distance; yyy-- {
				stock = append(stock, Stock{
					XX: x,
					YY: yyy,
				})
			}
		case DOWN:
			xx = x
			yy = y + distance
			for yyy := y + 1; yyy <= y+distance; yyy++ {
				stock = append(stock, Stock{
					XX: x,
					YY: yyy,
				})
			}
		case LEFT:
			xx = x - distance
			yy = y
			for xxx := x - 1; xxx >= y-distance; xxx-- {
				stock = append(stock, Stock{
					XX: xxx,
					YY: y,
				})
			}
		case RIGHT:
			xx = x + distance
			yy = y
			for xxx := x + 1; xxx <= y+distance; xxx++ {
				stock = append(stock, Stock{
					XX: xxx,
					YY: y,
				})
			}
		case JUMPUP:
			xx = -1
			yy = -1
			stock = nil
			count := 0
			for yyy := y - 1; yyy >= 0; yyy-- {
				if m.Board[yyy][x] != nil {
					count = count + 1
					if count == 2 {
						xx = x
						yy = yyy
						break
					}
				}
			}
			if count != 2 {
				return errors.New("wrong move")
			}
		case JUMPDOWN:
			xx = -1
			yy = -1
			stock = nil
			count := 0
			for yyy := y + 1; yyy <= 9; yyy++ {
				if m.Board[yyy][x] != nil {
					count = count + 1
					if count == 2 {
						xx = x
						yy = yyy
						break
					}
				}
			}
			if count != 2 {
				return errors.New("wrong move")
			}
		case JUMPLEFT:
			xx = -1
			yy = -1
			stock = nil
			count := 0
			for xxx := x - 1; xxx >= 0; xxx-- {
				if m.Board[y][xxx] != nil {
					count = count + 1
					if count == 2 {
						xx = xxx
						yy = y
						break
					}
				}
			}
			if count != 2 {
				return errors.New("wrong move")
			}
		case JUMPRIGHT:
			xx = -1
			yy = -1
			stock = nil
			count := 0
			for xxx := x + 1; xxx <= 8; xxx++ {
				if m.Board[y][xxx] != nil {
					count = count + 1
					if count == 2 {
						xx = xxx
						yy = y
						break
					}
				}
			}
			if count != 2 {
				return errors.New("wrong move")
			}
		default:
			return errors.New("wrong move")
		}
		if direction == UP || direction == DOWN || direction == LEFT || direction == RIGHT {
			if m.Board[yy][xx] != nil {
				return errors.New("wrong move")
			}
		}
	case 卒:
		if m.Action == red && direction == DOWN {
			return errors.New("wrong move")
		}
		if m.Action == black && direction == UP {
			return errors.New("wrong move")
		}
		if m.Action == red && y >= 5 && (direction == LEFT || direction == RIGHT) {
			return errors.New("wrong move")
		}
		if m.Action == black && y <= 4 && (direction == LEFT || direction == RIGHT) {
			return errors.New("wrong move")
		}
		switch direction {
		case UP:
			xx = x
			yy = y - 1
			stock = nil
		case DOWN:
			xx = x
			yy = y + 1
			stock = nil
		case LEFT:
			xx = x - 1
			yy = y
			stock = nil
		case RIGHT:
			xx = x + 1
			yy = y
			stock = nil
		default:
			return errors.New("wrong move")
		}

	}

	// stock
	if stock != nil {
		for _, s := range stock {
			if s.XX < 0 || s.YY < 0 || s.XX > 8 || s.YY > 9 { //ignore error stock
				continue
			}
			if m.Board[s.YY][s.XX] != nil {
				return errors.New("wrong move")
			}
		}
	}

	// move
	xfront := 0
	xrear := 8
	yfront := 0
	yrear := 9

	if xx < xfront || xx > xrear || yy < yfront || yy > yrear {
		return errors.New("wrong move")
	}
	if m.Board[yy][xx] != nil {
		if m.Board[yy][xx].Color == m.Action {
			return errors.New("wrong move")
		}
	}
	m.Board[yy][xx] = m.Board[y][x]
	m.Board[y][x] = nil
	if m.Action == red {
		m.Action = black
	} else if m.Action == black {
		m.Action = red
	}
	return nil
}
