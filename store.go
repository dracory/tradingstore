package tradingstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/dracory/neat"
	contractsorm "github.com/dracory/neat/contracts/database/orm"
	contractsschema "github.com/dracory/neat/contracts/database/schema"
	"github.com/dromara/carbon/v2"
)

// StoreInterface defines the interface for a store
type StoreInterface interface {
	MigrateDown(ctx context.Context, tx ...*sql.Tx) error
	MigrateUp(ctx context.Context, tx ...*sql.Tx) error
	DB() *sql.DB
	EnableDebug(bool)

	InstrumentCount(ctx context.Context, options InstrumentQueryInterface) (int64, error)
	InstrumentCreate(ctx context.Context, instrument InstrumentInterface) error
	InstrumentDelete(ctx context.Context, instrument InstrumentInterface) error
	InstrumentDeleteByID(ctx context.Context, id string) error
	InstrumentExists(ctx context.Context, options InstrumentQueryInterface) (bool, error)
	InstrumentFindByID(ctx context.Context, id string) (InstrumentInterface, error)
	InstrumentList(ctx context.Context, options InstrumentQueryInterface) ([]InstrumentInterface, error)
	InstrumentSoftDelete(ctx context.Context, instrument InstrumentInterface) error
	InstrumentSoftDeleteByID(ctx context.Context, id string) error
	InstrumentUpdate(ctx context.Context, instrument InstrumentInterface) error

	PriceCount(ctx context.Context, symbol string, exchange string, timeframe string, options PriceQueryInterface) (int64, error)
	PriceCreate(ctx context.Context, symbol string, exchange string, timeframe string, price PriceInterface) error
	PriceDelete(ctx context.Context, symbol string, exchange string, timeframe string, price PriceInterface) error
	PriceDeleteByID(ctx context.Context, symbol string, exchange string, timeframe string, priceID string) error
	PriceExists(ctx context.Context, symbol string, exchange string, timeframe string, options PriceQueryInterface) (bool, error)
	PriceFindByID(ctx context.Context, symbol string, exchange string, timeframe string, priceID string) (PriceInterface, error)
	PriceList(ctx context.Context, symbol string, exchange string, timeframe string, options PriceQueryInterface) ([]PriceInterface, error)
	PriceUpdate(ctx context.Context, symbol string, exchange string, timeframe string, price PriceInterface) error
}

var _ StoreInterface = (*storeImplementation)(nil)

// == TYPE =====================================================================

type storeImplementation struct {
	priceTableNamePrefix string
	instrumentTableName  string
	useMultipleExchanges bool
	automigrateEnabled   bool
	debugEnabled         bool
	logger               *slog.Logger
	db                   *neat.Database
}

// == CONSTRUCTOR ==============================================================

// NewStoreOptions define the options for creating a new tradingstore
type NewStoreOptions struct {
	PriceTableNamePrefix string
	InstrumentTableName  string
	UseMultipleExchanges bool
	DB                   *sql.DB
	AutomigrateEnabled   bool
	DebugEnabled         bool
}

// NewStore creates a new trading store
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	if opts.PriceTableNamePrefix == "" {
		return nil, errors.New("trading store: PriceTableNamePrefix is required")
	}

	if opts.InstrumentTableName == "" {
		return nil, errors.New("trading store: InstrumentTableName is required")
	}

	if opts.DB == nil {
		return nil, errors.New("trading store: DB is required")
	}

	neatDB, err := neat.NewFromSQLDB(opts.DB)
	if err != nil {
		return nil, err
	}

	store := &storeImplementation{
		priceTableNamePrefix: opts.PriceTableNamePrefix,
		instrumentTableName:  opts.InstrumentTableName,
		useMultipleExchanges: opts.UseMultipleExchanges,
		automigrateEnabled:   opts.AutomigrateEnabled,
		debugEnabled:         opts.DebugEnabled,
		logger:               slog.New(slog.NewTextHandler(os.Stdout, nil)),
		db:                   neatDB,
	}

	if store.automigrateEnabled {
		if err := store.MigrateUp(context.Background()); err != nil {
			return nil, err
		}
	}

	return store, nil
}

