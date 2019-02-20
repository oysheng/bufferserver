package orm

import "github.com/bufferserver/types"

type Balance struct {
	ID          uint64          `json:"-" gorm:"primary_key"`
	Address     string          `json:"address"`
	AssetID     string          `json:"asset"`
	Amount      int64           `json:"amount"`
	TxID        string          `json:"tx_id"`
	StatusFail  bool            `json:"status_fail"`
	IsConfirmed bool            `json:"is_confirmed"`
	CreatedAt   types.Timestamp `json:"create_at"`
}
