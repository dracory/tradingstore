package tradingstore

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/base/database"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// PriceCount returns the number of prices based on the given query options
func (store *Store) PriceCount(ctx context.Context, options PriceQueryInterface) (int64, error) {
	options.SetCountOnly(true)

	q, _, err := store.priceQuery(options)

	if err != nil {
		return -1, err
	}

	sqlStr, sqlParams, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, nil
	}

	store.logSql("count", sqlStr, sqlParams...)

	mapped, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr, sqlParams...)

	if err != nil {
		return -1, err
	}

	if len(mapped) < 1 {
		return -1, nil
	}

	countStr := mapped[0]["count"]

	i, err := strconv.ParseInt(countStr, 10, 64)

	if err != nil {
		return -1, err

	}

	return i, nil
}

// PriceExists returns true if a price exists based on the given query options
func (store *Store) PriceExists(ctx context.Context, options PriceQueryInterface) (bool, error) {
	count, err := store.PriceCount(ctx, options)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// PriceCreate creates a new price
func (store *Store) PriceCreate(ctx context.Context, price PriceInterface) error {

	data := price.Data()

	data[COLUMN_TIME] = price.GetTimeCarbon().ToDateTimeString(carbon.UTC)

	sqlStr, sqlParams, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.priceTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("create", sqlStr, sqlParams...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, sqlParams...)

	if err != nil {
		return err
	}

	price.MarkAsNotDirty()

	return nil
}

// PriceDelete deletes a price
func (store *Store) PriceDelete(ctx context.Context, price PriceInterface) error {
	if price == nil {
		return errors.New("price is nil")
	}

	return store.PriceDeleteByID(ctx, price.ID())
}

// PriceDeleteByID deletes a price by its ID
func (store *Store) PriceDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("price id is empty")
	}

	sqlStr, sqlParams, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.priceTableName).
		Prepared(true).
		Where(goqu.C("id").Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("delete", sqlStr, sqlParams...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, sqlParams...)

	return err
}

// PriceFindByID returns a price by its ID
func (store *Store) PriceFindByID(ctx context.Context, id string) (PriceInterface, error) {
	if id == "" {
		return nil, errors.New("price id is empty")
	}

	query := NewPriceQuery().SetID(id).SetLimit(1)

	list, err := store.PriceList(ctx, query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// PriceList returns a list of prices based on the given query options
func (store *Store) PriceList(ctx context.Context, options PriceQueryInterface) ([]PriceInterface, error) {
	q, columns, err := store.priceQuery(options)

	if err != nil {
		return []PriceInterface{}, err
	}

	q = q.Prepared(true).Select(columns...)

	sqlStr, sqlParams, errSql := q.ToSQL()

	if errSql != nil {
		return []PriceInterface{}, errSql
	}

	store.logSql("list", sqlStr, sqlParams...)

	modelMaps, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr, sqlParams...)
	if err != nil {
		return []PriceInterface{}, err
	}

	list := []PriceInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewPriceFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) PriceUpdate(ctx context.Context, price PriceInterface) error {
	if price == nil {
		return errors.New("price is nil")
	}

	dataChanged := price.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, sqlParams, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.priceTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(price.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("update", sqlStr, sqlParams...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, sqlParams...)

	price.MarkAsNotDirty()

	return err
}

// priceQuery returns a query for prices based on the given query options
func (store *Store) priceQuery(options PriceQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, nil, errors.New("price options is nil")
	}

	if err := options.Validate(); err != nil {
		return nil, nil, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.priceTableName)

	if options.HasID() {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if options.HasIDIn() {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn()))
	}

	if options.HasTime() {
		q = q.Where(goqu.C(COLUMN_TIME).Eq(options.Time()))
	}

	if options.HasTimeGte() && options.HasTimeLte() {
		q = q.Where(
			goqu.C(COLUMN_TIME).Gte(options.TimeGte()),
			goqu.C(COLUMN_TIME).Lte(options.TimeLte()),
		)
	} else if options.HasTimeGte() {
		q = q.Where(goqu.C(COLUMN_TIME).Gte(options.TimeGte()))
	} else if options.HasTimeLte() {
		q = q.Where(goqu.C(COLUMN_TIME).Lte(options.TimeLte()))
	}

	if !options.IsCountOnly() {
		if options.HasLimit() {
			q = q.Limit(cast.ToUint(options.Limit()))
		}

		if options.HasOffset() {
			q = q.Offset(cast.ToUint(options.Offset()))
		}
	}

	if options.HasOrderBy() {
		sort := lo.Ternary(options.HasSortDirection(), options.SortDirection(), sb.DESC)
		if strings.EqualFold(sort, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy()).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy()).Desc())
		}
	}

	columns = []any{}

	for _, column := range options.Columns() {
		columns = append(columns, column)
	}

	return q, columns, nil
}
