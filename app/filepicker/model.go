package filepicker

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle = lipgloss.NewStyle().PaddingLeft(4)
)

type directory string

func (d directory) FilterValue() string { return "" }

type directoryDelegate struct{}

func (d directoryDelegate) Height() int                             { return 1 }
func (d directoryDelegate) Spacing() int                            { return 0 }
func (d directoryDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d directoryDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	dir, ok := listItem.(directory)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", dir)

	fmt.Fprint(w, itemStyle.Render(str))
}

type Model struct {
	list list.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func NewFileModel() Model {
	const defaultWidth = 20
	const listHeight = 14
	l := list.New([]list.Item{
		directory("~/workspace"),
		directory("~/dotifles"),
	}, directoryDelegate{}, defaultWidth, listHeight)
	l.Title = "Fuzzy find directory"

	return Model{
		list: l,
	}
}
