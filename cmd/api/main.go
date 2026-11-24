package main

import (
	"inventory-api/internal/app"
	"log"
)

func main() {

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
