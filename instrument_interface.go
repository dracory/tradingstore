package tradingstore

import "github.com/dromara/carbon/v2"

type InstrumentInterface interface {
	// from dataobject
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	// methods

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
}
