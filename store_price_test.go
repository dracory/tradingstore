package tradingstore

import (
	"context"
	"testing"

	_ "modernc.org/sqlite"
)

func TestStorePriceCreate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		PriceTableName:     "forex_price_create",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	price := NewPrice().
		SetTime("2020-01-01 00:00:00").
		SetOpen("20.00").
		SetHigh("22.00").
		SetLow("18.00").
		SetClose("19.00").
		SetVolume("1000")

	ctx := context.Background()
	err = store.PriceCreate(ctx, price)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStorePriceFindByID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		PriceTableName:     "forex_price_create",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	price := NewPrice().
		SetTime("2020-01-01 00:00:00").
		SetOpen("20.12").
		SetHigh("22.23").
		SetLow("18.34").
		SetClose("19.45").
		SetVolume("1000")

	ctx := context.Background()
	err = store.PriceCreate(ctx, price)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	priceFound, errFind := store.PriceFindByID(ctx, price.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if priceFound == nil {
		t.Fatal("Price MUST NOT be nil")
	}

	if priceFound.ID() != price.ID() {
		t.Fatal("Price id MUST BE ", price.ID(), ", found: ", priceFound.ID())
	}

	if !priceFound.GetTimeCarbon().Eq(price.GetTimeCarbon()) {
		t.Fatal("Price time MUST BE ", price.GetTime(), ", found: ", priceFound.GetTime()[0:19])
	}

	if priceFound.GetOpen()[0:5] != price.GetOpen() {
		t.Fatal("Price open MUST BE ", price.GetOpen(), ", found: ", priceFound.GetOpen()[0:5])
	}

	if priceFound.GetHigh()[0:5] != price.GetHigh() {
		t.Fatal("Price high MUST BE ", price.GetHigh(), ", found: ", priceFound.GetHigh()[0:5])
	}

	if priceFound.GetLow()[0:5] != price.GetLow() {
		t.Fatal("Price low MUST BE ", price.GetLow(), ", found: ", priceFound.GetLow()[0:5])
	}

	if priceFound.GetClose()[0:5] != price.GetClose() {
		t.Fatal("Price close MUST BE ", price.GetClose(), ", found: ", priceFound.GetClose()[0:5])
	}

	if priceFound.GetVolume() != price.GetVolume() {
		t.Fatal("Price volume MUST BE ", price.GetVolume(), ", found: ", priceFound.GetVolume())
	}
}
