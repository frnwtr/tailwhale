package dockerx

import "testing"

func TestCacheSnapshotAndUpsert(t *testing.T){
    c := NewCache()
    c.ApplySnapshot([]Info{{ID:"1", Name:"a"}})
    if len(c.List()) != 1 { t.Fatal("expected 1") }
    c.Upsert(Info{ID:"2", Name:"b"})
    if len(c.List()) != 2 { t.Fatal("expected 2") }
    c.Upsert(Info{ID:"1", Name:"a2"})
    got := c.List()
    found := false
    for _, it := range got { if it.ID=="1" && it.Name=="a2" { found=true } }
    if !found { t.Fatal("upsert did not replace") }
    c.Remove("2")
    if len(c.List()) != 1 { t.Fatal("expected 1 after remove") }
}

