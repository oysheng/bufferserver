package orm

type Balance struct {
	ID      uint64 `json:"-" gorm:"primary_key"`
	Address string
	AssetID string
	Balance uint64
}
