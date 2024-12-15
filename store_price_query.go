package tradingstore

import "errors"

// type PriceQueryOptions struct {
// 	Columns   []string
// 	ID        string
// 	IDIn      string
// 	Time      string
// 	TimeGte   string
// 	TimeLte   string
// 	Offset    int
// 	Limit     int
// 	SortOrder string
// 	OrderBy   string
// 	CountOnly bool
// }

func NewPriceQuery() PriceQueryInterface {
	return &priceQueryImplementation{
		properties: make(map[string]any),
	}
}

type priceQueryImplementation struct {
	properties map[string]any
}

func (c *priceQueryImplementation) hasProperty(name string) bool {
	_, ok := c.properties[name]
	return ok
}

func (c *priceQueryImplementation) Validate() error {
	if c.HasID() && c.ID() == "" {
		return errors.New("price query. id cannot be empty")
	}

	if c.HasIDIn() && len(c.IDIn()) == 0 {
		return errors.New("price query. id_in cannot be empty")
	}

	if c.HasOrderBy() && c.OrderBy() == "" {
		return errors.New("price query. order_by cannot be empty")
	}

	if c.HasSortDirection() && c.SortDirection() == "" {
		return errors.New("price query. sort_direction cannot be empty")
	}

	if c.HasLimit() && c.Limit() <= 0 {
		return errors.New("price query. limit must be greater than 0")
	}

	if c.HasOffset() && c.Offset() < 0 {
		return errors.New("price query. offset must be greater than or equal to 0")
	}

	if c.HasTime() && c.Time() == "" {
		return errors.New("price query. time cannot be empty")
	}

	if c.HasTimeGte() && c.TimeGte() == "" {
		return errors.New("price query. time_gte cannot be empty")
	}

	if c.HasTimeLte() && c.TimeLte() == "" {
		return errors.New("price query. time_lte cannot be empty")
	}

	return nil
}

func (c *priceQueryImplementation) Columns() []string {
	if !c.hasProperty("columns") {
		return []string{}
	}

	return c.properties["columns"].([]string)
}

func (c *priceQueryImplementation) SetColumns(columns []string) PriceQueryInterface {
	c.properties["columns"] = columns

	return c
}

func (c *priceQueryImplementation) HasCountOnly() bool {
	return c.hasProperty("count_only")
}

func (c *priceQueryImplementation) IsCountOnly() bool {
	if !c.HasCountOnly() {
		return false
	}

	return c.properties["count_only"].(bool)
}

func (c *priceQueryImplementation) SetCountOnly(countOnly bool) PriceQueryInterface {
	c.properties["count_only"] = countOnly

	return c
}

func (c *priceQueryImplementation) HasID() bool {
	return c.hasProperty("id")
}

func (c *priceQueryImplementation) ID() string {
	if !c.hasProperty("id") {
		return ""
	}

	return c.properties["id"].(string)
}

func (c *priceQueryImplementation) SetID(id string) PriceQueryInterface {
	c.properties["id"] = id

	return c
}

func (c *priceQueryImplementation) HasIDIn() bool {
	return c.hasProperty("id_in")
}

func (c *priceQueryImplementation) IDIn() []string {
	if !c.hasProperty("id_in") {
		return []string{}
	}

	return c.properties["id_in"].([]string)
}

func (c *priceQueryImplementation) SetIDIn(idIn []string) PriceQueryInterface {
	c.properties["id_in"] = idIn

	return c
}

func (c *priceQueryImplementation) HasLimit() bool {
	return c.hasProperty("limit")
}

func (c *priceQueryImplementation) Limit() int {
	if !c.HasLimit() {
		return 0
	}

	return c.properties["limit"].(int)
}

func (c *priceQueryImplementation) SetLimit(limit int) PriceQueryInterface {
	c.properties["limit"] = limit

	return c
}

func (c *priceQueryImplementation) HasOffset() bool {
	return c.hasProperty("offset")
}

func (c *priceQueryImplementation) Offset() int {
	if !c.HasOffset() {
		return 0
	}

	return c.properties["offset"].(int)
}

func (c *priceQueryImplementation) SetOffset(offset int) PriceQueryInterface {
	c.properties["offset"] = offset

	return c
}

func (c *priceQueryImplementation) HasOrderBy() bool {
	return c.hasProperty("order_by")
}

func (c *priceQueryImplementation) OrderBy() string {
	if !c.HasOrderBy() {
		return ""
	}

	return c.properties["order_by"].(string)
}

func (c *priceQueryImplementation) SetOrderBy(orderBy string) PriceQueryInterface {
	c.properties["order_by"] = orderBy

	return c
}

func (c *priceQueryImplementation) HasSortDirection() bool {
	return c.hasProperty("sort_direction")
}

func (c *priceQueryImplementation) SortDirection() string {
	if !c.HasSortDirection() {
		return ""
	}

	return c.properties["sort_direction"].(string)
}

func (c *priceQueryImplementation) SetSortDirection(sortDirection string) PriceQueryInterface {
	c.properties["sort_direction"] = sortDirection

	return c
}

func (c *priceQueryImplementation) HasTime() bool {
	return c.hasProperty("time")
}

func (c *priceQueryImplementation) Time() string {
	if !c.HasTime() {
		return ""
	}

	return c.properties["time"].(string)
}

func (c *priceQueryImplementation) SetTime(time string) PriceQueryInterface {
	c.properties["time"] = time

	return c
}

func (c *priceQueryImplementation) HasTimeGte() bool {
	return c.hasProperty("time_gte")
}

func (c *priceQueryImplementation) TimeGte() string {
	if !c.HasTimeGte() {
		return ""
	}

	return c.properties["time_gte"].(string)
}

func (c *priceQueryImplementation) SetTimeGte(timeGte string) PriceQueryInterface {
	c.properties["time_gte"] = timeGte

	return c
}

func (c *priceQueryImplementation) HasTimeLte() bool {
	return c.hasProperty("time_lte")
}

func (c *priceQueryImplementation) TimeLte() string {
	if !c.HasTimeLte() {
		return ""
	}

	return c.properties["time_lte"].(string)
}

func (c *priceQueryImplementation) SetTimeLte(timeLte string) PriceQueryInterface {
	c.properties["time_lte"] = timeLte

	return c
}
