package tradingstore

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
)

// InstrumentInterface defines the interface for a financial instrument.
type InstrumentInterface interface {
	// setters and getters
	AssetClass() string
	SetAssetClass(assetClass string) InstrumentInterface

	Exchange() string
	SetExchange(exchange string) InstrumentInterface

	Description() string
	SetDescription(description string) InstrumentInterface

	ID() string
	SetID(id string) InstrumentInterface

	Meta(key string) (string, error)
	SetMeta(key string, value string) error
	DeleteMeta(key string) error

	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error

	Memo() string
	SetMemo(memo string) InstrumentInterface

	Name() string
	SetName(name string) InstrumentInterface

	Status() string
	SetStatus(status string) InstrumentInterface

	Symbol() string
	SetSymbol(symbol string) InstrumentInterface

	Timeframes() []string
	SetTimeframes(timeframes []string) InstrumentInterface

	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) InstrumentInterface

	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) InstrumentInterface

	SoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(softDeletedAt string) InstrumentInterface

	IsSoftDeleted() bool
}

var _ InstrumentInterface = (*instrumentImplementation)(nil)

// == TYPE =====================================================================

type instrumentImplementation struct {
	orm.ShortID

	SymbolField      string `db:"symbol"`
	ExchangeField    string `db:"exchange"`
	AssetClassField  string `db:"asset_class"`
	NameField        string `db:"name"`
	StatusField      string `db:"status"`
	DescriptionField string `db:"description"`
	MemoField        string `db:"memo"`
	MetasField       string `db:"metas"`
	TimeframesField  string `db:"timeframes"`

	CreatedAtField orm.CreatedAt
	UpdatedAtField orm.UpdatedAt
	soft_delete.SoftDeletesMaxDate
}

// == CONSTRUCTORS =============================================================

// NewInstrument creates a new instrument with default values.
func NewInstrument() InstrumentInterface {
	o := &instrumentImplementation{}
	o.SetID(neatuid.GenerateShortID())
	o.SetName("")
	o.SetStatus(INSTRUMENT_STATUS_DRAFT)
	o.SetAssetClass(ASSET_CLASS_UNKNOWN)
	o.SetExchange("")
	o.SetDescription("")
	o.SetMemo("")
	o.SetMetas(map[string]string{})
	o.SetTimeframes([]string{})
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetSoftDeletedAt(MAX_DATETIME)
	return o
}

// NewInstrumentFromExistingData hydrates an instrument from a raw column map.
func NewInstrumentFromExistingData(data map[string]string) InstrumentInterface {
	o := &instrumentImplementation{}
	o.SetID(data[COLUMN_ID])
	o.SetSymbol(data[COLUMN_SYMBOL])
	o.SetExchange(data[COLUMN_EXCHANGE])
	o.SetAssetClass(data[COLUMN_ASSET_CLASS])
	o.SetDescription(data[COLUMN_DESCRIPTION])
	o.SetMemo(data[COLUMN_MEMO])
	if v, ok := data[COLUMN_METAS]; ok {
		o.MetasField = v
	}
	if v, ok := data[COLUMN_TIMEFRAMES]; ok {
		o.TimeframesField = v
	}
	if v, ok := data[COLUMN_NAME]; ok {
		o.NameField = v
	}
	if v, ok := data[COLUMN_STATUS]; ok {
		o.StatusField = v
	}
	if v, ok := data[COLUMN_CREATED_AT]; ok {
		o.SetCreatedAt(v)
	}
	if v, ok := data[COLUMN_UPDATED_AT]; ok {
		o.SetUpdatedAt(v)
	}
	if v, ok := data[COLUMN_SOFT_DELETED_AT]; ok {
		o.SetSoftDeletedAt(v)
	}
	return o
}

// == METHODS ==================================================================

// IsSoftDeleted returns true if the instrument has been soft deleted.
func (instrument *instrumentImplementation) IsSoftDeleted() bool {
	return instrument.SoftDeletesMaxDate.IsSoftDeleted()
}

// == SETTERS & GETTERS ========================================================

func (instrument *instrumentImplementation) ID() string {
	return instrument.ShortID.ID
}

