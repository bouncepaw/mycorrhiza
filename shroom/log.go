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
func rejectUnattachLog(h hyphae.Hypha, u *user.User, errmsg string) {
	log.Printf("Reject unattach ‘%s’ by @%s: %s\n", h.CanonicalName(), u.Name, errmsg)
}
func rejectEditLog(h hyphae.Hypha, u *user.User, errmsg string) {
	log.Printf("Reject edit ‘%s’ by @%s: %s\n", h.CanonicalName(), u.Name, errmsg)
}
func rejectAttachLog(h hyphae.Hypha, u *user.User, errmsg string) {
	log.Printf("Reject attach ‘%s’ by @%s: %s\n", h.CanonicalName(), u.Name, errmsg)
}
