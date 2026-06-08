package main

import (
	"context"
	"fmt"
	"os"

	"github.com/suveshmoza/orbit/internal/benchmark"
	"github.com/suveshmoza/orbit/internal/config"
)

func main() {

	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	stats, err := benchmark.Run(context.Background(), cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	benchmark.PrintStats(stats)

}
