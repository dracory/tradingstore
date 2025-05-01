package tradingstore_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/dracory/tradingstore" // Import the package under test
)

// Constants for date formats and default values used in tests
const (
	testDateFormat = "2006-01-02 15:04:05"
	maxDateString  = "9999-12-31 23:59:59" // Corresponds to carbon.MaxValue().ToDateString()
)

// Helper to parse time and fail test on error
func mustParseTime(t *testing.T, layout, value string) time.Time {
	t.Helper()
	parsedTime, err := time.Parse(layout, value)
	if err != nil {
		t.Fatalf("Failed to parse time string %q with layout %q: %v", value, layout, err)
	}
	return parsedTime.UTC() // Ensure UTC
}

func TestNewInstrument(t *testing.T) {
	instrument := tradingstore.NewInstrument()

	if instrument == nil {
		t.Fatal("NewInstrument should not return nil")
	}
	if instrument.ID() == "" {
		t.Error("ID should be set")
	}
	if got := instrument.Description(); got != "" {
		t.Errorf("Expected Description to be empty, got %q", got)
	}
	if got := instrument.Memo(); got != "" {
		t.Errorf("Expected Memo to be empty, got %q", got)
	}

	metas, err := instrument.Metas()
	if err != nil {
		t.Errorf("Metas() returned an unexpected error for a new instrument: %v", err)
	}
	if len(metas) != 0 {
		t.Errorf("Expected Metas to be empty, got %v", metas)
	}
	if got := instrument.Timeframes(); len(got) != 0 {
		t.Errorf("Expected Timeframes to be empty, got %v", got)
	}

	// Check timestamps are recent and in UTC
	now := time.Now().UTC()
	createdAtStr := instrument.CreatedAt()
	updatedAtStr := instrument.UpdatedAt()
	softDeletedAtStr := instrument.SoftDeletedAt()

	// Check string format is today's date
	expectedDateStr := now.Format(testDateFormat)
	if createdAtStr != expectedDateStr {
		t.Errorf("Expected CreatedAt string %q, got %q", expectedDateStr, createdAtStr)
	}
	if updatedAtStr != expectedDateStr {
		t.Errorf("Expected UpdatedAt string %q, got %q", expectedDateStr, updatedAtStr)
	}

	// Check Carbon conversion is recent
	createdAt := instrument.CreatedAtCarbon().StdTime()
	updatedAt := instrument.UpdatedAtCarbon().StdTime()

	if diff := now.Sub(createdAt); diff < 0 || diff > 2*time.Second {
		t.Errorf("CreatedAt (%v) is not within 2 seconds of now (%v)", createdAt, now)
	}

	if diff := now.Sub(updatedAt); diff < 0 || diff > 2*time.Second {
		t.Errorf("UpdatedAt (%v) is not within 2 seconds of now (%v)", updatedAt, now)
	}

	// Check SoftDeletedAt
	if softDeletedAtStr != maxDateString {
		t.Errorf("Expected SoftDeletedAt string %q, got %q", maxDateString, softDeletedAtStr)
	}
	softDeletedAt := instrument.SoftDeletedAtCarbon().StdTime()
	expectedSoftDeletedAt := mustParseTime(t, testDateFormat, maxDateString)
	if !softDeletedAt.Equal(expectedSoftDeletedAt) {
		t.Errorf("Expected SoftDeletedAtCarbon to be %v, got %v", expectedSoftDeletedAt, softDeletedAt)
	}
}

