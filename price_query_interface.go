package tradingstore

type PriceQueryInterface interface {
	Validate() error

	IsColumnsSet() bool
	Columns() []string
	SetColumns(columns []string) PriceQueryInterface

	IsCountOnlySet() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) PriceQueryInterface

	IsTimeSet() bool
	Time() string
	SetTime(createdAt string) PriceQueryInterface

	IsTimeGteSet() bool
	TimeGte() string
	SetTimeGte(createdAtGte string) PriceQueryInterface

	IsTimeLteSet() bool
	TimeLte() string
	SetTimeLte(createdAtLte string) PriceQueryInterface

	IsIDSet() bool
	ID() string
	SetID(id string) PriceQueryInterface

	IsIDInSet() bool
	IDIn() []string
	SetIDIn(idIn []string) PriceQueryInterface

	IsLimitSet() bool
	Limit() int
	SetLimit(limit int) PriceQueryInterface

	IsOffsetSet() bool
	Offset() int
	SetOffset(offset int) PriceQueryInterface

	IsOrderBySet() bool
	OrderBy() string
	SetOrderBy(orderBy string) PriceQueryInterface

	IsOrderDirectionSet() bool
	OrderDirection() string
	SetOrderDirection(orderDirection string) PriceQueryInterface

	hasProperty(name string) bool
}
