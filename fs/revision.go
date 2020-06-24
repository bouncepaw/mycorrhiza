package fs

type Revision struct {
	Id         int      `json:"-"`
	FullName   string   `json:"-"`
	Tags       []string `json:"tags"`
	ShortName  string   `json:"name"`
	Comment    string   `json:"comment"`
	Author     string   `json:"author"`
	Time       int      `json:"time"`
	TextMime   string   `json:"text_mime"`
	BinaryMime string   `json:"binary_mime"`
	TextPath   string   `json:"-"`
	BinaryPath string   `json:"-"`
	TextName   string   `json:"text_name"`
	BinaryName string   `json:"binary_name"`
}
