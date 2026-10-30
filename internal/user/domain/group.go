package userdomain

type (
	GroupKind  string
	GroupLevel int
)

const (
	GroupAnon                GroupKind = "anon"
	GroupReader              GroupKind = "reader"
	GroupEditor              GroupKind = "editor"
	GroupLegacyTrustedEditor GroupKind = "trusted"
	GroupModerator           GroupKind = "moderator"
	GroupAdmin               GroupKind = "admin"
)

func (k GroupKind) IsValid() bool {
	switch k {
	case GroupAnon,
		GroupReader,
		GroupEditor,
		GroupLegacyTrustedEditor,
		GroupModerator,
		GroupAdmin:
		return true
	}
	return false
}

func (k GroupKind) String() string {
	return string(k)
}

func (k GroupKind) Level() GroupLevel {
	switch k {
	case GroupAnon, GroupReader:
		return 0
	case GroupEditor:
		return 1
	case GroupLegacyTrustedEditor:
		return 2
	case GroupModerator:
		return 3
	case GroupAdmin:
		return 4
	}
	return 0
}
