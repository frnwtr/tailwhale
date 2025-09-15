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
    ctx, cancel := context.WithCancel(context.Background())
    ch := make(chan servicesMsg, 4)
    m.updates = ch
    m.ctx = ctx
    m.cancel = cancel
    return tea.Batch(startWatch(m.host, m.tailnet, ctx, ch), fetchServices(m.host, m.tailnet), listen(ch))
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
type servicesMsg struct{ items []core.Service }
type errMsg struct{ err error }

func startWatch(host, tailnet string, ctx context.Context, ch chan servicesMsg) tea.Cmd {
    return func() tea.Msg {
        go func() {
            provider := dockerx.NewProvider()
            orch := core.Orchestrator{Provider: provider, Host: host, Tailnet: tailnet}
            _ = orch.Watch(ctx, 5*time.Second, func(svcs []core.Service, _ traefik.TLSConfig) {
                select {
                case ch <- servicesMsg{items: svcs}:
                case <-ctx.Done():
                }
            })
        }()
        return nil
    }
}

func listen(ch <-chan servicesMsg) tea.Cmd { return func() tea.Msg { return <-ch } }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            if m.cancel != nil { m.cancel() }
            return m, tea.Quit
        case "up", "k":
            if m.selected > 0 { m.selected-- }
        case "down", "j":
            if m.selected < len(m.items)-1 { m.selected++ }
        case "r":
            return m, fetchServices(m.host, m.tailnet)
        }
    case listMsg:
        m.loading = false
        m.items = msg.items
        if m.selected >= len(m.items) { m.selected = len(m.items) - 1 }
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
    header := "TailWhale TUI — services (↑/k, ↓/j, r to refresh, q to quit)\n\n"
    if len(m.items) == 0 {
        return header + "No services found.\n"
    }
    s := header
    for i, it := range m.items {
        cursor := "  "
        if i == m.selected { cursor = "> " }
        s += fmt.Sprintf("%s%s → %s\n", cursor, it.Name, it.Host)
    }
    sel := m.items[m.selected]
    s += "\nDetails:\n"
    s += fmt.Sprintf("  Name: %s\n  Host: %s\n  Mode: %d\n  Ports: %v\n", sel.Name, sel.Host, sel.Mode, sel.Ports)
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
