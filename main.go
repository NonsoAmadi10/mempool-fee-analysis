package main

import (
	"log"

	"github.com/NonsoAmadi10/mempool-fee/app"
)

func main() {

	err := app.App().Listen("0.0.0.0:4000")

	if err != nil {
		log.Fatal(err)
	}
}