// == PUBLIC METHODS ===========================================================

// DB returns the underlying database connection
func (store *storeImplementation) DB() *sql.DB {
	db, _ := store.db.DB()
	return db
}

// EnableDebug enables or disables the debug mode
func (store *storeImplementation) EnableDebug(debug bool) {
	store.debugEnabled = debug
	if debug {
		store.db.EnableDebug()
	} else {
		store.db.DisableDebug()
	}
}

// MigrateUp creates the trading store tables
func (store *storeImplementation) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	if !store.db.Schema().HasTable(store.instrumentTableName) {
		err := store.db.Schema().Create(store.instrumentTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, 40)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_NAME, 100)
			table.String(COLUMN_STATUS, 20)
			table.String(COLUMN_ASSET_CLASS, 40)
			table.String(COLUMN_SYMBOL, 10)
			table.String(COLUMN_EXCHANGE, 50)
			table.String(COLUMN_TIMEFRAMES, 100)
			table.Text(COLUMN_DESCRIPTION)
			table.Text(COLUMN_MEMO)
			table.LongText(COLUMN_METAS)
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
			table.DateTime(COLUMN_SOFT_DELETED_AT)
			table.Index(COLUMN_SYMBOL)
			table.Index(COLUMN_EXCHANGE)
			table.Index(COLUMN_ASSET_CLASS)
			table.Index(COLUMN_SOFT_DELETED_AT)
		})
		if err != nil {
			return err
		}
	}

	instruments, err := store.InstrumentList(ctx, InstrumentQuery().SetSoftDeletedIncluded(true))
	if err != nil {
		return err
	}

	for _, instrument := range instruments {
		timeframes := instrument.Timeframes()
		for _, timeframe := range timeframes {
			tableName := store.PriceTableName(instrument.Symbol(), instrument.Exchange(), timeframe)
			if !store.db.Schema().HasTable(tableName) {
				if err := store.createPriceTable(tableName); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// createPriceTable creates a single price table with the standard OHLCV schema,
// including an index and unique constraint on the time column.
func (store *storeImplementation) createPriceTable(tableName string) error {
	return store.db.Schema().Create(tableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 40)
		table.Primary(COLUMN_ID)
		table.Decimal(COLUMN_OPEN).Total(20).Places(8)
		table.Decimal(COLUMN_HIGH).Total(20).Places(8)
		table.Decimal(COLUMN_LOW).Total(20).Places(8)
		table.Decimal(COLUMN_CLOSE).Total(20).Places(8)
		table.BigInteger(COLUMN_VOLUME)
		table.DateTime(COLUMN_TIME)
		table.Index(COLUMN_TIME)
		table.Unique(COLUMN_TIME)
	})
}

// MigrateDown drops the trading store tables
func (store *storeImplementation) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	instruments, err := store.InstrumentList(ctx, InstrumentQuery().SetSoftDeletedIncluded(true))
	if err != nil {
		return err
	}

	for _, instrument := range instruments {
		timeframes := instrument.Timeframes()
		for _, timeframe := range timeframes {
			tableName := store.PriceTableName(instrument.Symbol(), instrument.Exchange(), timeframe)
			if store.db.Schema().HasTable(tableName) {
				if err := store.db.Schema().DropIfExists(tableName); err != nil {
					return err
				}
			}
		}
	}

	if store.db.Schema().HasTable(store.instrumentTableName) {
		if err := store.db.Schema().DropIfExists(store.instrumentTableName); err != nil {
			return err
		}
	}

	return nil
}

// PriceTableName returns the dynamic price table name for a symbol/exchange/timeframe
func (store *storeImplementation) PriceTableName(symbol string, exchange string, timeframe string) string {
	priceTableName := store.priceTableNamePrefix
	if exchange != "" && store.useMultipleExchanges {
		return priceTableName + strings.ToLower(symbol) + "_" + strings.ToLower(exchange) + "_" + strings.ToLower(timeframe)
	}
	return priceTableName + strings.ToLower(symbol) + "_" + strings.ToLower(timeframe)
}

// == INSTRUMENT METHODS =======================================================

