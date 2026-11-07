package web

import (
	"embed"

	"github.com/bouncepaw/mycorrhiza/web/newtmpl"
)

//go:embed views/*.html
var fs embed.FS

var (
	pageOrphans, pageBacklinks, pageUserList, pageChangePassword                                *newtmpl.Page
	pageHyphaDelete, pageHyphaEdit, pageHyphaEmpty, pageHypha                                   *newtmpl.Page
	pageRevision, pageMedia                                                                     *newtmpl.Page
	pageAuthLock, pageAuthLogin, pageAuthLogout, pageAuthRegister                               *newtmpl.Page
	pageCatPage, pageCatList, pageCatEdit                                                       *newtmpl.Page
	pageAdminPanel, pageAdminUserList, pageAdminNewUser, pageAdminEditUser, pageAdminDeleteUser *newtmpl.Page
)

func initPages() {
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

	pageAdminPanel = newtmpl.NewPage(fs, "views/admin-panel.html")
	pageAdminUserList = newtmpl.NewPage(fs, "views/admin-user-list.html")
	pageAdminNewUser = newtmpl.NewPage(fs, "views/admin-new-user.html")
	pageAdminEditUser = newtmpl.NewPage(fs, "views/admin-edit-user.html")
	pageAdminDeleteUser = newtmpl.NewPage(fs, "views/admin-delete-user.html")

}
