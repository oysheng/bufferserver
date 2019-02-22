package synchron

import (
	"time"

	"github.com/bytom/errors"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	"github.com/bufferserver/config"
	"github.com/bufferserver/database/orm"
	"github.com/bufferserver/service"
)

type browserKeeper struct {
	cfg     *config.Config
	db      *gorm.DB
	service *service.Service
}

func NewBrowserKeeper(cfg *config.Config, db *gorm.DB) *browserKeeper {
	service := service.NewService(cfg.Updater.BlockCenter.URL)
	return &browserKeeper{
		cfg:     cfg,
		db:      db,
		service: service,
	}
}

func (b *browserKeeper) Run() {
	ticker := time.NewTicker(time.Duration(b.cfg.Updater.Browser.SyncSeconds) * time.Second)
	for ; true; <-ticker.C {
		if err := b.syncBrowser(); err != nil {
			log.WithField("err", err).Errorf("fail on bytom browser")
		}
	}
}

func (b *browserKeeper) syncBrowser() error {
	var balances []*orm.Balance
	if err := b.db.Model(&orm.Balance{}).Where("status_fail = false").Where("is_confirmed = false").Find(&balances).Error; err != nil {
		return errors.Wrap(err, "query balances")
	}

	expireTime := time.Duration(b.cfg.Updater.Browser.ExpirationHours) * time.Hour
	for _, balance := range balances {
		if balance.TxID == "" {
			if err := b.db.Delete(&orm.Balance{ID: balance.ID}).Error; err != nil {
				return errors.Wrap(err, "delete without TxID balance record")
			}
			continue
		}

		res, err := b.service.GetTransactionStatus(balance.TxID)
		if err != nil {
			log.WithField("err", err).Errorf("fail on query transaction [%s] from bytom browser", balance.TxID)
			continue
		}

		if res.Height == 0 {
			if time.Now().Unix()-balance.CreatedAt.Unix() > int64(expireTime) {
				if err := b.db.Delete(&orm.Balance{ID: balance.ID}).Error; err != nil {
					return errors.Wrap(err, "delete expiration balance record")
				}
			}
			continue
		}

		if err := b.db.Model(&orm.Balance{}).Where(&orm.Balance{ID: balance.ID}).Update("status_fail", res.StatusFail).Update("is_confirmed", true).Error; err != nil {
			return errors.Wrap(err, "update balance")
		}
	}
	return nil
}
