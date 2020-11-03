package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bouncepaw/mycorrhiza/gemtext"
	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/util"
)

func init() {
	gemtext.HyphaExists = func(hyphaName string) bool {
		_, hyphaExists := HyphaStorage[hyphaName]
		return hyphaExists
	}
	gemtext.HyphaAccess = func(hyphaName string) (rawText, binaryBlock string, err error) {
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
	textType   TextType
	binaryPath string
	binaryType BinaryType
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

func renamingPairs(hyphaNames []string, replaceName func(string) string) map[string]string {
	renameMap := make(map[string]string)
	for _, hn := range hyphaNames {
		if hd, ok := HyphaStorage[hn]; ok {
			if hd.textPath != "" {
				renameMap[hd.textPath] = replaceName(hd.textPath)
			}
			if hd.binaryPath != "" {
				renameMap[hd.binaryPath] = replaceName(hd.binaryPath)
			}
		}
	}
	return renameMap
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
		hyphaNames = findHyphaeToRename(hyphaName, recursive)
		renameMap  = renamingPairs(hyphaNames, replaceName)
		renameMsg  = "Rename ‘%s’ to ‘%s’"
		hop        = history.Operation(history.TypeRenameHypha)
	)
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
func binaryHtmlBlock(hyphaName string, d *HyphaData) string {
	switch d.binaryType {
	case BinaryJpeg, BinaryGif, BinaryPng, BinaryWebp, BinarySvg, BinaryIco:
		return fmt.Sprintf(`
		<div class="binary-container binary-container_with-img">
			<a href="/page/%s"><img src="/binary/%s"/></a>
		</div>`, hyphaName)
	case BinaryOgg, BinaryWebm, BinaryMp4:
		return fmt.Sprintf(`
		<div class="binary-container binary-container_with-video">
			<video>
				<source src="/binary/%[1]s"/>
				<p>Your browser does not support video. See video's <a href="/binary/%[1]s">direct url</a></p>
			</video>
		`, hyphaName)
	case BinaryMp3:
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
		// If this hypha looks like it can be a hypha path, go deeper
		if node.IsDir() && isCanonicalName(node.Name()) {
			Index(filepath.Join(path, node.Name()))
		}

		hyphaPartFilename := filepath.Join(path, node.Name())
		skip, hyphaName, isText, mimeId := DataFromFilename(hyphaPartFilename)
		if !skip {
			var (
				hyphaData *HyphaData
				ok        bool
			)
			if hyphaData, ok = HyphaStorage[hyphaName]; !ok {
				hyphaData = &HyphaData{}
				HyphaStorage[hyphaName] = hyphaData
			}
			if isText {
				hyphaData.textPath = hyphaPartFilename
				hyphaData.textType = TextType(mimeId)
			} else {
				hyphaData.binaryPath = hyphaPartFilename
				hyphaData.binaryType = BinaryType(mimeId)
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
