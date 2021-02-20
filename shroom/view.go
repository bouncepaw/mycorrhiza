package shroom

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/markup"
	"github.com/bouncepaw/mycorrhiza/util"
)

// FetchTextPart tries to read text file of the given hypha. If there is no file, empty string is returned.
func FetchTextPart(h *hyphae.Hypha) (string, error) {
	if h.TextPath == "" {
		return "", nil
	}
	text, err := ioutil.ReadFile(h.TextPath)
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return string(text), nil
}

// binaryHtmlBlock creates an html block for binary part of the hypha.
func BinaryHtmlBlock(h *hyphae.Hypha) string {
	switch filepath.Ext(h.BinaryPath) {
	case ".jpg", ".gif", ".png", ".webp", ".svg", ".ico":
		return fmt.Sprintf(`
		<div class="binary-container binary-container_with-img">
			<a href="/binary/%[1]s"><img src="/binary/%[1]s"/></a>
		</div>`, h.Name)
	case ".ogg", ".webm", ".mp4":
		return fmt.Sprintf(`
		<div class="binary-container binary-container_with-video">
			<video controls>
				<source src="/binary/%[1]s"/>
				<p>Your browser does not support video. <a href="/binary/%[1]s">Download video</a></p>
			</video>
		`, h.Name)
	case ".mp3":
		return fmt.Sprintf(`
		<div class="binary-container binary-container_with-audio">
			<audio controls>
				<source src="/binary/%[1]s"/>
				<p>Your browser does not support audio. <a href="/binary/%[1]s">Download audio</a></p>
			</audio>
		`, h.Name)
	default:
		return fmt.Sprintf(`
		<div class="binary-container binary-container_with-nothing">
			<p><a href="/binary/%s">Download media</a></p>
		</div>
		`, h.Name)
	}
}

func SetHeaderLinks() {
	if userLinksHypha := hyphae.ByName(util.HeaderLinksHypha); !userLinksHypha.Exists {
		util.SetDefaultHeaderLinks()
	} else {
		contents, err := ioutil.ReadFile(userLinksHypha.TextPath)
		if err != nil || len(contents) == 0 {
			util.SetDefaultHeaderLinks()
		} else {
			text := string(contents)
			util.ParseHeaderLinks(text, markup.Rocketlink)
		}
	}
}
