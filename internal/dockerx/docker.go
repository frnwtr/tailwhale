package dockerx

// Info represents a subset of container metadata we care about.
type Info struct {
    ID      string
    Name    string
    Labels  map[string]string
    Ports   []int
    Running bool
    // Event carries a recent event action (e.g., start, stop, destroy) when originating from a watcher.
    Event   string
}

// Watcher emits container events (start/stop/label changes).
type Watcher interface {
    // Next blocks until an event or error occurs. ok=false on closed.
    Next() (Info, bool, error)
    // Close releases watcher resources.
    Close() error
}

// Provider discovers containers (snapshot) and returns a Watcher for changes.
type Provider interface {
    List() ([]Info, error)
    Watch() (Watcher, error)
}
