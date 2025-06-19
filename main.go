package main

import (
	"fmt"
	"groupie-tracker/handlers"
)

func main() {
	fmt.Println("Starting server on http://localhost:8080")
	handlers.StartServer()
}