// InstrumentCount returns the number of instruments based on the given query options
func (store *storeImplementation) InstrumentCount(ctx context.Context, options InstrumentQueryInterface) (int64, error) {
	if err := options.Validate(); err != nil {
		return 0, err
	}

	q := store.buildInstrumentQuery(options)

	var count int64
	err := q.Table(store.instrumentTableName).Count(&count)
	return count, err
}

// InstrumentExists returns true if an instrument exists based on the given query options
func (store *storeImplementation) InstrumentExists(ctx context.Context, options InstrumentQueryInterface) (bool, error) {
	count, err := store.InstrumentCount(ctx, options)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// InstrumentCreate creates a new instrument and automatically creates its
// price tables for each configured timeframe.
func (store *storeImplementation) InstrumentCreate(ctx context.Context, instrument InstrumentInterface) error {
	if instrument.CreatedAt() == "" {
		instrument.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}
	if instrument.UpdatedAt() == "" {
		instrument.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}
	if instrument.SoftDeletedAt() == "" {
		instrument.SetSoftDeletedAt(neat.MaxDateTime)
	}

	row, err := store.instrumentToMap(instrument)
	if err != nil {
		return err
	}

	if err := store.db.Query().Table(store.instrumentTableName).Create(row); err != nil {
		return err
	}

	// Auto-create price tables for each timeframe of the new instrument
	for _, timeframe := range instrument.Timeframes() {
		tableName := store.PriceTableName(instrument.Symbol(), instrument.Exchange(), timeframe)
		if !store.db.Schema().HasTable(tableName) {
			if err := store.createPriceTable(tableName); err != nil {
				return err
			}
		}
	}

	return nil
}

// InstrumentDelete deletes an instrument
func (store *storeImplementation) InstrumentDelete(ctx context.Context, instrument InstrumentInterface) error {
	if instrument == nil {
		return errors.New("instrument is nil")
	}
	return store.InstrumentDeleteByID(ctx, instrument.ID())
}

// InstrumentDeleteByID deletes an instrument by its ID
func (store *storeImplementation) InstrumentDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("instrument id is empty")
	}
	_, err := store.db.Query().
		Table(store.instrumentTableName).
		Where(COLUMN_ID+" = ?", id).
		Delete()
	return err
}

// InstrumentFindByID returns an instrument by its ID
func (store *storeImplementation) InstrumentFindByID(ctx context.Context, id string) (InstrumentInterface, error) {
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
func (store *storeImplementation) InstrumentList(ctx context.Context, options InstrumentQueryInterface) ([]InstrumentInterface, error) {
	if err := options.Validate(); err != nil {
		return []InstrumentInterface{}, err
	}

	q := store.buildInstrumentQuery(options)

	if options.IsColumnsSet() && len(options.Columns()) > 0 && options.Columns()[0] != "*" {
		q = q.Select(options.Columns())
	}

	type instrumentRow struct {
		ID            string    `db:"id"`
		Symbol        string    `db:"symbol"`
		Exchange      string    `db:"exchange"`
		AssetClass    string    `db:"asset_class"`
		Name          string    `db:"name"`
		Status        string    `db:"status"`
		Description   string    `db:"description"`
		Memo          string    `db:"memo"`
		Metas         string    `db:"metas"`
		Timeframes    string    `db:"timeframes"`
		CreatedAt     time.Time `db:"created_at"`
		UpdatedAt     time.Time `db:"updated_at"`
		SoftDeletedAt time.Time `db:"soft_deleted_at"`
	}

	var rows []instrumentRow
	if err := q.Table(store.instrumentTableName).Get(&rows); err != nil {
		return []InstrumentInterface{}, err
	}

	list := make([]InstrumentInterface, 0, len(rows))
	for _, r := range rows {
		instrument := &instrumentImplementation{}
		instrument.SetID(r.ID)
		instrument.SymbolField = r.Symbol
		instrument.ExchangeField = r.Exchange
		instrument.AssetClassField = r.AssetClass
		instrument.NameField = r.Name
		instrument.StatusField = r.Status
		instrument.DescriptionField = r.Description
		instrument.MemoField = r.Memo
		instrument.MetasField = r.Metas
		instrument.TimeframesField = r.Timeframes
		instrument.CreatedAtField.CreatedAt = r.CreatedAt
		instrument.UpdatedAtField.UpdatedAt = r.UpdatedAt
		instrument.SoftDeletesMaxDate.SoftDeletedAt = r.SoftDeletedAt
		list = append(list, instrument)
	}

	return list, nil
}

// InstrumentSoftDelete soft deletes an instrument
func (store *storeImplementation) InstrumentSoftDelete(ctx context.Context, instrument InstrumentInterface) error {
	if instrument == nil {
		return errors.New("instrument is nil")
	}
	return store.InstrumentSoftDeleteByID(ctx, instrument.ID())
}

// InstrumentSoftDeleteByID soft deletes an instrument by ID
func (store *storeImplementation) InstrumentSoftDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("instrument id is empty")
	}

	row := map[string]any{
		COLUMN_SOFT_DELETED_AT: carbon.Now(carbon.UTC).StdTime(),
		COLUMN_UPDATED_AT:      carbon.Now(carbon.UTC).StdTime(),
	}

	_, err := store.db.Query().
		Table(store.instrumentTableName).
		Where(COLUMN_ID+" = ?", id).
		Update(row)
	return err
}

