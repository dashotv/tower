package app

import (
	"context"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

// PlexPinToUsers ensures users are created from athorized pins
type PlexPinToUsers struct {
	minion.WorkerDefaults[*PlexPinToUsers]
}

func (j *PlexPinToUsers) Kind() string { return "PlexPinToUsers" }
func (j *PlexPinToUsers) Work(ctx context.Context, job *minion.Job[*PlexPinToUsers]) error {
	// app.Log.Debugf("creating users from authenticated pins")

	pins, err := app.DB.Pin.Query().Run()
	if err != nil {
		return fae.Wrap(err, "querying pins")
	}

	check := map[string]bool{}
	app.Log.Debugf("ranging pins")
	for _, p := range pins {
		if p.Token == "" {
			continue
		}

		if check[p.Token] {
			continue
		}

		check[p.Token] = true
		// app.Log.Debugf("find user by token %s", p.Token)
		// resp, err := app.DB.User.Query().Where("token", p.Token).Run()
		// if err != nil {
		// 	return fae.Wrap(err, "querying user")
		// }
		// if len(resp) > 0 {
		// 	// users exists
		// 	continue
		// }

		plexUser, err := app.Plex.GetUser(p.Token)
		if err != nil {
			return fae.Wrap(err, "getting plex user")
		}

		list, err := app.DB.User.Query().Where("email", plexUser.Email).Run()
		if err != nil {
			return fae.Wrap(err, "querying user")
		}
		if len(list) > 0 {
			if list[0].Token != p.Token {
				list[0].Token = p.Token
				err = app.DB.User.Save(list[0])
				if err != nil {
					return fae.Wrap(err, "saving user")
				}
			}
			continue
		}

		// create user
		user := &User{
			Email: plexUser.Email,
			Name:  plexUser.Title,
			Thumb: plexUser.Thumb,
			Home:  plexUser.Home,
			Admin: false,
			Token: p.Token,
		}
		err = app.DB.User.Save(user)
		if err != nil {
			return fae.Wrap(err, "saving user")
		}
	}

	if err := app.Workers.Enqueue(&PlexUserUpdates{}); err != nil {
		return fae.Wrap(err, "enqueuing worker")
	}

	return nil
}
