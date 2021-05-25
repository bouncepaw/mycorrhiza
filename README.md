# 🍄 MycorrhizaWiki 1.2
A wiki engine.

[Main wiki](https://mycorrhiza.lesarbr.es)

## Building
See [the guide](https://mycorrhiza.lesarbr.es/hypha/guide/deployment) on the wiki.

## Usage
```
mycorrhiza [OPTIONS...] WIKI_PATH

WIKI_PATH must be a path to git repository which you want to be a wiki.

Options:
  -config-path string
        Path to a configuration file. Leave empty if you don't want to use it.
  -print-example-config
        If true, print an example configuration file contents and exit. You can save the output to a file and base your own configuration on it.
```

## Features
* Wiki pages (called hyphae) are written in mycomarkup
* Edit pages through html forms, graphical preview, toolbar that helps you use the markup
* Responsive design, dark theme (synced with system theme)
* Works in text browsers
* Everything is stored as simple files, no database required. You can run a wiki on almost any directory and get something to work with
* Page trees; links to previous and next pages
* Changes are saved to git
* List of hyphae page
* History page
* Random page
* Recent changes page; RSS, Atom and JSON feeds available
* Hyphae can be deleted (while still preserving history)
* Hyphae can be renamed (recursive renaming of subhyphae is also supported)
* Light on resources
* Authorization with pre-set credentials, registration
* Basic Gemini protocol support

## Contributing
Help is always needed. We have a [tg chat](https://t.me/mycorrhizadev) where some development is coordinated. You can also sponsor on [boosty](https://boosty.to/bouncepaw). Feel free to open an issue or contact directly.

You can view list of all planned features on [our kanban board](https://github.com/bouncepaw/mycorrhiza/projects/1).
