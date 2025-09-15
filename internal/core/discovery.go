package core

import (
    "sort"

    "github.com/frnwtr/tailwhale/internal/dockerx"
)

// Discover returns the list of services to expose based on container labels.
func Discover(p dockerx.Provider, host, tailnet string) ([]Service, error) {
    list, err := p.List()
    if err != nil {
        return nil, err
    }
    var out []Service
    for _, c := range list {
        if c.Labels[LabelEnable] != "true" {
            continue
        }
        mode := ParseMode(c.Labels[LabelMode])
        svc := Service{
            ID:      c.ID,
            Name:    c.Name,
            Ports:   c.Ports,
            Exposed: true,
            Mode:    mode,
        }
        if h := c.Labels[LabelHost]; h != "" {
            svc.HostAlias = h
            svc.Host = h
        } else {
            svc.Host = HostnameFor(mode, NameInput{Container: c.Name, Host: host, Tailnet: tailnet})
        }
        out = append(out, svc)
    }
    sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
    return out, nil
}

