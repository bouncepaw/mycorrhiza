# The Structure
Here I am, doing new stuff before finishing the old stuff.

See https://github.com/bouncepaw/mycorrhiza/issues/57 for the discussion.

## The idea
Instead of letting users figure everything out by themselves, we think of the best (in our opinion) file layout and force it onto the users and remove the possibility of configuring it.

## What is inside the Structure
### Root
The whole wiki is inside one directory or inside one of its subdirectories. Only the Mycorrhiza binary itself and Git (and possible future runtime dependencies) might be outside that directory. That directory (called _root directory_) can have any name.

### Subdirectories
* `wiki.git` is a valid Git repository. If it is not present or is not a valid Git repository, the engine shall fail to work. When the Wizard is implemented, the engine will offer to make the Git repository.
* `cache` contains temporary files such as user token caches. Wiki administrators can safely delete this directory and expect the wiki to continue working. In the future, stuff like pre-rendered HTML can be stored here.
* All other subdirectories are ignored.

### User configuration
* `registered-users.json` contains a JSON array of all registered users. The engine will edit this file, and the administrators should not edit by themselves, unless they really want to.
* `fixed-users.json` contains a JSON array of all fixed users. Wiki administrators will edit this file by themselves.

### Wiki configuration
* `config.ini` is the main configuration file.

### Customisation
* `favicon.ico` is the Favicon as you know it.
* `common.css` redefines the built-in CSS, the Common style.
* `custom.css` is sent to the user after the Common style.

### Meta
* `README.txt` contains a short description of the files that can be inside the Structure. A small reminder for the administrators.