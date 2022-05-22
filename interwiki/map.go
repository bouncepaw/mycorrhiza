package interwiki

// Map is an interwiki map
type Map struct {
	list   []*Wiki
	byName map[string]*Wiki
}

var theMap Map

func init() {
	theMap.list = []*Wiki{}
	theMap.byName = map[string]*Wiki{}
}
