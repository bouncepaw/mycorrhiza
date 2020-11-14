# 🍄 MycorrhizaWiki 0.11
A wiki engine.

Features planned for this release:
* [ ] Authorization
  * [x] User groups: `anon`, `editor`, `trusted`, `moderator`, `admin`
  * [ ] Login page
  * [ ] Rights
* [ ] Mycomarkup improvements
  * [x] Strike-through syntax
  * [x] Formatting in headings
  * [ ] Fix empty line codeblock bug #26
  * [ ] `img{}` improvements
  * [ ] ...

## Building
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

Options:
  -auth-method string
        What auth method to use. Variants: "none", "fixed" (default "none")
  -fixed-credentials-path string
        Used when -auth-method=fixed. Path to file with user credentials. (default "mycocredentials.json")
  -home string
        The home page (default "home")
  -port string
        Port to serve the wiki at (default "1737")
  -title string
        How to call your wiki in the navititle (default "🍄")
  -user-tree string
        Hypha which is a superhypha of all user pages (default "u")
```

## Features
* Edit pages through html forms
* Responsive design
* Works in text browsers
* Wiki pages (called hyphae) are written in mycomarkup
* Everything is stored as simple files, no database required. You can run a wiki on almost any directory and get something to work with.
* Page trees
* Changes are saved to git
* List of hyphae page
* History page
* Random page
* Recent changes page
* Hyphae can be deleted (while still preserving history)
* Hyphae can be renamed (recursive renaming of subhyphae is also supported)
* Light on resources: I run a home wiki on this engine 24/7 at an [Orange π Lite](http://www.orangepi.org/orangepilite/).

## Contributing
Help is always needed. We have a [tg chat](https://t.me/mycorrhizadev) where some development is coordinated. Feel free to open an issue or contact me.

## Future plans
* Tagging system
* Better history viewing