func TestNewInstrumentFromExistingData(t *testing.T) {
	testID := "instr_123"
	testSymbol := "TEST"
	testExchange := "TESTEX"                           // Assuming Exchange is a string type
	testAssetClass := tradingstore.ASSET_CLASS_UNKNOWN // Use exported constant
	testDesc := "Test Description"
	testMemo := "Test Memo"
	testCreatedAt := "2023-01-01 12:34:56"
	testUpdatedAt := "2023-01-02 12:34:56"
	testSoftDeletedAt := maxDateString // Use constant for consistency
	testTimeframes := "1min,5min,15min"
	testMetasJSON := `{"key1":"value1","key2":"value2"}`

	data := map[string]string{
		tradingstore.COLUMN_ID:              testID,
		tradingstore.COLUMN_SYMBOL:          testSymbol,
		tradingstore.COLUMN_EXCHANGE:        testExchange,
		tradingstore.COLUMN_ASSET_CLASS:     testAssetClass,
		tradingstore.COLUMN_DESCRIPTION:     testDesc,
		tradingstore.COLUMN_MEMO:            testMemo,
		tradingstore.COLUMN_CREATED_AT:      testCreatedAt,
		tradingstore.COLUMN_UPDATED_AT:      testUpdatedAt,
		tradingstore.COLUMN_SOFT_DELETED_AT: testSoftDeletedAt,
		tradingstore.COLUMN_TIMEFRAMES:      testTimeframes,
		tradingstore.COLUMN_METAS:           testMetasJSON, // Note: COLUMN_METAS is used in implementation
	}

	instrument := tradingstore.NewInstrumentFromExistingData(data)

	if instrument == nil {
		t.Fatal("NewInstrumentFromExistingData should not return nil")
	}
	if got := instrument.ID(); got != testID {
		t.Errorf("Expected ID %q, got %q", testID, got)
	}
	if got := instrument.Symbol(); got != testSymbol {
		t.Errorf("Expected Symbol %q, got %q", testSymbol, got)
	}
	if got := instrument.Exchange(); got != testExchange {
		t.Errorf("Expected Exchange %q, got %q", testExchange, got)
	}
	if got := instrument.AssetClass(); got != testAssetClass {
		t.Errorf("Expected AssetClass %q, got %q", testAssetClass, got)
	}
	if got := instrument.Description(); got != testDesc {
		t.Errorf("Expected Description %q, got %q", testDesc, got)
	}
	if got := instrument.Memo(); got != testMemo {
		t.Errorf("Expected Memo %q, got %q", testMemo, got)
	}
	if got := instrument.CreatedAt(); got != testCreatedAt {
		t.Errorf("Expected CreatedAt %q, got %q", testCreatedAt, got)
	}
	if got := instrument.UpdatedAt(); got != testUpdatedAt {
		t.Errorf("Expected UpdatedAt %q, got %q", testUpdatedAt, got)
	}
	if got := instrument.SoftDeletedAt(); got != testSoftDeletedAt {
		t.Errorf("Expected SoftDeletedAt %q, got %q", testSoftDeletedAt, got)
	}

	expectedTimeframes := []string{"1min", "5min", "15min"}
	if got := instrument.Timeframes(); !reflect.DeepEqual(got, expectedTimeframes) {
		t.Errorf("Expected Timeframes %v, got %v", expectedTimeframes, got)
	}

	metas, err := instrument.Metas()
	if err != nil {
		t.Errorf("Metas() returned an unexpected error: %v", err)
	}
	expectedMetas := map[string]string{"key1": "value1", "key2": "value2"}
	if !reflect.DeepEqual(metas, expectedMetas) {
		t.Errorf("Expected Metas %v, got %v", expectedMetas, metas)
	}

	// Test time conversions (using Carbon methods which should internally parse)
	expectedCreatedAtTime := mustParseTime(t, testDateFormat, testCreatedAt)
	if got := instrument.CreatedAtCarbon().StdTime(); !got.Equal(expectedCreatedAtTime) {
		t.Errorf("Expected CreatedAtCarbon %v, got %v", expectedCreatedAtTime, got)
	}
	expectedUpdatedAtTime := mustParseTime(t, testDateFormat, testUpdatedAt)
	if got := instrument.UpdatedAtCarbon().StdTime(); !got.Equal(expectedUpdatedAtTime) {
		t.Errorf("Expected UpdatedAtCarbon %v, got %v", expectedUpdatedAtTime, got)
	}
	expectedSoftDeletedAtTime := mustParseTime(t, testDateFormat, testSoftDeletedAt)
	if got := instrument.SoftDeletedAtCarbon().StdTime(); !got.Equal(expectedSoftDeletedAtTime) {
		t.Errorf("Expected SoftDeletedAtCarbon %v, got %v", expectedSoftDeletedAtTime, got)
	}
}

