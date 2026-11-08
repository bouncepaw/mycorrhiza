package web

import (
	"github.com/bouncepaw/mycorrhiza/web/newtmpl"
)

var (
	pageOrphans, pageBacklinks, pageUserList, pageChangePassword                                *newtmpl.Page
	pageHyphaDelete, pageHyphaEdit, pageHyphaEmpty, pageHypha                                   *newtmpl.Page
	pageRevision, pageMedia                                                                     *newtmpl.Page
	pageAuthLock, pageAuthLogin, pageAuthLogout, pageAuthRegister                               *newtmpl.Page
	pageCatPage, pageCatList, pageCatEdit                                                       *newtmpl.Page
	pageAdminPanel, pageAdminUserList, pageAdminNewUser, pageAdminEditUser, pageAdminDeleteUser *newtmpl.Page
)

func initPages() {
	pageOrphans = newtmpl.NewPage("orphans.html")
	pageBacklinks = newtmpl.NewPage("backlinks.html")
	pageUserList = newtmpl.NewPage("user-list.html")
	pageChangePassword = newtmpl.NewPage("change-password.html")
	pageHyphaDelete = newtmpl.NewPage("hypha-delete.html")
	pageHyphaEdit = newtmpl.NewPage("hypha-edit.html")
	pageHypha = newtmpl.NewPage("hypha.html")
	pageRevision = newtmpl.NewPage("hypha-revision.html")
	pageMedia = newtmpl.NewPage("hypha-media.html")

	pageAuthLock = newtmpl.NewPage("auth-telegram.html", "auth-lock.html")
	pageAuthLogin = newtmpl.NewPage("auth-telegram.html", "auth-login.html")
	pageAuthLogout = newtmpl.NewPage("auth-logout.html")
	pageAuthRegister = newtmpl.NewPage("auth-telegram.html", "auth-register.html")

	pageCatPage = newtmpl.NewPage("cat-page.html")
	pageCatEdit = newtmpl.NewPage("cat-edit.html")
	pageCatList = newtmpl.NewPage("cat-list.html")

	pageAdminPanel = newtmpl.NewPage("admin-panel.html")
	pageAdminUserList = newtmpl.NewPage("admin-user-list.html")
	pageAdminNewUser = newtmpl.NewPage("admin-new-user.html")
	pageAdminEditUser = newtmpl.NewPage("admin-edit-user.html")
	pageAdminDeleteUser = newtmpl.NewPage("admin-delete-user.html")

}
