package tradingstore

import (
	"strings"

	"github.com/gouniverse/sb"
)

func (store *Store) PriceTableName(symbol string, exchange string, timeframe string) string {
	priceTableName := store.priceTableNamePrefix

	if exchange != "" {
		return priceTableName + strings.ToLower(symbol) + "_" + strings.ToLower(exchange) + "_" + strings.ToLower(timeframe)
	}

	return priceTableName + strings.ToLower(symbol) + "_" + strings.ToLower(timeframe)
}

func (store *Store) sqlTablePriceCreate(symbol string, exchange string, timeframe string) string {
	builder := sb.NewBuilder(sb.DatabaseDriverName(store.db)).
		Table(store.PriceTableName(symbol, exchange, timeframe)).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_OPEN,
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   10,
			Decimals: 8,
			Nullable: false,
		}).
		Column(sb.Column{
			Name:     COLUMN_HIGH,
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   10,
			Decimals: 8,
			Nullable: false,
		}).
		Column(sb.Column{
			Name:     COLUMN_LOW,
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   10,
			Decimals: 8,
			Nullable: false,
		}).
		Column(sb.Column{
			Name:     COLUMN_CLOSE,
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   10,
			Decimals: 8,
			Nullable: false,
		}).
		Column(sb.Column{
			Name:     COLUMN_VOLUME,
			Type:     sb.COLUMN_TYPE_INTEGER,
			Length:   10,
			Nullable: false,
		}).
		Column(sb.Column{
			Name:     COLUMN_TIME,
			Type:     sb.COLUMN_TYPE_DATETIME,
			Nullable: false,
		})

	// Create the table
	sql := builder.CreateIfNotExists()

	return sql
}

func (store *Store) sqlTableInstrumentCreate() string {
	builder := sb.NewBuilder(sb.DatabaseDriverName(store.db)).
		Table(store.instrumentTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_ASSET_CLASS,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   40,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_SYMBOL,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   10,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_EXCHANGE,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   50,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_TIMEFRAMES,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   100,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_DESCRIPTION,
			Type:     sb.COLUMN_TYPE_TEXT,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_CREATED_AT,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   50,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_UPDATED_AT,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   50,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_SOFT_DELETED_AT,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   50,
			Nullable: true,
		})

	// Create the table
	sql := builder.CreateIfNotExists()

	return sql
}

func (store *Store) sqlIndexesCreate() string {
	sql := ""

	// Create price indexes (using separate SQL statements)
	// sql += "\n\nCREATE INDEX IF NOT EXISTS idx_" + store.priceTableName + "_symbol ON " + store.priceTableName + " (" + COLUMN_SYMBOL + ");"
	// sql += "\nCREATE INDEX IF NOT EXISTS idx_" + store.priceTableName + "_exchange ON " + store.priceTableName + " (" + COLUMN_EXCHANGE + ");"
	// sql += "\nCREATE INDEX IF NOT EXISTS idx_" + store.priceTableName + "_asset_class ON " + store.priceTableName + " (" + COLUMN_ASSET_CLASS + ");"
	// sql += "\nCREATE INDEX IF NOT EXISTS idx_" + store.priceTableName + "_time ON " + store.priceTableName + " (" + COLUMN_TIME + ");"
	// sql += "\nCREATE INDEX IF NOT EXISTS idx_" + store.priceTableName + "_soft_deleted_at ON " + store.priceTableName + " (" + COLUMN_SOFT_DELETED_AT + ");"

	// Create instrument indexes (using separate SQL statements)
	sql += "\n\nCREATE INDEX IF NOT EXISTS idx_" + store.instrumentTableName + "_symbol ON " + store.instrumentTableName + " (" + COLUMN_SYMBOL + ");"
	sql += "\nCREATE INDEX IF NOT EXISTS idx_" + store.instrumentTableName + "_exchange ON " + store.instrumentTableName + " (" + COLUMN_EXCHANGE + ");"
	sql += "\nCREATE INDEX IF NOT EXISTS idx_" + store.instrumentTableName + "_asset_class ON " + store.instrumentTableName + " (" + COLUMN_ASSET_CLASS + ");"
	sql += "\nCREATE INDEX IF NOT EXISTS idx_" + store.instrumentTableName + "_soft_deleted_at ON " + store.instrumentTableName + " (" + COLUMN_SOFT_DELETED_AT + ");"

	return sql
}
