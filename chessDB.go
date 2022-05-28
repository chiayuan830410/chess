package main

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
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
