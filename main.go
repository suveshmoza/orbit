package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"codeberg.org/miekg/dns"
)

type ConfigFile struct {
	Config     AppConfig `json:"config"`
	DNSServers []Server  `json:"servers"`
}

type AppConfig struct {
	Samples     int      `json:"samples"`
	TimeoutMs   int      `json:"timeout_ms"`
	TestDomains []string `json:"test_domains"`
}

type Server struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func loadConfig(path string) (ConfigFile, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return ConfigFile{}, fmt.Errorf("read config: %w", err)
	}

	var cfg ConfigFile
	if err := json.Unmarshal(content, &cfg); err != nil {
		return ConfigFile{}, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}

func main() {

	cfg, err := loadConfig("config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	leaderboard := make(map[string]time.Duration)
	domains := cfg.Config.TestDomains
	for _, server := range cfg.DNSServers {

		client := dns.NewClient()
		results := make([]time.Duration, 0)

		for _, domain := range domains {
			msg := dns.NewMsg(domain, dns.TypeA)
			addr := net.JoinHostPort(server.Address, strconv.Itoa(server.Port))

			_, rtt, err := client.Exchange(context.Background(), msg, "udp", addr)
			if err != nil {
				fmt.Printf("Error exchanging: %v\n", err)
				os.Exit(1)
			}

			results = append(results, rtt)
		}

		var avg time.Duration
		for _, result := range results {
			avg += result
		}
		avg /= time.Duration(len(results))
		leaderboard[server.Name] = avg

	}

	fmt.Println(leaderboard)

}
