package api

import (
	"github.com/bytom/errors"
	"github.com/gin-gonic/gin"

	"github.com/blockcenter/api/common"
	"github.com/blockcenter/database/orm"
)

type ListBalanceReq struct {
	Address string `json:"address"`
	AssetID string `json:"asset"`
}

func (s *Server) ListBalances(c *gin.Context, req *ListBalanceReq) (*orm.Balance, error) {
	balance := &orm.Balance{Address: req.Address, AssetID: req.AssetID}
	if err := s.db.Master().Where(balance).First(balance).Error; err != nil {
		return nil, errors.Wrap(err, "db query balance")
	}

	return balance, nil
}

type ListUTXOReq struct {
	common.Display
}

type AttachUtxo struct {
	Hash   string `json:"hash"`
	Asset  string `json:"asset"`
	Amount uint64 `json:"amount"`
}

func (s *Server) ListUtxos(c *gin.Context, req *ListUTXOReq, page *common.PaginationQuery) ([]*AttachUtxo, error) {
	utxo := &orm.Utxo{}
	if asset, err := req.GetFilterString("asset"); err == nil {
		utxo.AssetID = asset
	}

	if cp, err := req.GetFilterString("program"); err == nil {
		utxo.ControlProgram = cp
	}

	var utxos []*orm.Utxo
	if err := s.db.Master().Where(utxo).Find(&utxos).Error; err != nil {
		return nil, err
	}

	var result []*AttachUtxo
	for _, u := range utxos {
		result = append(result, &AttachUtxo{
			Hash:   u.Hash,
			Asset:  u.AssetID,
			Amount: u.Amount,
		})
	}

	return result, nil
}
