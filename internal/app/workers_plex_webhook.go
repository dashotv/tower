package app

import (
	"context"
	"regexp"
	"strconv"
	"time"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
	"github.com/dashotv/tower/internal/plex"
)

var plexSupportedHooks = []string{"library.new", "media.scrobble"}
var resolutionRegex = regexp.MustCompile(`(?i)(\d{3,4})p`)

type PlexWebhook struct {
	minion.WorkerDefaults[*PlexWebhook]
	Payload *plex.WebhookPayload `bson:"payload" json:"payload"`
}

func (j *PlexWebhook) Kind() string { return "plex_webhook" }
func (j *PlexWebhook) Work(ctx context.Context, job *minion.Job[*PlexWebhook]) error {
	payload := job.Args.Payload

	// notifier.Log.Infof("plex", "event: %s", payload.Event)

	switch payload.Event {
	case "library.new":
		return j.LibraryNew(ctx, payload)
	case "media.scrobble":
		return j.MediaScrobble(ctx, payload)
	}
	return nil
}

func (j *PlexWebhook) LibraryNew(ctx context.Context, payload *plex.WebhookPayload) error {
	notifier.Log.Debugf("plex", "added: %s | %s | %s", payload.Metadata.GrandparentTitle, payload.Metadata.ParentTitle, payload.Metadata.Title)
	a := ContextApp(ctx)

	m, err := j.mediumFromMetadata(ctx, payload.Metadata)
	if err != nil {
		return fae.Wrap(err, "medium from metadata")
	}
	if m == nil {
		return nil
	}

	if err := a.DB.Medium.Save(m); err != nil {
		return fae.Wrap(err, "save")
	}

	return nil
}

func (j *PlexWebhook) MediaScrobble(ctx context.Context, payload *plex.WebhookPayload) error {
	notifier.Log.Debugf("plex", "scrobble: %s - %s | %s | %s", payload.Account.Title, payload.Metadata.GrandparentTitle, payload.Metadata.ParentTitle, payload.Metadata.Title)
	a := ContextApp(ctx)

	m, err := j.mediumFromMetadata(ctx, payload.Metadata)
	if err != nil {
		return fae.Wrap(err, "medium from metadata")
	}
	if m == nil {
		return fae.New("no medium found")
	}

	m.Downloaded = true
	m.Completed = true

	if err := a.DB.Medium.Save(m); err != nil {
		return fae.Wrap(err, "save medium")
	}

	w := &Watch{
		MediumID:  m.ID,
		Username:  payload.Account.Title,
		Player:    payload.Player.Title,
		WatchedAt: time.Now(),
	}
	if err := a.DB.Watch.Save(w); err != nil {
		return fae.Wrap(err, "save watch")
	}

	return nil
}

func (j *PlexWebhook) mediumFromMetadata(ctx context.Context, metadata *plex.WebhookPayloadMetadata) (*Medium, error) {
	a := ContextApp(ctx)

	resp, err := a.Plex.GetMetadataByKey(metadata.RatingKey)
	if err != nil {
		return nil, fae.Wrap(err, "get metadata")
	}
	if len(resp) == 0 || len(resp[0].Media) == 0 || len(resp[0].Media[0].Part) == 0 {
		a.Log.Named("plex_library_new").Warn("no media/part found in metadata")
		return nil, nil
	}

	m, err := a.DB.MediumByPlexMedia(resp[0].Media[0])
	if err != nil {
		return nil, fae.Wrap(err, "medium by plex media")
	}

	m.AddPathsByMetadata(resp[0])
	return m, nil
}

func mediumAddPath(m *Medium, file string, size int64, resolution int) *Path {
	p := m.AddPathByFullpath(file)
	if p == nil {
		return nil
	}

	p.Size = size
	p.Resolution = resolution
	return p
}

func metadataResolution(res string) int {
	if res == "" {
		return 0
	}
	if match := resolutionRegex.FindStringSubmatch(res); len(match) > 1 {
		r, _ := strconv.Atoi(match[1])
		return r
	}
	return 0
}
