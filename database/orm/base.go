package orm

type Base struct {
	ID             uint64 `json:"-" gorm:"primary_key"`
	AssetID        string
	ControlProgram string
}
