package orm

import "github.com/bufferserver/types"

type Balance struct {
	ID          uint64          `json:"-" gorm:"primary_key"`
	Address     string          `json:"address"`
	AssetID     string          `json:"asset"`
	Amount      int64           `json:"amount"`
	TxID        string          `json:"-"`
	StatusFail  bool            `json:"-"`
	IsConfirmed bool            `json:"-"`
	CreatedAt   types.Timestamp `json:"create_at"`
}
