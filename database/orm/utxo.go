package orm

import (
	"github.com/bufferserver/types"
)

type Utxo struct {
	ID             uint64
	Hash           string
	AssetID        string
	Amount         uint64
	ControlProgram string
	IsSpend        bool
	IsLocked       bool
	SubmitTime     types.Timestamp
	Duration       uint64
}
