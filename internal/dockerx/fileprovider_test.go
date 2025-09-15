package dockerx

import (
    "bytes"
    "testing"
)

func TestFileProviderList(t *testing.T){
    data := []byte(`[
        {"ID":"1","Name":"a","Labels":{"tailwhale.enable":"true"},"Ports":[80],"Running":true},
        {"ID":"2","Name":"b","Labels":{"tailwhale.enable":"false"},"Ports":[8080],"Running":true}
    ]`)
    items, err := decodeInfos(bytes.NewReader(data))
    if err != nil { t.Fatal(err) }
    if len(items) != 2 || items[0].Name != "a" { t.Fatalf("unexpected: %+v", items) }
}

