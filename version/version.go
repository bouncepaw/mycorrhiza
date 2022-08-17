package version

import (
	"fmt"
	"runtime/debug"
	"strconv"
)

// This is set through ldflags='-X ...' in the Makefile
var taggedRelease string = "Unknown"

func FormatVersion() string {
	var commitHash string = "Unknown"
	var dirty string = ""

	info, ok := debug.ReadBuildInfo()

	if ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				commitHash = setting.Value
			} else if setting.Key == "vcs.modified" {
				modified, err := strconv.ParseBool(setting.Value)
				if err == nil && modified {
					dirty = "-dirty"
				}
			}
		}
	}

	return fmt.Sprintf("Mycorrhiza Wiki %s+%s%s", taggedRelease, commitHash[:7], dirty)
}
