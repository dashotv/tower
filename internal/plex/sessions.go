package plex

func (p *Client) GetSessions() ([]*PlexSessionMetadata, error) {
	sessions := &PlexSessionContainer{}
	resp, err := p._server().SetResult(sessions).Get("/status/sessions")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, err
	}
	if sessions.MediaContainer.Size == 0 {
		return nil, nil
	}
	return sessions.MediaContainer.Metadata, nil
}

type PlexSessionContainer struct {
	MediaContainer struct {
		Size     int64                  `json:"size"`
		Metadata []*PlexSessionMetadata `json:"Metadata"`
	} `json:"MediaContainer"`
}

type PlexSessionMetadata struct {
	AddedAt               int64           `json:"addedAt"`
	Art                   string          `json:"art"`
	ContentRating         string          `json:"contentRating"`
	Duration              int64           `json:"duration"`
	GrandparentArt        string          `json:"grandparentArt"`
	GrandparentGUID       string          `json:"grandparentGuid"`
	GrandparentKey        string          `json:"grandparentKey"`
	GrandparentRatingKey  string          `json:"grandparentRatingKey"`
	GrandparentSlug       string          `json:"grandparentSlug"`
	GrandparentThumb      string          `json:"grandparentThumb"`
	GrandparentTitle      string          `json:"grandparentTitle"`
	GUID                  string          `json:"guid"`
	Index                 int64           `json:"index"`
	Key                   string          `json:"key"`
	LibrarySectionID      string          `json:"librarySectionID"`
	LibrarySectionKey     string          `json:"librarySectionKey"`
	LibrarySectionTitle   string          `json:"librarySectionTitle"`
	OriginalTitle         string          `json:"originalTitle"`
	OriginallyAvailableAt string          `json:"originallyAvailableAt"`
	ParentGUID            string          `json:"parentGuid"`
	ParentIndex           int64           `json:"parentIndex"`
	ParentKey             string          `json:"parentKey"`
	ParentRatingKey       string          `json:"parentRatingKey"`
	ParentTitle           string          `json:"parentTitle"`
	RatingKey             string          `json:"ratingKey"`
	SessionKey            string          `json:"sessionKey"`
	Summary               string          `json:"summary"`
	Thumb                 string          `json:"thumb"`
	Title                 string          `json:"title"`
	Type                  string          `json:"type"`
	UpdatedAt             int64           `json:"updatedAt"`
	ViewOffset            int64           `json:"viewOffset"`
	Year                  int64           `json:"year"`
	Media                 []*SessionMedia `json:"Media"`
	// Role                  []*Role           `json:"Role"`
	User             *User             `json:"User"`
	Player           *SessionPlayer    `json:"Player"`
	Session          *Session          `json:"Session"`
	TranscodeSession *TranscodeSession `json:"TranscodeSession"`
}

type SessionMedia struct {
	AudioProfile          string         `json:"audioProfile"`
	ID                    string         `json:"id"`
	VideoProfile          string         `json:"videoProfile"`
	AudioChannels         int64          `json:"audioChannels"`
	AudioCodec            string         `json:"audioCodec"`
	Bitrate               int64          `json:"bitrate"`
	Container             string         `json:"container"`
	Duration              int64          `json:"duration"`
	Height                int64          `json:"height"`
	OptimizedForStreaming bool           `json:"optimizedForStreaming"`
	Protocol              string         `json:"protocol"`
	VideoCodec            string         `json:"videoCodec"`
	VideoFrameRate        string         `json:"videoFrameRate"`
	VideoResolution       string         `json:"videoResolution"`
	Width                 int64          `json:"width"`
	Selected              bool           `json:"selected"`
	Part                  []*SessionPart `json:"Part"`
}

type SessionPart struct {
	AudioProfile          string    `json:"audioProfile"`
	ID                    string    `json:"id"`
	VideoProfile          string    `json:"videoProfile"`
	Bitrate               int64     `json:"bitrate"`
	Container             string    `json:"container"`
	Duration              int64     `json:"duration"`
	Height                int64     `json:"height"`
	OptimizedForStreaming bool      `json:"optimizedForStreaming"`
	Protocol              string    `json:"protocol"`
	Width                 int64     `json:"width"`
	Decision              string    `json:"decision"`
	Selected              bool      `json:"selected"`
	Stream                []*Stream `json:"Stream"`
}

