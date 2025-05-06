package tradingstore

import (
	"context"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestStorePriceCreate(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	price := NewPrice().
		SetTime("2020-01-01 00:00:00").
		SetOpen("20.00").
		SetHigh("22.00").
		SetLow("18.00").
		SetClose("19.00").
		SetVolume("1000")

	ctx := context.Background()
	err = store.PriceCreate(ctx, "AAPL", "NASDAQ", "1min", price)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStorePriceFindByID(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	price := NewPrice().
		SetTime("2020-01-01 00:00:00").
		SetOpen("20.12").
		SetHigh("22.23").
		SetLow("18.34").
		SetClose("19.45").
		SetVolume("1000")

	ctx := context.Background()
	err = store.PriceCreate(ctx, "AAPL", "NASDAQ", "1min", price)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	priceFound, errFind := store.PriceFindByID(ctx, "AAPL", "NASDAQ", "1min", price.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if priceFound == nil {
		t.Fatal("Price MUST NOT be nil")
	}

	if priceFound.ID() != price.ID() {
		t.Fatal("Price id MUST BE ", price.ID(), ", found: ", priceFound.ID())
	}

	if !priceFound.TimeCarbon().Eq(price.TimeCarbon()) {
		t.Fatal("Price time MUST BE ", price.Time(), ", found: ", priceFound.Time())
	}

	if priceFound.OpenFloat() != price.OpenFloat() {
		t.Fatal("Price open MUST BE ", price.OpenFloat(), ", found: ", priceFound.OpenFloat())
	}

	if priceFound.HighFloat() != price.HighFloat() {
		t.Fatal("Price high MUST BE ", price.HighFloat(), ", found: ", priceFound.HighFloat())
	}

	if priceFound.LowFloat() != price.LowFloat() {
		t.Fatal("Price low MUST BE ", price.LowFloat(), ", found: ", priceFound.LowFloat())
	}

	if priceFound.CloseFloat() != price.CloseFloat() {
		t.Fatal("Price close MUST BE ", price.CloseFloat(), ", found: ", priceFound.CloseFloat())
	}

	if priceFound.Volume() != price.Volume() {
		t.Fatal("Price volume MUST BE ", price.Volume(), ", found: ", priceFound.Volume())
	}
}

func TestStorePriceDelete(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	price := NewPrice().
		SetTime("2020-01-01 00:00:00").
		SetOpen("20.00").
		SetHigh("22.00").
		SetLow("18.00").
		SetClose("19.00").
		SetVolume("1000")

	ctx := context.Background()

	// Create the price first
	err = store.PriceCreate(ctx, "AAPL", "NASDAQ", "1min", price)
	if err != nil {
		t.Fatal("unexpected error creating price:", err)
	}

	// Verify it exists
	priceFound, errFind := store.PriceFindByID(ctx, "AAPL", "NASDAQ", "1min", price.ID())
	if errFind != nil {
		t.Fatal("unexpected error finding price:", errFind)
	}
	if priceFound == nil {
		t.Fatal("Price should exist")
	}

	// Delete the price
	err = store.PriceDelete(ctx, "AAPL", "NASDAQ", "1min", price)
	if err != nil {
		t.Fatal("unexpected error deleting price:", err)
	}

	// Verify it was deleted
	priceDeleted, errFindDeleted := store.PriceFindByID(ctx, "AAPL", "NASDAQ", "1min", price.ID())
	if errFindDeleted != nil {
		t.Fatal("unexpected error finding deleted price:", errFindDeleted)
	}
	if priceDeleted != nil {
		t.Fatal("Price should be deleted")
	}
}

func TestStorePriceDeleteByID(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	price := NewPrice().
		SetTime("2020-01-01 00:00:00").
		SetOpen("20.00").
		SetHigh("22.00").
		SetLow("18.00").
		SetClose("19.00").
		SetVolume("1000")

	ctx := context.Background()

	// Create the price first
	err = store.PriceCreate(ctx, "AAPL", "NASDAQ", "1min", price)
	if err != nil {
		t.Fatal("unexpected error creating price:", err)
	}

	// Verify it exists
	exists, errExists := store.PriceExists(ctx, "AAPL", "NASDAQ", "1min", NewPriceQuery().SetID(price.ID()))
	if errExists != nil {
		t.Fatal("unexpected error checking if price exists:", errExists)
	}
	if !exists {
		t.Fatal("Price should exist")
	}

	// Delete the price by ID
	err = store.PriceDeleteByID(ctx, "AAPL", "NASDAQ", "1min", price.ID())
	if err != nil {
		t.Fatal("unexpected error deleting price by ID:", err)
	}

	// Verify it was deleted
	existsAfterDelete, errExistsAfterDelete := store.PriceExists(ctx, "AAPL", "NASDAQ", "1min", NewPriceQuery().SetID(price.ID()))
	if errExistsAfterDelete != nil {
		t.Fatal("unexpected error checking if price exists after delete:", errExistsAfterDelete)
	}
	if existsAfterDelete {
		t.Fatal("Price should be deleted")
	}
}

func TestStorePriceExists(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create prices with different timestamps
	price1 := NewPrice().
		SetTime("2020-01-01 00:00:00").
		SetOpen("20.00").
		SetHigh("22.00").
		SetLow("18.00").
		SetClose("19.00").
		SetVolume("1000")

	price2 := NewPrice().
		SetTime("2020-01-02 00:00:00").
		SetOpen("19.00").
		SetHigh("21.00").
		SetLow("17.00").
		SetClose("18.00").
		SetVolume("2000")

	err = store.PriceCreate(ctx, "AAPL", "NASDAQ", "1min", price1)
	if err != nil {
		t.Fatal("unexpected error creating price1:", err)
	}

	err = store.PriceCreate(ctx, "AAPL", "NASDAQ", "1min", price2)
	if err != nil {
		t.Fatal("unexpected error creating price2:", err)
	}

	// Test exists by ID
	exists, errExists := store.PriceExists(ctx, "AAPL", "NASDAQ", "1min", NewPriceQuery().SetID(price1.ID()))
	if errExists != nil {
		t.Fatal("unexpected error checking if price exists by ID:", errExists)
	}
	if !exists {
		t.Fatal("Price should exist by ID")
	}

	// Test exists by time
	exists, errExists = store.PriceExists(ctx, "AAPL", "NASDAQ", "1min", NewPriceQuery().SetTime("2020-01-02 00:00:00"))
	if errExists != nil {
		t.Fatal("unexpected error checking if price exists by time:", errExists)
	}
	if !exists {
		t.Fatal("Price should exist by time")
	}

	// Test exists by time range
	exists, errExists = store.PriceExists(ctx, "AAPL", "NASDAQ", "1min", NewPriceQuery().
		SetTimeGte("2020-01-01 00:00:00").
		SetTimeLte("2020-01-02 23:59:59"))
	if errExists != nil {
		t.Fatal("unexpected error checking if price exists by time range:", errExists)
	}
	if !exists {
		t.Fatal("Price should exist within time range")
	}

	// Test not exists with non-existent ID
	exists, errExists = store.PriceExists(ctx, "AAPL", "NASDAQ", "1min", NewPriceQuery().SetID("non-existent-id"))
	if errExists != nil {
		t.Fatal("unexpected error checking if non-existent price exists:", errExists)
	}
	if exists {
		t.Fatal("Price should not exist with non-existent ID")
	}
}

func TestStorePriceCount(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create prices with different timestamps
	price1 := NewPrice().
		SetTime("2020-01-01 00:00:00").
		SetOpen("20.00").
		SetHigh("22.00").
		SetLow("18.00").
		SetClose("19.00").
		SetVolume("1000")

	price2 := NewPrice().
		SetTime("2020-01-02 00:00:00").
		SetOpen("19.00").
		SetHigh("21.00").
		SetLow("17.00").
		SetClose("18.00").
		SetVolume("2000")

	price3 := NewPrice().
		SetTime("2020-01-03 00:00:00").
		SetOpen("18.00").
		SetHigh("20.00").
		SetLow("16.00").
		SetClose("17.00").
		SetVolume("3000")

	err = store.PriceCreate(ctx, "AAPL", "NASDAQ", "1min", price1)
	if err != nil {
		t.Fatal("unexpected error creating price1:", err)
	}

	err = store.PriceCreate(ctx, "AAPL", "NASDAQ", "1min", price2)
	if err != nil {
		t.Fatal("unexpected error creating price2:", err)
	}

	err = store.PriceCreate(ctx, "AAPL", "NASDAQ", "1min", price3)
	if err != nil {
		t.Fatal("unexpected error creating price3:", err)
	}

	// Test count all
	count, errCount := store.PriceCount(ctx, "AAPL", "NASDAQ", "1min", NewPriceQuery())
	if errCount != nil {
		t.Fatal("unexpected error counting all prices:", errCount)
	}
	if count != 3 {
		t.Fatal("Count should be 3, got:", count)
	}

	// Test count by time range (partial)
	count, errCount = store.PriceCount(ctx, "AAPL", "NASDAQ", "1min", NewPriceQuery().
		SetTimeGte("2020-01-01 00:00:00").
		SetTimeLte("2020-01-02 00:00:00"))
	if errCount != nil {
		t.Fatal("unexpected error counting prices by time range:", errCount)
	}
	if count != 2 {
		t.Fatal("Count should be 2 for partial range, got:", count)
	}

	// Test count with ID filter (single result)
	count, errCount = store.PriceCount(ctx, "AAPL", "NASDAQ", "1min", NewPriceQuery().SetID(price1.ID()))
	if errCount != nil {
		t.Fatal("unexpected error counting prices by ID:", errCount)
	}
	if count != 1 {
		t.Fatal("Count should be 1 for ID filter, got:", count)
	}

	// Test count with multiple IDs
	count, errCount = store.PriceCount(ctx, "AAPL", "NASDAQ", "1min", NewPriceQuery().SetIDIn([]string{price1.ID(), price3.ID()}))
	if errCount != nil {
		t.Fatal("unexpected error counting prices by multiple IDs:", errCount)
	}
	if count != 2 {
		t.Fatal("Count should be 2 for multiple IDs filter, got:", count)
	}
}

func TestStorePriceList(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create prices with different timestamps
	price1 := NewPrice().
		SetTime("2020-01-01 00:00:00").
		SetOpen("20.00").
		SetHigh("22.00").
		SetLow("18.00").
		SetClose("19.00").
		SetVolume("1000")

	err = store.PriceCreate(ctx, "AAPL", "NASDAQ", "1min", price1)
	if err != nil {
		t.Fatal("unexpected error creating price1:", err)
	}

	time.Sleep(1 * time.Second)

	price2 := NewPrice().
		SetTime("2020-01-02 00:00:00").
		SetOpen("19.00").
		SetHigh("21.00").
		SetLow("17.00").
		SetClose("18.00").
		SetVolume("2000")

	err = store.PriceCreate(ctx, "AAPL", "NASDAQ", "1min", price2)
	if err != nil {
		t.Fatal("unexpected error creating price2:", err)
	}

	time.Sleep(1 * time.Second)

	price3 := NewPrice().
		SetTime("2020-01-03 00:00:00").
		SetOpen("18.00").
		SetHigh("20.00").
		SetLow("16.00").
		SetClose("17.00").
		SetVolume("3000")

	err = store.PriceCreate(ctx, "AAPL", "NASDAQ", "1min", price3)
	if err != nil {
		t.Fatal("unexpected error creating price3:", err)
	}

	// Test list all
	prices, errList := store.PriceList(ctx, "AAPL", "NASDAQ", "1min", NewPriceQuery())
	if errList != nil {
		t.Fatal("unexpected error listing all prices:", errList)
	}
	if len(prices) != 3 {
		t.Fatal("Should list 3 prices, got:", len(prices))
	}

	// Test list with limit
	prices, errList = store.PriceList(ctx, "AAPL", "NASDAQ", "1min", NewPriceQuery().SetLimit(2))
	if errList != nil {
		t.Fatal("unexpected error listing prices with limit:", errList)
	}
	if len(prices) != 2 {
		t.Fatal("Should list 2 prices with limit, got:", len(prices))
	}

	// Test list with time range
	prices, errList = store.PriceList(ctx, "AAPL", "NASDAQ", "1min", NewPriceQuery().
		SetTimeGte("2020-01-02 00:00:00").
		SetTimeLte("2020-01-03 00:00:00"))
	if errList != nil {
		t.Fatal("unexpected error listing prices by time range:", errList)
	}
	if len(prices) != 2 {
		t.Fatal("Should list 2 prices in time range, got:", len(prices))
	}

	// Test list with ordering (ascending)
	prices, errList = store.PriceList(ctx, "AAPL", "NASDAQ", "1min", NewPriceQuery().
		SetOrderBy(COLUMN_TIME).
		SetOrderDirection("ASC"))
	if errList != nil {
		t.Fatal("unexpected error listing prices with ordering:", errList)
	}
	if len(prices) != 3 {
		t.Fatal("Should list 3 prices with ordering, got:", len(prices))
	}
	if !prices[0].TimeCarbon().Eq(price1.TimeCarbon()) {
		t.Fatal("First price should be oldest, got:", prices[0].Time())
	}

	// Test list with ordering (descending)
	prices, errList = store.PriceList(ctx, "AAPL", "NASDAQ", "1min", NewPriceQuery().
		SetOrderBy(COLUMN_TIME).
		SetOrderDirection("DESC"))
	if errList != nil {
		t.Fatal("unexpected error listing prices with descending ordering:", errList)
	}
	if !prices[0].TimeCarbon().Eq(price3.TimeCarbon()) {
		t.Fatal("First price should be newest, got:", prices[0].Time())
	}

	// Test list with offset
	prices, errList = store.PriceList(ctx, "AAPL", "NASDAQ", "1min", NewPriceQuery().
		SetOrderBy(COLUMN_TIME).
		SetOrderDirection("ASC").
		SetOffset(1).
		SetLimit(10))
	if errList != nil {
		t.Fatal("unexpected error listing prices with offset:", errList)
	}
	if len(prices) != 2 {
		t.Fatal("Should list 2 prices with offset, got:", len(prices))
	}
	if !prices[0].TimeCarbon().Eq(price2.TimeCarbon()) {
		t.Fatal("First price with offset should be second oldest, got:", prices[0].Time())
	}
}

func TestStorePriceUpdate(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create a price
	price := NewPrice().
		SetTime("2020-01-01 00:00:00").
		SetOpen("20.00").
		SetHigh("22.00").
		SetLow("18.00").
		SetClose("19.00").
		SetVolume("1000")

	err = store.PriceCreate(ctx, "AAPL", "NASDAQ", "1min", price)
	if err != nil {
		t.Fatal("unexpected error creating price:", err)
	}

	// Update the price values
	price.SetOpen("21.00")
	price.SetHigh("23.00")
	price.SetLow("19.00")
	price.SetClose("20.00")
	price.SetVolume("1500")

	// Save the updates
	err = store.PriceUpdate(ctx, "AAPL", "NASDAQ", "1min", price)
	if err != nil {
		t.Fatal("unexpected error updating price:", err)
	}

	// Retrieve the updated price
	updatedPrice, errFind := store.PriceFindByID(ctx, "AAPL", "NASDAQ", "1min", price.ID())
	if errFind != nil {
		t.Fatal("unexpected error finding updated price:", errFind)
	}

	if updatedPrice == nil {
		t.Fatal("Updated price should exist")
	}

	// Check if the values were updated correctly
	if updatedPrice.OpenFloat() != 21.0 {
		t.Fatal("Price open should be updated to 21.00, got:", updatedPrice.OpenFloat())
	}

	if updatedPrice.HighFloat() != 23.0 {
		t.Fatal("Price high should be updated to 23.00, got:", updatedPrice.HighFloat())
	}

	if updatedPrice.LowFloat() != 19.0 {
		t.Fatal("Price low should be updated to 19.00, got:", updatedPrice.LowFloat())
	}

	if updatedPrice.CloseFloat() != 20.0 {
		t.Fatal("Price close should be updated to 20.00, got:", updatedPrice.CloseFloat())
	}

	if updatedPrice.VolumeFloat() != 1500.0 {
		t.Fatal("Price volume should be updated to 1500, got:", updatedPrice.VolumeFloat())
	}

	// Check that the ID and time remain unchanged
	if updatedPrice.ID() != price.ID() {
		t.Fatal("Price ID should remain unchanged")
	}

	if !updatedPrice.TimeCarbon().Eq(price.TimeCarbon()) {
		t.Fatal("Price time should remain unchanged")
	}
}
