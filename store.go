package tradingstore

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/sb"
	"github.com/mingrammer/cfmt"
	"github.com/samber/lo"
)

// const DISCOUNT_TABLE_NAME = "shop_discount"

var _ StoreInterface = (*Store)(nil) // verify it extends the interface

type Store struct {
	priceTableName     string
	db                 *sql.DB
	dbDriverName       string
	automigrateEnabled bool
	debugEnabled       bool
}

// AutoMigrate auto migrate
func (store *Store) AutoMigrate() error {
	sql := store.sqlTablePriceCreate()

	_, err := store.db.Exec(sql)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// EnableDebug - enables the debug option
func (st *Store) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

func (st *Store) PriceCount(options PriceQueryOptions) (int64, error) {
	options.CountOnly = true
	q := st.priceQuery(options)

	sqlStr, _, errSql := q.Limit(1).Select(goqu.COUNT(goqu.Star()).As("count")).ToSQL()

	if errSql != nil {
		return -1, nil
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	db := sb.NewDatabase(st.db, st.dbDriverName)
	mapped, err := db.SelectToMapString(sqlStr)
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

func (store *Store) PriceExists(options PriceQueryOptions) (bool, error) {
	count, err := store.PriceCount(options)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (store *Store) PriceCreate(price *Price) error {

	data := price.Data()

	data["time"] = price.TimeCarbon().ToDateTimeString(carbon.UTC)

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.priceTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		cfmt.Infoln(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, params...)

	if err != nil {
		return err
	}

	price.MarkAsNotDirty()

	return nil
}

func (store *Store) PriceDelete(price *Price) error {
	if price == nil {
		return errors.New("price is nil")
	}

	return store.PriceDeleteByID(price.ID())
}

func (store *Store) PriceDeleteByID(id string) error {
	if id == "" {
		return errors.New("price id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.priceTableName).
		Prepared(true).
		Where(goqu.C("id").Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, params...)

	return err
}

func (store *Store) PriceFindByID(id string) (*Price, error) {
	if id == "" {
		return nil, errors.New("price id is empty")
	}

	list, err := store.PriceList(PriceQueryOptions{
		ID:    id,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return &list[0], nil
	}

	return nil, nil
}

func (store *Store) PriceList(options PriceQueryOptions) ([]Price, error) {
	q := store.priceQuery(options)

	if len(options.Columns) > 0 {
		q = q.Select(options.Columns[0])
		if len(options.Columns) > 1 {
			for _, column := range options.Columns[1:] {
				q = q.SelectAppend(goqu.C(column))
			}
		}
	} else {
		q = q.Select(goqu.Star())
	}

	sqlStr, _, errSql := q.ToSQL()

	if errSql != nil {
		return []Price{}, errSql
	}

	if store.debugEnabled {
		cfmt.Infoln(sqlStr)
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)
	modelMaps, err := db.SelectToMapString(sqlStr)
	if err != nil {
		return []Price{}, err
	}

	list := []Price{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewPriceFromExistingData(modelMap)
		list = append(list, *model)
	})

	return list, nil
}

func (store *Store) PriceUpdate(price *Price) error {
	if price == nil {
		return errors.New("price is nil")
	}

	// price.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := price.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.priceTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(price.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		cfmt.Infoln(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, params...)

	price.MarkAsNotDirty()

	return err
}

func (store *Store) priceQuery(options PriceQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).From(store.priceTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if options.Time != "" {
		q = q.Where(goqu.C(COLUMN_TIME).Eq(options.Time))
	}

	if options.TimeGte != "" {
		q = q.Where(goqu.C(COLUMN_TIME).Gte(options.TimeGte))
	}

	if options.TimeLte != "" {
		q = q.Where(goqu.C(COLUMN_TIME).Lte(options.TimeLte))
	}

	if !options.CountOnly {
		if options.Limit > 0 {
			q = q.Limit(uint(options.Limit))
		}

		if options.Offset > 0 {
			q = q.Offset(uint(options.Offset))
		}
	}

	sortOrder := "desc"
	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.SortOrder != "" {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy).Desc())
		}
	}

	return q
}

type PriceQueryOptions struct {
	Columns   []string
	ID        string
	IDIn      string
	Time      string
	TimeGte   string
	TimeLte   string
	Offset    int
	Limit     int
	SortOrder string
	OrderBy   string
	CountOnly bool
}
