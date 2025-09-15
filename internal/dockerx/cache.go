package dockerx

import "sync"

// Cache maintains a snapshot of containers keyed by ID.
type Cache struct {
    mu   sync.RWMutex
    data map[string]Info
}

func NewCache() *Cache { return &Cache{data: make(map[string]Info)} }

// ApplySnapshot replaces the cache with the provided list.
func (c *Cache) ApplySnapshot(items []Info) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data = make(map[string]Info, len(items))
    for _, it := range items {
        c.data[it.ID] = it
    }
}

// Upsert applies a single Info update (insert or replace).
func (c *Cache) Upsert(it Info) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.data == nil { c.data = make(map[string]Info) }
    c.data[it.ID] = it
}

// Remove deletes an entry by ID.
func (c *Cache) Remove(id string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.data, id)
}

// List returns a copy of all entries.
func (c *Cache) List() []Info {
    c.mu.RLock()
    defer c.mu.RUnlock()
    out := make([]Info, 0, len(c.data))
    for _, it := range c.data { out = append(out, it) }
    return out
}

