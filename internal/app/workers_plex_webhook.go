package app

import (
	"context"
	"time"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
	"github.com/dashotv/tower/internal/plex"
)

var plexSupportedHooks = []string{"library.new", "media.scrobble"}

type PlexWebhook struct {
	minion.WorkerDefaults[*PlexWebhook]
	Payload *plex.WebhookPayload `bson:"payload" json:"payload"`
}

func (j *PlexWebhook) Kind() string { return "plex_webhook" }
func (j *PlexWebhook) Work(ctx context.Context, job *minion.Job[*PlexWebhook]) error {
	payload := job.Args.Payload

	switch payload.Event {
	case "library.new":
		return j.LibraryNew(ctx, payload)
	case "media.scrobble":
		return j.MediaScrobble(ctx, payload)
	}
	return nil
}

func (j *PlexWebhook) LibraryNew(ctx context.Context, payload *plex.WebhookPayload) error {
	a := ContextApp(ctx)

	m, err := j.mediumFromMetadata(ctx, payload.Metadata)
	if err != nil {
		return fae.Wrap(err, "medium from metadata")
	}
	if m == nil {
		return fae.Errorf("medium not found: %s", payload.Metadata.RatingKey)
	}

	if err := a.DB.Medium.Save(m); err != nil {
		return fae.Wrap(err, "save")
	}

	return nil
}

func (j *PlexWebhook) MediaScrobble(ctx context.Context, payload *plex.WebhookPayload) error {
	a := ContextApp(ctx)

	m, err := j.mediumFromMetadata(ctx, payload.Metadata)
	if err != nil {
		return fae.Wrap(err, "medium from metadata")
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
		a.Log.Named("plex_library_new").Warn("no media found in metadata")
		return nil, nil
	}

	kind, name, file, ext, err := pathParts(resp[0].Media[0].Part[0].File)
	if err != nil {
		return nil, fae.Wrap(err, "path parts")
	}

	m, ok, err := a.DB.MediumBy(kind, name, file, ext)
	if err != nil {
		return nil, fae.Wrap(err, "medium by")
	}
	if !ok {
		return nil, fae.New("medium not found")
	}

	p := m.AddPathByFullpath(file)
	if p == nil {
		return nil, fae.Errorf("failed to add path: %s", file)
	}

	return m, nil
}
