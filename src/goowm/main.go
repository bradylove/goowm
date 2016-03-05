package main

import (
	"fmt"
	"goowm/config"
	"goowm/windowmanager"
	"log"
)

func main() {
	conf, err := config.Load("config", "$HOME/.config/goowm")
	if err != nil {
		panic(fmt.Errorf("Error loading config: %s", err))
	}

	wm, err := windowmanager.New(conf)
	if err != nil {
		panic(err)
	}

	log.Println("Press them buttons!")
	wm.Run()
}
