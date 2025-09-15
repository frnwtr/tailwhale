package fsx

import (
    "bytes"
    "os"
    "path/filepath"
    "testing"
)

func TestWriteFileAtomic(t *testing.T){
    dir := t.TempDir()
    path := filepath.Join(dir, "a", "b", "file.txt")
    if err := WriteFileAtomic(path, []byte("hello"), 0o644); err != nil { t.Fatal(err) }
    got, err := os.ReadFile(path)
    if err != nil { t.Fatal(err) }
    if !bytes.Equal(got, []byte("hello")) { t.Fatalf("wrong content: %q", string(got)) }
}

