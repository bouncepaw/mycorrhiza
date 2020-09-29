package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/bouncepaw/mycorrhiza/gemtext"
	"github.com/bouncepaw/mycorrhiza/history"
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
	return history.Operation(history.TypeDeleteHypha).
		WithFilesRemoved(hd.textPath, hd.binaryPath).
		WithMsg(fmt.Sprintf("Delete ‘%s’", hyphaName)).
		WithSignature("anon").
		Apply()
}

// binaryHtmlBlock creates an html block for binary part of the hypha.
func binaryHtmlBlock(hyphaName string, d *HyphaData) string {
	switch d.binaryType {
	case BinaryJpeg, BinaryGif, BinaryPng, BinaryWebp, BinarySvg, BinaryIco:
		return fmt.Sprintf(`
		<div class="binary-container binary-container_with-img">
			<img src="/binary/%s"/>
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
