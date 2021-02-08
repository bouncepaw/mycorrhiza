package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/markup"
	"github.com/bouncepaw/mycorrhiza/mimetype"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
)

func init() {
	markup.HyphaExists = func(hyphaName string) bool {
		_, hyphaExists := HyphaStorage[hyphaName]
		return hyphaExists
	}
	markup.HyphaAccess = func(hyphaName string) (rawText, binaryBlock string, err error) {
		if hyphaData, ok := HyphaStorage[hyphaName]; ok {
			rawText, err = FetchTextPart(hyphaData)
			if hyphaData.BinaryPath != "" {
				binaryBlock = binaryHtmlBlock(hyphaName, hyphaData)
			}
		} else {
			err = errors.New("Hypha " + hyphaName + " does not exist")
		}
		return
	}
	markup.HyphaIterate = IterateHyphaNamesWith
	markup.HyphaImageForOG = func(hyphaName string) string {
		if hd, isOld := GetHyphaData(hyphaName); isOld && hd.BinaryPath != "" {
			return util.URL + "/binary/" + hyphaName
		}
		return util.URL + "/favicon.ico"
	}
}

// GetHyphaData finds a hypha addressed by `hyphaName` and returns its `hyphaData`. `hyphaData` is set to a zero value if this hypha does not exist. `isOld` is false if this hypha does not exist.
func GetHyphaData(hyphaName string) (hyphaData *HyphaData, isOld bool) {
	hyphaData, isOld = HyphaStorage[hyphaName]
	if hyphaData == nil {
		hyphaData = &HyphaData{}
	}
	return
}

// HyphaData represents a hypha's meta information: binary and text parts rooted paths and content types.
type HyphaData hyphae.Hypha

// uploadHelp is a helper function for UploadText and UploadBinary
func uploadHelp(hop *history.HistoryOp, hyphaName, ext string, data []byte, u *user.User) *history.HistoryOp {
	var (
		hyphaData, isOld = GetHyphaData(hyphaName)
		fullPath         = filepath.Join(WikiDir, hyphaName+ext)
		originalFullPath = &hyphaData.TextPath
	)
	if hop.Type == history.TypeEditBinary {
		originalFullPath = &hyphaData.BinaryPath
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0777); err != nil {
		return hop.WithError(err)
	}

	if err := ioutil.WriteFile(fullPath, data, 0644); err != nil {
		return hop.WithError(err)
	}

	if isOld && *originalFullPath != fullPath && *originalFullPath != "" {
		if err := history.Rename(*originalFullPath, fullPath); err != nil {
			return hop.WithError(err)
		}
		log.Println("Move", *originalFullPath, "to", fullPath)
	}
	// New hyphae must be added to the hypha storage
	if !isOld {
		HyphaStorage[hyphaName] = hyphaData
		hyphae.IncrementCount()
	}
	*originalFullPath = fullPath
	if isOld && hop.Type == history.TypeEditText && !history.FileChanged(fullPath) {
		return hop.Abort()
	}
	return hop.WithFiles(fullPath).
		WithUser(u).
		Apply()
}

// UploadText loads a new text part from `textData` for hypha `hyphaName`.
func UploadText(hyphaName, textData string, u *user.User) *history.HistoryOp {
	return uploadHelp(
		history.
			Operation(history.TypeEditText).
			WithMsg(fmt.Sprintf("Edit ‘%s’", hyphaName)),
		hyphaName, ".myco", []byte(textData), u)
}

// UploadBinary loads a new binary part from `file` for hypha `hyphaName` with `hd`. The contents have the specified `mime` type. It must be marked if the hypha `isOld`.
func UploadBinary(hyphaName, mime string, file multipart.File, u *user.User) *history.HistoryOp {
	var (
		hop       = history.Operation(history.TypeEditBinary).WithMsg(fmt.Sprintf("Upload binary part for ‘%s’ with type ‘%s’", hyphaName, mime))
		data, err = ioutil.ReadAll(file)
	)
	if err != nil {
		return hop.WithError(err).Apply()
	}
	return uploadHelp(hop, hyphaName, mimetype.ToExtension(mime), data, u)
}

// DeleteHypha deletes hypha and makes a history record about that.
func (hd *HyphaData) DeleteHypha(hyphaName string, u *user.User) *history.HistoryOp {
	hop := history.Operation(history.TypeDeleteHypha).
		WithFilesRemoved(hd.TextPath, hd.BinaryPath).
		WithMsg(fmt.Sprintf("Delete ‘%s’", hyphaName)).
		WithUser(u).
		Apply()
	if len(hop.Errs) == 0 {
		delete(HyphaStorage, hyphaName)
		hyphae.DecrementCount()
	}
	return hop
}

// UnattachHypha unattaches hypha and makes a history record about that.
func (hd *HyphaData) UnattachHypha(hyphaName string, u *user.User) *history.HistoryOp {
	hop := history.Operation(history.TypeUnattachHypha).
		WithFilesRemoved(hd.BinaryPath).
		WithMsg(fmt.Sprintf("Unattach ‘%s’", hyphaName)).
		WithUser(u).
		Apply()
	if len(hop.Errs) == 0 {
		hd, ok := HyphaStorage[hyphaName]
		if ok {
			if hd.BinaryPath != "" {
				hd.BinaryPath = ""
			}
			// If nothing is left of the hypha
			if hd.TextPath == "" {
				delete(HyphaStorage, hyphaName)
				hyphae.DecrementCount()
			}
		}
	}
	return hop
}

func findHyphaeToRename(hyphaName string, recursive bool) []string {
	hyphae := []string{hyphaName}
	if recursive {
		hyphae = append(hyphae, util.FindSubhyphae(hyphaName, IterateHyphaNamesWith)...)
	}
	return hyphae
}

