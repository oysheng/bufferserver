package synchron

import (
	"time"

	"github.com/bytom/errors"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	"github.com/bufferserver/api/common"
	"github.com/bufferserver/database/orm"
	"github.com/bufferserver/service"
)

func BlockCenterKeeper(db *gorm.DB, node *service.Node, duration time.Duration) {
	ticker := time.NewTicker(duration)
	for ; true; <-ticker.C {
		if err := syncBlockCenter(db, node); err != nil {
			log.WithField("err", err).Errorf("fail on blockcenter")
		}
	}
}

func syncBlockCenter(db *gorm.DB, node *service.Node) error {
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

		if err := UpdateOrSaveUTXO(db, base.AssetID, base.ControlProgram, resUTXOs); err != nil {
			return err
		}
	}
	return nil
}

func UpdateOrSaveUTXO(db *gorm.DB, asset string, program string, bcUTXOs []*service.AttachUtxo) error {
	utxoMap := make(map[string]bool)
	for _, butxo := range bcUTXOs {
		utxo := orm.Utxo{Hash: butxo.Hash}
		utxoMap[butxo.Hash] = true
		if err := db.Where(utxo).First(&utxo).Error; err != nil && err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, "query utxo")
		} else if err == gorm.ErrRecordNotFound {
			utxo := &orm.Utxo{
				Hash:           butxo.Hash,
				AssetID:        butxo.Asset,
				Amount:         butxo.Amount,
				ControlProgram: program,
				IsSpend:        false,
				IsLocked:       false,
				Duration:       uint64(60),
			}

			if err := db.Save(utxo).Error; err != nil {
				return errors.Wrap(err, "save utxo")
			}
			continue
		}

		if time.Now().Unix()-utxo.SubmitTime.Unix() < int64(utxo.Duration) {
			continue
		}

		if err := db.Model(&orm.Utxo{}).Where(&orm.Utxo{Hash: butxo.Hash}).Where("is_spend = false").Update("is_locked", false).Error; err != nil {
			return errors.Wrap(err, "update utxo unlocked")
		}
	}

	var utxos []*orm.Utxo
	if err := db.Model(&orm.Utxo{}).Where(&orm.Utxo{AssetID: asset, ControlProgram: program}).Where("is_spend = false").Find(&utxos).Error; err != nil {
		return errors.Wrap(err, "list unspent utxos")
	}

	for _, u := range utxos {
		if _, ok := utxoMap[u.Hash]; ok {
			continue
		}

		if err := db.Model(&orm.Utxo{}).Where(&orm.Utxo{Hash: u.Hash}).Update("is_spend", true).Error; err != nil {
			return errors.Wrap(err, "update utxo spent")
		}
	}
	return nil
}
