package benchmark

import (
	"fmt"
	"slices"
	"time"
)

type ResultStats struct {
	Mean   time.Duration
	Median time.Duration
	Lowest time.Duration
}

func Compute(results []time.Duration) ResultStats {
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

func PrintStats(stats map[string][]time.Duration) {
	for server, results := range stats {
		s := Compute(results)
		fmt.Printf("%s Mean: [%v], Median: [%v], Lowest: [%v]\n", server, s.Mean, s.Median, s.Lowest)
	}
}
