package shroom

import (
	"log"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/user"
)

func rejectDeleteLog(h hyphae.Hypha, u *user.User, errmsg string) {
	log.Printf("Reject delete ‘%s’ by @%s: %s\n", h.CanonicalName(), u.Name, errmsg)
}
func rejectRenameLog(h hyphae.Hypha, u *user.User, errmsg string) {
	log.Printf("Reject rename ‘%s’ by @%s: %s\n", h.CanonicalName(), u.Name, errmsg)
}
func rejectRemoveMediaLog(h hyphae.Hypha, u *user.User, errmsg string) {
	log.Printf("Reject remove media ‘%s’ by @%s: %s\n", h.CanonicalName(), u.Name, errmsg)
}
func rejectEditLog(h hyphae.Hypha, u *user.User, errmsg string) {
	log.Printf("Reject edit ‘%s’ by @%s: %s\n", h.CanonicalName(), u.Name, errmsg)
}
func rejectUploadMediaLog(h hyphae.Hypha, u *user.User, errmsg string) {
	log.Printf("Reject upload media ‘%s’ by @%s: %s\n", h.CanonicalName(), u.Name, errmsg)
}
