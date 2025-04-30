package tradingstore

import (
	"context"
	"database/sql"
)

// StoreInterface defines the interface for a store
type StoreInterface interface {
	// AutoMigrateInstruments automatically creates the schema if it does not exist
	AutoMigrateInstruments(ctx context.Context) error

	// AutoMigratePrices automatically creates the price tables if they do not exist
	// It will create a price table for each instrument and each timeframe
	// You will need to call this method when you create a new instrument
	AutoMigratePrices(ctx context.Context) error

	// DB returns the underlying sql.DB connection
	DB() *sql.DB

	// EnableDebug enables debug mode
	EnableDebug(bool)

	// InstrumentCount returns the number of instruments that match the criteria
	InstrumentCount(ctx context.Context, options InstrumentQueryInterface) (int64, error)

	// InstrumentCreate creates a new instrument in the database
	InstrumentCreate(ctx context.Context, instrument InstrumentInterface) error

	// InstrumentDelete deletes an instrument
	InstrumentDelete(ctx context.Context, instrument InstrumentInterface) error

	// InstrumentDeleteByID deletes an instrument by ID
	InstrumentDeleteByID(ctx context.Context, id string) error

	// InstrumentExists checks if an instrument exists by checking a number of criteria
	InstrumentExists(ctx context.Context, options InstrumentQueryInterface) (bool, error)

	// InstrumentFindByID finds an instrument by its ID
	InstrumentFindByID(ctx context.Context, id string) (InstrumentInterface, error)

	// InstrumentList returns a list of instruments from the database based on criteria
	InstrumentList(ctx context.Context, options InstrumentQueryInterface) ([]InstrumentInterface, error)

	// InstrumentSoftDelete soft deletes an instrument
	InstrumentSoftDelete(ctx context.Context, instrument InstrumentInterface) error

	// InstrumentSoftDeleteByID soft deletes an instrument by ID
	InstrumentSoftDeleteByID(ctx context.Context, id string) error

	// InstrumentUpdate updates an instrument
	InstrumentUpdate(ctx context.Context, instrument InstrumentInterface) error

	// PriceCount returns the number of prices that match the criteria
	PriceCount(ctx context.Context, symbol string, exchange string, timeframe string, options PriceQueryInterface) (int64, error)

	// PriceCreate creates a new price in the database
	PriceCreate(ctx context.Context, symbol string, exchange string, timeframe string, price PriceInterface) error

	// PriceDelete deletes a price
	PriceDelete(ctx context.Context, symbol string, exchange string, timeframe string, price PriceInterface) error

	// PriceDeleteByID deletes a price by ID
	PriceDeleteByID(ctx context.Context, symbol string, exchange string, timeframe string, priceID string) error

	// PriceExists checks if a price exists by checking a number of criteria
	PriceExists(ctx context.Context, symbol string, exchange string, timeframe string, options PriceQueryInterface) (bool, error)

	// PriceFindByID finds a price by its ID
	PriceFindByID(ctx context.Context, symbol string, exchange string, timeframe string, priceID string) (PriceInterface, error)

	// PriceList returns a list of prices from the database based on criteria
	PriceList(ctx context.Context, symbol string, exchange string, timeframe string, options PriceQueryInterface) ([]PriceInterface, error)

	// PriceUpdate updates a price
	PriceUpdate(ctx context.Context, symbol string, exchange string, timeframe string, price PriceInterface) error
}
