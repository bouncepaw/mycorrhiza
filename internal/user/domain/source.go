package userdomain

type Source string

const (
	SourceLocal    Source = "local"
	SourceTelegram Source = "telegram"
)

func (s Source) IsValid() bool {
	switch s {
	case SourceLocal, SourceTelegram:
		return true
	}
	return false
}

func (s Source) String() string {
	return string(s)
}
