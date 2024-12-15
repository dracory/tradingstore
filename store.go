package tradingstore

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/gouniverse/base/database"
)

// ============================================================================
// == INTERFACE
// ============================================================================

var _ StoreInterface = (*Store)(nil) // verify it extends the interface

type Store struct {
	// priceTableName is the name of the price table
	priceTableName string

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

// AutoMigrate auto migrate
func (store *Store) AutoMigrate() error {
	sql := store.sqlTablePriceCreate()

	_, err := store.db.Exec(sql)

	if err != nil {
		return err
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
