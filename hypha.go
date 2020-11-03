package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/markup"
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
			if hyphaData.binaryPath != "" {
				binaryBlock = binaryHtmlBlock(hyphaName, hyphaData)
			}
		} else {
			err = errors.New("Hypha " + hyphaName + " does not exist")
		}
		return
	}
}

// HyphaData represents a hypha's meta information: binary and text parts rooted paths and content types.
type HyphaData struct {
	textPath   string
	binaryPath string
}

// uploadHelp is a helper function for UploadText and UploadBinary
func (hd *HyphaData) uploadHelp(hop *history.HistoryOp, hyphaName, ext string, originalFullPath *string, isOld bool, data []byte) *history.HistoryOp {
	var (
		fullPath = filepath.Join(WikiDir, hyphaName+ext)
	)
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
	if !isOld {
		HyphaStorage[hyphaName] = hd
	}
	*originalFullPath = fullPath
	log.Printf("%v\n", *hd)
	return hop.WithFiles(fullPath).
		WithSignature("anon").
		Apply()
}

// UploadText loads a new text part from `textData` for hypha `hyphaName` with `hd`.  It must be marked if the hypha `isOld`.
func (hd *HyphaData) UploadText(hyphaName, textData string, isOld bool) *history.HistoryOp {
	hop := history.Operation(history.TypeEditText).WithMsg(fmt.Sprintf("Edit ‘%s’", hyphaName))
	return hd.uploadHelp(hop, hyphaName, ".myco", &hd.textPath, isOld, []byte(textData))
}

// UploadBinary loads a new binary part from `file` for hypha `hyphaName` with `hd`. The contents have the specified `mime` type. It must be marked if the hypha `isOld`.
func (hd *HyphaData) UploadBinary(hyphaName, mime string, file multipart.File, isOld bool) *history.HistoryOp {
	var (
		hop       = history.Operation(history.TypeEditBinary).WithMsg(fmt.Sprintf("Upload binary part for ‘%s’ with type ‘%s’", hyphaName, mime))
		data, err = ioutil.ReadAll(file)
	)
	if err != nil {
		return hop.WithError(err).Apply()
	}

	return hd.uploadHelp(hop, hyphaName, MimeToExtension(mime), &hd.binaryPath, isOld, data)
}

// DeleteHypha deletes hypha and makes a history record about that.
func (hd *HyphaData) DeleteHypha(hyphaName string) *history.HistoryOp {
	hop := history.Operation(history.TypeDeleteHypha).
		WithFilesRemoved(hd.textPath, hd.binaryPath).
		WithMsg(fmt.Sprintf("Delete ‘%s’", hyphaName)).
		WithSignature("anon").
		Apply()
	if len(hop.Errs) == 0 {
		delete(HyphaStorage, hyphaName)
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
			if hd.textPath != "" {
				renameMap[hd.textPath] = replaceName(hd.textPath)
			}
			if hd.binaryPath != "" {
				renameMap[hd.binaryPath] = replaceName(hd.binaryPath)
			}
		}
	}
	return renameMap, nil
}

// word Data is plural here
func relocateHyphaData(hyphaNames []string, replaceName func(string) string) {
	for _, hyphaName := range hyphaNames {
		if hd, ok := HyphaStorage[hyphaName]; ok {
			hd.textPath = replaceName(hd.textPath)
			hd.binaryPath = replaceName(hd.binaryPath)
			HyphaStorage[replaceName(hyphaName)] = hd
			delete(HyphaStorage, hyphaName)
		}
	}
}

// RenameHypha renames hypha from old name `hyphaName` to `newName` and makes a history record about that. If `recursive` is `true`, its subhyphae will be renamed the same way.
func (hd *HyphaData) RenameHypha(hyphaName, newName string, recursive bool) *history.HistoryOp {
	var (
		replaceName = func(str string) string {
			return strings.Replace(str, hyphaName, newName, 1)
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
		WithSignature("anon").
		Apply()
	if len(hop.Errs) == 0 {
		relocateHyphaData(hyphaNames, replaceName)
	}
	return hop
}

// binaryHtmlBlock creates an html block for binary part of the hypha.
func binaryHtmlBlock(hyphaName string, hd *HyphaData) string {
	switch filepath.Ext(hd.binaryPath) {
	case ".jpg", ".gif", ".png", ".webp", ".svg", ".ico":
		return fmt.Sprintf(`
		<div class="binary-container binary-container_with-img">
			<a href="/page/%[1]s"><img src="/binary/%[1]s"/></a>
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
			<p>This hypha's media cannot be rendered. Access it <a href="/binary/%s">directly</a></p>
		</div>
		`, hyphaName)
	}
}

// Index finds all hypha files in the full `path` and saves them to HyphaStorage. This function is recursive.
func Index(path string) {
	nodes, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, node := range nodes {
		// If this hypha looks like it can be a hypha path, go deeper. Do not touch the .git and static folders for they have an admnistrative importance!
		if node.IsDir() && isCanonicalName(node.Name()) && node.Name() != ".git" && node.Name() != "static" {
			Index(filepath.Join(path, node.Name()))
			continue
		}

		var (
			hyphaPartPath           = filepath.Join(path, node.Name())
			hyphaName, isText, skip = DataFromFilename(hyphaPartPath)
			hyphaData               *HyphaData
		)
		if !skip {
			// Reuse the entry for existing hyphae, create a new one for those that do not exist yet.
			if hd, ok := HyphaStorage[hyphaName]; ok {
				hyphaData = hd
			} else {
				hyphaData = &HyphaData{}
				HyphaStorage[hyphaName] = hyphaData
			}
			if isText {
				hyphaData.textPath = hyphaPartPath
			} else {
				// Notify the user about binary part collisions. It's a design decision to just use any of them, it's the user's fault that they have screwed up the folder structure, but the engine should at least let them know, right?
				if hyphaData.binaryPath != "" {
					log.Println("There is a file collision for binary part of a hypha:", hyphaData.binaryPath, "and", hyphaPartPath, "-- going on with the latter")
				}
				hyphaData.binaryPath = hyphaPartPath
			}
		}

	}
}

// FetchTextPart tries to read text file in the `d`. If there is no file, empty string is returned.
func FetchTextPart(d *HyphaData) (string, error) {
	if d.textPath == "" {
		return "", nil
	}
	_, err := os.Stat(d.textPath)
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	text, err := ioutil.ReadFile(d.textPath)
	if err != nil {
		return "", err
	}
	return string(text), nil
}
