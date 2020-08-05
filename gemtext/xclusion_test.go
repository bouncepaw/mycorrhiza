package gemtext

import (
	"testing"
)

func TestParseTransclusion(t *testing.T) {
	check := func(line string, expectedXclusion Transclusion) {
		if xcl := parseTransclusion(line); xcl != expectedXclusion {
			t.Error(line, "; got:", xcl, "wanted:", expectedXclusion)
		}
	}
	check("<=  ", Transclusion{"", -9, -9})
	check("<=hypha", Transclusion{"hypha", 0, 0})
	check("<=  hypha\t", Transclusion{"hypha", 0, 0})
	check("<= hypha :", Transclusion{"hypha", 0, 0})
	check("<= hypha : ..", Transclusion{"hypha", 0, 0})
	check("<= hypha : 3", Transclusion{"hypha", 3, 3})
	check("<= hypha : 3..", Transclusion{"hypha", 3, 0})
	check("<= hypha : ..3", Transclusion{"hypha", 0, 3})
	check("<= hypha : 3..4", Transclusion{"hypha", 3, 4})
}
