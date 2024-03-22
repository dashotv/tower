package importer

type Series struct {
	Title        string
	Description  string
	Airdate      string
	Status       string
	Language     string
	EpisodeOrder int
	Episodes     []*Episode
	Covers       []string
	Backgrounds  []string
}

type Episode struct {
	Title       string
	Description string
	Airdate     string
	Season      int
	Episode     int
	Absolute    int
}

const (
	EpisodeOrderUnknown = iota
	EpisodeOrderDefault
	EpisodeOrderDVD
	EpisodeOrderAbsolute
)

func episodeOrderString(order int) string {
	switch order {
	case EpisodeOrderDVD:
		return "dvd"
	case EpisodeOrderAbsolute:
		return "absolute"
	default:
		return "default"
	}
}
