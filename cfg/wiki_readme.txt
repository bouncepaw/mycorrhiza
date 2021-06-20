Hello! This is your new wiki directory. It contains various files that control
how wiki looks and works. Here's a brief breakdown of them.

wiki.git/
    This is a Git repository that holds all of your wiki's content. If you want
    to manually change something, add or remove files, do it there. Everything
    inside is publicly accessible, so don't place something private in there.

static/
    This directory can be used to serve additional static files, or to overwrite
    built-in ones. Here are some of the files you might want to place here:
        common.css:     replaces Mycorrhiza Wiki's default stylesheet
        custom.css:     adds additional styles to common.css
        favicon.ico:    sets the website favicon

cache/
    This directory contains temporary files such as user tokens. You can
    safely delete this directory, and nothing will break.

README.txt
    We're here! Feel free to remove this file if you don't need it.

config.ini
    This is the main configuration file that holds all your wiki's settings,
    such as it's name, special hyphae, network configuration and, optionally,
    additional JavaScripts.

fixed-users.json
    This file holds all fixed users' credentials. Not recommended to use.
    See https://mycorrhiza.wiki/hypha/feature/authorization/fixed.

registered-users.json
    This file holds all registered users' credentials, securely hashed and
    protected. This file is updated automatically, and you're discouraged from
    editing it by yourself. More user control features are coming in the future.
