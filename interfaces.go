package tradingstore

import (
	"context"
	"database/sql"

	"github.com/dromara/carbon/v2"
)

type StoreInterface interface {
	// AutoMigrate auto migrates the database schema
	AutoMigrate() error

	// EnableDebug enables or disables the debug mode
	EnableDebug(debug bool)

	// DB returns the underlying database connection
	DB() *sql.DB

	// == Price Methods =======================================================//

	// PriceCount returns the number of prices based on the given query options
	PriceCreate(ctx context.Context, price PriceInterface) error

	// PriceCount returns the number of prices based on the given query options
	PriceCount(ctx context.Context, options PriceQueryInterface) (int64, error)

	// PriceExists returns true if a price exists based on the given query options
	PriceExists(ctx context.Context, options PriceQueryInterface) (bool, error)

	// PriceFind returns a price based on the given query options
	PriceDelete(ctx context.Context, price PriceInterface) error

	// PriceFind returns a price based on the given query options
	PriceDeleteByID(ctx context.Context, priceID string) error

	// PriceFindByID returns a price by its ID
	PriceFindByID(ctx context.Context, price string) (PriceInterface, error)

	// PriceList returns a list of prices based on the given query options
	PriceList(ctx context.Context, options PriceQueryInterface) ([]PriceInterface, error)

	// PriceUpdate updates a price
	PriceUpdate(ctx context.Context, price PriceInterface) error
}

type PriceInterface interface {
	// from dataobject

	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	// methods

	// setters and getters
	ID() string
	SetID(id string) PriceInterface

	GetClose() string
	GetCloseFloat() float64
	SetClose(close string) PriceInterface

	GetHigh() string
	GetHighFloat() float64
	SetHigh(high string) PriceInterface

	GetLow() string
	GetLowFloat() float64
	SetLow(low string) PriceInterface

	GetOpen() string
	GetOpenFloat() float64
	SetOpen(open string) PriceInterface

	GetTime() string
	GetTimeCarbon() carbon.Carbon
	SetTime(time string) PriceInterface

	GetVolume() string
	GetVolumeFloat() float64
	SetVolume(volume string) PriceInterface
}

type PriceQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) PriceQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) PriceQueryInterface

	HasTime() bool
	Time() string
	SetTime(createdAt string) PriceQueryInterface

	HasTimeGte() bool
	TimeGte() string
	SetTimeGte(createdAtGte string) PriceQueryInterface

	HasTimeLte() bool
	TimeLte() string
	SetTimeLte(createdAtLte string) PriceQueryInterface

	HasID() bool
	ID() string
	SetID(id string) PriceQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) PriceQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) PriceQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) PriceQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) PriceQueryInterface

	HasSortDirection() bool
	SortDirection() string
	SetSortDirection(sortDirection string) PriceQueryInterface

	hasProperty(name string) bool
}
