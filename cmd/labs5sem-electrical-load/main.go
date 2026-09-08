package main

import (
	"log"

	"labs5sem-electrical-load/internal/api"
)

func main() {
	log.Println("Application start!")

	api.StartServer()

	log.Println("Application terminated!")
}
