package userdomain

type Source string

const (
	SourceLocal Source = "local"
)

func (s Source) IsValid() bool {
	switch s {
	case SourceLocal:
		return true
	}
	return false
}

func (s Source) String() string {
	return string(s)
}
