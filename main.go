package main

import (
	"fmt"
	"os"

	"github.com/suveshmoza/orbit/internal/config"
	"github.com/suveshmoza/orbit/internal/tui"
)

func main() {

	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := tui.Run(cfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
