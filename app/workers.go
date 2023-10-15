package app

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Minion struct {
	Concurrency int
	Queue       chan *Job
	Log         *zap.SugaredLogger
}

type MinionFunc func(id int, log *zap.SugaredLogger) error

func NewMinion(concurrency int) *Minion {
	return &Minion{
		Concurrency: concurrency,
		Log:         log.Named("minion"),
		Queue:       make(chan *Job, concurrency*concurrency),
	}
}

func (m *Minion) Start() {
	m.Log.Infof("starting minion (concurrency=%d)...", m.Concurrency)
	for w := 0; w < m.Concurrency; w++ {
		worker := &Worker{w, m.Log.Named(fmt.Sprintf("worker:%d", w)), m.Queue}
		go worker.Run()
	}
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
		log.Infof("starting %s", name)
		err := f(id, log)

		j.ProcessedAt = time.Now()
		if err != nil {
			log.Errorf("processing %s: %s", name, err)
			j.Error = errors.Wrap(err, "failed to run minion job").Error()
		}

		err = db.MinionJob.Update(j)
		if err != nil {
			log.Errorf("error %s: %s", name, err)
			return errors.Wrap(err, "failed to save minion job")
		}

		log.Infof("finished %s", name)
		return nil
	}
	m.Queue <- &Job{ID: j.ID.Hex(), Func: mf}
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
