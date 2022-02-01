# 🍄 Mycorrhiza Wiki

**Mycorrhiza Wiki** is a lightweight file-system wiki engine that uses Git for keeping history. [Main wiki](https://mycorrhiza.wiki)

<img src="https://mycorrhiza.wiki/binary/release/1.8/screenshot" alt="A screenshot of mycorrhiza.wiki's home page in the Safari browser" width="600">

## Features

* **No database required.** Everything is stored in plain files. It makes installation super easy, and you can modify the content directly by yourself.
* **Everything is hyphae.** A [hypha] is a unit of content such as a picture, video or a text article. Hyphae can [transclude][transclusion] and link each other resulting in a tight network of hypertext pages.
* **Hyphae are authored in [Mycomarkup],** a markup language that's designed to be unambigious yet easy to use.
* **Nesting of hyphae** is supported. A tree of related hyphae is shown on every page.
* **History of changes** for textual parts of hyphae. Every change is safely stored in [Git]. Web feeds for recent changes included!
* **Keyboard-driven navigation.** Press `?` to see the list of shortcuts.
* **Support for [authorization].**
* **[Open Graph] support.**
* **Optional [Telegram] authentication.**

Compare Mycorrhiza Wiki with other engines on [WikiMatrix](https://www.wikimatrix.org/show/mycorrhiza).

## Installing

See [the deployment guide](https://mycorrhiza.wiki/hypha/guide/deployment) on the wiki.

## Contributing

Help is always welcome! We have an IRC channel [#mycorrhiza on irc.libera.chat]
and a [Telegram chat] for discussions and development. You can also sponsor the
maintainer of Mycorrhiza, [@bouncepaw], on [Boosty]. If you want to contribute
with code, you can either open a pull request on GitHub or send a patch to the
[mailing list]. Feel free to open an issue on GitHub or contact us directly.

You can view the list of planned features at [our GitHub project kanban
board][kanban] or on the [roadmap page].

[hypha]: https://mycorrhiza.wiki/hypha/feature/hypha
[transclusion]: https://mycorrhiza.wiki/hypha/feature/transclusion
[authorization]: https://mycorrhiza.wiki/hypha/feature/authorization
[Mycomarkup]: https://mycorrhiza.wiki/help/en/mycomarkup
[Git]: https://mycorrhiza.wiki/hypha/integration/git
[Open Graph]: https://mycorrhiza.wiki/hypha/standard/opengraph
[Telegram]: https://mycorrhiza.wiki/help/en/telegram

[#mycorrhiza on irc.libera.chat]: irc://irc.libera.chat/#mycorrhiza
[Telegram chat]: https://t.me/mycorrhizadev
[@bouncepaw]: https://github.com/bouncepaw
[Boosty]: https://boosty.to/bouncepaw
[mailing list]: https://lists.sr.ht/~handlerug/mycorrhiza-devel
[kanban]: https://github.com/bouncepaw/mycorrhiza/projects/1
[roadmap page]: https://mycorrhiza.wiki/hypha/release/roadmap