func TestInstrumentSettersGetters(t *testing.T) {
	instrument := tradingstore.NewInstrument()

	// ID
	testID := "new_id_456"
	instrument.SetID(testID)
	if got := instrument.ID(); got != testID {
		t.Errorf("SetID/ID failed: expected %q, got %q", testID, got)
	}

	// Symbol
	testSymbol := "NEWSYMB"
	instrument.SetSymbol(testSymbol)
	if got := instrument.Symbol(); got != testSymbol {
		t.Errorf("SetSymbol/Symbol failed: expected %q, got %q", testSymbol, got)
	}

	// Exchange
	testExchange := "NEWEX"
	instrument.SetExchange(testExchange)
	if got := instrument.Exchange(); got != testExchange {
		t.Errorf("SetExchange/Exchange failed: expected %q, got %q", testExchange, got)
	}

	// AssetClass
	testAssetClass := tradingstore.ASSET_CLASS_ETF // Use exported constant
	instrument.SetAssetClass(testAssetClass)
	if got := instrument.AssetClass(); got != testAssetClass {
		t.Errorf("SetAssetClass/AssetClass failed: expected %q, got %q", testAssetClass, got)
	}

	// Description
	testDesc := "New Description"
	instrument.SetDescription(testDesc)
	if got := instrument.Description(); got != testDesc {
		t.Errorf("SetDescription/Description failed: expected %q, got %q", testDesc, got)
	}

	// Memo
	testMemo := "New Memo"
	instrument.SetMemo(testMemo)
	if got := instrument.Memo(); got != testMemo {
		t.Errorf("SetMemo/Memo failed: expected %q, got %q", testMemo, got)
	}

	// CreatedAt
	testCreatedAt := "2024-05-10 12:34:56"
	instrument.SetCreatedAt(testCreatedAt)
	if got := instrument.CreatedAt(); got != testCreatedAt {
		t.Errorf("SetCreatedAt/CreatedAt failed: expected %q, got %q", testCreatedAt, got)
	}
	expectedCreatedAtTime := mustParseTime(t, testDateFormat, testCreatedAt)
	if got := instrument.CreatedAtCarbon().StdTime(); !got.Equal(expectedCreatedAtTime) {
		t.Errorf("SetCreatedAt/CreatedAtCarbon failed: expected %v, got %v", expectedCreatedAtTime, got)
	}

	// UpdatedAt
	testUpdatedAt := "2024-05-11 12:34:56"
	instrument.SetUpdatedAt(testUpdatedAt)
	if got := instrument.UpdatedAt(); got != testUpdatedAt {
		t.Errorf("SetUpdatedAt/UpdatedAt failed: expected %q, got %q", testUpdatedAt, got)
	}
	expectedUpdatedAtTime := mustParseTime(t, testDateFormat, testUpdatedAt)
	if got := instrument.UpdatedAtCarbon().StdTime(); !got.Equal(expectedUpdatedAtTime) {
		t.Errorf("SetUpdatedAt/UpdatedAtCarbon failed: expected %v, got %v", expectedUpdatedAtTime, got)
	}

	// SoftDeletedAt
	testSoftDeletedAt := "2024-05-12 12:34:56"
	instrument.SetSoftDeletedAt(testSoftDeletedAt)
	if got := instrument.SoftDeletedAt(); got != testSoftDeletedAt {
		t.Errorf("SetSoftDeletedAt/SoftDeletedAt failed: expected %q, got %q", testSoftDeletedAt, got)
	}
	expectedSoftDeletedAtTime := mustParseTime(t, testDateFormat, testSoftDeletedAt)
	if got := instrument.SoftDeletedAtCarbon().StdTime(); !got.Equal(expectedSoftDeletedAtTime) {
		t.Errorf("SetSoftDeletedAt/SoftDeletedAtCarbon failed: expected %v, got %v", expectedSoftDeletedAtTime, got)
	}

	// Timeframes
	testTimeframes := []string{"1h", "4h", "1d"}
	instrument.SetTimeframes(testTimeframes)
	if got := instrument.Timeframes(); !reflect.DeepEqual(got, testTimeframes) {
		t.Errorf("SetTimeframes/Timeframes failed: expected %v, got %v", testTimeframes, got)
	}
	// Cannot test instrument.Get(COLUMN_TIMEFRAMES) from _test package

	// Test setting empty timeframes
	instrument.SetTimeframes([]string{})
	if got := instrument.Timeframes(); len(got) != 0 {
		t.Errorf("SetTimeframes (empty) failed: expected empty slice, got %v", got)
	}
	// Cannot test instrument.Get(COLUMN_TIMEFRAMES) from _test package

}

