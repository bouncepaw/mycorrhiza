# 🍄 Mycorrhiza Wiki 1.2
<img src="https://mycorrhiza.lesarbr.es/binary/release/1.2/screenshot" alt="A screenshot of Mycorrhiza Wiki home hypha in the Safari browser" width="500">

**Mycorrhiza Wiki** is a filesystem and git-based wiki engine.

[Main wiki](https://mycorrhiza.lesarbr.es)

## Usage
```
mycorrhiza [OPTIONS...] WIKI_PATH

WIKI_PATH must be a path to a git repository which you want to be a wiki.

Options:
  -config-path string
        Path to a configuration file. Leave empty if you don't want to use it.
  -print-example-config
        If true, print an example configuration file contents and exit. You can save the output to a file and base your own configuration on it.
```

## Features
* Wiki pages (called hyphae) are written in Mycomarkup
* Edit pages through HTML forms, a graphical preview and a toolbar that helps you use Mycomarkup
* Responsive design, dark theme (synced with system theme)
* Works in text browsers
* Everything is stored as simple files, no database required. You can run a wiki on almost any directory and get something to work with
* Page trees; links to previous and next pages
* Changes are saved to git
* List of hyphae page
* History page
* Random page
* Recent changes page; RSS, Atom and JSON feeds available
* Hyphae can be deleted while still preserving history
* Hyphae can be renamed (recursive renaming of subhyphae is also supported)
* Light on resources
* Registration
* Hotkeys (press `?` to see what hotkeys there are)

## Building
See [the guide](https://mycorrhiza.lesarbr.es/hypha/guide/deployment) on the wiki.

## Installing

### AUR
You can install Mycorrhiza Wiki from AUR using your favorite package manager on any Arch Linux-derivative distro (Arch, Manjaro, Garuda, etc):
```sh
# Build from sources
yay -S mycorrhiza
# OR
# Use pre-built binaries from the Releases page
yay -S mycorrhiza-bin
```

### Docker
You can run Mycorrhiza Wiki in Docker using Dockerfile provided by this repository. Clone the repo and build the image:
```sh
git clone https://github.com/bouncepaw/mycorrhiza/
docker build -t mycorrhiza .
```

Now you can create a new Mycorrhiza Wiki container using this command:
```sh
docker run -v /path/to/wiki:/wiki -p 1737:1737 mycorrhiza
```

Example:
```sh
cd /dev/shm
git clone https://github.com/bouncepaw/mycorrhiza/
docker build -t mycorrhiza .
git clone https://github.com/bouncepaw/example-wiki
docker run -v /dev/shm/example-wiki:/wiki -p 1737:1737 mycorrhiza
```

Example 2:
```sh
# ...
docker run -v /dev/shm/:/config -v /dev/shm/example-wiki:/wiki -p 80:1737 mycorrhiza -config-path /config/myconfig.ini /wiki
```

## Contributing
We always need help. We have a [Telegram chat](https://t.me/mycorrhizadev) where we coordinate development. You can also sponsor @bouncepaw on [Boosty](https://boosty.to/bouncepaw). Feel free to open an issue or contact us directly.

You can view list of many planned features on [our kanban board](https://github.com/bouncepaw/mycorrhiza/projects/1).
