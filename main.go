package main

import (
	"fmt"
	"os"
)

const (
	VERSION = "0.1"
)

func usage() {
	fmt.Printf(
		`herald-bot %s herald-bot is a Telegram bot that can do 
Usage:

Where:
      
`, VERSION)
	os.Exit(1)
}

func main() {
	h, err := NewHerald("139682629:AAE6123HxiwNlFMKr366mIGlZSYq_uV1Nn0")
	if err != nil {
		panic(err)
	}

	h.Run()
}
