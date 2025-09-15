//go:build docker

package dockerx

import (
    "context"

    "github.com/docker/docker/api/types"
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
    // Minimal skeleton: return a no-op watcher for now.
    // TODO: hook into Docker events and emit Info changes.
    return &FakeWatcher{}, nil
}

func trimSlash(s string) string {
    if len(s) > 0 && s[0] == '/' { return s[1:] }
    return s
}

