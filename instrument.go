package tradingstore

import (
	"encoding/json"
	"strings"

	"github.com/dromara/carbon/v2"
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

	// Default values
	o.SetName("")
	o.SetStatus(INSTRUMENT_STATUS_DRAFT)
	o.SetAssetClass(ASSET_CLASS_UNKNOWN)
	o.SetExchange("")
	o.SetDescription("")
	o.SetMemo("")
	o.SetMetas(map[string]string{})
	o.SetTimeframes([]string{})
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	o.SetSoftDeletedAt(carbon.MaxValue().ToDateTimeString())

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

func (instrument *instrumentImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(instrument.CreatedAt(), carbon.UTC)
}

func (instrument *instrumentImplementation) SetCreatedAt(createdAt string) InstrumentInterface {
	instrument.Set(COLUMN_CREATED_AT, createdAt)
	return instrument
}

func (instrument *instrumentImplementation) Memo() string {
	return instrument.Get(COLUMN_MEMO)
}

func (instrument *instrumentImplementation) SetMemo(memo string) InstrumentInterface {
	instrument.Set(COLUMN_MEMO, memo)
	return instrument
}

func (instrument *instrumentImplementation) Meta(key string) (string, error) {
	metas, err := instrument.Metas()
	if err != nil {
		return "", err
	}

	value, ok := metas[key]
	if !ok {
		return "", nil
	}

	return value, nil
}

func (instrument *instrumentImplementation) SetMeta(key string, value string) error {
	metas, err := instrument.Metas()
	if err != nil {
		return err
	}

	metas[key] = value
	return instrument.SetMetas(metas)
}

func (instrument *instrumentImplementation) DeleteMeta(key string) error {
	metas, err := instrument.Metas()
	if err != nil {
		return err
	}

	delete(metas, key)
	return instrument.SetMetas(metas)
}

func (instrument *instrumentImplementation) Metas() (map[string]string, error) {
	metasStr := instrument.Get(COLUMN_METAS)
	if metasStr == "" {
		return map[string]string{}, nil
	}

	var metas map[string]string
	err := json.Unmarshal([]byte(metasStr), &metas)
	if err != nil {
		return map[string]string{}, err
	}

	return metas, nil
}

func (instrument *instrumentImplementation) SetMetas(metas map[string]string) error {
	metasBytes, err := json.Marshal(metas)
	if err != nil {
		return err
	}

	instrument.Set(COLUMN_METAS, string(metasBytes))
	return nil
}

func (instrument *instrumentImplementation) Name() string {
	return instrument.Get(COLUMN_NAME)
}

func (instrument *instrumentImplementation) SetName(name string) InstrumentInterface {
	instrument.Set(COLUMN_NAME, name)
	return instrument
}

func (instrument *instrumentImplementation) SoftDeletedAt() string {
	return instrument.Get(COLUMN_SOFT_DELETED_AT)
}

func (instrument *instrumentImplementation) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(instrument.SoftDeletedAt(), carbon.UTC)
}

func (instrument *instrumentImplementation) SetSoftDeletedAt(softDeletedAt string) InstrumentInterface {
	instrument.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return instrument
}

func (instrument *instrumentImplementation) Status() string {
	return instrument.Get(COLUMN_STATUS)
}

func (instrument *instrumentImplementation) SetStatus(status string) InstrumentInterface {
	instrument.Set(COLUMN_STATUS, status)
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

func (instrument *instrumentImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(instrument.UpdatedAt(), carbon.UTC)
}

func (instrument *instrumentImplementation) SetUpdatedAt(updatedAt string) InstrumentInterface {
	instrument.Set(COLUMN_UPDATED_AT, updatedAt)
	return instrument
}
