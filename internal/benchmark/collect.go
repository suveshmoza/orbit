package benchmark

import "time"

type Event struct {
	Server string
	Domain string
	Sample int
	RTT    time.Duration
	Err    error
}

func ApplyEvent(results map[string]ServerResult, e Event) {
	r := results[e.Server]
	if e.Err != nil {
		r.Failed++
	} else {
		r.RTTs = append(r.RTTs, e.RTT)
		r.Passed++
	}
	results[e.Server] = r
}

func InitResults(servers []string, expectedPerServer int) map[string]ServerResult {
	results := make(map[string]ServerResult, len(servers))
	for _, name := range servers {
		results[name] = ServerResult{Expected: expectedPerServer}
	}
	return results
}
