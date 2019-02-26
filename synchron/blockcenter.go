package synchron

import (
	"time"

	"github.com/bytom/errors"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	"github.com/bufferserver/api/common"
	"github.com/bufferserver/config"
	"github.com/bufferserver/database/orm"
	"github.com/bufferserver/service"
)

type blockCenterKeeper struct {
	cfg     *config.Config
	db      *gorm.DB
	service *service.Service
}

func NewBlockCenterKeeper(cfg *config.Config, db *gorm.DB) *blockCenterKeeper {
	service := service.NewService(cfg.Updater.BlockCenter.URL)
	return &blockCenterKeeper{
		cfg:     cfg,
		db:      db,
		service: service,
	}
}

func (b *blockCenterKeeper) Run() {
	ticker := time.NewTicker(time.Duration(b.cfg.Updater.BlockCenter.SyncSeconds) * time.Second)
	for ; true; <-ticker.C {
		if err := b.syncBlockCenter(); err != nil {
			log.WithField("err", err).Errorf("fail on bytom blockcenter")
		}
	}
}

func (b *blockCenterKeeper) syncBlockCenter() error {
	var bases []*orm.Base
	if err := b.db.Find(&bases).Error; err != nil {
		return errors.Wrap(err, "query bases")
	}

	filter := make(map[string]interface{})
	for _, base := range bases {
		filter["asset"] = base.AssetID
		filter["script"] = base.ControlProgram
		filter["unconfirmed"] = true
		req := &common.Display{Filter: filter}
		resUTXOs, err := b.service.ListBlockCenterUTXOs(req)
		if err != nil {
			return errors.Wrap(err, "list blockcenter utxos")
		}

		if err := b.updateOrSaveUTXO(base.AssetID, base.ControlProgram, resUTXOs); err != nil {
			return err
		}

		if err := b.updateUTXOStatus(base.AssetID, base.ControlProgram, resUTXOs); err != nil {
			return err
		}
	}

	if err := b.delIrrelevantUTXO(); err != nil {
		return err
	}

	return nil
}

func (b *blockCenterKeeper) updateOrSaveUTXO(asset string, program string, bcUTXOs []*service.AttachUtxo) error {
	for _, butxo := range bcUTXOs {
		utxo := orm.Utxo{Hash: butxo.Hash}
		if err := b.db.Where(utxo).First(&utxo).Error; err != nil && err != gorm.ErrRecordNotFound {
			return errors.Wrap(err, "query utxo")
		} else if err == gorm.ErrRecordNotFound {
			utxo := &orm.Utxo{
				Hash:           butxo.Hash,
				AssetID:        butxo.Asset,
				Amount:         butxo.Amount,
				ControlProgram: program,
				IsSpend:        false,
				IsConfirmed:    butxo.IsConfirmed,
				IsLocked:       false,
				Duration:       uint64(600),
			}

			if err := b.db.Save(utxo).Error; err != nil {
				return errors.Wrap(err, "save utxo")
			}
			continue
		}

		if butxo.IsConfirmed {
			if err := b.db.Model(&orm.Utxo{}).Where(&orm.Utxo{Hash: butxo.Hash}).Where("is_confirmed = false").Update("is_confirmed", true).Error; err != nil {
				return errors.Wrap(err, "update utxo confirmed")
			}
		}

		if time.Now().Unix()-utxo.SubmitTime.Unix() < int64(utxo.Duration) {
			continue
		}

		if err := b.db.Model(&orm.Utxo{}).Where(&orm.Utxo{Hash: butxo.Hash}).Where("is_locked = true").Update("is_locked", false).Error; err != nil {
			return errors.Wrap(err, "update utxo unlocked")
		}
	}

	return nil
}

func (b *blockCenterKeeper) updateUTXOStatus(asset string, program string, bcUTXOs []*service.AttachUtxo) error {
	utxoMap := make(map[string]bool)
	for _, butxo := range bcUTXOs {
		utxoMap[butxo.Hash] = true
	}

	var utxos []*orm.Utxo
	if err := b.db.Model(&orm.Utxo{}).Where(&orm.Utxo{AssetID: asset, ControlProgram: program}).Where("is_spend = false").Find(&utxos).Error; err != nil {
		return errors.Wrap(err, "list unspent utxos")
	}

	for _, u := range utxos {
		if _, ok := utxoMap[u.Hash]; ok {
			continue
		}

		if err := b.db.Model(&orm.Utxo{}).Where(&orm.Utxo{Hash: u.Hash}).Update("is_spend", true).Error; err != nil {
			return errors.Wrap(err, "update utxo spent")
		}
	}

	return nil
}

func (b *blockCenterKeeper) delIrrelevantUTXO() error {
	var utxos []*orm.Utxo
	query := b.db.Joins("left join bases on (utxos.control_program = bases.control_program and utxos.asset_id = bases.asset_id)").Where("bases.id is null")
	if err := query.Find(&utxos).Error; err == gorm.ErrRecordNotFound {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "query utxo not in base")
	}

	for _, u := range utxos {
		if err := b.db.Delete(&orm.Utxo{}, "hash = ? ", u.Hash).Error; err != nil {
			return errors.Wrap(err, "delete irrelevant utxo")
		}
	}

	return nil
}
