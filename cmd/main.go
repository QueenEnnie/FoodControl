package main

import (
	"context"
	"food-control/internal/db"
)


func main() {
	db := db.New()

	err := db.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	println("DB connected")
}
