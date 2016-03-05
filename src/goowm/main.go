package main

import (
	"goowm/windowmanager"
	"log"
	"os"
)

func main() {
	d := os.Getenv("DISPLAY")
	if d == "" {
		d = ":0"
	}

	wm, err := windowmanager.New(d)
	if err != nil {
		panic(err)
	}

	log.Println("Press them buttons!")
	wm.Run()
}
