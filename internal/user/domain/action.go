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
	ActionText:              0,
	ActionBacklinks:         0,
	ActionHistory:           0,
	ActionMedia:             1,
	ActionEdit:              1,
	ActionUploadMedia:       1,
	ActionRename:            1,
	ActionUploadText:        1,
	ActionAddToCat:          1,
	ActionRemoveFromCat:     1,
	ActionRemoveMedia:       2,
	ActionUpdateHeaderLinks: 3,
	ActionDelete:            3,
	ActionReindex:           4,
	ActionAdmin:             4,
	ActionAdminShutdown:     4,
}