func TestInstrumentMetaMethods(t *testing.T) {
	instrument := tradingstore.NewInstrument()

	// Initial state
	metas, err := instrument.Metas()
	if err != nil {
		t.Fatalf("Initial Metas() failed: %v", err)
	}
	if len(metas) != 0 {
		t.Fatalf("Initial Metas() expected empty map, got %v", metas)
	}
	val, err := instrument.Meta("key1")
	if err != nil {
		t.Fatalf("Initial Meta('key1') failed: %v", err)
	}
	if val != "" {
		t.Fatalf("Initial Meta('key1') expected empty string, got %q", val)
	}

	// SetMeta
	err = instrument.SetMeta("key1", "value1")
	if err != nil {
		t.Fatalf("SetMeta('key1', 'value1') failed: %v", err)
	}
	val, err = instrument.Meta("key1")
	if err != nil {
		t.Fatalf("Meta('key1') after set failed: %v", err)
	}
	if val != "value1" {
		t.Fatalf("Meta('key1') after set expected 'value1', got %q", val)
	}

	// Check internal storage
	// Cannot directly check internal storage via instrument.Get(COLUMN_METAS) from _test package.
	// Rely on testing public methods like Meta() and Metas().

	// Set another meta
	err = instrument.SetMeta("key2", "value2")
	if err != nil {
		t.Fatalf("SetMeta('key2', 'value2') failed: %v", err)
	}

	// Metas
	metas, err = instrument.Metas()
	if err != nil {
		t.Fatalf("Metas() after adding key2 failed: %v", err)
	}
	expectedMetas := map[string]string{"key1": "value1", "key2": "value2"}
	if !reflect.DeepEqual(metas, expectedMetas) {
		t.Fatalf("Metas() mismatch after adding key2: expected %v, got %v", expectedMetas, metas)
	}

	// SetMetas (overwrite)
	newMetas := map[string]string{"key3": "value3", "key4": "value4"}
	err = instrument.SetMetas(newMetas)
	if err != nil {
		t.Fatalf("SetMetas failed: %v", err)
	}
	metas, err = instrument.Metas()
	if err != nil {
		t.Fatalf("Metas() after SetMetas failed: %v", err)
	}
	if !reflect.DeepEqual(metas, newMetas) {
		t.Fatalf("Metas() mismatch after SetMetas: expected %v, got %v", newMetas, metas)
	}

	// DeleteMeta
	err = instrument.DeleteMeta("key3")
	if err != nil {
		t.Fatalf("DeleteMeta('key3') failed: %v", err)
	}
	val, err = instrument.Meta("key3")
	if err != nil {
		t.Fatalf("Meta('key3') after delete failed: %v", err)
	}
	if val != "" {
		t.Fatalf("Meta('key3') after delete expected empty string, got %q", val)
	}
	metas, err = instrument.Metas()
	if err != nil {
		t.Fatalf("Metas() after DeleteMeta failed: %v", err)
	}
	expectedAfterDelete := map[string]string{"key4": "value4"}
	if !reflect.DeepEqual(metas, expectedAfterDelete) {
		t.Fatalf("Metas() mismatch after DeleteMeta: expected %v, got %v", expectedAfterDelete, metas)
	}

	// Delete non-existent key
	err = instrument.DeleteMeta("nonexistent")
	if err != nil {
		t.Fatalf("DeleteMeta('nonexistent') failed: %v", err) // Should not error
	}
	metas, err = instrument.Metas()
	if err != nil {
		t.Fatalf("Metas() after deleting non-existent key failed: %v", err)
	}
	if !reflect.DeepEqual(metas, expectedAfterDelete) { // Should be unchanged
		t.Fatalf("Metas() mismatch after deleting non-existent key: expected %v, got %v", expectedAfterDelete, metas)
	}
}
