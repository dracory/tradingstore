package tradingstore

// InstrumentQuery is a shortcut to create a new instrument query
func InstrumentQuery() InstrumentQueryInterface {
	return NewInstrumentQuery()
}

// NewInstrumentQuery creates a new instrument query
func NewInstrumentQuery() InstrumentQueryInterface {
	return &instrumentQueryImplementation{}
}

// instrumentQueryImplementation implements the InstrumentQueryInterface
type instrumentQueryImplementation struct {
	// options
	countOnly bool

	// limit offset
	limit  int
	offset int

	// orderBy
	orderBy       string
	sortDirection string

	// id
	id   string
	idIn []string

	// symbol
	symbol     string
	symbolLike string

	// exchange
	exchange string

	// assetClass
	assetClass string
}

var _ InstrumentQueryInterface = (*instrumentQueryImplementation)(nil) // verify interface is implemented

// Validate validates the query options
func (iq *instrumentQueryImplementation) Validate() error {
	return nil
}

// Columns returns the columns to select
func (iq *instrumentQueryImplementation) Columns() []string {
	return []string{"*"}
}

// SetCountOnly sets the count only option
func (iq *instrumentQueryImplementation) SetCountOnly(countOnly bool) InstrumentQueryInterface {
	iq.countOnly = countOnly
	return iq
}

// IsCountOnly returns true if the count only option is set
func (iq *instrumentQueryImplementation) IsCountOnly() bool {
	return iq.countOnly
}

// SetLimit sets the limit
func (iq *instrumentQueryImplementation) SetLimit(limit int) InstrumentQueryInterface {
	iq.limit = limit
	return iq
}

// HasLimit returns true if the limit is set
func (iq *instrumentQueryImplementation) HasLimit() bool {
	return iq.limit > 0
}

// Limit returns the limit
func (iq *instrumentQueryImplementation) Limit() int {
	return iq.limit
}

// SetOffset sets the offset
func (iq *instrumentQueryImplementation) SetOffset(offset int) InstrumentQueryInterface {
	iq.offset = offset
	return iq
}

// HasOffset returns true if the offset is set
func (iq *instrumentQueryImplementation) HasOffset() bool {
	return iq.offset > 0
}

// Offset returns the offset
func (iq *instrumentQueryImplementation) Offset() int {
	return iq.offset
}

// SetOrderBy sets the order by
func (iq *instrumentQueryImplementation) SetOrderBy(orderBy string) InstrumentQueryInterface {
	iq.orderBy = orderBy
	return iq
}

// HasOrderBy returns true if the order by is set
func (iq *instrumentQueryImplementation) HasOrderBy() bool {
	return iq.orderBy != ""
}

// OrderBy returns the order by
func (iq *instrumentQueryImplementation) OrderBy() string {
	return iq.orderBy
}

// SetSortDirection sets the sort direction
func (iq *instrumentQueryImplementation) SetSortDirection(sortDirection string) InstrumentQueryInterface {
	iq.sortDirection = sortDirection
	return iq
}

// HasSortDirection returns true if the sort direction is set
func (iq *instrumentQueryImplementation) HasSortDirection() bool {
	return iq.sortDirection != ""
}

// SortDirection returns the sort direction
func (iq *instrumentQueryImplementation) SortDirection() string {
	return iq.sortDirection
}

// SetID sets the id
func (iq *instrumentQueryImplementation) SetID(id string) InstrumentQueryInterface {
	iq.id = id
	return iq
}

// HasID returns true if the id is set
func (iq *instrumentQueryImplementation) HasID() bool {
	return iq.id != ""
}

// ID returns the id
func (iq *instrumentQueryImplementation) ID() string {
	return iq.id
}

// SetIDIn sets the id in
func (iq *instrumentQueryImplementation) SetIDIn(ids []string) InstrumentQueryInterface {
	iq.idIn = ids
	return iq
}

// HasIDIn returns true if the id in is set
func (iq *instrumentQueryImplementation) HasIDIn() bool {
	return len(iq.idIn) > 0
}

// IDIn returns the id in
func (iq *instrumentQueryImplementation) IDIn() []string {
	return iq.idIn
}

// SetSymbol sets the symbol
func (iq *instrumentQueryImplementation) SetSymbol(symbol string) InstrumentQueryInterface {
	iq.symbol = symbol
	return iq
}

// HasSymbol returns true if the symbol is set
func (iq *instrumentQueryImplementation) HasSymbol() bool {
	return iq.symbol != ""
}

// Symbol returns the symbol
func (iq *instrumentQueryImplementation) Symbol() string {
	return iq.symbol
}

// SetSymbolLike sets the symbol like
func (iq *instrumentQueryImplementation) SetSymbolLike(symbolLike string) InstrumentQueryInterface {
	iq.symbolLike = symbolLike
	return iq
}

// HasSymbolLike returns true if the symbol like is set
func (iq *instrumentQueryImplementation) HasSymbolLike() bool {
	return iq.symbolLike != ""
}

// SymbolLike returns the symbol like
func (iq *instrumentQueryImplementation) SymbolLike() string {
	return iq.symbolLike
}

// SetExchange sets the exchange
func (iq *instrumentQueryImplementation) SetExchange(exchange string) InstrumentQueryInterface {
	iq.exchange = exchange
	return iq
}

// HasExchange returns true if the exchange is set
func (iq *instrumentQueryImplementation) HasExchange() bool {
	return iq.exchange != ""
}

// Exchange returns the exchange
func (iq *instrumentQueryImplementation) Exchange() string {
	return iq.exchange
}

// SetAssetClass sets the asset class
func (iq *instrumentQueryImplementation) SetAssetClass(assetClass string) InstrumentQueryInterface {
	iq.assetClass = assetClass
	return iq
}

// HasAssetClass returns true if the asset class is set
func (iq *instrumentQueryImplementation) HasAssetClass() bool {
	return iq.assetClass != ""
}

// AssetClass returns the asset class
func (iq *instrumentQueryImplementation) AssetClass() string {
	return iq.assetClass
}
