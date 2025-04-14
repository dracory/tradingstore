package tradingstore

import (
	"database/sql"
	"errors"
	"os"

	_ "modernc.org/sqlite"
)

func initDB(filepath string) *sql.DB {
	os.Remove(filepath) // remove database
	dsn := filepath + "?parseTime=true&loc=UTC&_loc=UTC"
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func initStore() (StoreInterface, error) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                  db,
		PriceTableName:      "price_create",
		InstrumentTableName: "instrument",
		AutomigrateEnabled:  true,
	})

	if err != nil {
		return nil, err
	}

	if store == nil {
		return nil, errors.New("unexpected nil store")
	}

	return store, nil
}