// InstrumentUpdate updates an instrument
func (store *storeImplementation) InstrumentUpdate(ctx context.Context, instrument InstrumentInterface) error {
	if instrument == nil {
		return errors.New("instrument is nil")
	}

	instrument.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	row, err := store.instrumentToMap(instrument)
	if err != nil {
		return err
	}
	// ID and created_at are not updatable
	delete(row, COLUMN_ID)
	delete(row, COLUMN_CREATED_AT)

	_, err = store.db.Query().
		Table(store.instrumentTableName).
		Where(COLUMN_ID+" = ?", instrument.ID()).
		Update(row)
	return err
}

// == PRICE METHODS ============================================================

// PriceCount returns the number of prices based on the given query options
func (store *storeImplementation) PriceCount(ctx context.Context, symbol string, exchange string, timeframe string, options PriceQueryInterface) (int64, error) {
	if err := options.Validate(); err != nil {
		return 0, err
	}

	q := store.buildPriceQuery(symbol, exchange, timeframe, options)

	var count int64
	err := q.Table(store.PriceTableName(symbol, exchange, timeframe)).Count(&count)
	return count, err
}

// PriceExists returns true if a price exists based on the given query options
func (store *storeImplementation) PriceExists(ctx context.Context, symbol string, exchange string, timeframe string, options PriceQueryInterface) (bool, error) {
	count, err := store.PriceCount(ctx, symbol, exchange, timeframe, options)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// PriceCreate creates a new price
func (store *storeImplementation) PriceCreate(ctx context.Context, symbol string, exchange string, timeframe string, price PriceInterface) error {
	row := map[string]any{
		COLUMN_ID:     price.ID(),
		COLUMN_OPEN:   price.Open(),
		COLUMN_HIGH:   price.High(),
		COLUMN_LOW:    price.Low(),
		COLUMN_CLOSE:  price.Close(),
		COLUMN_VOLUME: price.Volume(),
		COLUMN_TIME:   price.TimeCarbon().ToDateTimeString(carbon.UTC),
	}

	return store.db.Query().Table(store.PriceTableName(symbol, exchange, timeframe)).Create(row)
}

// PriceDelete deletes a price
func (store *storeImplementation) PriceDelete(ctx context.Context, symbol string, exchange string, timeframe string, price PriceInterface) error {
	if price == nil {
		return errors.New("price is nil")
	}
	return store.PriceDeleteByID(ctx, symbol, exchange, timeframe, price.ID())
}

// PriceDeleteByID deletes a price by its ID
func (store *storeImplementation) PriceDeleteByID(ctx context.Context, symbol string, exchange string, timeframe string, id string) error {
	if id == "" {
		return errors.New("price id is empty")
	}
	_, err := store.db.Query().
		Table(store.PriceTableName(symbol, exchange, timeframe)).
		Where(COLUMN_ID+" = ?", id).
		Delete()
	return err
}

// PriceFindByID returns a price by its ID
func (store *storeImplementation) PriceFindByID(ctx context.Context, symbol string, exchange string, timeframe string, priceID string) (PriceInterface, error) {
	if priceID == "" {
		return nil, errors.New("price id is empty")
	}

	query := NewPriceQuery().SetID(priceID).SetLimit(1)
	list, err := store.PriceList(ctx, symbol, exchange, timeframe, query)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

// PriceList returns a list of prices based on the given query options
func (store *storeImplementation) PriceList(ctx context.Context, symbol string, exchange string, timeframe string, options PriceQueryInterface) ([]PriceInterface, error) {
	if err := options.Validate(); err != nil {
		return []PriceInterface{}, err
	}

	q := store.buildPriceQuery(symbol, exchange, timeframe, options)

	if options.IsColumnsSet() && len(options.Columns()) > 0 && options.Columns()[0] != "*" {
		q = q.Select(options.Columns())
	}

	type priceRow struct {
		ID     string    `db:"id"`
		Open   string    `db:"open"`
		High   string    `db:"high"`
		Low    string    `db:"low"`
		Close  string    `db:"close"`
		Volume string    `db:"volume"`
		Time   time.Time `db:"time"`
	}

	var rows []priceRow
	if err := q.Table(store.PriceTableName(symbol, exchange, timeframe)).Get(&rows); err != nil {
		return []PriceInterface{}, err
	}

	list := make([]PriceInterface, 0, len(rows))
	for _, r := range rows {
		price := &priceImplementation{}
		price.SetID(r.ID)
		price.OpenField = r.Open
		price.HighField = r.High
		price.LowField = r.Low
		price.CloseField = r.Close
		price.VolumeField = r.Volume
		price.TimeField = carbon.CreateFromStdTime(r.Time).ToIso8601ZuluString()
		list = append(list, price)
	}

	return list, nil
}

// PriceUpdate updates a price
func (store *storeImplementation) PriceUpdate(ctx context.Context, symbol string, exchange string, timeframe string, price PriceInterface) error {
	if price == nil {
		return errors.New("price is nil")
	}

	row := map[string]any{
		COLUMN_OPEN:   price.Open(),
		COLUMN_HIGH:   price.High(),
		COLUMN_LOW:    price.Low(),
		COLUMN_CLOSE:  price.Close(),
		COLUMN_VOLUME: price.Volume(),
		COLUMN_TIME:   price.TimeCarbon().ToDateTimeString(carbon.UTC),
	}

	_, err := store.db.Query().
		Table(store.PriceTableName(symbol, exchange, timeframe)).
		Where(COLUMN_ID+" = ?", price.ID()).
		Update(row)
	return err
}

// == PRIVATE METHODS ==========================================================

// instrumentToMap serializes an instrument to a column map using only the
// InstrumentInterface methods (no concrete type assertions). Metas are
// JSON-marshaled and timeframes are comma-joined to match the DB schema.
func (store *storeImplementation) instrumentToMap(instrument InstrumentInterface) (map[string]any, error) {
	metas, err := instrument.Metas()
	if err != nil {
		return nil, err
	}
	metasBytes, err := json.Marshal(metas)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		COLUMN_ID:              instrument.ID(),
		COLUMN_SYMBOL:          instrument.Symbol(),
		COLUMN_EXCHANGE:        instrument.Exchange(),
		COLUMN_ASSET_CLASS:     instrument.AssetClass(),
		COLUMN_NAME:            instrument.Name(),
		COLUMN_STATUS:          instrument.Status(),
		COLUMN_DESCRIPTION:     instrument.Description(),
		COLUMN_MEMO:            instrument.Memo(),
		COLUMN_METAS:           string(metasBytes),
		COLUMN_TIMEFRAMES:      strings.Join(instrument.Timeframes(), ","),
		COLUMN_CREATED_AT:      instrument.CreatedAtCarbon().StdTime(),
		COLUMN_UPDATED_AT:      instrument.UpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT: instrument.SoftDeletedAtCarbon().StdTime(),
	}, nil
}

// buildInstrumentQuery returns a query for instruments based on the given query options
func (store *storeImplementation) buildInstrumentQuery(options InstrumentQueryInterface) contractsorm.Query {
	// Use Model() to enable neat's automatic soft delete handling via SoftDeletesMaxDate
	q := store.db.Query().Model(&instrumentImplementation{})

	if options == nil {
		return q
	}

	if options.IsAssetClassSet() && options.AssetClass() != "" {
		q = q.Where(COLUMN_ASSET_CLASS+" = ?", options.AssetClass())
	}

	if options.IsExchangeSet() && options.Exchange() != "" {
		q = q.Where(COLUMN_EXCHANGE+" = ?", options.Exchange())
	}

	if options.IsIDSet() && options.ID() != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID())
	}

	if options.IsIDInSet() && len(options.IDIn()) > 0 {
		args := make([]any, len(options.IDIn()))
		for i, id := range options.IDIn() {
			args[i] = id
		}
		q = q.WhereIn(COLUMN_ID, args)
	}

	if options.IsStatusSet() && options.Status() != "" {
		q = q.Where(COLUMN_STATUS+" = ?", options.Status())
	}

	if options.IsSymbolSet() && options.Symbol() != "" {
		q = q.Where(COLUMN_SYMBOL+" = ?", options.Symbol())
	}

	if options.IsSymbolLikeSet() && options.SymbolLike() != "" {
		q = q.Where(COLUMN_SYMBOL+" LIKE ?", "%"+options.SymbolLike()+"%")
	}

	if !options.IsCountOnly() {
		if options.IsLimitSet() && options.Limit() > 0 {
			q = q.Limit(options.Limit())
		}

		if options.IsOffsetSet() && options.Offset() > 0 {
			q = q.Offset(options.Offset())
		}
	}

	if options.IsOrderBySet() && options.OrderBy() != "" {
		sortOrder := "desc"
		if options.IsOrderDirectionSet() && options.OrderDirection() != "" {
			sortOrder = options.OrderDirection()
		}
		q = q.OrderBy(options.OrderBy(), sortOrder)
	}

	// Handle soft delete filtering via neat's automatic handling (SoftDeletesMaxDate)
	if options.IsSoftDeletedIncluded() && options.SoftDeletedIncluded() {
		q = q.WithSoftDeleted()
	}

	return q
}

