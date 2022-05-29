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
		Password: "", // no password set
		DB:       0,  // use default DB
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

func (c *ChessDB) SetBoard2Redis(hash string, moment MomentResult) {
	err := c.RDB.Set(c.ctx, hash, c.MarshalBinary(moment), 0).Err()
	if err != nil {
		panic(err)
	}
}
func (c *ChessDB) GetBoard2Redis(hash string) (moment MomentResult) {
	b, err := c.RDB.Get(c.ctx, hash).Result()
	if err != nil {
		panic(err)
	}
	moment = c.UnMarshalBinary(b)
	return moment
}
func (c *ChessDB) MarshalBinary(moment MomentResult) []byte {
	b, err := json.Marshal(moment)
	if err != nil {
		panic(err)
	}
	return b
}

func (c *ChessDB) UnMarshalBinary(data string) (moment MomentResult) {

	b := []byte(data)
	err := json.Unmarshal(b, &moment)
	if err != nil {
		panic(err)
	}
	return moment
}
