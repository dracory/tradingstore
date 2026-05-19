package tradingstore

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/dracory/database"
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

// MigrateUp creates the trading store tables
func (store *Store) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	var txToUse *sql.Tx
	if len(tx) > 0 {
		txToUse = tx[0]
	}

	// Create instrument table
	sql := store.sqlTableInstrumentCreate()
	var errExec error
	if txToUse != nil {
		_, errExec = txToUse.ExecContext(ctx, sql)
	} else {
		_, errExec = store.db.ExecContext(ctx, sql)
	}
	if errExec != nil {
		return errExec
	}

	// Create price tables for each instrument and timeframe
	instruments, err := store.InstrumentList(ctx, InstrumentQuery())
	if err != nil {
		return err
	}

	for _, instrument := range instruments {
		timeframes := instrument.Timeframes()
		for _, timeframe := range timeframes {
			sql := store.sqlTablePriceCreate(instrument.Symbol(), instrument.Exchange(), timeframe)
			if txToUse != nil {
				_, errExec = txToUse.ExecContext(ctx, sql)
			} else {
				_, errExec = store.db.ExecContext(ctx, sql)
			}
			if errExec != nil {
				return errExec
			}
		}
	}

	return nil
}

// MigrateDown drops the trading store tables
func (store *Store) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	var txToUse *sql.Tx
	if len(tx) > 0 {
		txToUse = tx[0]
	}

	// Drop price tables for each instrument and timeframe
	instruments, err := store.InstrumentList(ctx, InstrumentQuery())
	if err != nil {
		return err
	}

	for _, instrument := range instruments {
		timeframes := instrument.Timeframes()
		for _, timeframe := range timeframes {
			sql, err := store.sqlTablePriceDrop(instrument.Symbol(), instrument.Exchange(), timeframe)
			if err != nil {
				return err
			}
			var errExec error
			if txToUse != nil {
				_, errExec = txToUse.ExecContext(ctx, sql)
			} else {
				_, errExec = store.db.ExecContext(ctx, sql)
			}
			if errExec != nil {
				return errExec
			}
		}
	}

	// Drop instrument table
	sql, err := store.sqlTableInstrumentDrop()
	if err != nil {
		return err
	}
	var errExec error
	if txToUse != nil {
		_, errExec = txToUse.ExecContext(ctx, sql)
	} else {
		_, errExec = store.db.ExecContext(ctx, sql)
	}
	if errExec != nil {
		return errExec
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
