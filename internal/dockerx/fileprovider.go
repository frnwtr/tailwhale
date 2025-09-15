package dockerx

import (
    "encoding/json"
    "errors"
    "io"
    "os"
)

// FileProvider implements Provider by reading a static JSON file of []Info.
// It is intended for testing and offline development.
type FileProvider struct{
    Path string
}

func (p *FileProvider) List() ([]Info, error){
    if p.Path == "" { return nil, errors.New("file provider: empty path") }
    f, err := os.Open(p.Path)
    if err != nil { return nil, err }
    defer f.Close()
    return decodeInfos(f)
}

func (p *FileProvider) Watch() (Watcher, error){
    // No events for file provider; caller can poll List periodically.
    return &FakeWatcher{}, nil
}

func decodeInfos(r io.Reader) ([]Info, error){
    dec := json.NewDecoder(r)
    var items []Info
    if err := dec.Decode(&items); err != nil { return nil, err }
    return items, nil
}