func (instrument *instrumentImplementation) SetID(id string) InstrumentInterface {
	instrument.ShortID.ID = id
	return instrument
}

func (instrument *instrumentImplementation) Symbol() string {
	return instrument.SymbolField
}

func (instrument *instrumentImplementation) SetSymbol(symbol string) InstrumentInterface {
	instrument.SymbolField = symbol
	return instrument
}

func (instrument *instrumentImplementation) Exchange() string {
	return instrument.ExchangeField
}

func (instrument *instrumentImplementation) SetExchange(exchange string) InstrumentInterface {
	instrument.ExchangeField = exchange
	return instrument
}

func (instrument *instrumentImplementation) AssetClass() string {
	return instrument.AssetClassField
}

func (instrument *instrumentImplementation) SetAssetClass(assetClass string) InstrumentInterface {
	instrument.AssetClassField = assetClass
	return instrument
}

func (instrument *instrumentImplementation) Description() string {
	return instrument.DescriptionField
}

func (instrument *instrumentImplementation) SetDescription(description string) InstrumentInterface {
	instrument.DescriptionField = description
	return instrument
}

func (instrument *instrumentImplementation) CreatedAt() string {
	if instrument.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(instrument.CreatedAtField.CreatedAt).ToDateTimeString()
}

func (instrument *instrumentImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(instrument.CreatedAtField.CreatedAt)
}

func (instrument *instrumentImplementation) SetCreatedAt(createdAt string) InstrumentInterface {
	if createdAt == "" {
		instrument.CreatedAtField.CreatedAt = time.Time{}
		return instrument
	}
	instrument.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return instrument
}

func (instrument *instrumentImplementation) Memo() string {
	return instrument.MemoField
}

func (instrument *instrumentImplementation) SetMemo(memo string) InstrumentInterface {
	instrument.MemoField = memo
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
	if instrument.MetasField == "" {
		return map[string]string{}, nil
	}
	var metas map[string]string
	err := json.Unmarshal([]byte(instrument.MetasField), &metas)
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
	instrument.MetasField = string(metasBytes)
	return nil
}

func (instrument *instrumentImplementation) Name() string {
	return instrument.NameField
}

func (instrument *instrumentImplementation) SetName(name string) InstrumentInterface {
	instrument.NameField = name
	return instrument
}

func (instrument *instrumentImplementation) SoftDeletedAt() string {
	if instrument.SoftDeletesMaxDate.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(instrument.SoftDeletesMaxDate.SoftDeletedAt).ToDateTimeString()
}

func (instrument *instrumentImplementation) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(instrument.SoftDeletesMaxDate.SoftDeletedAt)
}

func (instrument *instrumentImplementation) SetSoftDeletedAt(softDeletedAt string) InstrumentInterface {
	if softDeletedAt == "" {
		instrument.SoftDeletesMaxDate.SoftDeletedAt = time.Time{}
		return instrument
	}
	instrument.SoftDeletesMaxDate.SoftDeletedAt = carbon.Parse(softDeletedAt, carbon.UTC).StdTime()
	return instrument
}

func (instrument *instrumentImplementation) Status() string {
	return instrument.StatusField
}

func (instrument *instrumentImplementation) SetStatus(status string) InstrumentInterface {
	instrument.StatusField = status
	return instrument
}

func (instrument *instrumentImplementation) Timeframes() []string {
	if instrument.TimeframesField == "" {
		return []string{}
	}
	return strings.Split(instrument.TimeframesField, ",")
}

func (instrument *instrumentImplementation) SetTimeframes(timeframes []string) InstrumentInterface {
	instrument.TimeframesField = strings.Join(timeframes, ",")
	return instrument
}

func (instrument *instrumentImplementation) UpdatedAt() string {
	if instrument.UpdatedAtField.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(instrument.UpdatedAtField.UpdatedAt).ToDateTimeString()
}

func (instrument *instrumentImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(instrument.UpdatedAtField.UpdatedAt)
}

func (instrument *instrumentImplementation) SetUpdatedAt(updatedAt string) InstrumentInterface {
	if updatedAt == "" {
		instrument.UpdatedAtField.UpdatedAt = time.Time{}
		return instrument
	}
	instrument.UpdatedAtField.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return instrument
}
