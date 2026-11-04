package interwikidomain

type WikiEngine string

const (
	Mycorrhiza WikiEngine = "mycorrhiza"
	Betula     WikiEngine = "betula"
	Agora      WikiEngine = "agora"
	// Generic is any website.
	Generic WikiEngine = "generic"
)

func (we WikiEngine) Valid() bool {
	switch we {
	case Mycorrhiza, Betula, Agora, Generic:
		return true
	}
	return false
}

// WikiDTO is an entry in the interwiki map.
type WikiDTO struct {
	// Name is the name of the wiki, and is also one of the possible prefices.
	Name string `json:"name"`

	// Aliases are alternative prefices you can use instead of Name. This slice can be empty.
	Aliases []string `json:"aliases,omitempty"`

	// URL is the address of the wiki.
	URL string `json:"url"`

	// LinkHrefFormat is a format string for interwiki links. See Mycomarkup internal docs hidden deep inside for more information.
	//
	// This field is optional. If it is not set, it is derived from other data. See the code.
	LinkHrefFormat string `json:"link_href_format"`

	ImgSrcFormat string `json:"img_src_format"`

	// Engine is the engine of the wiki. Invalid values will result in a start-up error.
	Engine WikiEngine `json:"engine"`
}
