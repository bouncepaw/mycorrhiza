package util

import (
	"testing"
)

func TestWikilink(t *testing.T) {
	atHypha := ":example/test"
	results := map[string]string{
		"foo":      "/foo",
		"::foo":    "/:example/foo",
		":bar/foo": "/:bar/foo",
		"/baz":     "/:example/test/baz",
		"./baz":    "/:example/test/baz",
		"../qux":   "/:example/qux",
	}
	for link, expect := range results {
		if res := Wikilink(link, atHypha); expect != res {
			t.Errorf("%s → %s; expected %s", link, res, expect)
		}
	}
}
