package tradingstore

import (
	"context"
	"database/sql"
)

// StoreInterface defines the interface for a store
type StoreInterface interface {
	// AutoMigrate automatically creates the schema if it does not exist
	AutoMigrate() error

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

	// InstrumentUpdate updates an instrument
	InstrumentUpdate(ctx context.Context, instrument InstrumentInterface) error

	// PriceCount returns the number of prices that match the criteria
	PriceCount(ctx context.Context, options PriceQueryInterface) (int64, error)

	// PriceCreate creates a new price in the database
	PriceCreate(ctx context.Context, price PriceInterface) error

	// PriceDelete deletes a price
	PriceDelete(ctx context.Context, price PriceInterface) error

	// PriceDeleteByID deletes a price by ID
	PriceDeleteByID(ctx context.Context, id string) error

	// PriceExists checks if a price exists by checking a number of criteria
	PriceExists(ctx context.Context, options PriceQueryInterface) (bool, error)

	// PriceFindByID finds a price by its ID
	PriceFindByID(ctx context.Context, id string) (PriceInterface, error)

	// PriceList returns a list of prices from the database based on criteria
	PriceList(ctx context.Context, options PriceQueryInterface) ([]PriceInterface, error)

	// PriceUpdate updates a price
	PriceUpdate(ctx context.Context, price PriceInterface) error
}
