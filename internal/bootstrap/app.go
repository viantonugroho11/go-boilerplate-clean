package bootstrap

// RunApp memuat config (global), wiring terisolasi (DB, Redis, Kafka, routes), lalu jalankan HTTP server sampai signal.
func RunApp() error {
	if err := LoadConfig(); err != nil {
		return err
	}

	db, err := initDB()
	if err != nil {
		return err
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	userService, closeDeps, err := wireUserService(db)
	if err != nil {
		return err
	}
	defer closeDeps()

	redisClient, err := initRedis()
	if err != nil {
		return err
	}
	defer redisClient.Close()

	e := newEcho(userService)
	return runHTTP(e)
}
