package tradingstore

import (
	"strings"

	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/uid"
)

// == CLASS ====================================================================

// Instrument represents a financial instrument data object
type instrumentImplementation struct {
	dataobject.DataObject
}

// == CONSTRUCTORS =============================================================

func NewInstrument() InstrumentInterface {
	o := (&instrumentImplementation{}).
		SetID(uid.HumanUid())

	return o
}

func NewInstrumentFromExistingData(data map[string]string) InstrumentInterface {
	o := &instrumentImplementation{}
	o.Hydrate(data)
	return o
}

var _ InstrumentInterface = (*instrumentImplementation)(nil)

// == SETTERS & GETTERS ========================================================

func (instrument *instrumentImplementation) ID() string {
	return instrument.Get(COLUMN_ID)
}

func (instrument *instrumentImplementation) SetID(id string) InstrumentInterface {
	instrument.Set(COLUMN_ID, id)
	return instrument
}

func (instrument *instrumentImplementation) Symbol() string {
	return instrument.Get(COLUMN_SYMBOL)
}

func (instrument *instrumentImplementation) SetSymbol(symbol string) InstrumentInterface {
	instrument.Set(COLUMN_SYMBOL, symbol)
	return instrument
}

func (instrument *instrumentImplementation) Exchange() string {
	return instrument.Get(COLUMN_EXCHANGE)
}

func (instrument *instrumentImplementation) SetExchange(exchange string) InstrumentInterface {
	instrument.Set(COLUMN_EXCHANGE, exchange)
	return instrument
}

func (instrument *instrumentImplementation) AssetClass() string {
	return instrument.Get(COLUMN_ASSET_CLASS)
}

func (instrument *instrumentImplementation) SetAssetClass(assetClass string) InstrumentInterface {
	instrument.Set(COLUMN_ASSET_CLASS, assetClass)
	return instrument
}

func (instrument *instrumentImplementation) Description() string {
	return instrument.Get(COLUMN_DESCRIPTION)
}

func (instrument *instrumentImplementation) SetDescription(description string) InstrumentInterface {
	instrument.Set(COLUMN_DESCRIPTION, description)
	return instrument
}

func (instrument *instrumentImplementation) CreatedAt() string {
	return instrument.Get(COLUMN_CREATED_AT)
}

func (instrument *instrumentImplementation) SetCreatedAt(createdAt string) InstrumentInterface {
	instrument.Set(COLUMN_CREATED_AT, createdAt)
	return instrument
}

func (instrument *instrumentImplementation) SoftDeletedAt() string {
	return instrument.Get(COLUMN_SOFT_DELETED_AT)
}

func (instrument *instrumentImplementation) SetSoftDeletedAt(softDeletedAt string) InstrumentInterface {
	instrument.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return instrument
}

func (instrument *instrumentImplementation) Timeframes() []string {
	timeframes := instrument.Get(COLUMN_TIMEFRAMES)
	if timeframes == "" {
		return []string{}
	}
	return strings.Split(timeframes, ",")
}

func (instrument *instrumentImplementation) SetTimeframes(timeframes []string) InstrumentInterface {
	instrument.Set(COLUMN_TIMEFRAMES, strings.Join(timeframes, ","))
	return instrument
}

func (instrument *instrumentImplementation) UpdatedAt() string {
	return instrument.Get(COLUMN_UPDATED_AT)
}

func (instrument *instrumentImplementation) SetUpdatedAt(updatedAt string) InstrumentInterface {
	instrument.Set(COLUMN_UPDATED_AT, updatedAt)
	return instrument
}
