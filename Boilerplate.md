# Boilerplate in Mycorrhiza codebase

Being programmed by Go, mostly by Bouncepaw, the codebase contains a lot of boilerplate. This document is an attempt to describe how it is done.

## Modules

Mycorrhiza is arranged in quite many packages. They are thematic. For example, package `backlinks` has all things backlinks, including the storage, the views, and exported functions (not many, ideally). Such packages can be called modules, if you want.

## Views

Views are the biggest source of similar code. Before the transition from QTPL to Go's standard templates, this boilerplate energy was split differently, but was not instantly obvious. The current approach does not really introduce new boilerplate energy, but it does focus it, resulting in actual boilerplate. I hope you get the idea.

All related views are part of one module.

Views come in multiple parts.

The first part is the template itself. Call template files like that: `view_user_list.html`, prefixed with `view_`. The boilerplate is as follows.

```html
{{define "title"}}{{end}}
{{define "body"}}
{{end}}
```

More often than not, you will want to make template `title` a different template in the same file. See existing files for inspiration.

The code that makes those templates runnable lies in one file. This is the second part. It contains the following:.

The Russian translation is a `string` variable called `ruTranslation`; we currently have no other translations, but they are to be called like `frTranslation`, `eoTranslation`, et cetera.


```go
var (
	ruTranslation = `
{{define "one thing"}}...{{end}}
{{define "other thing"}}...{{end}}
`
	...
)
```

Chains are collection of different language variants of the same template. Declare them and then assign them in a function, which you call somewhere (not just `init`!).

```go
var (
	ruTranslation              = `...`
	chainStuff, chainAddStuff viewutil.Chain
)

func initViews() {
	chainStuff = viewutil.CopyEnRuWith(fs, "view_stuff.html", ruTranslation)
	chainAddStuff = viewutil.CopyEnRuWith(fs, "view_add_stuff.html", ruTranslation)
}
```

Then every view has a runner and its own datatype.

```go
//...

type dataStuff struct {
	*viewutil.BaseData
	StuffName string
}

func viewStuff(meta viewutil.Meta, stuffName string) {
	viewutil.ExecutePage(meta, chainStuff, dataUserList{
		BaseData: &viewutil.BaseData{},
		StuffName: stuffName,
	})
}
```

Sometimes, two datatypes of different views are the same, it is ok to just share one, but name it so that it mentions both.

Avoid any logic in those runners. Keep them as boilerplate as they are. You rarely need to fill the `BaseData` field. Do it if you need to. If you don't, `viewutil.ExecutePage` will do its best to guess. We name the fields with capital letters.

### Troubleshooting

* Declared the chains?
* Assigned the chains?
* Assigned the chains, sure?
* Used the correct data type?
* `*viewutil.BaseData` is there?
* Used the correct chain?