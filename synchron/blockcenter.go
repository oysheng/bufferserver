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
	"github.com/bufferserver/types"
)

func BlockCenterKeeper(cfg *config.Config, db *gorm.DB, cache *database.RedisDB, node *service.Node, duration time.Duration) {
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
			return errors.Wrap(err, "list blockcenter utxos")
		}

		for _, utxo := range resUTXOs {
			u := orm.Utxo{Hash: utxo.Hash}
			if err := db.Where(u).First(&u).Error; err != nil && err != gorm.ErrRecordNotFound {
				return errors.Wrap(err, "query utxo")
			} else if err == gorm.ErrRecordNotFound {
				butxo := &orm.Utxo{Hash: utxo.Hash, AssetID: utxo.Asset, Amount: utxo.Amount, ControlProgram: base.ControlProgram, IsSpend: false, IsLocked: false, Duration: uint64(60)}
				if err := db.Save(butxo).Error; err != nil {
					return errors.Wrap(err, "save utxo")
				}
				continue
			}

			currentTime := time.Now()
			if (u.SubmitTime != types.Timestamp{}) && (currentTime.Unix()-u.SubmitTime.Unix()) > int64(u.Duration) {
				if err := db.Model(&orm.Utxo{}).Where(&orm.Utxo{Hash: utxo.Hash}).Update("is_locked", false).Error; err != nil {
					return errors.Wrap(err, "update utxo unlocked")
				}
			}
		}
	}
	return nil
}
