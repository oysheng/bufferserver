package main

//import (
//	"sync"
//	"time"
//
//	log "github.com/sirupsen/logrus"
//
//	"github.com/blockcenter/config"
//	"github.com/blockcenter/database"
//	"github.com/blockcenter/service"
//)
//
//func main() {
//	cfg := config.NewConfig()
//	db, err := database.NewMySQLDB(cfg.MySQL, cfg.Updater.MySQLConnCfg)
//	if err != nil {
//		log.WithField("err", err).Panic("initialize mysql db error")
//	}
//
//	cache, err := database.NewRedisDB(cfg.Redis)
//	if err != nil {
//		log.WithField("err", err).Panic("initialize redis db error")
//	}
//
//	market := service.NewMarket(cfg.Updater.Market.ExRateServerURL, cfg.Updater.Market.PriceServerURL)
//	marketSyncFreq := time.Duration(cfg.Updater.Market.SyncSeconds) * time.Second
//	go synchron.MarketKeeper(db.Master(), cache, market, marketSyncFreq, "btm")
//
//	go synchron.UnconfirmedTxKeeper(cfg, db, cache, "btm")
//
//	node := service.NewNode(cfg.Coin.Btm.Upstream.URL)
//	blockSyncFreq := time.Duration(cfg.Updater.SyncSeconds) * time.Second
//	go synchron.BlockKeeper(cfg, db.Master(), cache, node, blockSyncFreq, "btm")
//
//	// keep the main func running in case of terminating goroutines
//	var wg sync.WaitGroup
//	wg.Add(1)
//	wg.Wait()
//}
