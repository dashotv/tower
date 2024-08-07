package app

import (
	"fmt"
	"time"

	runic "github.com/dashotv/runic/client"
)

func onRunicReleases(a *Application, msg *runic.Release) error {
	log := a.Log.Named("runic.releases")
	// log.Infof("received: '%s' %02dx%02d", msg.Title, msg.Season, msg.Episode)

	if !a.Config.ProcessRunicEvents {
		// log.Warnf("skipping: runic events disabled")
		return nil
	}

	if a.Want == nil {
		log.Warnf("want not initialized, waiting...")
		for i := 0; i < 10 && a.Want == nil; i++ {
			time.Sleep(1 * time.Second)
		}
		if a.Want == nil {
			log.Errorf("want not initialized")
			return nil
		}
	}

	medium := a.Want.Release(msg)
	if medium == nil {
		// log.Debugf("skipping: [%s] %s (%d) %dx%d: not wanted", msg.Type, msg.Title, msg.Year, msg.Season, msg.Episode)
		return nil
	}

	log.Debugf("found: %s s%02de%02d", msg.Title, msg.Season, msg.Episode)
	var d *Download

	downloads, err := a.DB.Download.Query().Where("medium_id", medium.ID).Run()
	if err != nil {
		return err
	}

	switch len(downloads) {
	case 0:
		d = &Download{MediumID: medium.ID}
	case 1:
		if downloads[0].Status != "searching" {
			log.Warnf("skipping: %s s%02de%02d: download exists", msg.Title, msg.Season, msg.Episode)
			return nil
		}
		d = downloads[0]
	default:
		log.Warnf("skipping: %s s%02de%02d: multiple download exists", msg.Title, msg.Season, msg.Episode)
		return nil
	}

	d.SetRelease(msg)
	if !a.Config.Production {
		d.Status = "reviewing"
	}
	if err := a.DB.Download.Save(d); err != nil {
		return err
	}

	medium.Downloaded = true
	if err := a.DB.Medium.Save(medium); err != nil {
		return err
	}

	notifier.Info("EventFound", fmt.Sprintf("%s (%d) S%02dE%02d", msg.Title, msg.Year, msg.Season, msg.Episode))
	return nil
}
