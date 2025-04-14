package tradingstore

type InstrumentQueryInterface interface {
	// Validation
	Validate() error

	// Columns
	Columns() []string

	// Options
	SetCountOnly(countOnly bool) InstrumentQueryInterface
	IsCountOnly() bool

	// Limit, Offset
	SetLimit(limit int) InstrumentQueryInterface
	HasLimit() bool
	Limit() int
	SetOffset(offset int) InstrumentQueryInterface
	HasOffset() bool
	Offset() int

	// Order By
	SetOrderBy(orderBy string) InstrumentQueryInterface
	HasOrderBy() bool
	OrderBy() string
	SetSortDirection(sortDirection string) InstrumentQueryInterface
	HasSortDirection() bool
	SortDirection() string

	// ID
	SetID(id string) InstrumentQueryInterface
	HasID() bool
	ID() string
	SetIDIn(ids []string) InstrumentQueryInterface
	HasIDIn() bool
	IDIn() []string

	// Symbol
	SetSymbol(symbol string) InstrumentQueryInterface
	HasSymbol() bool
	Symbol() string
	SetSymbolLike(symbolLike string) InstrumentQueryInterface
	HasSymbolLike() bool
	SymbolLike() string

	// Exchange
	SetExchange(exchange string) InstrumentQueryInterface
	HasExchange() bool
	Exchange() string

	// Asset Class
	SetAssetClass(assetClass string) InstrumentQueryInterface
	HasAssetClass() bool
	AssetClass() string
}
