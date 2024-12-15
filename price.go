package tradingstore

import (
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/uid"
	"github.com/spf13/cast"
)

// == CLASS ====================================================================

// Price represents an OHLCV data object, for storing the pricing data in the database
type Price struct {
	dataobject.DataObject
}

// == CONSTRUCTORS =============================================================

func NewPrice() PriceInterface {
	o := (&Price{}).
		SetID(uid.HumanUid())

	return o
}

func NewPriceFromExistingData(data map[string]string) PriceInterface {
	o := &Price{}
	o.Hydrate(data)
	return o
}

// == METHODS ==================================================================

// == SETTERS & GETTERS ========================================================

func (price *Price) GetClose() string {
	return price.Get(COLUMN_CLOSE)
}

func (price *Price) GetCloseFloat() float64 {
	return cast.ToFloat64(price.GetClose())
}

func (price *Price) SetClose(close string) PriceInterface {
	price.Set(COLUMN_CLOSE, close)
	return price
}

func (price *Price) GetHigh() string {
	return price.Get(COLUMN_HIGH)
}

func (price *Price) GetHighFloat() float64 {
	return cast.ToFloat64(price.GetHigh())
}

func (price *Price) SetHigh(high string) PriceInterface {
	price.Set(COLUMN_HIGH, high)
	return price
}

func (price *Price) ID() string {
	return price.Get(COLUMN_ID)
}

func (price *Price) SetID(id string) PriceInterface {
	price.Set(COLUMN_ID, id)
	return price
}

// Time returns the time as a Iso8601 formatted string.
//
// Parameters:
// - none
//
// Returns:
// - string: the time in ISO8601 format
func (price *Price) GetTime() string {
	return price.Get(COLUMN_TIME)
}

func (price *Price) GetTimeCarbon() carbon.Carbon {
	return carbon.Parse(price.GetTime(), carbon.UTC)
}

// SetTime sets the time for a Price, must be in UTC.
// The time is stored as an ISO8601 formatted string.
//
// Parameters:
// - timeUtc: time in UTC format
//
// Returns:
// - *Price: the Price
func (price *Price) SetTime(timeUtc string) PriceInterface {
	timeUtc = carbon.Parse(timeUtc, carbon.UTC).ToIso8601ZuluString()
	price.Set(COLUMN_TIME, timeUtc)
	return price
}

func (price *Price) GetLow() string {
	return price.Get(COLUMN_LOW)
}

func (price *Price) GetLowFloat() float64 {
	return cast.ToFloat64(price.GetLow())
}

func (price *Price) SetLow(low string) PriceInterface {
	price.Set(COLUMN_LOW, low)
	return price
}

func (price *Price) GetOpen() string {
	return price.Get(COLUMN_OPEN)
}

func (price *Price) GetOpenFloat() float64 {
	return cast.ToFloat64(price.GetOpen())
}

func (price *Price) SetOpen(open string) PriceInterface {
	price.Set(COLUMN_OPEN, open)
	return price
}

func (price *Price) GetVolume() string {
	return price.Get(COLUMN_VOLUME)
}

func (price *Price) GetVolumeFloat() float64 {
	return cast.ToFloat64(price.GetVolume())
}

func (price *Price) SetVolume(volume string) PriceInterface {
	price.Set(COLUMN_VOLUME, volume)
	return price
}
