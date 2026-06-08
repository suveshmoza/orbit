package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/suveshmoza/orbit/internal/config"
)

func Run(cfg config.ConfigFile) error {
	p := tea.NewProgram(New(cfg))
	_, err := p.Run()
	return err
}
