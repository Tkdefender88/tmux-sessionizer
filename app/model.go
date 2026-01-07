package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"gitlab.com/Tkdefender88/tmux-sessionizer/app/filepicker"
)

type Model struct {
	filePicker filepicker.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}
