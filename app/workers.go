package app

import (
	"fmt"
	"time"

	"github.com/dashotv/minion"
	"github.com/pkg/errors"
)

var workers *minion.Minion

var jobs = map[string]minion.Func{
	"PopularReleases": PopularReleases,
	"CleanPlexPins":   CleanPlexPins,
	"CleanJobs":       CleanJobs,
	// "DownloadsProcess": DownloadsProcess,
	// "UserWatchlist":   UserWatchlist,
}

func setupWorkers() error {
	workers = minion.New(cfg.Minion.Concurrency)

	for n, f := range jobs {
		workers.Register(n, wrapJob(n, f))
	}

	if cfg.Cron {
		// every 5 seconds
		// if _, err := workers.Schedule("*/5 * * * * *", "DownloadsProcess"); err != nil {
		// 	return err
		// }
		// if _, err := workers.Schedule("*/5 * * * * *", "CausingErrors"); err != nil {
		// 	return err
		// }

		// every 5 minutes
		if _, err := workers.Schedule("0 */5 * * * *", "PopularReleases"); err != nil {
			return err
		}
		// every 15 minutes
		// if _, err := workers.Schedule("0 */15 * * * *", "ProcessFeeds"); err != nil {
		// 	return err
		// }

		// every day at 3am (11am UTC)
		if _, err := workers.Schedule("0 0 11 * * *", "CleanPlexPins"); err != nil {
			return err
		}
		// every day at 3am (11am UTC)
		if _, err := workers.Schedule("0 0 11 * * *", "CleanJobs"); err != nil {
			return err
		}
		// every 15 minutes
		// if _, err := workers.Schedule("0 */15 * * * *", "UserWatchlist"); err != nil {
		// 	return err
		// }
	}

	return nil
}

func CleanPlexPins() error {
	list, err := db.Pin.Query().
		GreaterThan("created_at", time.Now().UTC().AddDate(0, 0, -1)).
		Run()
	if err != nil {
		return errors.Wrap(err, "querying pins")
	}

	for _, p := range list {
		err := db.Pin.Delete(p)
		if err != nil {
			return errors.Wrap(err, "deleting pin")
		}
	}

	return nil
}

func CleanJobs() error {
	list, err := db.MinionJob.Query().
		GreaterThan("created_at", time.Now().UTC().AddDate(0, 0, -1)).
		Run()
	if err != nil {
		return errors.Wrap(err, "querying jobs")
	}

	for _, j := range list {
		err := db.MinionJob.Delete(j)
		if err != nil {
			return errors.Wrap(err, "deleting job")
		}
	}

	return nil
}

func PopularReleases() error {
	limit := 25
	intervals := map[string]int{
		"daily":   1,
		"weekly":  7,
		"monthly": 30,
	}

	for f, i := range intervals {
		for _, t := range releaseTypes {
			date := time.Now().AddDate(0, 0, -i)

			results, err := db.ReleasesPopularQuery(t, date, limit)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("popular releases %s %s", f, t))
			}

			cache.Set(fmt.Sprintf("releases_popular_%s_%s", f, t), results)
		}
	}

	return nil
}

func UserWatchlist() error {
	return nil
}

func CausingErrors() error {
	log.Info("causing error")
	return nil
}

func DownloadsProcess() error {
	log.Info("processing downloads")
	return nil
}

func ProcessFeeds() error {
	log.Info("processing feeds")
	return db.ProcessFeeds()
}

func wrapJob(name string, f func() error) func() error {
	return func() error {
		j := &MinionJob{Name: name}

		err := db.MinionJob.Save(j)
		if err != nil {
			return errors.Wrap(err, "saving job")
		}

		start := time.Now()
		err = f()
		if err != nil {
			j.Error = err.Error()
		}

		j.ProcessedAt = time.Now()
		duration := time.Since(start)
		j.Duration = duration.Seconds()

		err = db.MinionJob.Update(j)
		if err != nil {
			return errors.Wrap(err, "updating job")
		}

		log.Infof("job:%s: %s", name, duration)
		return nil
	}
}
