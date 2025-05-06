package tradingstore

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/base/database"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// InstrumentCount returns the number of instruments based on the given query options
func (store *Store) InstrumentCount(ctx context.Context, options InstrumentQueryInterface) (int64, error) {
	options.SetCountOnly(true)

	q, _, err := store.instrumentQuery(options)

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

// InstrumentExists returns true if an instrument exists based on the given query options
func (store *Store) InstrumentExists(ctx context.Context, options InstrumentQueryInterface) (bool, error) {
	count, err := store.InstrumentCount(ctx, options)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// InstrumentCreate creates a new instrument
func (store *Store) InstrumentCreate(ctx context.Context, instrument InstrumentInterface) error {
	data := instrument.Data()

	sqlStr, sqlParams, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.instrumentTableName).
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

	instrument.MarkAsNotDirty()

	return nil
}

// InstrumentDelete deletes an instrument
func (store *Store) InstrumentDelete(ctx context.Context, instrument InstrumentInterface) error {
	if instrument == nil {
		return errors.New("instrument is nil")
	}

	return store.InstrumentDeleteByID(ctx, instrument.ID())
}

// InstrumentDeleteByID deletes an instrument by its ID
func (store *Store) InstrumentDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("instrument id is empty")
	}

	sqlStr, sqlParams, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.instrumentTableName).
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

// InstrumentFindByID returns an instrument by its ID
func (store *Store) InstrumentFindByID(ctx context.Context, id string) (InstrumentInterface, error) {
	if id == "" {
		return nil, errors.New("instrument id is empty")
	}

	query := NewInstrumentQuery().SetID(id).SetLimit(1)

	list, err := store.InstrumentList(ctx, query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// InstrumentList returns a list of instruments based on the given query options
func (store *Store) InstrumentList(ctx context.Context, options InstrumentQueryInterface) ([]InstrumentInterface, error) {
	q, columns, err := store.instrumentQuery(options)

	if err != nil {
		return []InstrumentInterface{}, err
	}

	q = q.Prepared(true).Select(columns...)

	sqlStr, sqlParams, errSql := q.ToSQL()

	if errSql != nil {
		return []InstrumentInterface{}, errSql
	}

	store.logSql("list", sqlStr, sqlParams...)

	modelMaps, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr, sqlParams...)
	if err != nil {
		return []InstrumentInterface{}, err
	}

	list := []InstrumentInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewInstrumentFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

// InstrumentSoftDelete soft deletes an instrument
func (store *Store) InstrumentSoftDelete(ctx context.Context, instrument InstrumentInterface) error {
	if instrument == nil {
		return errors.New("instrument is nil")
	}

	return store.InstrumentSoftDeleteByID(ctx, instrument.ID())
}

// InstrumentSoftDeleteByID soft deletes an instrument by ID
func (store *Store) InstrumentSoftDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("instrument id is empty")
	}

	sqlStr, sqlParams, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.instrumentTableName).
		Prepared(true).
		Set(goqu.Record{COLUMN_SOFT_DELETED_AT: time.Now().Format(time.RFC3339)}).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("soft delete", sqlStr, sqlParams...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, sqlParams...)

	return err
}

// InstrumentUpdate updates an instrument
func (store *Store) InstrumentUpdate(ctx context.Context, instrument InstrumentInterface) error {
	if instrument == nil {
		return errors.New("instrument is nil")
	}

	dataChanged := instrument.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, sqlParams, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.instrumentTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(instrument.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("update", sqlStr, sqlParams...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, sqlParams...)

	instrument.MarkAsNotDirty()

	return err
}

// instrumentQuery returns a query for instruments based on the given query options
func (store *Store) instrumentQuery(options InstrumentQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, nil, errors.New("instrument options is nil")
	}

	if err := options.Validate(); err != nil {
		return nil, nil, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.instrumentTableName)

	if options.IsAssetClassSet() {
		q = q.Where(goqu.C(COLUMN_ASSET_CLASS).Eq(options.AssetClass()))
	}

	if options.IsExchangeSet() {
		q = q.Where(goqu.C(COLUMN_EXCHANGE).Eq(options.Exchange()))
	}

	if options.IsIDSet() {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if options.IsIDInSet() {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn()))
	}

	if options.IsStatusSet() {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status()))
	}

	if options.IsSymbolSet() {
		q = q.Where(goqu.C(COLUMN_SYMBOL).Eq(options.Symbol()))
	}

	if options.IsSymbolLikeSet() {
		q = q.Where(goqu.C(COLUMN_SYMBOL).Like("%" + options.SymbolLike() + "%"))
	}

	if !options.IsCountOnly() {
		if options.IsLimitSet() {
			q = q.Limit(cast.ToUint(options.Limit()))
		}

		if options.IsOffsetSet() {
			q = q.Offset(cast.ToUint(options.Offset()))
		}
	}

	if options.IsOrderBySet() {
		sort := lo.Ternary(options.IsOrderDirectionSet(), options.OrderDirection(), sb.DESC)
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
