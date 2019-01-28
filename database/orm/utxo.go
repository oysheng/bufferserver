package orm

import "time"

type Utxo struct {
	ID             uint64
	Hash           string
	AssetID        string
	Amount         uint64
	ControlProgram string
	IsSpend        bool
	IsLocked       bool
	SubmitTime     *time.Time
	Duration       uint64
}
