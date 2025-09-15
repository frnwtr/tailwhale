//go:build tui

package main

import (
    "context"
    "fmt"
    "os"
    "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/frnwtr/tailwhale/internal/core"
    "github.com/frnwtr/tailwhale/internal/dockerx"
)

type model struct {
    loading bool
    err     error
    items   []core.Service
    host    string
    tailnet string
}

func initialModel() model {
    return model{loading: true, host: "host", tailnet: "tn"}
}

func (m model) Init() tea.Cmd {
    return fetchServices(m.host, m.tailnet)
}

func fetchServices(host, tailnet string) tea.Cmd {
    return func() tea.Msg {
        p := dockerx.NewProvider()
        svcs, err := core.Discover(p, host, tailnet)
        if err != nil {
            return errMsg{err}
        }
        return listMsg{svcs}
    }
}

type listMsg struct{ items []core.Service }
type errMsg struct{ err error }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "r":
            return m, fetchServices(m.host, m.tailnet)
        }
    case listMsg:
        m.loading = false
        m.items = msg.items
    case errMsg:
        m.loading = false
        m.err = msg.err
    case tea.WindowSizeMsg:
        // ignore for now
    }
    return m, nil
}

func (m model) View() string {
    if m.loading {
        return "TailWhale TUI — loading... (q to quit)\n"
    }
    if m.err != nil {
        return fmt.Sprintf("TailWhale TUI — error: %v (r to retry, q to quit)\n", m.err)
    }
    s := "TailWhale TUI — services (r to refresh, q to quit)\n\n"
    if len(m.items) == 0 {
        s += "No services found.\n"
        return s
    }
    for _, it := range m.items {
        s += fmt.Sprintf("• %s → %s\n", it.Name, it.Host)
    }
    return s
}

func main() {
    p := tea.NewProgram(initialModel(), tea.WithContext(context.Background()), tea.WithAltScreen())
    if err := p.Start(); err != nil {
        fmt.Fprintln(os.Stderr, "tui error:", err)
        os.Exit(1)
    }
    _ = time.Second
}

