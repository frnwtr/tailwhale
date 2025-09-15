package core

// ExposureMode defines how services are exposed.
type ExposureMode int

const (
    ModeA ExposureMode = iota // Host + Traefik
    ModeB                     // Per-container sidecar
    ModeC                     // Funnel on Traefik
)

// Service represents a container/service that may be exposed.
type Service struct {
    ID        string
    Name      string
    Host      string
    Ports     []int
    Exposed   bool
    Mode      ExposureMode
    HostAlias string // optional override
}

// NameInput contains data to compute a hostname.
type NameInput struct {
    Container string
    Host      string
    Tailnet   string
}

// HostnameFor returns the hostname for a service depending on the exposure mode.
func HostnameFor(mode ExposureMode, in NameInput) string {
    switch mode {
    case ModeA:
        // <container>.<host>.<tailnet>.ts.net
        if in.Container == "" || in.Host == "" || in.Tailnet == "" {
            return ""
        }
        return in.Container + "." + in.Host + "." + in.Tailnet + ".ts.net"
    case ModeB:
        // <container>.<tailnet>.ts.net
        if in.Container == "" || in.Tailnet == "" {
            return ""
        }
        return in.Container + "." + in.Tailnet + ".ts.net"
    case ModeC:
        // https://<host>.ts.net
        if in.Host == "" {
            return ""
        }
        return "https://" + in.Host + ".ts.net"
    default:
        return ""
    }
}

