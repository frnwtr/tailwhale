package fsx

import (
    "io"
    "os"
    "path/filepath"
)

// WriteFileAtomic writes data to path atomically via a temp file and rename.
func WriteFileAtomic(path string, data []byte, perm os.FileMode) error {
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0o755); err != nil { return err }
    f, err := os.CreateTemp(dir, ".tmp-*")
    if err != nil { return err }
    tmp := f.Name()
    defer func(){ _ = os.Remove(tmp) }()
    if _, err := f.Write(data); err != nil { f.Close(); return err }
    if err := f.Sync(); err != nil { f.Close(); return err }
    if err := f.Close(); err != nil { return err }
    if err := os.Chmod(tmp, perm); err != nil { return err }
    return os.Rename(tmp, path)
}

// CopyAtomic copies from r to path atomically.
func CopyAtomic(path string, r io.Reader, perm os.FileMode) error {
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0o755); err != nil { return err }
    f, err := os.CreateTemp(dir, ".tmp-*")
    if err != nil { return err }
    tmp := f.Name()
    defer func(){ _ = os.Remove(tmp) }()
    if _, err := io.Copy(f, r); err != nil { f.Close(); return err }
    if err := f.Sync(); err != nil { f.Close(); return err }
    if err := f.Close(); err != nil { return err }
    if err := os.Chmod(tmp, perm); err != nil { return err }
    return os.Rename(tmp, path)
}

