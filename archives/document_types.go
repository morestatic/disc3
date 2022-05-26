package archives

type DocumentType string

const (
	Unknown  = iota
	Artists  = "artists"
	Labels   = "labels"
	Releases = "releases"
	Masters  = "masters"
)

func (dt DocumentType) String() string {
	switch dt {
	case Artists:
		return "artists"
	case Labels:
		return "labels"
	case Releases:
		return "releases"
	case Masters:
		return "masters"
	default:
		return "unknown"
	}
}

func (dt DocumentType) Singular() string {
	switch dt {
	case Artists:
		return "artist"
	case Labels:
		return "label"
	case Releases:
		return "release"
	case Masters:
		return "master"
	default:
		return "unknown"
	}
}

func (dt DocumentType) ShortForm() string {
	switch dt {
	case Artists:
		return "A"
	case Labels:
		return "L"
	case Releases:
		return "R"
	case Masters:
		return "M"
	default:
		return "U"
	}
}
