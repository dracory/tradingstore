package tradingstore

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/dracory/base/database"
)

// ============================================================================
// == INTERFACE
// ============================================================================

var _ StoreInterface = (*Store)(nil) // verify it extends the interface

type Store struct {
	// priceTableNamePrefix is the prefix of the price table
	priceTableNamePrefix string

	// instrumentTableName is the name of the instrument table
	instrumentTableName string

	// useMultipleExchanges enables or disables the use of multiple exchanges
	// if true, a price table will be created for each exchange, i.e price_eurusd_binance_1min
	// if false, a price table will be created for the default exchange, i.e price_eurusd_1min
	useMultipleExchanges bool

	// db is the underlying database connection
	db *sql.DB

	// dbDriverName is the name of the database driver
	dbDriverName string

	// automigrateEnabled enables auto migrate
	automigrateEnabled bool

	// debugEnabled enables or disables the debug mode
	debugEnabled bool

	// sqlLogger is the sql logger used when debug mode is enabled
	sqlLogger *slog.Logger
}

// ============================================================================
// == PUBLIC METHODS
// ============================================================================

// AutoMigrateInstruments auto migrates the instrument table
func (store *Store) AutoMigrateInstruments(ctx context.Context) error {
	sql := store.sqlTableInstrumentCreate()

	_, err := store.db.Exec(sql)

	if err != nil {
		return err
	}

	return nil
}

// AutoMigratePrices auto migrates the price tables
// It will create a price table for each instrument and each timeframe
// You will need to call this method when you create a new instrument
func (store *Store) AutoMigratePrices(ctx context.Context) error {
	instruments, err := store.InstrumentList(ctx, InstrumentQuery())

	if err != nil {
		return err
	}

	sqls := []string{}
	for _, instrument := range instruments {
		timeframes := instrument.Timeframes()

		for _, timeframe := range timeframes {
			sql := store.sqlTablePriceCreate(instrument.Symbol(), instrument.Exchange(), timeframe)
			sqls = append(sqls, sql)
		}
	}

	for _, sql := range sqls {
		_, err := store.db.Exec(sql)

		if err != nil {
			return err
		}
	}

	return nil
}

// DB returns the underlying database connection
func (st *Store) DB() *sql.DB {
	return st.db
}

// EnableDebug - enables the debug option
func (st *Store) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

// ============================================================================
// == PRIVATE METHODS
// ============================================================================

// logSql logs sql to the sql logger, if debug mode is enabled
func (store *Store) logSql(sqlOperationType string, sql string, params ...interface{}) {
	if !store.debugEnabled {
		return
	}

	if store.sqlLogger != nil {
		store.sqlLogger.Debug("sql: "+sqlOperationType, slog.String("sql", sql), slog.Any("params", params))
	}
}

// toQuerableContext converts the context to a QueryableContext
func (store *Store) toQuerableContext(ctx context.Context) database.QueryableContext {
	if database.IsQueryableContext(ctx) {
		return ctx.(database.QueryableContext)
	}

	return database.Context(ctx, store.db)
}
