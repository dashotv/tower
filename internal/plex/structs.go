package plex

type Country struct {
	Tag string `json:"tag"`
}

type Field struct {
	Locked bool   `json:"locked"`
	Name   string `json:"name"`
}

type Location struct {
	Path string `json:"path"`
}

type Media struct {
	ID                    int64   `json:"id"`
	Duration              int64   `json:"duration"`
	Bitrate               int64   `json:"bitrate"`
	Width                 int64   `json:"width"`
	Height                int64   `json:"height"`
	AspectRatio           float64 `json:"aspectRatio"`
	AudioChannels         int64   `json:"audioChannels"`
	AudioCodec            string  `json:"audioCodec"`
	VideoCodec            string  `json:"videoCodec"`
	VideoResolution       string  `json:"videoResolution"`
	Container             string  `json:"container"`
	VideoFrameRate        string  `json:"videoFrameRate"`
	AudioProfile          *string `json:"audioProfile,omitempty"`
	VideoProfile          string  `json:"videoProfile"`
	Part                  []Part  `json:"Part"`
	OptimizedForStreaming *int64  `json:"optimizedForStreaming,omitempty"`
	Has64BitOffsets       *bool   `json:"has64bitOffsets,omitempty"`
}

func (m *Media) GetVideoResolution() int {
	switch m.VideoResolution {
	case "480", "sd", "576":
		return 480
	case "720":
		return 720
	case "1080", "2k":
		return 1080
	case "2160", "4k":
		return 2160
	default:
		return 0
	}
}

type Part struct {
	ID                    int64   `json:"id"`
	Key                   string  `json:"key"`
	Duration              int64   `json:"duration"`
	File                  string  `json:"file"`
	Size                  int64   `json:"size"`
	AudioProfile          *string `json:"audioProfile,omitempty"`
	Container             string  `json:"container"`
	VideoProfile          string  `json:"videoProfile"`
	Has64BitOffsets       *bool   `json:"has64bitOffsets,omitempty"`
	OptimizedForStreaming *bool   `json:"optimizedForStreaming,omitempty"`
}

type Style string

const (
	Shelf Style = "shelf"
)

//
// type Hub struct {
// 	Title         string        `json:"title"`
// 	Type          string        `json:"type"`
// 	HubIdentifier string        `json:"hubIdentifier"`
// 	Context       string        `json:"context"`
// 	Size          int64         `json:"size"`
// 	More          bool          `json:"more"`
// 	Style         Style         `json:"style"`
// 	Metadata      []HubMetadata `json:"Metadata,omitempty"`
// 	Directory     []Directory   `json:"Directory,omitempty"`
// }
//
// type Directory struct {
// 	Key                 string `json:"key"`
// 	LibrarySectionID    int64  `json:"librarySectionID"`
// 	LibrarySectionKey   string `json:"librarySectionKey"`
// 	LibrarySectionTitle string `json:"librarySectionTitle"`
// 	LibrarySectionType  int64  `json:"librarySectionType"`
// 	Reason              string `json:"reason"`
// 	ReasonID            int64  `json:"reasonID"`
// 	ReasonTitle         string `json:"reasonTitle"`
// 	Score               string `json:"score"`
// 	Type                string `json:"type"`
// 	ID                  int64  `json:"id"`
// 	Filter              string `json:"filter"`
// 	Tag                 string `json:"tag"`
// 	TagType             int64  `json:"tagType"`
// 	Thumb               string `json:"thumb"`
// 	Art                 string `json:"art"`
// 	Count               int64  `json:"count"`
// 	GUID                string `json:"guid"`
// 	Summary             string `json:"summary"`
// }

