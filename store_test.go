package tradingstore

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strings"

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
		DB:                   db,
		PriceTableNamePrefix: "price_",
		InstrumentTableName:  "instrument",
		UseMultipleExchanges: true,
		AutomigrateEnabled:   true,
	})

	if err != nil {
		return nil, err
	}

	if store == nil {
		return nil, errors.New("unexpected nil store")
	}

	err = seedInstruments(store)
	if err != nil {
		return nil, err
	}

	err = store.AutoMigratePrices(context.Background())
	if err != nil {
		return nil, err
	}

	return store, nil
}

func seedInstruments(store StoreInterface) error {
	timeframes := TIMEFRAME_1_MINUTE + "," + TIMEFRAME_5_MINUTES + "," + TIMEFRAME_15_MINUTES + "," + TIMEFRAME_30_MINUTES + "," + TIMEFRAME_1_HOUR + "," + TIMEFRAME_4_HOURS + "," + TIMEFRAME_1_DAY

	data := []map[string]string{
		{
			"symbol":      "AAPL",
			"exchange":    "NASDAQ",
			"asset_class": ASSET_CLASS_STOCK,
			"timeframes":  timeframes,
		},
		{
			"symbol":      "MSFT",
			"exchange":    "NASDAQ",
			"asset_class": ASSET_CLASS_STOCK,
			"timeframes":  timeframes,
		},
	}

	for _, data := range data {
		instrument := NewInstrument().
			SetSymbol(data["symbol"]).
			SetExchange(data["exchange"]).
			SetAssetClass(data["asset_class"]).
			SetTimeframes(strings.Split(data["timeframes"], ","))

		err := store.InstrumentCreate(context.Background(), instrument)
		if err != nil {
			return err
		}
	}

	return nil
}
