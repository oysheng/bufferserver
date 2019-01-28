package api

import (
	"github.com/bytom/errors"
	"github.com/gin-gonic/gin"

	"github.com/bufferserver/api/common"
	"github.com/bufferserver/database/orm"
)

type ListBalanceReq struct {
	Address string `json:"address"`
	AssetID string `json:"asset"`
}

type ListBalanceResp struct {
	Address string `json:"address"`
	AssetID string `json:"asset"`
	Amount  uint64 `json:"amount"`
}

func (s *Server) ListBalances(c *gin.Context, req *ListBalanceReq) (*ListBalanceResp, error) {
	balance := &orm.Balance{Address: req.Address, AssetID: req.AssetID}
	if err := s.db.Master().Where(balance).First(balance).Error; err != nil {
		return nil, errors.Wrap(err, "db query balance")
	}

	return &ListBalanceResp{
		Address: balance.Address,
		AssetID: balance.AssetID,
		Amount:  balance.Balance,
	}, nil
}

type ListUTXOsResp struct {
	Hash   string `json:"hash"`
	Asset  string `json:"asset"`
	Amount uint64 `json:"amount"`
}

func (s *Server) ListUtxos(c *gin.Context, req *common.AssetProgram) ([]*ListUTXOsResp, error) {
	utxo := &orm.Utxo{AssetID: req.Asset, ControlProgram: req.Program}
	var utxos []*orm.Utxo
	if err := s.db.Master().Where(utxo).Find(&utxos).Where("is_locked = false").Error; err != nil {
		return nil, err
	}

	var result []*ListUTXOsResp
	for _, u := range utxos {
		result = append(result, &ListUTXOsResp{
			Hash:   u.Hash,
			Asset:  u.AssetID,
			Amount: u.Amount,
		})
	}

	return result, nil
}
