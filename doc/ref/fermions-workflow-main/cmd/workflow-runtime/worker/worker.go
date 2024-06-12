package main

import (
	"log"

	"github.com/galgotech/fermions-workflow/internal/cmd"
)

func main() {
	err := cmd.Worker()
	if err != nil {
		log.Fatal(err)
	}
}
