package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func scatter(n int, fn func(i int) error) error {
	errors := make(chan error, n)

	var i int
	for i = 0; i < n; i++ {
		go func(i int) { errors <- fn(i) }(i)
	}

	var err, innerErr error
	for i = 0; i < cap(errors); i++ {
		if innerErr = <-errors; innerErr != nil {
			err = innerErr
		}
	}

	return err
}

// DB is a logical database with multiple underlying physical databases
// forming a single master multiple slaves topology.
// Reads and writes are automatically directed to the correct physical db.
type DB struct {
	pdb   *gorm.DB // Physical databases
	count uint64   // Monotonically incrementing counter on each query
}

// Open concurrently opens each underlying physical db.
func open(driverName, masterDSN string) (*DB, error) {
	conn := masterDSN
	db := &DB{}
	err := scatter(1, func(i int) (err error) {
		db.pdb, err = gorm.Open(driverName, conn)
		return err
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) logMode(logMode bool) error {
	return scatter(1, func(i int) error {
		db.pdb.LogMode(logMode)
		return nil
	})
}

// Ping verifies if a connection to each physical database is still alive,
// establishing a connection if necessary.
func (db *DB) Ping() error {
	return scatter(1, func(i int) error {
		return db.pdb.DB().Ping()
	})
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool for each underlying physical db.
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns then the
// new MaxIdleConns will be reduced to match the MaxOpenConns limit
// If n <= 0, no idle connections are retained.
func (db *DB) SetMaxIdleConns(n int) {
	db.pdb.DB().SetMaxIdleConns(n)
}

// SetMaxOpenConns sets the maximum number of open connections
// to each physical database.
// If MaxIdleConns is greater than 0 and the new MaxOpenConns
// is less than MaxIdleConns, then MaxIdleConns will be reduced to match
// the new MaxOpenConns limit. If n <= 0, then there is no limit on the number
// of open connections. The default is 0 (unlimited).
func (db *DB) SetMaxOpenConns(n int) {
	db.pdb.DB().SetMaxOpenConns(n)
}

// Master returns the master physical database
func (db *DB) Master() *gorm.DB {
	return db.pdb
}
