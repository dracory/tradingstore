package tradingstore

type InstrumentInterface interface {
	// from dataobject
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	// methods

	// setters and getters
	ID() string
	SetID(id string) InstrumentInterface

	GetSymbol() string
	SetSymbol(symbol string) InstrumentInterface

	GetExchange() string
	SetExchange(exchange string) InstrumentInterface

	GetAssetClass() string
	SetAssetClass(assetClass string) InstrumentInterface

	GetDescription() string
	SetDescription(description string) InstrumentInterface

	GetTimeframes() []string
	SetTimeframes(timeframes []string) InstrumentInterface
}
