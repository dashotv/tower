package app

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var minion *Minion

func setupWorkers() error {
	minion = NewMinion(cfg.Minion.Concurrency)
	return nil
}

type Minion struct {
	Concurrency int
	Queue       chan *Job
	Cron        *cron.Cron
	Log         *zap.SugaredLogger
}

type MinionFunc func(id int, log *zap.SugaredLogger) error

func NewMinion(concurrency int) *Minion {
	return &Minion{
		Concurrency: concurrency,
		Log:         log.Named("minion"),
		Queue:       make(chan *Job, concurrency*concurrency),
		Cron:        cron.New(cron.WithSeconds()),
	}
}

func (m *Minion) Start() error {
	m.Log.Infof("starting minion (concurrency=%d)...", m.Concurrency)

	if cfg.Cron {
		// every 5 seconds DownloadsProcess
		// if err := m.AddCron("*/5 * * * * *", "DownloadsProcess", s.DownloadsProcess); err != nil {
		// 	return err
		// }
		// if err := m.AddCron("*/5 * * * * *", "CausingErrors", s.CausingErrors); err != nil {
		// 	return err
		// }

		// every 5 minutes
		if err := m.AddCron("0 */5 * * * *", "PopularReleases", m.PopularReleases); err != nil {
			return err
		}
		// every 15 minutes
		// if err := m.AddCron("0 */15 * * * *", "ProcessFeeds", s.ProcessFeeds); err != nil {
		// 	return err
		// }

		// every day at 3am
		if err := m.AddCron("0 0 3 * * *", "CleanPlexPins", m.CleanPlexPins); err != nil {
			return err
		}
		// every day at 3am
		if err := m.AddCron("0 0 3 * * *", "CleanJobs", m.CleanJobs); err != nil {
			return err
		}
	}

	for w := 0; w < m.Concurrency; w++ {
		worker := &Worker{w, m.Log.Named(fmt.Sprintf("worker:%d", w)), m.Queue}
		go worker.Run()
	}

	go func() {
		m.Cron.Start()
	}()

	return nil
}

func (m *Minion) Add(name string, f MinionFunc) error {
	j := &MinionJob{
		Name: name,
	}

	err := db.MinionJob.Save(j)
	if err != nil {
		return errors.Wrap(err, "failed to save minion job")
	}

	mf := func(id int, log *zap.SugaredLogger) error {
		log.Infof("starting %s: %s", name, j.ID.Hex())
		err := f(id, log)

		j.ProcessedAt = time.Now()
		if err != nil {
			log.Errorf("processing %s: %s: %s", name, j.ID.Hex(), err)
			j.Error = errors.Wrap(err, "failed to run minion job").Error()
		}

		err = db.MinionJob.Update(j)
		if err != nil {
			log.Errorf("error %s: %s: %s", name, j.ID.Hex(), err)
			return errors.Wrap(err, "failed to save minion job")
		}

		log.Infof("finished %s: %s", name, j.ID.Hex())
		return nil
	}

	m.Queue <- &Job{ID: j.ID.Hex(), Func: mf}
	return nil
}

func (m *Minion) AddCron(spec, name string, f MinionFunc) error {
	_, err := m.Cron.AddFunc(spec, func() {
		minion.Add(name, func(id int, log *zap.SugaredLogger) error {
			return f(id, log)
		})
	})

	if err != nil {
		return errors.Wrap(err, "adding cron function")
	}

	return nil
}

type Job struct {
	ID   string // reference to MinionJob in db
	Func MinionFunc
}

func (j *Job) Run(id int, log *zap.SugaredLogger) error {
	return j.Func(id, log)
}

type Worker struct {
	ID    int
	Log   *zap.SugaredLogger
	Queue chan *Job
}

func (w *Worker) Run() {
	for j := range w.Queue {
		j.Run(w.ID, w.Log)
	}
}

func (m *Minion) CausingErrors(id int, log *zap.SugaredLogger) error {
	log.Info("causing error")
	return nil
}
func (m *Minion) DownloadsProcess(id int, log *zap.SugaredLogger) error {
	log.Info("processing downloads")
	return nil
}

func (m *Minion) ProcessFeeds(id int, log *zap.SugaredLogger) error {
	m.Log.Info("processing feeds")
	return db.ProcessFeeds()
}

func (m *Minion) CleanPlexPins(id int, log *zap.SugaredLogger) error {
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

func (m *Minion) CleanJobs(id int, log *zap.SugaredLogger) error {
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

func (m *Minion) PopularReleases(id int, log *zap.SugaredLogger) error {
	limit := 25
	intervals := map[string]int{
		"daily":   1,
		"weekly":  7,
		"monthly": 30,
	}

	start := time.Now()
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

	diff := time.Since(start)
	log.Infof("PopularReleases: took %s", diff)

	return nil
}
