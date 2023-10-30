package app

var initialized bool
var db *Connector

type SetupFunc func() error

func Start() error {
	err := setup(
		setupConfig,
		setupLogger,
		setupDb,
		setupServer,
		setupCache,
		setupWorkers,
		setupEvents,
		setupPlex,
	)
	if err != nil {
		return err
	}

	initialized = true
	log.Info("initialized: ", initialized)
	log.Debugf("config: %+v", cfg)

	return server.Start()
}

func setup(fs ...SetupFunc) error {
	for _, f := range fs {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}

func setupDb() (err error) {
	db, err = NewConnector()
	if err != nil {
		return err
	}

	return nil
}