// buildPriceQuery returns a query for prices based on the given query options
func (store *storeImplementation) buildPriceQuery(symbol string, exchange string, timeframe string, options PriceQueryInterface) contractsorm.Query {
	q := store.db.Query()

	if options == nil {
		return q
	}

	if options.IsIDSet() && options.ID() != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID())
	}

	if options.IsIDInSet() && len(options.IDIn()) > 0 {
		args := make([]any, len(options.IDIn()))
		for i, id := range options.IDIn() {
			args[i] = id
		}
		q = q.WhereIn(COLUMN_ID, args)
	}

	if options.IsTimeSet() && options.Time() != "" {
		q = q.Where(COLUMN_TIME+" = ?", options.Time())
	}

	if options.IsTimeGteSet() && options.TimeGte() != "" {
		q = q.Where(COLUMN_TIME+" >= ?", options.TimeGte())
	}

	if options.IsTimeLteSet() && options.TimeLte() != "" {
		q = q.Where(COLUMN_TIME+" <= ?", options.TimeLte())
	}

	if !options.IsCountOnly() {
		if options.IsLimitSet() && options.Limit() > 0 {
			q = q.Limit(options.Limit())
		}

		if options.IsOffsetSet() && options.Offset() > 0 {
			q = q.Offset(options.Offset())
		}
	}

	if options.IsOrderBySet() && options.OrderBy() != "" {
		sortOrder := "asc"
		if options.IsOrderDirectionSet() && options.OrderDirection() != "" {
			sortOrder = options.OrderDirection()
		}
		q = q.OrderBy(options.OrderBy(), sortOrder)
	} else {
		q = q.OrderBy(COLUMN_TIME, "asc")
	}

	return q
}