type Stream struct {
	Bitrate              int64    `json:"bitrate"`
	Codec                string   `json:"codec"`
	Default              bool     `json:"default"`
	DisplayTitle         string   `json:"displayTitle"`
	ExtendedDisplayTitle string   `json:"extendedDisplayTitle"`
	FrameRate            *float64 `json:"frameRate,omitempty"`
	Height               *int64   `json:"height,omitempty"`
	ID                   string   `json:"id"`
	Language             string   `json:"language"`
	LanguageCode         string   `json:"languageCode"`
	LanguageTag          string   `json:"languageTag"`
	StreamType           int64    `json:"streamType"`
	Width                *int64   `json:"width,omitempty"`
	Decision             string   `json:"decision"`
	Location             string   `json:"location"`
	AudioChannelLayout   *string  `json:"audioChannelLayout,omitempty"`
	BitrateMode          *string  `json:"bitrateMode,omitempty"`
	Channels             *int64   `json:"channels,omitempty"`
	Profile              *string  `json:"profile,omitempty"`
	SamplingRate         *int64   `json:"samplingRate,omitempty"`
	Selected             *bool    `json:"selected,omitempty"`
	Burn                 *string  `json:"burn,omitempty"`
	Title                *string  `json:"title,omitempty"`
}

type SessionPlayer struct {
	Address             string `json:"address"`
	Device              string `json:"device"`
	MachineIdentifier   string `json:"machineIdentifier"`
	Model               string `json:"model"`
	Platform            string `json:"platform"`
	PlatformVersion     string `json:"platformVersion"`
	Product             string `json:"product"`
	Profile             string `json:"profile"`
	RemotePublicAddress string `json:"remotePublicAddress"`
	State               string `json:"state"`
	Title               string `json:"title"`
	Version             string `json:"version"`
	Local               bool   `json:"local"`
	Relayed             bool   `json:"relayed"`
	Secure              bool   `json:"secure"`
	UserID              int64  `json:"userID"`
}

//
// type Role struct {
// 	Filter string `json:"filter"`
// 	ID     string `json:"id"`
// 	Role   string `json:"role"`
// 	Tag    string `json:"tag"`
// 	TagKey string `json:"tagKey"`
// 	Thumb  string `json:"thumb"`
// }

type Session struct {
	ID        string `json:"id"`
	Bandwidth int64  `json:"bandwidth"`
	Location  string `json:"location"`
}

type TranscodeSession struct {
	Key                      string  `json:"key"`
	Throttled                bool    `json:"throttled"`
	Complete                 bool    `json:"complete"`
	Progress                 float64 `json:"progress"`
	Size                     int64   `json:"size"`
	Speed                    float64 `json:"speed"`
	Error                    bool    `json:"error"`
	Duration                 int64   `json:"duration"`
	Remaining                int64   `json:"remaining"`
	Context                  string  `json:"context"`
	SourceVideoCodec         string  `json:"sourceVideoCodec"`
	SourceAudioCodec         string  `json:"sourceAudioCodec"`
	VideoDecision            string  `json:"videoDecision"`
	AudioDecision            string  `json:"audioDecision"`
	SubtitleDecision         string  `json:"subtitleDecision"`
	Protocol                 string  `json:"protocol"`
	Container                string  `json:"container"`
	VideoCodec               string  `json:"videoCodec"`
	AudioCodec               string  `json:"audioCodec"`
	AudioChannels            int64   `json:"audioChannels"`
	TranscodeHwRequested     bool    `json:"transcodeHwRequested"`
	TranscodeHwDecoding      string  `json:"transcodeHwDecoding"`
	TranscodeHwEncoding      string  `json:"transcodeHwEncoding"`
	TranscodeHwDecodingTitle string  `json:"transcodeHwDecodingTitle"`
	TranscodeHwFullPipeline  bool    `json:"transcodeHwFullPipeline"`
	TranscodeHwEncodingTitle string  `json:"transcodeHwEncodingTitle"`
	TimeStamp                float64 `json:"timeStamp"`
	MaxOffsetAvailable       float64 `json:"maxOffsetAvailable"`
	MinOffsetAvailable       float64 `json:"minOffsetAvailable"`
}

type User struct {
	ID    string `json:"id"`
	Thumb string `json:"thumb"`
	Title string `json:"title"`
}
