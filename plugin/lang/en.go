package lang

var EnglishMap = map[string]string{
	"edit hypha title template": "Edit %s at MycorrhizaWiki",
	"view hypha title template": "%s at MycorrhizaWiki",
	"this site runs myco wiki":  `<p>This website runs <a href="https://github.com/bouncepaw/mycorrhiza">MycorrhizaWiki</a></p>`,
	"generic error msg":         `<b>Sorry, something went wrong</b>`,

	"edit/text mime type":     "Text MIME-type",
	"edit/text mime type/tip": "We support <code>text/markdown</code>, <code>text/creole</code> and <code>text/gemini</code>",

	"edit/revision comment":     "Revision comment",
	"edit/revision comment/tip": "Please make your comment helpful",
	"edit/revision comment/new": "Create %s",
	"edit/revision comment/old": "Update %s",

	"edit/tags":     "Edit tags",
	"edit/tags/tip": "Tags are separated by commas, whitespace is ignored",

	"edit/upload file":     "Upload file",
	"edit/upload file/tip": "Only images are fully supported for now",

	"edit/box":              "Edit box",
	"edit/box/title":        "Edit %s",
	"edit/box/help pattern": "Describe %s here",

	"edit/cancel": "Cancel",

	"update ok/title": "Saved %s",
	"update ok/msg":   "Saved successfully. <a href='/%s'>Go back</a>",
}
