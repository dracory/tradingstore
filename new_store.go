package tradingstore

import (
	"database/sql"
	"errors"

	"github.com/gouniverse/sb"
)

// NewStoreOptions define the options for creating a new tradingstore
type NewStoreOptions struct {
	PriceTableName      string
	InstrumentTableName string

	// UseMultipleExchanges is used to create a new price table for each exchange
	// if false, the price table will be created with the exchange name as the table name (i.e. price_btc_usdt)
	// if true, the price table will be created with the exchange name as the table name (i.e. price_btc_binance_usdt)
	UseMultipleExchanges bool

	DB                 *sql.DB
	DbDriverName       string
	AutomigrateEnabled bool
	DebugEnabled       bool
}

// NewStore creates a new trading store
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	if opts.PriceTableName == "" {
		return nil, errors.New("trading store: PriceTableName is required")
	}

	if opts.InstrumentTableName == "" {
		return nil, errors.New("trading store: InstrumentTableName is required")
	}

	if opts.DB == nil {
		return nil, errors.New("trading store: DB is required")
	}

	if opts.DbDriverName == "" {
		opts.DbDriverName = sb.DatabaseDriverName(opts.DB)
	}

	store := &Store{
		priceTableName:      opts.PriceTableName,
		instrumentTableName: opts.InstrumentTableName,
		automigrateEnabled:  opts.AutomigrateEnabled,
		db:                  opts.DB,
		dbDriverName:        opts.DbDriverName,
		debugEnabled:        opts.DebugEnabled,
	}

	if store.automigrateEnabled {
		err := store.AutoMigrate()

		if err != nil {
			return nil, err
		}
	}

	return store, nil
}
