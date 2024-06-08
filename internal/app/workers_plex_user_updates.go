package app

import (
	"context"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

// PlexUserUpdates updates users from plex
type PlexUserUpdates struct {
	minion.WorkerDefaults[*PlexUserUpdates]
}

func (j *PlexUserUpdates) Kind() string { return "PlexUserUpdates" }
func (j *PlexUserUpdates) Work(ctx context.Context, job *minion.Job[*PlexUserUpdates]) error {
	// app.Log.Debugf("updating users")

	users, err := app.DB.User.Query().NotEqual("token", "").Run()
	if err != nil {
		return fae.Wrap(err, "querying users")
	}

	for _, u := range users {
		data, err := app.Plex.GetUser(u.Token)
		if err != nil {
			notifier.Log.Errorf("PlexUser", "getting plex user: %s: %v", u.Email, err)
			continue
		}

		u.Name = data.Username
		u.Email = data.Email
		u.Thumb = data.Thumb
		u.Home = data.Home
		u.Admin = data.HomeAdmin

		app.Log.Debugf("updating user %s", u.Name)
		err = app.DB.User.Update(u)
		if err != nil {
			return fae.Wrap(err, "updating user")
		}
	}

	if err := app.Workers.Enqueue(&PlexWatchlistUpdates{}); err != nil {
		return fae.Wrap(err, "enqueuing worker")
	}

	return nil
}
