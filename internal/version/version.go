package version

import (
	"regexp"
	"runtime/debug"
	"strconv"

	"github.com/bouncepaw/mycorrhiza/help"
)

// Long is the full version string, including VCS information, that looks like
// x.y.z+hash-dirty.
var Long string

// Short is the human-friendly x.y.z part of the long version string.
var Short string

var versionRegexp = regexp.MustCompile(`This is documentation for Mycorrhiza Wiki (.*)\. `)

func init() {
	if b, err := help.Get("en"); err == nil {
		matches := versionRegexp.FindSubmatch(b)
		if matches != nil {
			Short = string(matches[1])
		}
	}

	Long = Short
	info, ok := debug.ReadBuildInfo()
	if ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				val := setting.Value
				if len(val) > 7 {
					val = val[:7]
				}
				Long += "+" + val
			} else if setting.Key == "vcs.modified" {
				modified, err := strconv.ParseBool(setting.Value)
				if err == nil && modified {
					Long += "-dirty"
				}
			}
		}
	}
}
