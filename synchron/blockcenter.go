package synchron

import (
	"time"

	"github.com/bytom/errors"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	"github.com/bufferserver/api/common"
	"github.com/bufferserver/config"
	"github.com/bufferserver/database"
	"github.com/bufferserver/database/orm"
	"github.com/bufferserver/service"
)

func BlockKeeper(cfg *config.Config, db *gorm.DB, cache *database.RedisDB, node *service.Node, duration time.Duration) {
	ticker := time.NewTicker(duration)
	for ; true; <-ticker.C {
		if err := syncBlockCenter(db, cache, node); err != nil {
			log.WithField("err", err).Errorf("fail on blockcenter")
		}
	}
}

func syncBlockCenter(db *gorm.DB, cache *database.RedisDB, node *service.Node) error {
	var bases []*orm.Base
	if err := db.Find(&bases).Error; err != nil {
		return errors.Wrap(err, "query bases")
	}

	filter := make(map[string]interface{})
	for _, base := range bases {
		filter["asset"] = base.AssetID
		filter["script"] = base.ControlProgram
		req := &common.Display{Filter: filter}
		resUTXOs, err := node.ListBlockCenterUTXOs(req)
		if err != nil {
			return err
		}

		for _, utxo := range resUTXOs {
			if err := db.Save(utxo).Error; err != nil {
				return errors.Wrap(err, "update utxo")
			}
		}
	}
	return nil
}
