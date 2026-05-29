package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"codeberg.org/miekg/dns"
)

func main() {

	client := dns.NewClient()
	results := make([]time.Duration, 0)

	for range 10 {
		msg := dns.NewMsg("google.com", dns.TypeA)
		addr := net.JoinHostPort("1.1.1.1", "53")

		_, rtt, err := client.Exchange(context.Background(), msg, "udp", addr)
		if err != nil {
			fmt.Printf("Error exchanging: %v\n", err)
			os.Exit(1)
		}

		results = append(results, rtt)
	}

	var average time.Duration
	for _, result := range results {
		average += result
	}
	fmt.Println(average)
	average /= time.Duration(len(results))
	fmt.Printf("Average RTT: %v\n", average)
}
