package tradingstore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/sb"
)

// NewStoreOptions define the options for creating a new tradingstore
type NewStoreOptions struct {
	// PriceTableNamePrefix is the prefix of the price table
	PriceTableNamePrefix string

	// InstrumentTableName is the name of the instrument table
	InstrumentTableName string

	// UseMultipleExchanges is used to create a new price table for each exchange
	// if false, the price table will be created without the exchange name as the table name (i.e. price_btc_usdt)
	// if true, the price table will be created with the exchange name as the table name (i.e. price_btc_binance_usdt)
	UseMultipleExchanges bool

	// DB is the underlying database connection
	DB *sql.DB

	// DbDriverName is the name of the database driver
	DbDriverName string

	// AutomigrateEnabled is used to auto migrate the instrument table
	// Note: You will need to call AutoMigratePrices after creating a new instrument
	AutomigrateEnabled bool

	// DebugEnabled is used to enable debug mode
	DebugEnabled bool
}

// NewStore creates a new trading store
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	if opts.PriceTableNamePrefix == "" {
		return nil, errors.New("trading store: PriceTableNamePrefix is required")
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
		priceTableNamePrefix: opts.PriceTableNamePrefix,
		instrumentTableName:  opts.InstrumentTableName,
		useMultipleExchanges: opts.UseMultipleExchanges,
		automigrateEnabled:   opts.AutomigrateEnabled,
		db:                   opts.DB,
		dbDriverName:         opts.DbDriverName,
		debugEnabled:         opts.DebugEnabled,
	}

	if store.automigrateEnabled {
		err := store.AutoMigrateInstruments(context.Background())

		if err != nil {
			return nil, err
		}
	}

	return store, nil
}
