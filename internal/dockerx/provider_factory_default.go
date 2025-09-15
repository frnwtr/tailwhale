package dockerx

// NewProvider returns the default container provider for this build.
// Without build tags, this returns a no-op FakeProvider suitable for tests.
func NewProvider() Provider {
    return &FakeProvider{}
}

