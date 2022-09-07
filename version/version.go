package version

import (
	"runtime/debug"
	"strconv"
)

// This is set through ldflags='-X ...' in the Makefile
var taggedRelease string = "unknown"

func FormatVersion() string {
	var commitHash string = ""
	var dirty string = ""

	info, ok := debug.ReadBuildInfo()

	if ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				commitHash = "+" + setting.Value
				if len(commitHash) > 8 {
					commitHash = commitHash[:8]
				}
			} else if setting.Key == "vcs.modified" {
				modified, err := strconv.ParseBool(setting.Value)
				if err == nil && modified {
					dirty = "-dirty"
				}
			}
		}
	}

	return taggedRelease + commitHash + dirty
}
