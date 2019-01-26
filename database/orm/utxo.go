package orm

import (
	"github.com/blockcenter/types"
	"time"
)

type Utxo struct {
	ID             uint64
	Hash           string
	AssetID        string
	Amount         uint64
	SourceID       string
	SourcePos      uint64
	ControlProgram string
	IsSpend        bool
	SubmitTime     types.Timestamp
	Duration       time.Duration
}
