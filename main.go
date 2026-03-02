package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")
	s.InputField = lipgloss.
		NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.NormalBorder()).
		Padding(1).
		Width(80)
	return s
}

type model struct {
	styles    *Styles
	index     int
	questions []Question
	width     int
	height    int
	done      bool
}

type Question struct {
	question string
	answer   string
	input    Input
}

func NewQuestion(q string) Question {
	return Question{question: q}
}

func newShortQuestion(question string) Question {
	q := NewQuestion(question)
	model := NewShortAnswerField()
	q.input = model
	return q
}

func newLongQuestion(question string) Question {
	q := NewQuestion(question)
	model := NewLongAnswerField()
	q.input = model
	return q
}

func New(questions []Question) *model {
	styles := DefaultStyles()
	return &model{
		questions: questions,
		styles:    styles,
	}
}

func (m model) Init() tea.Cmd {
	return m.questions[m.index].input.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	current := &m.questions[m.index]
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.index == len(m.questions)-1 {
				m.done = true
			}
			current.answer = current.input.Value()
			m.Next()
			return m, current.input.Blur
		}
	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	current := m.questions[m.index]
	if m.done {
		var output string
		for _, q := range m.questions {
			output += fmt.Sprintf("%s: %s\n", q.question, q.answer)
		}
		return output
	}
	if m.width == 0 {
		return "loading..."
	}
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			current.question,
			m.styles.InputField.Render(current.input.View())),
	)
}

func (m *model) Next() {
	if m.index < len(m.questions)-1 {
		m.index++
	} else {
		m.index = 0
	}
}

func main() {
	questions := []Question{
		newShortQuestion("what is your name?"),
		newShortQuestion("what is your favorite editor?"),
		newLongQuestion("what is your favorite quote?")}
	m := New(questions)
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	p := tea.NewProgram(*m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
