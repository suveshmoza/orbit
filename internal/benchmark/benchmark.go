package benchmark

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"codeberg.org/miekg/dns"
	"github.com/suveshmoza/orbit/internal/config"
)

func Run(ctx context.Context, cfg config.ConfigFile) (map[string][]time.Duration, error) {
	stats := make(map[string][]time.Duration)
	domains := cfg.Config.TestDomains

	for _, server := range cfg.DNSServers {
		client := dns.NewClient()

		addr := net.JoinHostPort(server.Address, strconv.Itoa(server.Port))
		results := make([]time.Duration, 0, len(domains)*cfg.Config.Samples)

		for _, domain := range domains {
			for range cfg.Config.Samples {
				msg := dns.NewMsg(domain, dns.TypeA)

				_, rtt, err := client.Exchange(ctx, msg, "udp", addr)
				if err != nil {
					return nil, fmt.Errorf("error exchanging %s via %s: %w", domain, server.Name, err)
				}

				results = append(results, rtt)
			}
		}
		stats[server.Name] = results
	}
	return stats, nil
}
