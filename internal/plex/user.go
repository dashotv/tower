package plex

import "github.com/dashotv/fae"

type ApiUser struct {
	ID       int64  `json:"id" xml:"id,attr"`
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

func (p *Client) GetUser(token string) (*ApiUser, error) {
	user := &ApiUser{}
	resp, err := p._plextv().SetResult(user).
		SetHeader("X-Plex-Token", token).Get("/user")
	if err != nil {
		return nil, fae.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return nil, fae.Errorf("getting user: %s", resp.Status())
	}

	return user, nil
}

func (p *Client) GetServicesUser(token string) (*ApiUser, error) {
	user := &ApiUser{}
	resp, err := p._plextv().SetResult(user).
		SetHeader("X-Plex-Token", token).Get("/services/user")
	if err != nil {
		return nil, fae.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return nil, fae.Errorf("getting user: %s", resp.Status())
	}

	return user, nil
}
