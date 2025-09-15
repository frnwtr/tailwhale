//go:build docker

package dockerx

import (
    "context"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/events"
    "github.com/docker/docker/api/types/filters"
    "github.com/docker/docker/client"
)

// NewProvider returns a real Docker API provider (requires -tags docker).
func NewProvider() Provider {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return &FakeProvider{}
    }
    return &DockerProvider{cli: cli}
}

type DockerProvider struct{
    cli *client.Client
}

func (p *DockerProvider) List() ([]Info, error){
    ctx := context.Background()
    cs, err := p.cli.ContainerList(ctx, types.ContainerListOptions{All: true})
    if err != nil { return nil, err }
    out := make([]Info, 0, len(cs))
    for _, c := range cs {
        labels := map[string]string{}
        for k, v := range c.Labels { labels[k] = v }
        ports := make([]int, 0, len(c.Ports))
        for _, p := range c.Ports { ports = append(ports, int(p.PublicPort)) }
        name := ""
        if len(c.Names) > 0 { name = c.Names[0] }
        out = append(out, Info{ID: c.ID, Name: trimSlash(name), Labels: labels, Ports: ports, Running: c.State == "running"})
    }
    return out, nil
}

func (p *DockerProvider) Watch() (Watcher, error){
    ctx, cancel := context.WithCancel(context.Background())
    args := filters.NewArgs()
    args.Add("type", "container")
    msgs, errs := p.cli.Events(ctx, types.EventsOptions{Filters: args})
    w := &dockerWatcher{cli: p.cli, ctx: ctx, cancel: cancel, msgs: msgs, errs: errs, out: make(chan Info)}
    go w.loop()
    return w, nil
}

func trimSlash(s string) string {
    if len(s) > 0 && s[0] == '/' { return s[1:] }
    return s
}

type dockerWatcher struct{
    cli   *client.Client
    ctx   context.Context
    cancel context.CancelFunc
    msgs  <-chan events.Message
    errs  <-chan error
    out   chan Info
}

func (w *dockerWatcher) loop(){
    defer close(w.out)
    for {
        select {
        case <-w.ctx.Done():
            return
        case err := <-w.errs:
            _ = err // propagate as no event; Next will see closed channel
            return
        case m, ok := <-w.msgs:
            if !ok { return }
            // Filter by common actions
            switch m.Action {
            case "start", "stop", "die", "pause", "unpause", "update", "destroy":
                // Inspect container and emit Info
                id := m.Actor.ID
                if id == "" { continue }
                go func(cid string){
                    // separate context to avoid blocking main loop
                    info := Info{ID: cid}
                    // Best-effort inspect; ignore errors
                    if json, err := w.cli.ContainerInspect(w.ctx, cid); err == nil {
                        name := trimSlash(json.Name)
                        labels := map[string]string{}
                        for k, v := range json.Config.Labels { labels[k] = v }
                        ports := make([]int, 0, len(json.NetworkSettings.Ports))
                        for p := range json.NetworkSettings.Ports { if p.Int() > 0 { ports = append(ports, p.Int()) } }
                        info.Name = name
                        info.Labels = labels
                        info.Ports = ports
                        info.Running = json.State != nil && json.State.Running
                    }
                    select { case w.out <- info: case <-w.ctx.Done(): }
                }(id)
            }
        }
    }
}

func (w *dockerWatcher) Next() (Info, bool, error){
    i, ok := <-w.out
    return i, ok, nil
}

func (w *dockerWatcher) Close() error { w.cancel(); return nil }
