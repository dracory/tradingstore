package tradingstore

import (
	"database/sql"
	"errors"

	"github.com/gouniverse/sb"
)

// NewStoreOptions define the options for creating a new forexstore store
type NewStoreOptions struct {
	PriceTableName     string
	DB                 *sql.DB
	DbDriverName       string
	AutomigrateEnabled bool
	DebugEnabled       bool
}

// NewStore creates a new block store
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	if opts.PriceTableName == "" {
		return nil, errors.New("trade store: PriceTableName is required")
	}

	if opts.DB == nil {
		return nil, errors.New("trade store: DB is required")
	}

	if opts.DbDriverName == "" {
		opts.DbDriverName = sb.DatabaseDriverName(opts.DB)
	}

	store := &Store{
		priceTableName:     opts.PriceTableName,
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 opts.DB,
		dbDriverName:       opts.DbDriverName,
		debugEnabled:       opts.DebugEnabled,
	}

	if store.automigrateEnabled {
		err := store.AutoMigrate()

		if err != nil {
			return nil, err
		}
	}

	return store, nil
}
