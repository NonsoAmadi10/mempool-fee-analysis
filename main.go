package main

import (
	"log"

	//"github.com/NonsoAmadi10/mempool-fee/app"
	mempoolfee "github.com/NonsoAmadi10/mempool-fee/mempool-fee"
)

func main() {

	// err := app.App().Listen("0.0.0.0:4000")

	// if err != nil {
	// 	log.Fatal(err)
	// }

	fee := mempoolfee.GetBestFee()

	log.Println(fee)
}