type User struct {
	ID       string `json:"id" xml:"id,attr"`
	UUID     string `json:"uuid" xml:"uuid,attr"`
	Username string `json:"username" xml:"username,attr"`
	Title    string `json:"title" xml:"title,attr"`
	Email    string `json:"email" xml:"email,attr"`
	// FriendlyName      string      `json:"friendlyName"`
	// Locale            interface{} `json:"locale"`
	Confirmed bool  `json:"confirmed" xml:"confirmed,attr"`
	JoinedAt  int64 `json:"joinedAt" xml:"joinedAt,attr"`
	// EmailOnlyAuth     bool        `json:"emailOnlyAuth"`
	// HasPassword       bool        `json:"hasPassword"`
	// Protected         bool        `json:"protected"`
	Thumb string `json:"thumb" xml:"thumb,attr"`
	// AuthToken         string      `json:"authToken"`
	// MailingListStatus string      `json:"mailingListStatus"`
	// MailingListActive bool        `json:"mailingListActive"`
	// ScrobbleTypes     string      `json:"scrobbleTypes"`
	// Country           string      `json:"country"`
	// Pin               string      `json:"pin"`
	// Subscription      struct {
	// 	Active         bool     `json:"active"`
	// 	SubscribedAt   string   `json:"subscribedAt"`
	// 	Status         string   `json:"status"`
	// 	PaymentService string   `json:"paymentService"`
	// 	Plan           string   `json:"plan"`
	// 	Features       []string `json:"features"`
	// } `json:"subscription"`
	// SubscriptionDescription string `json:"subscriptionDescription"`
	// Restricted              bool   `json:"restricted"`
	// Anonymous               bool   `json:"anonymous"`
	Home bool `json:"home" xml:"home,attr"`
	// Guest                   bool   `json:"guest"`
	HomeSize  int64 `json:"homeSize" xml:"homeSize,attr"`
	HomeAdmin bool  `json:"homeAdmin" xml:"homeAdmin,attr"`
	// MaxHomeSize             int64  `json:"maxHomeSize"`
	// RememberExpiresAt       int64  `json:"rememberExpiresAt"`
	// Profile                 struct {
	// 	AutoSelectAudio              bool   `json:"autoSelectAudio"`
	// 	DefaultAudioLanguage         string `json:"defaultAudioLanguage"`
	// 	DefaultSubtitleLanguage      string `json:"defaultSubtitleLanguage"`
	// 	AutoSelectSubtitle           int64  `json:"autoSelectSubtitle"`
	// 	DefaultSubtitleAccessibility int64  `json:"defaultSubtitleAccessibility"`
	// 	DefaultSubtitleForced        int64  `json:"defaultSubtitleForced"`
	// } `json:"profile"`
	// Entitlements []string `json:"entitlements"`
	// Roles        []string `json:"roles"`
	// Services     []struct {
	// 	Identifier string  `json:"identifier"`
	// 	Endpoint   string  `json:"endpoint"`
	// 	Token      *string `json:"token"`
	// 	Secret     *string `json:"secret"`
	// 	Status     string  `json:"status"`
	// } `json:"services"`
	// AdsConsent           interface{} `json:"adsConsent"`
	// AdsConsentSetAt      interface{} `json:"adsConsentSetAt"`
	// AdsConsentReminderAt interface{} `json:"adsConsentReminderAt"`
	// ExperimentalFeatures bool        `json:"experimentalFeatures"`
	// TwoFactorEnabled     bool        `json:"twoFactorEnabled"`
	// BackupCodesCreated   bool        `json:"backupCodesCreated"`
}

type Director struct {
	ID        string  `json:"id"`
	Slug      string  `json:"slug"`
	Tag       string  `json:"tag"`
	Role      string  `json:"role"`
	Directory bool    `json:"directory"`
	Thumb     *string `json:"thumb,omitempty"`
}

type GUID struct {
	ID string `json:"id"`
}

type Genre struct {
	Filter      string  `json:"filter"`
	ID          string  `json:"id"`
	RatingKey   string  `json:"ratingKey"`
	Slug        string  `json:"slug"`
	Tag         string  `json:"tag"`
	Directory   bool    `json:"directory"`
	Context     string  `json:"context"`
	OriginalTag *string `json:"originalTag,omitempty"`
}

type Image struct {
	Alt  string `json:"alt"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Rating struct {
	Image string  `json:"image"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type Role struct {
	ID        string  `json:"id"`
	Order     int64   `json:"order"`
	Slug      string  `json:"slug"`
	Tag       string  `json:"tag"`
	Thumb     *string `json:"thumb,omitempty"`
	Role      *string `json:"role,omitempty"`
	Directory bool    `json:"directory"`
}

type Similar struct {
	GUID string `json:"guid"`
	Tag  string `json:"tag"`
}

type Account struct {
	ID    int64  `json:"id"`
	Thumb string `json:"thumb"`
	Title string `json:"title"`
}

type Player struct {
	Local         bool   `json:"local"`
	PublicAddress string `json:"publicAddress"`
	Title         string `json:"title"`
	UUID          string `json:"uuid"`
}

type Server struct {
	Title string `json:"title"`
	UUID  string `json:"uuid"`
}

type WebhookPayload struct {
	Event    string                  `json:"event"`
	User     bool                    `json:"user"`
	Owner    bool                    `json:"owner"`
	Account  *Account                `json:"Account"`
	Server   *Server                 `json:"Server"`
	Player   *Player                 `json:"Player"`
	Metadata *WebhookPayloadMetadata `json:"Metadata"`
}

type WebhookPayloadMetadata struct {
	LibrarySectionType   string `json:"librarySectionType"`
	RatingKey            string `json:"ratingKey"`
	Key                  string `json:"key"`
	ParentRatingKey      string `json:"parentRatingKey"`
	GrandparentRatingKey string `json:"grandparentRatingKey"`
	// GUID                 string `json:"Guid"`
	LibrarySectionID int64  `json:"librarySectionID"`
	Type             string `json:"type"`
	Title            string `json:"title"`
	GrandparentKey   string `json:"grandparentKey"`
	ParentKey        string `json:"parentKey"`
	GrandparentTitle string `json:"grandparentTitle"`
	ParentTitle      string `json:"parentTitle"`
	Summary          string `json:"summary"`
	Index            int64  `json:"index"`
	ParentIndex      int64  `json:"parentIndex"`
	RatingCount      int64  `json:"ratingCount"`
	Thumb            string `json:"thumb"`
	Art              string `json:"art"`
	ParentThumb      string `json:"parentThumb"`
	GrandparentThumb string `json:"grandparentThumb"`
	GrandparentArt   string `json:"grandparentArt"`
	AddedAt          int64  `json:"addedAt"`
	UpdatedAt        int64  `json:"updatedAt"`
}
