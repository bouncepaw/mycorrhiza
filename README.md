# 🍄 MycorrhizaWiki 0.10
A wiki engine.

This is the development branch for version 0.10. Features planned for this version:
* [x] New file structure
* [ ] Mycomarkup
  * [x] 6 headings
  * [x] Paragraph styling
  * [ ] Inline links
  * [ ] Better lists
  * [ ] Better quotes
* [x] CLI options
* [x] CSS improvements

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
  -home string
        The home page (default "home")
  -port string
        Port to serve the wiki at (default "1737")
  -title string
        How to call your wiki in the navititle (default "🍄")
```

## Features
* Edit pages through html forms
* Responsive design
* Works in text browsers
* Wiki pages (called hyphae) are in gemtext
* Everything is stored as simple files, no database required
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
* Authorization
* Better history viewing
