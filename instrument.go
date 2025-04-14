package tradingstore

import (
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/uid"
)

// CREATE TABLE financial_instruments (
// 	id INT AUTO_INCREMENT PRIMARY KEY,
// 	asset_class VARCHAR(10), -- 'CURRENCY', 'STOCK', 'INDEX'
// 	symbol VARCHAR(20),
// 	isin VARCHAR(12),
// 	exchange VARCHAR(20)
// );

// == CLASS ====================================================================

// Instrument represents a trading instrument data object for storing in the database
type Instrument struct {
	dataobject.DataObject
}

// == CONSTRUCTORS =============================================================

func NewInstrument() InstrumentInterface {
	o := (&Instrument{}).
		SetID(uid.HumanUid())

	return o
}

func NewInstrumentFromExistingData(data map[string]string) InstrumentInterface {
	o := &Instrument{}
	o.Hydrate(data)
	return o
}

var _ InstrumentInterface = (*Instrument)(nil)

// == SETTERS & GETTERS ========================================================

func (instrument *Instrument) ID() string {
	return instrument.Get(COLUMN_ID)
}

func (instrument *Instrument) SetID(id string) InstrumentInterface {
	instrument.Set(COLUMN_ID, id)
	return instrument
}

func (instrument *Instrument) GetSymbol() string {
	return instrument.Get(COLUMN_SYMBOL)
}

func (instrument *Instrument) SetSymbol(symbol string) InstrumentInterface {
	instrument.Set(COLUMN_SYMBOL, symbol)
	return instrument
}

func (instrument *Instrument) GetExchange() string {
	return instrument.Get(COLUMN_EXCHANGE)
}

func (instrument *Instrument) SetExchange(exchange string) InstrumentInterface {
	instrument.Set(COLUMN_EXCHANGE, exchange)
	return instrument
}

func (instrument *Instrument) GetAssetClass() string {
	return instrument.Get(COLUMN_ASSET_CLASS)
}

func (instrument *Instrument) SetAssetClass(assetClass string) InstrumentInterface {
	instrument.Set(COLUMN_ASSET_CLASS, assetClass)
	return instrument
}

func (instrument *Instrument) GetDescription() string {
	return instrument.Get("description")
}

func (instrument *Instrument) SetDescription(description string) InstrumentInterface {
	instrument.Set("description", description)
	return instrument
}

func (i *Instrument) CreatedAt() string {
	return i.Get(COLUMN_CREATED_AT)
}

func (i *Instrument) SetCreatedAt(createdAt string) *Instrument {
	i.Set(COLUMN_CREATED_AT, createdAt)
	return i
}

func (i *Instrument) DeletedAt() string {
	return i.Get(COLUMN_DELETED_AT)
}

func (i *Instrument) SetDeletedAt(deletedAt string) *Instrument {
	i.Set(COLUMN_DELETED_AT, deletedAt)
	return i
}

func (i *Instrument) UpdatedAt() string {
	return i.Get(COLUMN_UPDATED_AT)
}

func (i *Instrument) SetUpdatedAt(updatedAt string) *Instrument {
	i.Set(COLUMN_UPDATED_AT, updatedAt)
	return i
}
