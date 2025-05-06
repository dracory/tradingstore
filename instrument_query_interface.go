package tradingstore

type InstrumentQueryInterface interface {
	// Validation
	Validate() error

	// Asset Class
	SetAssetClass(assetClass string) InstrumentQueryInterface
	IsAssetClassSet() bool
	AssetClass() string

	// Exchange
	IsExchangeSet() bool
	Exchange() string
	SetExchange(exchange string) InstrumentQueryInterface

	// Columns
	IsColumnsSet() bool
	Columns() []string
	SetColumns(columns []string) InstrumentQueryInterface

	// Count Only
	SetCountOnly(countOnly bool) InstrumentQueryInterface
	IsCountOnly() bool

	// ID
	SetID(id string) InstrumentQueryInterface
	IsIDSet() bool
	ID() string

	SetIDIn(ids []string) InstrumentQueryInterface
	IsIDInSet() bool
	IDIn() []string

	// Limit
	IsLimitSet() bool
	Limit() int
	SetLimit(limit int) InstrumentQueryInterface

	// Offset
	IsOffsetSet() bool
	Offset() int
	SetOffset(offset int) InstrumentQueryInterface

	// Order By
	IsOrderBySet() bool
	OrderBy() string
	SetOrderBy(orderBy string) InstrumentQueryInterface

	// Order Direction
	IsOrderDirectionSet() bool
	OrderDirection() string
	SetOrderDirection(orderDirection string) InstrumentQueryInterface

	// Status
	SetStatus(status string) InstrumentQueryInterface
	IsStatusSet() bool
	Status() string

	// Symbol
	IsSymbolSet() bool
	Symbol() string
	SetSymbol(symbol string) InstrumentQueryInterface

	// Symbol Like
	IsSymbolLikeSet() bool
	SymbolLike() string
	SetSymbolLike(symbolLike string) InstrumentQueryInterface
}