func renamingPairs(hyphaNames []string, replaceName func(string) string) (map[string]string, error) {
	renameMap := make(map[string]string)
	for _, hn := range hyphaNames {
		if hd, ok := HyphaStorage[hn]; ok {
			if _, nameUsed := HyphaStorage[replaceName(hn)]; nameUsed {
				return nil, errors.New("Hypha " + replaceName(hn) + " already exists")
			}
			if hd.TextPath != "" {
				renameMap[hd.TextPath] = replaceName(hd.TextPath)
			}
			if hd.BinaryPath != "" {
				renameMap[hd.BinaryPath] = replaceName(hd.BinaryPath)
			}
		}
	}
	return renameMap, nil
}

// word Data is plural here
func relocateHyphaData(hyphaNames []string, replaceName func(string) string) {
	for _, hyphaName := range hyphaNames {
		if hd, ok := HyphaStorage[hyphaName]; ok {
			hd.TextPath = replaceName(hd.TextPath)
			hd.BinaryPath = replaceName(hd.BinaryPath)
			HyphaStorage[replaceName(hyphaName)] = hd
			delete(HyphaStorage, hyphaName)
		}
	}
}

// RenameHypha renames hypha from old name `hyphaName` to `newName` and makes a history record about that. If `recursive` is `true`, its subhyphae will be renamed the same way.
func RenameHypha(hyphaName, newName string, recursive bool, u *user.User) *history.HistoryOp {
	var (
		re          = regexp.MustCompile(`(?i)` + hyphaName)
		replaceName = func(str string) string {
			return re.ReplaceAllString(CanonicalName(str), newName)
		}
		hyphaNames     = findHyphaeToRename(hyphaName, recursive)
		renameMap, err = renamingPairs(hyphaNames, replaceName)
		renameMsg      = "Rename ‘%s’ to ‘%s’"
		hop            = history.Operation(history.TypeRenameHypha)
	)
	if err != nil {
		hop.Errs = append(hop.Errs, err)
		return hop
	}
	if recursive {
		renameMsg += " recursively"
	}
	hop.WithFilesRenamed(renameMap).
		WithMsg(fmt.Sprintf(renameMsg, hyphaName, newName)).
		WithUser(u).
		Apply()
	if len(hop.Errs) == 0 {
		relocateHyphaData(hyphaNames, replaceName)
	}
	return hop
}

// binaryHtmlBlock creates an html block for binary part of the hypha.
func binaryHtmlBlock(hyphaName string, hd *HyphaData) string {
	switch filepath.Ext(hd.BinaryPath) {
	case ".jpg", ".gif", ".png", ".webp", ".svg", ".ico":
		return fmt.Sprintf(`
		<div class="binary-container binary-container_with-img">
			<a href="/binary/%[1]s"><img src="/binary/%[1]s"/></a>
		</div>`, hyphaName)
	case ".ogg", ".webm", ".mp4":
		return fmt.Sprintf(`
		<div class="binary-container binary-container_with-video">
			<video>
				<source src="/binary/%[1]s"/>
				<p>Your browser does not support video. See video's <a href="/binary/%[1]s">direct url</a></p>
			</video>
		`, hyphaName)
	case ".mp3":
		return fmt.Sprintf(`
		<div class="binary-container binary-container_with-audio">
			<audio>
				<source src="/binary/%[1]s"/>
				<p>Your browser does not support audio. See audio's <a href="/binary/%[1]s">direct url</a></p>
			</audio>
		`, hyphaName)
	default:
		return fmt.Sprintf(`
		<div class="binary-container binary-container_with-nothing">
			<p>This hypha's media cannot be rendered. <a href="/binary/%s">Download it</a></p>
		</div>
		`, hyphaName)
	}
}

// FetchTextPart tries to read text file in the `d`. If there is no file, empty string is returned.
func FetchTextPart(d *HyphaData) (string, error) {
	if d.TextPath == "" {
		return "", nil
	}
	_, err := os.Stat(d.TextPath)
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	text, err := ioutil.ReadFile(d.TextPath)
	if err != nil {
		return "", err
	}
	return string(text), nil
}

func setHeaderLinks() {
	if userLinksHypha, ok := GetHyphaData(util.HeaderLinksHypha); !ok {
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

func HyphaToTemporaryWorkaround(h *hyphae.Hypha) *HyphaData {
	return &HyphaData{
		Name:       h.Name,
		TextPath:   h.TextPath,
		BinaryPath: h.BinaryPath,
	}
}

// MergeIn merges in content file paths from a different hypha object. Prints warnings sometimes.
func (h *HyphaData) MergeIn(oh *hyphae.Hypha) {
	if h.TextPath == "" && oh.TextPath != "" {
		h.TextPath = oh.TextPath
	}
	if oh.BinaryPath != "" {
		if h.BinaryPath != "" {
			log.Println("There is a file collision for binary part of a hypha:", h.BinaryPath, "and", oh.BinaryPath, "-- going on with the latter")
		}
		h.BinaryPath = oh.BinaryPath
	}
}

// Index finds all hypha files in the full `path` and saves them to HyphaStorage. This function is recursive.
func Index(path string) {
	ch := make(chan *hyphae.Hypha, 5)

	go func() {
		hyphae.Index(path, 0, ch)
		close(ch)
	}()

	for h := range ch {
		if oldHypha, ok := HyphaStorage[h.Name]; ok {
			oldHypha.MergeIn(h)
		} else {
			HyphaStorage[h.Name] = HyphaToTemporaryWorkaround(h)
			hyphae.IncrementCount()
		}
	}

}
