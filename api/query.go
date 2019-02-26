package api

import (
	"fmt"
	"time"

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
	balance := &orm.Balance{Address: req.Address, AssetID: req.AssetID, IsConfirmed: true}
	if err := s.db.Master().Model(&orm.Balance{}).Where(balance).Where("status_fail = false").Find(&balances).Error; err != nil {
		return nil, err
	}
	return balances, nil
}

type ListUTXOsReq struct {
	common.AssetProgram
	Confirmed bool          `json:"confirmed"`
	Sorter    common.Sorter `json:"sort"`
}

type ListUTXOsResp struct {
	Hash   string `json:"hash"`
	Asset  string `json:"asset"`
	Amount uint64 `json:"amount"`
}

func (s *Server) ListUtxos(c *gin.Context, req *ListUTXOsReq, page *common.PaginationQuery) ([]*ListUTXOsResp, error) {
	utxo := &orm.Utxo{AssetID: req.Asset, ControlProgram: req.Program}
	var utxos []*orm.Utxo
	query := s.db.Master().Where(utxo).Where("is_spend = false").Where("is_locked = false")
	if req.Confirmed {
		query = query.Where("is_confirmed = true")
	}

	if req.Sorter.By == "amount" {
		query = query.Order(fmt.Sprintf("amount %s", req.Sorter.Order))
	}

	if err := query.Offset(page.Start).Limit(page.Limit).Find(&utxos).Error; err != nil {
		return nil, err
	}

	// list locked UTXOs with more than 60 seconds from submit time while the count of unlock UTXOs is 0
	if len(utxos) == 0 {
		var lockUTXOs []*orm.Utxo
		query := s.db.Master().Where(utxo).Where("is_spend = false").Where("is_locked = true")
		if req.Sorter.By == "amount" {
			query = query.Order(fmt.Sprintf("amount %s", req.Sorter.Order))
		}

		if err := query.Offset(page.Start).Limit(page.Limit).Find(&lockUTXOs).Error; err != nil {
			return nil, err
		}

		for _, u := range lockUTXOs {
			if time.Now().Unix()-u.SubmitTime.Unix() >= 60 {
				utxos = append(utxos, u)
			}
		}
	}

	result := []*ListUTXOsResp{}
	for _, u := range utxos {
		result = append(result, &ListUTXOsResp{
			Hash:   u.Hash,
			Asset:  u.AssetID,
			Amount: u.Amount,
		})
	}

	return result, nil
}
