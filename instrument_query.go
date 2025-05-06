package tradingstore

import (
	"errors"
	"strings"
)

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
	// assetClass
	isAssetClassSet bool
	assetClass      string

	// options
	countOnly bool

	// columns
	isColumnsSet bool
	columns      []string

	// exchange
	isExchangeSet bool
	exchange      string

	// id
	id      string
	isIDSet bool

	// id in
	idIn      []string
	isIDInSet bool

	// limit
	limit      int
	isLimitSet bool

	// offset
	offset      int
	isOffsetSet bool

	// orderBy
	orderBy      string
	isOrderBySet bool

	// orderDirection
	orderDirection      string
	isOrderDirectionSet bool

	// status
	isStatusSet bool
	status      string

	// symbol
	isSymbolSet bool
	symbol      string

	// symbolLike
	isSymbolLikeSet bool
	symbolLike      string
}

var _ InstrumentQueryInterface = (*instrumentQueryImplementation)(nil) // verify interface is implemented

// Validate validates the query options
func (q *instrumentQueryImplementation) Validate() error {
	if q.IsAssetClassSet() && q.AssetClass() == "" {
		return errors.New("instrument query. asset class cannot be empty")
	}

	if q.IsExchangeSet() && q.Exchange() == "" {
		return errors.New("instrument query. exchange cannot be empty")
	}

	if q.IsIDSet() && q.ID() == "" {
		return errors.New("instrument query. id cannot be empty")
	}

	if q.IsIDInSet() && len(q.IDIn()) == 0 {
		return errors.New("instrument query. id in cannot be empty")
	}

	if q.IsOffsetSet() && q.Offset() < 0 {
		return errors.New("instrument query. offset cannot be negative")
	}

	if q.IsLimitSet() && q.Limit() < 0 {
		return errors.New("instrument query. limit cannot be negative")
	}

	if q.IsOrderBySet() && q.OrderBy() == "" {
		return errors.New("instrument query. order by cannot be empty")
	}

	if q.IsOrderDirectionSet() && q.OrderDirection() != "asc" && q.OrderDirection() != "desc" {
		return errors.New("instrument query. order direction must be 'asc' or 'desc'")
	}

	if q.IsStatusSet() && q.Status() == "" {
		return errors.New("instrument query. status cannot be empty")
	}

	if q.IsSymbolSet() && q.Symbol() == "" {
		return errors.New("instrument query. symbol cannot be empty")
	}

	if q.IsSymbolLikeSet() && q.SymbolLike() == "" {
		return errors.New("instrument query. symbol like cannot be empty")
	}

	return nil
}

// IsAssetClassSet returns true if the asset class is set
func (iq *instrumentQueryImplementation) IsAssetClassSet() bool {
	return iq.isAssetClassSet
}

// AssetClass returns the asset class
func (iq *instrumentQueryImplementation) AssetClass() string {
	if iq.IsAssetClassSet() {
		return iq.assetClass
	}
	return ""
}

// SetAssetClass sets the asset class
func (iq *instrumentQueryImplementation) SetAssetClass(assetClass string) InstrumentQueryInterface {
	iq.assetClass = assetClass
	iq.isAssetClassSet = true
	return iq
}

// IsColumnsSet returns true if the columns are set
func (iq *instrumentQueryImplementation) IsColumnsSet() bool {
	return iq.isColumnsSet
}

// Columns returns the columns to select
func (iq *instrumentQueryImplementation) Columns() []string {
	if iq.IsColumnsSet() {
		return iq.columns
	}

	return []string{"*"}
}

// SetColumns sets the columns to select
func (iq *instrumentQueryImplementation) SetColumns(columns []string) InstrumentQueryInterface {
	iq.isColumnsSet = true
	iq.columns = columns
	return iq
}

// SetCountOnly sets the count only option
func (iq *instrumentQueryImplementation) SetCountOnly(countOnly bool) InstrumentQueryInterface {
	iq.countOnly = countOnly
	// not needed as bool is already set: iq.isCountOnlySet = true
	return iq
}

// IsCountOnly returns true if the count only option is set
func (iq *instrumentQueryImplementation) IsCountOnly() bool {
	return iq.countOnly
}

// IsExchangeSet returns true if the exchange is set
func (iq *instrumentQueryImplementation) IsExchangeSet() bool {
	return iq.isExchangeSet
}

// Exchange returns the exchange
func (iq *instrumentQueryImplementation) Exchange() string {
	if iq.IsExchangeSet() {
		return iq.exchange
	}
	return ""
}

// SetExchange sets the exchange
func (iq *instrumentQueryImplementation) SetExchange(exchange string) InstrumentQueryInterface {
	iq.exchange = exchange
	iq.isExchangeSet = true
	return iq
}

// IsIDSet returns true if the id is set
func (iq *instrumentQueryImplementation) IsIDSet() bool {
	return iq.isIDSet
}

