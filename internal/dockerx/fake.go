package dockerx

type FakeProvider struct{
    Items []Info
}

func (f *FakeProvider) List() ([]Info, error){
    // return a copy
    out := make([]Info, len(f.Items))
    copy(out, f.Items)
    return out, nil
}

func (f *FakeProvider) Watch() (Watcher, error){
    return &FakeWatcher{}, nil
}

type FakeWatcher struct{}

func (w *FakeWatcher) Next() (Info, bool, error){
    return Info{}, false, nil
}
func (w *FakeWatcher) Close() error { return nil }

