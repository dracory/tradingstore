package tradingstore

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
