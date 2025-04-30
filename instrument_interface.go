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
