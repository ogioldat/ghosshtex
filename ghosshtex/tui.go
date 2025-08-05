package ghosshtex

import (
	"io"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	viewport viewport.Model
	textarea textarea.Model
	err      error
}

func NewEditorTUI(output io.Writer) *tea.Program {
	ta := textarea.New()
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent("Hej, tu bubbletea")

	initialModel := model{
		textarea: ta,
		viewport: vp,
		err:      nil,
	}

	return tea.NewProgram(initialModel, tea.WithOutput(output))
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return "View from TUI"
}
