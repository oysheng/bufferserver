package orm

type Balance struct {
	ID      uint64 `json:"-" gorm:"primary_key"`
	Address string `json:"address"`
	AssetID string `json:"asset"`
	Amount  int64  `json:"amount"`
}
