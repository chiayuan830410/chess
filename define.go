package main

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

type Board [10][9]*Piece

type Moment struct {
	Board  Board
	Action int // red: 1, black: 2
}

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
