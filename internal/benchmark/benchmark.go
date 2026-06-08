package benchmark

import (
	"context"
	"net"
	"strconv"
	"sync"
	"time"

	"codeberg.org/miekg/dns"
	"github.com/suveshmoza/orbit/internal/config"
)

type ServerResult struct {
	RTTs     []time.Duration
	Passed   int
	Failed   int
	Expected int
}

func RunStreaming(ctx context.Context, cfg config.ConfigFile, events chan<- Event) {
	client := dns.NewClient()
	domains := cfg.Config.TestDomains
	samples := cfg.Config.Samples
	timeoutMs := cfg.Config.TimeoutMs

	var wg sync.WaitGroup

	for _, server := range cfg.DNSServers {
		wg.Add(1)
		go func(server config.Server) {
			defer wg.Done()
			runServer(ctx, client, server, domains, samples, timeoutMs, func(e Event) {
				sendEvent(ctx, events, e)
			})
		}(server)
	}

	wg.Wait()
}

func sendEvent(ctx context.Context, events chan<- Event, e Event) {
	select {
	case events <- e:
	case <-ctx.Done():
	}
}

func runServer(
	ctx context.Context,
	client *dns.Client,
	server config.Server,
	domains []string,
	samples int,
	timeoutMs int,
	report func(Event),
) {
	addr := net.JoinHostPort(server.Address, strconv.Itoa(server.Port))

	for _, domain := range domains {
		for s := range samples {
			if ctx.Err() != nil {
				return
			}

			qctx := ctx
			var cancel context.CancelFunc
			if timeoutMs > 0 {
				qctx, cancel = context.WithTimeout(ctx, time.Duration(timeoutMs)*time.Millisecond)
			}

			msg := dns.NewMsg(domain, dns.TypeA)
			_, rtt, err := client.Exchange(qctx, msg, "udp", addr)
			if cancel != nil {
				cancel()
			}

			report(Event{
				Server: server.Name,
				Domain: domain,
				Sample: s + 1,
				RTT:    rtt,
				Err:    err,
			})
		}
	}
}
