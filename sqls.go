package tradingstore

import (
	"github.com/gouniverse/sb"
)

func (store *Store) sqlTablePriceCreate() string {
	sql := sb.NewBuilder(sb.DatabaseDriverName(store.db)).
		Table(store.priceTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name: COLUMN_TIME,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name:     COLUMN_HIGH,
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   10,
			Decimals: 8,
		}).
		Column(sb.Column{
			Name:     COLUMN_LOW,
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   10,
			Decimals: 8,
		}).
		Column(sb.Column{
			Name:     COLUMN_OPEN,
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   10,
			Decimals: 8,
		}).
		Column(sb.Column{
			Name:     COLUMN_CLOSE,
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   10,
			Decimals: 8,
		}).
		Column(sb.Column{
			Name:   COLUMN_VOLUME,
			Type:   sb.COLUMN_TYPE_INTEGER,
			Length: 12,
		}).
		CreateIfNotExists()

	return sql
}

func (store *Store) sqlTableSymbolCreate() string {
	sql := sb.NewBuilder(sb.DatabaseDriverName(store.db)).
		Table(store.priceTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   COLUMN_ASSET_CLASS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   COLUMN_SYMBOL,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 10,
		}).
		Column(sb.Column{
			Name:   COLUMN_EXCHANGE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:     COLUMN_CREATED_AT,
			Type:     sb.COLUMN_TYPE_DATETIME,
			Nullable: false,
		}).
		Column(sb.Column{
			Name:     COLUMN_UPDATED_AT,
			Type:     sb.COLUMN_TYPE_DATETIME,
			Nullable: false,
		}).
		Column(sb.Column{
			Name:     COLUMN_DELETED_AT,
			Type:     sb.COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		CreateIfNotExists()

	return sql
}
