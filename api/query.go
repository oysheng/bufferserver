package api

import (
	"github.com/gin-gonic/gin"

	"github.com/bufferserver/api/common"
	"github.com/bufferserver/database/orm"
)

type ListBalanceReq struct {
	Address string `json:"address"`
	AssetID string `json:"asset"`
}

func (s *Server) ListBalances(c *gin.Context, req *ListBalanceReq) ([]*orm.Balance, error) {
	var balances []*orm.Balance
	balance := &orm.Balance{Address: req.Address, AssetID: req.AssetID}
	if err := s.db.Master().Model(&orm.Balance{}).Where(balance).Find(&balances).Error; err != nil {
		return nil, err
	}
	return balances, nil
}

type ListUTXOsResp struct {
	Hash   string `json:"hash"`
	Asset  string `json:"asset"`
	Amount uint64 `json:"amount"`
}

func (s *Server) ListUtxos(c *gin.Context, req *common.AssetProgram) ([]*ListUTXOsResp, error) {
	utxo := &orm.Utxo{AssetID: req.Asset, ControlProgram: req.Program}
	var utxos []*orm.Utxo
	if err := s.db.Master().Where(utxo).Where("is_spend = false").Where("is_locked = false").Find(&utxos).Error; err != nil {
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
