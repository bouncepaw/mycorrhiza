package web

import (
	"embed"
	"text/template"

	"github.com/bouncepaw/mycorrhiza/web/newtmpl"
	"github.com/bouncepaw/mycorrhiza/web/viewutil"
)

//go:embed views/*.html
var fs embed.FS

var pageOrphans, pageBacklinks, pageUserList, pageChangePassword *newtmpl.Page
var pageHyphaDelete, pageHyphaEdit, pageHyphaEmpty, pageHypha *newtmpl.Page
var pageRevision, pageMedia *newtmpl.Page
var pageAuthLock, pageAuthLogin, pageAuthLogout, pageAuthRegister *newtmpl.Page
var pageCatPage, pageCatList, pageCatEdit *newtmpl.Page

var panelTemplate, listTemplate, newUserTemplate, editUserTemplate, deleteUserTemplate *template.Template

func initPages() {

	panelTemplate = viewutil.CopyWith(fs, "views/admin-panel.html")
	listTemplate = viewutil.CopyWith(fs, "views/admin-user-list.html")
	newUserTemplate = viewutil.CopyWith(fs, "views/admin-new-user.html")
	editUserTemplate = viewutil.CopyWith(fs, "views/admin-edit-user.html")
	deleteUserTemplate = viewutil.CopyWith(fs, "views/admin-delete-user.html")

	pageOrphans = newtmpl.NewPage(fs, "views/orphans.html")
	pageBacklinks = newtmpl.NewPage(fs, "views/backlinks.html")
	pageUserList = newtmpl.NewPage(fs, "views/user-list.html")
	pageChangePassword = newtmpl.NewPage(fs, "views/change-password.html")
	pageHyphaDelete = newtmpl.NewPage(fs, "views/hypha-delete.html")
	pageHyphaEdit = newtmpl.NewPage(fs, "views/hypha-edit.html")
	pageHypha = newtmpl.NewPage(fs, "views/hypha.html")
	pageRevision = newtmpl.NewPage(fs, "views/hypha-revision.html")
	pageMedia = newtmpl.NewPage(fs, "views/hypha-media.html")

	pageAuthLock = newtmpl.NewPage(fs, "views/auth-telegram.html", "views/auth-lock.html")

	pageAuthLogin = newtmpl.NewPage(fs, "views/auth-telegram.html", "views/auth-login.html")

	pageAuthLogout = newtmpl.NewPage(fs, "views/auth-logout.html")

	pageAuthRegister = newtmpl.NewPage(fs, "views/auth-telegram.html", "views/auth-register.html")

	pageCatPage = newtmpl.NewPage(fs, "views/cat-page.html")

	pageCatEdit = newtmpl.NewPage(fs, "views/cat-edit.html")

	pageCatList = newtmpl.NewPage(fs, "views/cat-list.html")
}
