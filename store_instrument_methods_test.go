package tradingstore

import (
	"context"
	"testing"

	_ "modernc.org/sqlite"
)

// clearInstruments removes all existing instruments from the database
func clearInstruments(t *testing.T, store StoreInterface) {
	ctx := context.Background()
	instruments, err := store.InstrumentList(ctx, NewInstrumentQuery())
	if err != nil {
		t.Fatal("Error listing instruments to clear:", err)
	}

	for _, instrument := range instruments {
		err := store.InstrumentDelete(ctx, instrument)
		if err != nil {
			t.Fatal("Error deleting instrument:", err)
		}
	}
}

func TestStoreInstrumentCreate(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Clear any existing instruments
	clearInstruments(t, store)

	instrument := NewInstrument().
		SetSymbol("AAPL").
		SetExchange("NASDAQ").
		SetAssetClass(ASSET_CLASS_STOCK)

	ctx := context.Background()
	err = store.InstrumentCreate(ctx, instrument)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreInstrumentFindByID(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Clear any existing instruments
	clearInstruments(t, store)

	instrument := NewInstrument().
		SetSymbol("MSFT").
		SetExchange("NASDAQ").
		SetAssetClass(ASSET_CLASS_STOCK)

	ctx := context.Background()
	err = store.InstrumentCreate(ctx, instrument)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	instrumentFound, errFind := store.InstrumentFindByID(ctx, instrument.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if instrumentFound == nil {
		t.Fatal("Instrument MUST NOT be nil")
	}

	if instrumentFound.ID() != instrument.ID() {
		t.Fatal("Instrument id MUST BE ", instrument.ID(), ", found: ", instrumentFound.ID())
	}

	if instrumentFound.Symbol() != instrument.Symbol() {
		t.Fatal("Instrument symbol MUST BE ", instrument.Symbol(), ", found: ", instrumentFound.Symbol())
	}

	if instrumentFound.Exchange() != instrument.Exchange() {
		t.Fatal("Instrument exchange MUST BE ", instrument.Exchange(), ", found: ", instrumentFound.Exchange())
	}

	if instrumentFound.AssetClass() != instrument.AssetClass() {
		t.Fatal("Instrument asset class MUST BE ", instrument.AssetClass(), ", found: ", instrumentFound.AssetClass())
	}
}

func TestStoreInstrumentDelete(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Clear any existing instruments
	clearInstruments(t, store)

	instrument := NewInstrument().
		SetSymbol("GOOG").
		SetExchange("NASDAQ").
		SetAssetClass(ASSET_CLASS_STOCK)

	ctx := context.Background()

	// Create the instrument first
	err = store.InstrumentCreate(ctx, instrument)
	if err != nil {
		t.Fatal("unexpected error creating instrument:", err)
	}

	// Verify it exists
	instrumentFound, errFind := store.InstrumentFindByID(ctx, instrument.ID())
	if errFind != nil {
		t.Fatal("unexpected error finding instrument:", errFind)
	}
	if instrumentFound == nil {
		t.Fatal("Instrument should exist")
	}

	// Delete the instrument
	err = store.InstrumentDelete(ctx, instrument)
	if err != nil {
		t.Fatal("unexpected error deleting instrument:", err)
	}

	// Verify it was deleted
	instrumentDeleted, errFindDeleted := store.InstrumentFindByID(ctx, instrument.ID())
	if errFindDeleted != nil {
		t.Fatal("unexpected error finding deleted instrument:", errFindDeleted)
	}
	if instrumentDeleted != nil {
		t.Fatal("Instrument should be deleted")
	}
}

func TestStoreInstrumentDeleteByID(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Clear any existing instruments
	clearInstruments(t, store)

	instrument := NewInstrument().
		SetSymbol("AMZN").
		SetExchange("NASDAQ").
		SetAssetClass(ASSET_CLASS_STOCK)

	ctx := context.Background()

	// Create the instrument first
	err = store.InstrumentCreate(ctx, instrument)
	if err != nil {
		t.Fatal("unexpected error creating instrument:", err)
	}

	// Verify it exists
	exists, errExists := store.InstrumentExists(ctx, NewInstrumentQuery().SetID(instrument.ID()))
	if errExists != nil {
		t.Fatal("unexpected error checking if instrument exists:", errExists)
	}
	if !exists {
		t.Fatal("Instrument should exist")
	}

	// Delete the instrument by ID
	err = store.InstrumentDeleteByID(ctx, instrument.ID())
	if err != nil {
		t.Fatal("unexpected error deleting instrument by ID:", err)
	}

	// Verify it was deleted
	existsAfterDelete, errExistsAfterDelete := store.InstrumentExists(ctx, NewInstrumentQuery().SetID(instrument.ID()))
	if errExistsAfterDelete != nil {
		t.Fatal("unexpected error checking if instrument exists after delete:", errExistsAfterDelete)
	}
	if existsAfterDelete {
		t.Fatal("Instrument should be deleted")
	}
}

func TestStoreInstrumentExists(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Clear any existing instruments
	clearInstruments(t, store)

	ctx := context.Background()

	// Create instruments with different properties
	instrument1 := NewInstrument().
		SetSymbol("TSLA").
		SetExchange("NASDAQ").
		SetAssetClass(ASSET_CLASS_STOCK)

	instrument2 := NewInstrument().
		SetSymbol("JPM").
		SetExchange("NYSE").
		SetAssetClass(ASSET_CLASS_STOCK)

	err = store.InstrumentCreate(ctx, instrument1)
	if err != nil {
		t.Fatal("unexpected error creating instrument1:", err)
	}

	err = store.InstrumentCreate(ctx, instrument2)
	if err != nil {
		t.Fatal("unexpected error creating instrument2:", err)
	}

	// Test exists by ID
	exists, errExists := store.InstrumentExists(ctx, NewInstrumentQuery().SetID(instrument1.ID()))
	if errExists != nil {
		t.Fatal("unexpected error checking if instrument exists by ID:", errExists)
	}
	if !exists {
		t.Fatal("Instrument should exist by ID")
	}

	// Test exists by symbol
	exists, errExists = store.InstrumentExists(ctx, NewInstrumentQuery().SetSymbol("JPM"))
	if errExists != nil {
		t.Fatal("unexpected error checking if instrument exists by symbol:", errExists)
	}
	if !exists {
		t.Fatal("Instrument should exist by symbol")
	}

	// Test exists by exchange
	exists, errExists = store.InstrumentExists(ctx, NewInstrumentQuery().SetExchange("NYSE"))
	if errExists != nil {
		t.Fatal("unexpected error checking if instrument exists by exchange:", errExists)
	}
	if !exists {
		t.Fatal("Instrument should exist by exchange")
	}

	// Test exists by asset class
	exists, errExists = store.InstrumentExists(ctx, NewInstrumentQuery().SetAssetClass(ASSET_CLASS_STOCK))
	if errExists != nil {
		t.Fatal("unexpected error checking if instrument exists by asset class:", errExists)
	}
	if !exists {
		t.Fatal("Instrument should exist by asset class")
	}

	// Test exists by symbol like
	exists, errExists = store.InstrumentExists(ctx, NewInstrumentQuery().SetSymbolLike("SL"))
	if errExists != nil {
		t.Fatal("unexpected error checking if instrument exists by symbol like:", errExists)
	}
	if !exists {
		t.Fatal("Instrument should exist by symbol like")
	}

	// Test not exists with non-existent ID
	exists, errExists = store.InstrumentExists(ctx, NewInstrumentQuery().SetID("non-existent-id"))
	if errExists != nil {
		t.Fatal("unexpected error checking if non-existent instrument exists:", errExists)
	}
	if exists {
		t.Fatal("Instrument should not exist with non-existent ID")
	}

	// Test not exists with non-existent symbol
	exists, errExists = store.InstrumentExists(ctx, NewInstrumentQuery().SetSymbol("NONEXIST"))
	if errExists != nil {
		t.Fatal("unexpected error checking if instrument with non-existent symbol exists:", errExists)
	}
	if exists {
		t.Fatal("Instrument should not exist with non-existent symbol")
	}
}

func TestStoreInstrumentCount(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Clear any existing instruments
	clearInstruments(t, store)

	ctx := context.Background()

	// Create instruments with different properties
	instrument1 := NewInstrument().
		SetSymbol("AAPL").
		SetExchange("NASDAQ").
		SetAssetClass(ASSET_CLASS_STOCK)

	instrument2 := NewInstrument().
		SetSymbol("MSFT").
		SetExchange("NASDAQ").
		SetAssetClass(ASSET_CLASS_STOCK)

	instrument3 := NewInstrument().
		SetSymbol("GBP/USD").
		SetExchange("FOREX").
		SetAssetClass(ASSET_CLASS_CURRENCY)

	err = store.InstrumentCreate(ctx, instrument1)
	if err != nil {
		t.Fatal("unexpected error creating instrument1:", err)
	}

	err = store.InstrumentCreate(ctx, instrument2)
	if err != nil {
		t.Fatal("unexpected error creating instrument2:", err)
	}

	err = store.InstrumentCreate(ctx, instrument3)
	if err != nil {
		t.Fatal("unexpected error creating instrument3:", err)
	}

	// Test count all
	count, errCount := store.InstrumentCount(ctx, NewInstrumentQuery())
	if errCount != nil {
		t.Fatal("unexpected error counting all instruments:", errCount)
	}
	if count != 3 {
		t.Fatal("Count should be 3, got:", count)
	}

	// Test count by exchange (partial)
	count, errCount = store.InstrumentCount(ctx, NewInstrumentQuery().SetExchange("NASDAQ"))
	if errCount != nil {
		t.Fatal("unexpected error counting instruments by exchange:", errCount)
	}
	if count != 2 {
		t.Fatal("Count should be 2 for NASDAQ exchange, got:", count)
	}

	// Test count with ID filter (single result)
	count, errCount = store.InstrumentCount(ctx, NewInstrumentQuery().SetID(instrument1.ID()))
	if errCount != nil {
		t.Fatal("unexpected error counting instruments by ID:", errCount)
	}
	if count != 1 {
		t.Fatal("Count should be 1 for ID filter, got:", count)
	}

	// Test count by asset class
	count, errCount = store.InstrumentCount(ctx, NewInstrumentQuery().SetAssetClass(ASSET_CLASS_CURRENCY))
	if errCount != nil {
		t.Fatal("unexpected error counting instruments by asset class:", errCount)
	}
	if count != 1 {
		t.Fatal("Count should be 1 for CURRENCY asset class, got:", count)
	}

	// Test count with multiple IDs
	count, errCount = store.InstrumentCount(ctx, NewInstrumentQuery().SetIDIn([]string{instrument1.ID(), instrument3.ID()}))
	if errCount != nil {
		t.Fatal("unexpected error counting instruments by multiple IDs:", errCount)
	}
	if count != 2 {
		t.Fatal("Count should be 2 for multiple IDs filter, got:", count)
	}
}

func TestStoreInstrumentList(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Clear any existing instruments
	clearInstruments(t, store)

	ctx := context.Background()

	// Create instruments with different properties
	instrument1 := NewInstrument().
		SetSymbol("AAPL").
		SetExchange("NASDAQ").
		SetAssetClass(ASSET_CLASS_STOCK)

	instrument2 := NewInstrument().
		SetSymbol("MSFT").
		SetExchange("NASDAQ").
		SetAssetClass(ASSET_CLASS_STOCK)

	instrument3 := NewInstrument().
		SetSymbol("GBP/USD").
		SetExchange("FOREX").
		SetAssetClass(ASSET_CLASS_CURRENCY)

	err = store.InstrumentCreate(ctx, instrument1)
	if err != nil {
		t.Fatal("unexpected error creating instrument1:", err)
	}

	err = store.InstrumentCreate(ctx, instrument2)
	if err != nil {
		t.Fatal("unexpected error creating instrument2:", err)
	}

	err = store.InstrumentCreate(ctx, instrument3)
	if err != nil {
		t.Fatal("unexpected error creating instrument3:", err)
	}

	// Test list all
	instruments, errList := store.InstrumentList(ctx, NewInstrumentQuery())
	if errList != nil {
		t.Fatal("unexpected error listing all instruments:", errList)
	}
	if len(instruments) != 3 {
		t.Fatal("Should list 3 instruments, got:", len(instruments))
	}

	// Test list with limit
	instruments, errList = store.InstrumentList(ctx, NewInstrumentQuery().SetLimit(2))
	if errList != nil {
		t.Fatal("unexpected error listing instruments with limit:", errList)
	}
	if len(instruments) != 2 {
		t.Fatal("Should list 2 instruments with limit, got:", len(instruments))
	}

	// Test list with exchange filter
	instruments, errList = store.InstrumentList(ctx, NewInstrumentQuery().SetExchange("NASDAQ"))
	if errList != nil {
		t.Fatal("unexpected error listing instruments by exchange:", errList)
	}
	if len(instruments) != 2 {
		t.Fatal("Should list 2 instruments in NASDAQ exchange, got:", len(instruments))
	}

	// Test list with asset class filter
	instruments, errList = store.InstrumentList(ctx, NewInstrumentQuery().SetAssetClass(ASSET_CLASS_CURRENCY))
	if errList != nil {
		t.Fatal("unexpected error listing instruments by asset class:", errList)
	}
	if len(instruments) != 1 {
		t.Fatal("Should list 1 instrument with CURRENCY asset class, got:", len(instruments))
	}
	if instruments[0].AssetClass() != ASSET_CLASS_CURRENCY {
		t.Fatal("Asset class should be CURRENCY, got:", instruments[0].AssetClass())
	}

	// Test list with symbol filter
	instruments, errList = store.InstrumentList(ctx, NewInstrumentQuery().SetSymbol("AAPL"))
	if errList != nil {
		t.Fatal("unexpected error listing instruments by symbol:", errList)
	}
	if len(instruments) != 1 {
		t.Fatal("Should list 1 instrument with symbol AAPL, got:", len(instruments))
	}
	if instruments[0].Symbol() != "AAPL" {
		t.Fatal("Symbol should be AAPL, got:", instruments[0].Symbol())
	}

	// Test list with symbol like filter
	instruments, errList = store.InstrumentList(ctx, NewInstrumentQuery().SetSymbolLike("MS"))
	if errList != nil {
		t.Fatal("unexpected error listing instruments by symbol like:", errList)
	}
	if len(instruments) != 1 {
		t.Fatal("Should list 1 instrument with symbol like MS, got:", len(instruments))
	}
	if instruments[0].Symbol() != "MSFT" {
		t.Fatal("Symbol should be MSFT, got:", instruments[0].Symbol())
	}

	// Test list with ordering (ascending by symbol)
	instruments, errList = store.InstrumentList(ctx, NewInstrumentQuery().
		SetOrderBy(COLUMN_SYMBOL).
		SetSortDirection("ASC"))
	if errList != nil {
		t.Fatal("unexpected error listing instruments with ordering:", errList)
	}
	if len(instruments) != 3 {
		t.Fatal("Should list 3 instruments with ordering, got:", len(instruments))
	}
	if instruments[0].Symbol() != "AAPL" {
		t.Fatal("First instrument should be AAPL, got:", instruments[0].Symbol())
	}

	// Test list with offset and limit
	instruments, errList = store.InstrumentList(ctx, NewInstrumentQuery().
		SetOrderBy(COLUMN_SYMBOL).
		SetSortDirection("ASC").
		SetOffset(1).
		SetLimit(10))
	if errList != nil {
		t.Fatal("unexpected error listing instruments with offset:", errList)
	}
	if len(instruments) != 2 {
		t.Fatal("Should list 2 instruments with offset, got:", len(instruments))
	}
	// Can't reliably check for specific symbols here since we don't know the exact order
}

func TestStoreInstrumentUpdate(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Clear any existing instruments
	clearInstruments(t, store)

	ctx := context.Background()

	// Create an instrument
	instrument := NewInstrument().
		SetSymbol("AAPL").
		SetExchange("NASDAQ").
		SetAssetClass(ASSET_CLASS_STOCK)

	err = store.InstrumentCreate(ctx, instrument)
	if err != nil {
		t.Fatal("unexpected error creating instrument:", err)
	}

	// Update the instrument values
	instrument.SetSymbol("AAPL.US")
	instrument.SetExchange("NYSE")
	instrument.SetAssetClass(ASSET_CLASS_ETF)

	// Save the updates
	err = store.InstrumentUpdate(ctx, instrument)
	if err != nil {
		t.Fatal("unexpected error updating instrument:", err)
	}

	// Retrieve the updated instrument
	updatedInstrument, errFind := store.InstrumentFindByID(ctx, instrument.ID())
	if errFind != nil {
		t.Fatal("unexpected error finding updated instrument:", errFind)
	}

	if updatedInstrument == nil {
		t.Fatal("Updated instrument should exist")
	}

	// Check if the values were updated correctly
	if updatedInstrument.Symbol() != "AAPL.US" {
		t.Fatal("Instrument symbol should be updated to AAPL.US, got:", updatedInstrument.Symbol())
	}

	if updatedInstrument.Exchange() != "NYSE" {
		t.Fatal("Instrument exchange should be updated to NYSE, got:", updatedInstrument.Exchange())
	}

	if updatedInstrument.AssetClass() != ASSET_CLASS_ETF {
		t.Fatal("Instrument asset class should be updated to ETF, got:", updatedInstrument.AssetClass())
	}

	// Check that the ID remains unchanged
	if updatedInstrument.ID() != instrument.ID() {
		t.Fatal("Instrument ID should remain unchanged")
	}
}
