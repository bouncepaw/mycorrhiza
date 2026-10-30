package userdomain

import "github.com/bouncepaw/mycorrhiza/util/fp"

type Action string

const (
	ActionText              Action = "text"
	ActionBacklinks         Action = "backlinks"
	ActionHistory           Action = "history"
	ActionMedia             Action = "media"
	ActionEdit              Action = "edit"
	ActionUploadMedia       Action = "upload-binary"
	ActionRename            Action = "rename"
	ActionUploadText        Action = "upload-text"
	ActionAddToCat          Action = "add-to-category"
	ActionRemoveFromCat     Action = "remove-from-category"
	ActionRemoveMedia       Action = "remove-media"
	ActionUpdateHeaderLinks Action = "update-header-links"
	ActionDelete            Action = "delete"
	ActionReindex           Action = "reindex"
	ActionAdmin             Action = "admin"
	ActionAdminShutdown     Action = "admin/shutdown"
)

func (a Action) String() string {
	return string(a)
}

func (a Action) MinimalLevel() GroupLevel {
	lvl, ok := minimalRights[a]
	return fp.If(ok, lvl, 5)
}

var minimalRights = map[Action]GroupLevel{
	// TODO: use consts above
	"text":                 0,
	"backlinks":            0,
	"history":              0,
	"media":                1,
	"edit":                 1,
	"upload-binary":        1,
	"rename":               1,
	"upload-text":          1,
	"add-to-category":      1,
	"remove-from-category": 1,
	"remove-media":         2,
	"update-header-links":  3,
	"delete":               3,
	"reindex":              4,
	"admin":                4,
	"admin/shutdown":       4,
}