// ID returns the id
func (iq *instrumentQueryImplementation) ID() string {
	if iq.IsIDSet() {
		return iq.id
	}
	return ""
}

// SetID sets the id
func (iq *instrumentQueryImplementation) SetID(id string) InstrumentQueryInterface {
	iq.id = id
	iq.isIDSet = true
	return iq
}

// IsIDInSet returns true if the id in is set
func (iq *instrumentQueryImplementation) IsIDInSet() bool {
	return len(iq.idIn) > 0
}

// IDIn returns the id in
func (iq *instrumentQueryImplementation) IDIn() []string {
	if iq.IsIDInSet() {
		return iq.idIn
	}
	return []string{}
}

// SetIDIn sets the id in
func (iq *instrumentQueryImplementation) SetIDIn(ids []string) InstrumentQueryInterface {
	iq.idIn = ids
	iq.isIDInSet = true
	return iq
}

// IsLimitSet returns true if the limit is set
func (iq *instrumentQueryImplementation) IsLimitSet() bool {
	return iq.isLimitSet
}

// Limit returns the limit
func (iq *instrumentQueryImplementation) Limit() int {
	if iq.IsLimitSet() {
		return iq.limit
	}
	return -1
}

// SetLimit sets the limit
func (iq *instrumentQueryImplementation) SetLimit(limit int) InstrumentQueryInterface {
	iq.limit = limit
	iq.isLimitSet = true
	return iq
}

// IsOffsetSet returns true if the offset is set
func (iq *instrumentQueryImplementation) IsOffsetSet() bool {
	return iq.isOffsetSet
}

// Offset returns the offset
func (iq *instrumentQueryImplementation) Offset() int {
	if iq.IsOffsetSet() {
		return iq.offset
	}
	return -1
}

// SetOffset sets the offset
func (iq *instrumentQueryImplementation) SetOffset(offset int) InstrumentQueryInterface {
	iq.offset = offset
	iq.isOffsetSet = true
	return iq
}

// IsOrderBySet returns true if the order by is set
func (iq *instrumentQueryImplementation) IsOrderBySet() bool {
	return iq.isOrderBySet
}

// OrderBy returns the order by
func (iq *instrumentQueryImplementation) OrderBy() string {
	if iq.IsOrderBySet() {
		return iq.orderBy
	}
	return ""
}

// SetOrderBy sets the order by
func (iq *instrumentQueryImplementation) SetOrderBy(orderBy string) InstrumentQueryInterface {
	iq.orderBy = orderBy
	iq.isOrderBySet = true
	return iq
}

// IsOrderDirectionSet returns true if the order direction is set
func (iq *instrumentQueryImplementation) IsOrderDirectionSet() bool {
	return iq.isOrderDirectionSet
}

// OrderDirection returns the order direction
func (iq *instrumentQueryImplementation) OrderDirection() string {
	if iq.IsOrderDirectionSet() {
		return iq.orderDirection
	}
	return ""
}

// SetOrderDirection sets the order direction
func (iq *instrumentQueryImplementation) SetOrderDirection(orderDirection string) InstrumentQueryInterface {
	iq.orderDirection = strings.ToLower(orderDirection)
	iq.isOrderDirectionSet = true
	return iq
}

// IsStatusSet returns true if the status is set
func (iq *instrumentQueryImplementation) IsStatusSet() bool {
	return iq.isStatusSet
}

// Status returns the status
func (iq *instrumentQueryImplementation) Status() string {
	if iq.IsStatusSet() {
		return iq.status
	}
	return ""
}

// SetStatus sets the status
func (iq *instrumentQueryImplementation) SetStatus(status string) InstrumentQueryInterface {
	iq.status = status
	iq.isStatusSet = true
	return iq
}

// IsSymbolSet returns true if the symbol is set
func (iq *instrumentQueryImplementation) IsSymbolSet() bool {
	return iq.isSymbolSet
}

// Symbol returns the symbol
func (iq *instrumentQueryImplementation) Symbol() string {
	if iq.IsSymbolSet() {
		return iq.symbol
	}
	return ""
}

// SetSymbol sets the symbol
func (iq *instrumentQueryImplementation) SetSymbol(symbol string) InstrumentQueryInterface {
	iq.symbol = symbol
	iq.isSymbolSet = true
	return iq
}

// IsSymbolLikeSet returns true if the symbol like is set
func (iq *instrumentQueryImplementation) IsSymbolLikeSet() bool {
	return iq.isSymbolLikeSet
}

// SymbolLike returns the symbol like
func (iq *instrumentQueryImplementation) SymbolLike() string {
	if iq.IsSymbolLikeSet() {
		return iq.symbolLike
	}
	return ""
}

// SetSymbolLike sets the symbol like
func (iq *instrumentQueryImplementation) SetSymbolLike(symbolLike string) InstrumentQueryInterface {
	iq.symbolLike = symbolLike
	iq.isSymbolLikeSet = true
	return iq
}
