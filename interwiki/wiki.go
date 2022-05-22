package interwiki

// WikiEngine is an enumeration of supported interwiki targets.
type WikiEngine int

const (
	// Mycorrhiza is a Mycorrhiza wiki. This is the default value.
	Mycorrhiza WikiEngine = iota
	// Generic is any website.
	Generic
)

// Wiki is an entry in the interwiki map.
type Wiki struct {
	// Names is a slice of link prefices that correspond to this wiki.
	Names []string `json:"names"`

	// URL is the address of the wiki.
	URL string `json:"url"`

	// LinkFormat is a format string for incoming interwiki links. The format strings should look like this:
	//     http://wiki.example.org/view/%s
	// where %s is where text will be inserted. No other % instructions are supported yet. They will be added once we learn of their use cases.
	//
	// This field is optional. For Generic wikis, it is automatically set to <URL>/%s; for Mycorrhiza wikis, it is automatically set to <URL>/hypha/%s.
	LinkFormat string `json:"link_format"`

	// Description is a plain-text description of the wiki.
	Description string `json:"description"`

	// Engine is the engine of the wiki. This field is not set in JSON.
	Engine WikiEngine `json:"-"`

	// EngineString is a string name of the engine. It is then converted to Engine. Supported values are:
	//     * mycorrhiza -> Mycorrhiza
	//     * generic -> Generic
	// All other values will result in an error.
	EngineString string `json:"engine"`
}

func (w *Wiki) canonize() {

}
