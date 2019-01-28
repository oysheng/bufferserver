package main

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bufferserver/config"
	"github.com/bufferserver/database"
	"github.com/bufferserver/service"
	"github.com/bufferserver/synchron"
)

func main() {
	cfg := config.NewConfig()
	db, err := database.NewMySQLDB(cfg.MySQL, cfg.Updater.MySQLConnCfg)
	if err != nil {
		log.WithField("err", err).Panic("initialize mysql db error")
	}

	cache, err := database.NewRedisDB(cfg.Redis)
	if err != nil {
		log.WithField("err", err).Panic("initialize redis db error")
	}

	node := service.NewNode(cfg.Updater.URL)
	blockSyncFreq := time.Duration(cfg.Updater.SyncSeconds) * time.Second
	go synchron.BlockCenterKeeper(cfg, db.Master(), cache, node, blockSyncFreq)

	// keep the main func running in case of terminating goroutines
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
