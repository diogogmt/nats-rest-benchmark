package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting nats client...")
	go startNats()
	fmt.Println("Starting rest server...")
	startRest()
}
