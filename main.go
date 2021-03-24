package main

import (
	"log"

	"github.com/aaronellington/projectl/pkg/projectl"
)

func main() {
	app := &projectl.App{}
	if err := app.Execute(); err != nil {
		log.Fatal(err)
	}
}
