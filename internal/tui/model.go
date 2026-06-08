package tui

import (
	"context"
	"fmt"
	"time"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/suveshmoza/orbit/internal/benchmark"
	"github.com/suveshmoza/orbit/internal/config"
)

type phase int

type doneMsg struct{}

const (
	phaseRunning phase = iota
	phaseDone

	doneDelay       = 150 * time.Millisecond
	eventBufferSize = 64
)

type eventMsg struct {
	event  benchmark.Event
	closed bool
}

type model struct {
	phase     phase
	total     int
	samples   int
	completed int
	current   string
	results   map[string]benchmark.ServerResult
	events    <-chan benchmark.Event
	cancel    context.CancelFunc
	spinner   spinner.Model
	progress  progress.Model
}

func New(cfg config.ConfigFile) tea.Model {
	ctx, cancel := context.WithCancel(context.Background())
	events := make(chan benchmark.Event, eventBufferSize)

	go func() {
		benchmark.RunStreaming(ctx, cfg, events)
		close(events)
	}()

	expected := len(cfg.Config.TestDomains) * cfg.Config.Samples
	serverNames := make([]string, len(cfg.DNSServers))
	for i, s := range cfg.DNSServers {
		serverNames[i] = s.Name
	}

	s := spinner.New(
		spinner.WithSpinner(spinner.Dot),
	)

	p := progress.New(
		progress.WithWidth(20),
		progress.WithoutPercentage(),
		progress.WithColors(lipgloss.Color("42")),
		progress.WithFillCharacters('█', '░'),
	)

	return model{
		phase:    phaseRunning,
		total:    expected * len(cfg.DNSServers),
		samples:  cfg.Config.Samples,
		results:  benchmark.InitResults(serverNames, expected),
		events:   events,
		cancel:   cancel,
		spinner:  s,
		progress: p,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(waitForEvent(m.events), m.spinner.Tick)
}

func waitForEvent(ch <-chan benchmark.Event) tea.Cmd {
	return func() tea.Msg {
		e, ok := <-ch
		if !ok {
			return eventMsg{closed: true}
		}
		return eventMsg{event: e}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.cancel()
			return m, tea.Quit
		}

	case spinner.TickMsg:
		if m.phase != phaseRunning {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case eventMsg:
		if msg.closed {
			m.phase = phaseDone
			return m, tea.Tick(doneDelay, func(t time.Time) tea.Msg { return doneMsg{} })
		}

		e := msg.event
		m.current = fmt.Sprintf("%s · %s · sample %d/%d", e.Server, e.Domain, e.Sample, m.samples)
		benchmark.ApplyEvent(m.results, e)
		m.completed++

		if m.phase == phaseRunning {
			return m, waitForEvent(m.events)
		}

	case doneMsg:
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() tea.View {
	return tea.NewView(render(m))
}
