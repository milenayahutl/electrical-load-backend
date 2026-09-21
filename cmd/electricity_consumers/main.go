package main

import (
	"log"

	"electricity_consumers/internal/api"
)

func main() {
	log.Println("Application start!")

	api.StartServer()

	log.Println("Application terminated!")
}
