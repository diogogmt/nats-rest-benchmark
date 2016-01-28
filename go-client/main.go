package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("Starting nats client...")
	go startNats()
	fmt.Println("Starting rest server...")
	go startRest()
	runtime.Goexit()
}
