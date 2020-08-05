# mycorrhiza wiki
A wiki engine inspired by fungi. Not production-ready.

This branch is devoted to version 0.8.
* [ ] Tree generation
  * [x] Basic generation
  * [ ] Generation that takes non-existent hyphae into account¹
* [ ] History
  * [ ] Saving all changes to git
  * [ ] Make it possible to see any previous version
  * [ ] A nice UI for that

¹ Current algorithm won't detect `a/b/c` as a child of `a` if `a/b` does not exist.

## Current features
* Edit pages through html forms
* Responsive design
* Works in text browsers
* Pages (called hyphae) can be in gemtext.
* Everything is stored as simple files, no database required

## Future features
* Tags
* Authorization
* History view
* Granular user rights

## Installation
I guess you can just clone this repo and run `make` to play around with the default wiki.

