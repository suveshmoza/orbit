package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"slices"
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

type ResultStats struct {
	Mean   time.Duration
	Median time.Duration
	Lowest time.Duration
}

func generateResultStats(results []time.Duration) ResultStats {
	if len(results) == 0 {
		return ResultStats{}
	}

	var sum time.Duration
	lowest := results[0]
	for _, r := range results {
		sum += r
		if r < lowest {
			lowest = r
		}
	}
	mean := sum / time.Duration(len(results))

	sorted := slices.Clone(results)
	slices.Sort(sorted)

	n := len(sorted)
	var median time.Duration
	if n%2 == 1 {
		median = sorted[n/2]
	} else {
		median = (sorted[n/2-1] + sorted[n/2]) / 2
	}

	return ResultStats{
		Mean:   mean,
		Median: median,
		Lowest: lowest,
	}
}

func main() {

	cfg, err := loadConfig("config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	stats := make(map[string][]time.Duration)
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
		stats[server.Name] = results
	}

	for server, results := range stats {
		s := generateResultStats(results)
		fmt.Printf("Server [%s] Mean: [%v], Median: [%v], Lowest: [%v]\n", server, s.Mean, s.Median, s.Lowest)
	}

}
