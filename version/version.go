package version

// These are set through ldflags='-X ...' in the Makefile
var (
	TaggedRelease string
	CommitHash    string
)
