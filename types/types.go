package types

import (
	"fmt"
	"strconv"
	"time"
)

type Input struct {
	Script  string `json:"script"`
	Address string `json:"address"`
	Asset   string `json:"asset"`
	Amount  uint64 `json:"amount"`
}

type Output struct {
	UtxoID  string `json:"utxo_id"`
	Script  string `json:"script"`
	Address string `json:"address"`
	Asset   string `json:"asset"`
	Amount  uint64 `json:"amount"`
}

type Tx struct {
	Hash                string     `json:"hash"`
	StatusFail          bool       `json:"status_fail"`
	Size                uint64     `json:"size"`
	SubmissionTimestamp uint64     `json:"submission_timestamp"`
	BlockHeight         uint64     `json:"block_height,omitempty"`
	BlockTimestamp      uint64     `json:"block_timestamp,omitempty"`
	Memo                string     `json:"memo"`
	Inputs              []*Input   `json:"inputs"`
	Outputs             []*Output  `json:"outputs"`
	Fee                 uint64     `json:"fee"`
	Balances            []*Balance `json:"balances"`
}

type Balance struct {
	Asset  string `json:"asset"`
	Amount string `json:"amount"`
}

type Timestamp time.Time

func (t *Timestamp) Unix() int64 {
	return time.Time(*t).Unix()
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t).Unix()
	stamp := fmt.Sprint(ts)
	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	*t = Timestamp(time.Unix(int64(ts), 0))
	return nil
}
