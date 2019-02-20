package main

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/bufferserver/config"
	"github.com/bufferserver/database"
	"github.com/bufferserver/synchron"
)

func main() {
	cfg := config.NewConfig()
	db, err := database.NewMySQLDB(cfg.MySQL, cfg.Updater.BlockCenter.MySQLConnCfg)
	if err != nil {
		log.WithField("err", err).Panic("initialize mysql db error")
	}

	go synchron.NewBlockCenterKeeper(cfg, db.Master()).Run()
	go synchron.NewBrowserKeeper(cfg, db.Master()).Run()

	// keep the main func running in case of terminating goroutines
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
