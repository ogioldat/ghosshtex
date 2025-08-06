package ghosshtex

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	textarea textarea.Model
	err      error
}

type errMsg error

func NewEditorTUI(output io.Writer, input io.Reader) *tea.Program {
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
		err:      nil,
	}

	return tea.NewProgram(
		initialModel,
		tea.WithOutput(output),
		tea.WithInput(input),
	)
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyCtrlD, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"Welcome to text editor\n\n%s\n\n%s",
		m.textarea.View(),
		"(esc to quit)",
	) + "\n"
}
