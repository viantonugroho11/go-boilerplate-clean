package bootstrap

// RunApp loads config, wires dependencies, and runs the HTTP server until signal.
func RunApp() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	db, err := initDB(cfg)
	if err != nil {
		return err
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	userService, closeUser, err := wireUserService(cfg, db)
	if err != nil {
		return err
	}
	defer closeUser()

	redisClient, err := initRedis(cfg)
	if err != nil {
		return err
	}
	defer redisClient.Close()

	sampleService, closeSample, err := wireSampleService(cfg, db)
	if err != nil {
		return err
	}
	defer closeSample()

	e := newEcho(userService, sampleService)
	return runHTTP(cfg, e)
}
