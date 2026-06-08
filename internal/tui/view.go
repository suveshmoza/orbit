package tui

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/suveshmoza/orbit/internal/benchmark"
)

var (
	titleStyle   = lipgloss.NewStyle().Bold(true)
	mutedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	bestStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	footerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	headerStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("252"))
	dividerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

func render(m model) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Orbit — DNS Benchmark"))
	b.WriteString("\n")

	if m.phase == phaseRunning {
		b.WriteString(dividerStyle.Render(strings.Repeat("─", 40)))
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render(fmt.Sprintf("%s Running: %s", m.spinner.View(), m.current)))
		b.WriteString("\n\n")
		pct := 0.0
		if m.total > 0 {
			pct = float64(m.completed) / float64(m.total)
		}
		b.WriteString(m.progress.ViewAs(pct))
		fmt.Fprintf(&b, " %d/%d", m.completed, m.total)
		b.WriteString("\n")
	}

	bestMedian := bestMedianServer(m.results, m.phase == phaseDone)
	b.WriteString(renderResultsTable(m.results, bestMedian))
	b.WriteString("\n")
	if m.phase == phaseRunning {
		b.WriteString(footerStyle.Render("q quit\n"))
	}

	return b.String()
}

func renderResultsTable(results map[string]benchmark.ServerResult, bestMedian string) string {
	var rows [][]string
	bestRowIdx := -1
	rowIdx := 0

	for _, server := range sortedServers(results) {
		result := results[server]
		if result.Passed+result.Failed == 0 {
			continue
		}

		passed := fmt.Sprintf("%d/%d", result.Passed, result.Expected)
		mean, median, lowest := "n/a", "n/a", "n/a"
		if result.Passed > 0 {
			s := benchmark.Compute(result.RTTs)
			mean = formatDuration(s.Mean)
			median = formatDuration(s.Median)
			lowest = formatDuration(s.Lowest)
		}

		rows = append(rows, []string{server, passed, mean, median, lowest})
		if server == bestMedian {
			bestRowIdx = rowIdx
		}
		rowIdx++
	}

	if len(rows) == 0 {
		return ""
	}

	cellStyle := lipgloss.NewStyle().Padding(0, 1)

	t := table.New().
		Headers("Server", "Passed", "Mean", "Median", "Lowest").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch row {
			case table.HeaderRow:
				return headerStyle.Padding(0, 1)
			case bestRowIdx:
				return bestStyle.Padding(0, 1)
			default:
				return cellStyle
			}
		})

	return t.String()
}

func sortedServers(results map[string]benchmark.ServerResult) []string {
	return slices.Sorted(maps.Keys(results))
}

func bestMedianServer(results map[string]benchmark.ServerResult, done bool) string {
	if !done {
		return ""
	}

	var best string
	var bestMedian time.Duration
	for server, result := range results {
		if result.Passed == 0 {
			continue
		}
		median := benchmark.Compute(result.RTTs).Median
		if best == "" || median < bestMedian {
			best = server
			bestMedian = median
		}
	}
	return best
}

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return d.Round(time.Millisecond).String()
}
