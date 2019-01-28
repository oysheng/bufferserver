package api

import (
	"time"

	"github.com/bytom/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/bufferserver/api/common"
	"github.com/bufferserver/database/orm"
)

func (s *Server) UpdateBase(c *gin.Context, req *common.AssetProgram) error {
	base := &orm.Base{AssetID: req.Asset, ControlProgram: req.Program}
	if err := s.db.Master().Where(base).First(base).Error; err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "db query base")
	} else if err == gorm.ErrRecordNotFound {
		if err := s.db.Master().Save(base).Error; err != nil {
			return errors.Wrap(err, "save base")
		}
	}

	return nil
}

type UpdateUTXOsReq struct {
	Hash string `json:"hash"`
}

func (s *Server) UpdateUtxo(c *gin.Context, req *UpdateUTXOsReq) error {
	utxo := &orm.Utxo{Hash: req.Hash}
	if err := s.db.Master().Where(utxo).First(utxo).Error; err != nil {
		return errors.Wrap(err, "db query utxo")
	}

	if err := s.db.Master().Model(&orm.Utxo{}).Where(&orm.Utxo{Hash: utxo.Hash}).Update("submit_time", time.Now()).Update("is_locked", true).Error; err != nil {
		return errors.Wrap(err, "update utxo locked")
	}
	return nil
}

type UpdateBalanceReq struct {
	Address string `json:"address"`
	AssetID string `json:"asset"`
	Amount  int64  `json:"amount"`
}

func (s *Server) UpdateBalance(c *gin.Context, req *UpdateBalanceReq) error {
	balance := &orm.Balance{Address: req.Address, AssetID: req.AssetID, Amount: req.Amount}
	if err := s.db.Master().Save(balance).Error; err != nil {
		return errors.Wrap(err, "insert balance record")
	}

	return nil
}
