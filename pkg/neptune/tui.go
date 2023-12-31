package neptune

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jedib0t/go-pretty/v6/table"
)

var (
	currentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	donePkgNameStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("140"))
	doneStyle           = lipgloss.NewStyle().Margin(1, 2)
	checkMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	helpStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
	baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
)

const (
	padding  = 2
	maxWidth = 80
)

type progressMsg float64

type progressErrMsg struct{ err error }
type progressDone bool

func finalPause() tea.Cmd {
	return tea.Tick(time.Millisecond*750, func(_ time.Time) tea.Msg {
		return nil
	})
}

type model struct {
	pw       *progressWriter
	s        spinner.Model
	progress progress.Model
	width    int
	height   int
	err      error
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.s.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil
	case progressErrMsg:
		m.err = msg.err
		return m, tea.Quit
	case progressMsg:
		var cmds []tea.Cmd
		if msg >= 1.0 {
			cmds = append(cmds, tea.Sequence(finalPause(), tea.Quit))
			cmds = append(cmds, tea.Printf("%s %s", checkMark, soundsInfo[sounds[Xindex]]))
			return m, tea.Batch(cmds...)
		}
		cmds = append(cmds, m.progress.SetPercent(float64(msg)))
		return m, tea.Batch(cmds...)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.s, cmd = m.s.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	default:
		return m, nil
	}
}

func (m model) View() string {
	if m.err != nil {
		return "Error downloading: " + m.err.Error() + "\n"
	}
	n := len(sounds)
	w := lipgloss.Width(fmt.Sprintf("%d", n))
	pkgCount := fmt.Sprintf(" %*d/%*d", w, Xindex, w, n)
	spin := m.s.View() + " "
	prog := m.progress.View()
	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+pkgCount))
	pkgName := currentPkgNameStyle.Render(soundsInfo[sounds[Xindex]])
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Installing " + pkgName)
	pad := strings.Repeat(" ", padding)

	return spin + info + pad + prog + pkgCount
}

// Print Table with Available Sounds
func PrintTableWithAliens(list []string) {
	alien := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("♪ ")
	// Generate a random color for each row
	colors := generateRandomColors(len(list))
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")).Render("Available Sound Keys"))
	for i, str := range list {
		row := table.Row{alien+lipgloss.NewStyle().Foreground(lipgloss.Color(colors[i])).Render(str)}
		t.AppendRow(row)
	}
	t.Render()
}

// generate random colors
func generateRandomColors(count int) []string {
	colors := make([]string, count)
	for i := 0; i < count; i++ {
		color := fmt.Sprintf("#%06x", rand.Intn(0xffffff))
		colors[i] = color
	}
	return colors
}
