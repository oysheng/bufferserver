package types

import (
	"fmt"
	"math/big"
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

func (tx *Tx) CalcBalances(addresses []string) {
	// create the map for each time, and thus no concurrent access.
	addressMap := make(map[string]bool)
	for _, address := range addresses {
		addressMap[address] = true
	}

	balanceMap := make(map[string]*big.Int)
	for _, input := range tx.Inputs {
		if _, ok := addressMap[input.Address]; !ok {
			continue
		}

		balance, ok := balanceMap[input.Asset]
		if !ok {
			balance = big.NewInt(0)
			balanceMap[input.Asset] = balance
		}

		balance.Sub(balance, new(big.Int).SetUint64(input.Amount))
	}

	for _, output := range tx.Outputs {
		if _, ok := addressMap[output.Address]; !ok {
			continue
		}

		balance, ok := balanceMap[output.Asset]
		if !ok {
			balance = big.NewInt(0)
			balanceMap[output.Asset] = balance
		}

		balance.Add(balance, new(big.Int).SetUint64(output.Amount))
	}

	for asset, amount := range balanceMap {
		balance := &Balance{
			Asset:  asset,
			Amount: amount.String(),
		}
		tx.Balances = append(tx.Balances, balance)
	}
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
