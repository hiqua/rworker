package main

import "fmt"
import "time"

func main() {
	fmt.Println("Starting server...")
	for {
		// Lazy server!
		time.Sleep(3600 * time.Second)
	}
}
