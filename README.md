# 🍄 MycorrhizaWiki 0.12
A wiki engine.

[Main wiki](https://mycorrhiza.lesarbr.es)

## Building
Also see [detailed instructions](https://mycorrhiza.lesarbr.es/page/deploy) on wiki.
```sh
git clone --recurse-submodules https://github.com/bouncepaw/mycorrhiza
cd mycorrhiza
make
# That make will:
# * run the default wiki. You can edit it right away.
# * create an executable called `mycorrhiza`. Run it with path to your wiki.
```

## Usage
```
mycorrhiza [OPTIONS...] WIKI_PATH

WIKI_PATH must be a path to git repository which you want to be a wiki.

Options:
  -auth-method string
        What auth method to use. Variants: "none", "fixed" (default "none")
  -fixed-credentials-path string
        Used when -auth-method=fixed. Path to file with user credentials. (default "mycocredentials.json")
  -header-links-hypha string
        Optional hypha that overrides the header links
  -home string
        The home page (default "home")
  -icon string
        What to show in the navititle in the beginning, before the colon (default "🍄")
  -name string
        What is the name of your wiki (default "wiki")
  -port string
        Port to serve the wiki at (default "1737")
  -url string
        URL at which your wiki can be found. Used to generate feeds (default "http://0.0.0.0:$port")
  -user-hypha string
        Hypha which is a superhypha of all user pages (default "u")
```

## Features
* Edit pages through html forms, graphical preview
* Responsive design, dark theme (synced with system theme)
* Works in text browsers
* Wiki pages (called hyphae) are written in mycomarkup
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
* Authorization with pre-set credentials

## Contributing
Help is always needed. We have a [tg chat](https://t.me/mycorrhizadev) where some development is coordinated. You can also sponsor on [boosty](https://boosty.to/bouncepaw). Feel free to open an issue or contact directly.

You can view list of all planned features on [our kanban board](https://github.com/bouncepaw/mycorrhiza/projects/1).
