package version

import (
	"regexp"
	"runtime/debug"
	"strconv"

	"github.com/bouncepaw/mycorrhiza/help"
)

var tag = "unknown"
var versionRegexp = regexp.MustCompile(`This is documentation for \*\*Mycorrhiza Wiki\*\* (.*).`)

func init() {
	if b, err := help.Get("en"); err == nil {
		matches := versionRegexp.FindSubmatch(b)
		if matches != nil {
			tag = "v" + string(matches[1])
		}
	}
}

func FormatVersion() string {
	var commit, dirty string
	info, ok := debug.ReadBuildInfo()

	if ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				commit = "+" + setting.Value
				if len(commit) > 8 {
					commit = commit[:8]
				}
			} else if setting.Key == "vcs.modified" {
				modified, err := strconv.ParseBool(setting.Value)
				if err == nil && modified {
					dirty = "-dirty"
				}
			}
		}
	}

	return tag + commit + dirty
}
