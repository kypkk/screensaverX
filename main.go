package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

// exitKey is the only key combination that ends the screensaver. Change this
// constant to use a different combination — see README.md for the strings
// bubbletea recognizes and the terminal-compatibility caveats.
const exitKey = "ctrl+q"

func main() {
	// If we're inside tmux, hand off control to the outer shell and exit.
	if handleTmux() {
		return
	}
	runScreensaver()
}

func runScreensaver() {
	// Make Ctrl+C / Ctrl+Z / Ctrl+\ / kill PID etc. impotent.
	blockExitSignals()

	m := InitialModel()
	// New Program / Initial Model. WithoutSignalHandler stops bubbletea from
	// installing its own SIGINT/SIGTERM handler that would race ours.
	p := tea.NewProgram(m,
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
		tea.WithoutSignalHandler(),
	)
	// Run
	final, err := p.Run()
	if err != nil {
		os.Exit(1)
	}
	// Wipe the terminal (screen + scrollback) the same way `clear` does, so
	// no prior shell history is visible after the break.
	fmt.Print("\033[2J\033[3J\033[H")
	if fm, ok := final.(model); ok {
		fmt.Printf("Break complete. Total time: %s\n", formatDuration(time.Since(fm.start)))
	}
}

// tickMsg fires every second so the on-screen timer can re-render.
type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Model: The BubbleTea App State
type model struct {
	title   string
	start   time.Time
	elapsed time.Duration
	width   int
	height  int
}

// Initial a new model
func InitialModel() model {
	return model{
		title: "Currently taking a break......",
		start: time.Now(),
	}
}

// Init: kick off the event loop
func (m model) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle(m.title), tickCmd())
}

// Update: handle Msgs
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// setting the window height and width for the layout
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	// keep the live timer ticking
	case tickMsg:
		m.elapsed = time.Since(m.start)
		return m, tickCmd()
	// Ignore mouse entirely.
	case tea.MouseMsg:
		return m, nil
	// Only the exit key combo can quit; everything else is swallowed.
	case tea.KeyMsg:
		if msg.String() == exitKey {
			return m, tea.Quit
		}
		return m, nil
	}
	return m, nil
}

// View: return a string based on the state of our model
func (m model) View() string {
	body := lipgloss.JoinVertical(
		lipgloss.Center,
		m.title,
		"",
		formatDuration(m.elapsed),
	)
	s := fmt.Sprintln(lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("36")).
		Padding(0, 4).Align(lipgloss.Center).Render(body))

	s = lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		s,
	)
	return s
}

// formatDuration renders a duration as HH:MM:SS, rounded to whole seconds.
func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	d = d.Round(time.Second)
	h := int(d / time.Hour)
	m := int(d % time.Hour / time.Minute)
	s := int(d % time.Minute / time.Second)
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
