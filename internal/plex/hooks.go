package plex

type HookData struct {
	Payload HookPayload `json:"payload"`
}

type HookPayload struct {
	Event    string              `json:"event"`
	User     bool                `json:"user"`
	Owner    bool                `json:"owner"`
	Account  Account             `json:"Account"`
	Server   Server              `json:"Server"`
	Player   Player              `json:"Player"`
	Metadata HookPayloadMetadata `json:"Metadata"`
}

type HookPayloadMetadata struct {
	LibrarySectionType   string `json:"librarySectionType"`
	RatingKey            string `json:"ratingKey"`
	Key                  string `json:"key"`
	ParentRatingKey      string `json:"parentRatingKey"`
	GrandparentRatingKey string `json:"grandparentRatingKey"`
	GUID                 string `json:"guid"`
	LibrarySectionID     int64  `json:"librarySectionID"`
	Type                 string `json:"type"`
	Title                string `json:"title"`
	GrandparentKey       string `json:"grandparentKey"`
	ParentKey            string `json:"parentKey"`
	GrandparentTitle     string `json:"grandparentTitle"`
	ParentTitle          string `json:"parentTitle"`
	Summary              string `json:"summary"`
	Index                int64  `json:"index"`
	ParentIndex          int64  `json:"parentIndex"`
	RatingCount          int64  `json:"ratingCount"`
	Thumb                string `json:"thumb"`
	Art                  string `json:"art"`
	ParentThumb          string `json:"parentThumb"`
	GrandparentThumb     string `json:"grandparentThumb"`
	GrandparentArt       string `json:"grandparentArt"`
	AddedAt              int64  `json:"addedAt"`
	UpdatedAt            int64  `json:"updatedAt"`
}
