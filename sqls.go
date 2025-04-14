package tradingstore

import "github.com/gouniverse/sb"

func (store *Store) sqlTablePriceCreate() string {
	builder := sb.NewBuilder(sb.DatabaseDriverName(store.db)).
		Table(store.priceTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     50,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_ASSET_CLASS,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   50,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_SYMBOL,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   50,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_EXCHANGE,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   50,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_OPEN,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   50,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_HIGH,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   50,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_LOW,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   50,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_CLOSE,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   50,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_VOLUME,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   50,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_TIME,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   50,
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
			Name:     COLUMN_DELETED_AT,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   50,
			Nullable: true,
		})

	// Create the table
	sql := builder.CreateIfNotExists()

	// Create indexes (using separate SQL statements)
	sql += "\n\nCREATE INDEX IF NOT EXISTS idx_" + store.priceTableName + "_symbol ON " + store.priceTableName + " (" + COLUMN_SYMBOL + ");"
	sql += "\nCREATE INDEX IF NOT EXISTS idx_" + store.priceTableName + "_exchange ON " + store.priceTableName + " (" + COLUMN_EXCHANGE + ");"
	sql += "\nCREATE INDEX IF NOT EXISTS idx_" + store.priceTableName + "_asset_class ON " + store.priceTableName + " (" + COLUMN_ASSET_CLASS + ");"
	sql += "\nCREATE INDEX IF NOT EXISTS idx_" + store.priceTableName + "_time ON " + store.priceTableName + " (" + COLUMN_TIME + ");"
	sql += "\nCREATE INDEX IF NOT EXISTS idx_" + store.priceTableName + "_deleted_at ON " + store.priceTableName + " (" + COLUMN_DELETED_AT + ");"

	return sql
}

func (store *Store) sqlTableInstrumentCreate() string {
	builder := sb.NewBuilder(sb.DatabaseDriverName(store.db)).
		Table(store.instrumentTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     50,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_ASSET_CLASS,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   50,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_SYMBOL,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   50,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     COLUMN_EXCHANGE,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   50,
			Nullable: true,
		}).
		Column(sb.Column{
			Name:     "description",
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
			Name:     COLUMN_DELETED_AT,
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   50,
			Nullable: true,
		})

	// Create the table
	sql := builder.CreateIfNotExists()

	// Create indexes (using separate SQL statements)
	sql += "\n\nCREATE INDEX IF NOT EXISTS idx_" + store.instrumentTableName + "_symbol ON " + store.instrumentTableName + " (" + COLUMN_SYMBOL + ");"
	sql += "\nCREATE INDEX IF NOT EXISTS idx_" + store.instrumentTableName + "_exchange ON " + store.instrumentTableName + " (" + COLUMN_EXCHANGE + ");"
	sql += "\nCREATE INDEX IF NOT EXISTS idx_" + store.instrumentTableName + "_asset_class ON " + store.instrumentTableName + " (" + COLUMN_ASSET_CLASS + ");"
	sql += "\nCREATE INDEX IF NOT EXISTS idx_" + store.instrumentTableName + "_deleted_at ON " + store.instrumentTableName + " (" + COLUMN_DELETED_AT + ");"

	return sql
}
