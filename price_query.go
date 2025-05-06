package tradingstore

import "errors"

// PriceQuery is a shortcut for NewPriceQuery
func PriceQuery() PriceQueryInterface {
	return NewPriceQuery()
}

// NewPriceQuery creates a new price query
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
	if c.IsIDSet() && c.ID() == "" {
		return errors.New("price query. id cannot be empty")
	}

	if c.IsIDInSet() && len(c.IDIn()) == 0 {
		return errors.New("price query. id_in cannot be empty")
	}

	if c.IsOrderBySet() && c.OrderBy() == "" {
		return errors.New("price query. order_by cannot be empty")
	}

	if c.IsOrderDirectionSet() && c.OrderDirection() == "" {
		return errors.New("price query. order_direction cannot be empty")
	}

	if c.IsLimitSet() && c.Limit() <= 0 {
		return errors.New("price query. limit must be greater than 0")
	}

	if c.IsOffsetSet() && c.Offset() < 0 {
		return errors.New("price query. offset must be greater than or equal to 0")
	}

	if c.IsTimeSet() && c.Time() == "" {
		return errors.New("price query. time cannot be empty")
	}

	if c.IsTimeGteSet() && c.TimeGte() == "" {
		return errors.New("price query. time_gte cannot be empty")
	}

	if c.IsTimeLteSet() && c.TimeLte() == "" {
		return errors.New("price query. time_lte cannot be empty")
	}

	return nil
}

func (c *priceQueryImplementation) IsColumnsSet() bool {
	return c.hasProperty("columns")
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

func (c *priceQueryImplementation) IsCountOnlySet() bool {
	return c.hasProperty("count_only")
}

func (c *priceQueryImplementation) IsCountOnly() bool {
	if !c.IsCountOnlySet() {
		return false
	}

	return c.properties["count_only"].(bool)
}

func (c *priceQueryImplementation) SetCountOnly(countOnly bool) PriceQueryInterface {
	c.properties["count_only"] = countOnly

	return c
}

func (c *priceQueryImplementation) IsIDSet() bool {
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

func (c *priceQueryImplementation) IsIDInSet() bool {
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

func (c *priceQueryImplementation) IsLimitSet() bool {
	return c.hasProperty("limit")
}

func (c *priceQueryImplementation) Limit() int {
	if !c.IsLimitSet() {
		return 0
	}

	return c.properties["limit"].(int)
}

func (c *priceQueryImplementation) SetLimit(limit int) PriceQueryInterface {
	c.properties["limit"] = limit

	return c
}

func (c *priceQueryImplementation) IsOffsetSet() bool {
	return c.hasProperty("offset")
}

func (c *priceQueryImplementation) Offset() int {
	if !c.IsOffsetSet() {
		return 0
	}

	return c.properties["offset"].(int)
}

func (c *priceQueryImplementation) SetOffset(offset int) PriceQueryInterface {
	c.properties["offset"] = offset

	return c
}

func (c *priceQueryImplementation) IsOrderBySet() bool {
	return c.hasProperty("order_by")
}

func (c *priceQueryImplementation) OrderBy() string {
	if !c.IsOrderBySet() {
		return ""
	}

	return c.properties["order_by"].(string)
}

func (c *priceQueryImplementation) SetOrderBy(orderBy string) PriceQueryInterface {
	c.properties["order_by"] = orderBy

	return c
}

func (c *priceQueryImplementation) IsOrderDirectionSet() bool {
	return c.hasProperty("order_direction")
}

func (c *priceQueryImplementation) OrderDirection() string {
	if !c.IsOrderDirectionSet() {
		return ""
	}

	return c.properties["order_direction"].(string)
}

func (c *priceQueryImplementation) SetOrderDirection(orderDirection string) PriceQueryInterface {
	c.properties["order_direction"] = orderDirection

	return c
}

func (c *priceQueryImplementation) IsTimeSet() bool {
	return c.hasProperty("time")
}

func (c *priceQueryImplementation) Time() string {
	if !c.IsTimeSet() {
		return ""
	}

	return c.properties["time"].(string)
}

func (c *priceQueryImplementation) SetTime(time string) PriceQueryInterface {
	c.properties["time"] = time

	return c
}

func (c *priceQueryImplementation) IsTimeGteSet() bool {
	return c.hasProperty("time_gte")
}

func (c *priceQueryImplementation) TimeGte() string {
	if !c.IsTimeGteSet() {
		return ""
	}

	return c.properties["time_gte"].(string)
}

func (c *priceQueryImplementation) SetTimeGte(timeGte string) PriceQueryInterface {
	c.properties["time_gte"] = timeGte

	return c
}

func (c *priceQueryImplementation) IsTimeLteSet() bool {
	return c.hasProperty("time_lte")
}

func (c *priceQueryImplementation) TimeLte() string {
	if !c.IsTimeLteSet() {
		return ""
	}

	return c.properties["time_lte"].(string)
}

func (c *priceQueryImplementation) SetTimeLte(timeLte string) PriceQueryInterface {
	c.properties["time_lte"] = timeLte

	return c
}
