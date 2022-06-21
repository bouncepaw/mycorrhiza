package interwiki

var (
	listOfEntries []*Wiki
	entriesByName map[string]*Wiki
)

func init() {
	listOfEntries = []*Wiki{}
	entriesByName = map[string]*Wiki{}
}
