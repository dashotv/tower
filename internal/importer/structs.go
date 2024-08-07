package importer

type Series struct {
	ID          int64
	Title       string
	Description string
	Airdate     string
	Status      string
	Language    string
}

type Episode struct {
	ID          int64
	Title       string
	Description string
	Airdate     string
	Season      int
	Episode     int
	Absolute    int
}

type Movie struct {
	ID          int64
	ImdbID      string
	Title       string
	Description string
	Airdate     string
	Poster      string
	Backdrop    string
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
