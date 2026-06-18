package tradingstore

import (
	"github.com/dracory/neat/database/orm"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

// PriceInterface defines the interface for an OHLCV price record.
type PriceInterface interface {
	ID() string
	SetID(id string) PriceInterface

	Close() string
	CloseFloat() float64
	SetClose(close string) PriceInterface

	High() string
	HighFloat() float64
	SetHigh(high string) PriceInterface

	Low() string
	LowFloat() float64
	SetLow(low string) PriceInterface

	Open() string
	OpenFloat() float64
	SetOpen(open string) PriceInterface

	Time() string
	TimeCarbon() *carbon.Carbon
	SetTime(time string) PriceInterface

	Volume() string
	VolumeFloat() float64
	SetVolume(volume string) PriceInterface
}

var _ PriceInterface = (*priceImplementation)(nil)

// == TYPE =====================================================================

type priceImplementation struct {
	orm.ShortID

	OpenField   string `db:"open"`
	HighField   string `db:"high"`
	LowField    string `db:"low"`
	CloseField  string `db:"close"`
	VolumeField string `db:"volume"`
	TimeField   string `db:"time"`
}

// == CONSTRUCTORS =============================================================

// NewPrice creates a new price with default values.
func NewPrice() PriceInterface {
	o := &priceImplementation{}
	o.SetID(neatuid.GenerateShortID())
	return o
}

// NewPriceFromExistingData hydrates a price from a raw column map.
func NewPriceFromExistingData(data map[string]string) PriceInterface {
	o := &priceImplementation{}
	o.SetID(data[COLUMN_ID])
	o.SetOpen(data[COLUMN_OPEN])
	o.SetHigh(data[COLUMN_HIGH])
	o.SetLow(data[COLUMN_LOW])
	o.SetClose(data[COLUMN_CLOSE])
	o.SetVolume(data[COLUMN_VOLUME])
	if v, ok := data[COLUMN_TIME]; ok {
		o.SetTime(v)
	}
	return o
}

// == SETTERS & GETTERS ========================================================

func (price *priceImplementation) Close() string {
	return price.CloseField
}

func (price *priceImplementation) CloseFloat() float64 {
	return cast.ToFloat64(price.Close())
}

func (price *priceImplementation) SetClose(close string) PriceInterface {
	price.CloseField = close
	return price
}

func (price *priceImplementation) High() string {
	return price.HighField
}

func (price *priceImplementation) HighFloat() float64 {
	return cast.ToFloat64(price.High())
}

func (price *priceImplementation) SetHigh(high string) PriceInterface {
	price.HighField = high
	return price
}

func (price *priceImplementation) Low() string {
	return price.LowField
}

func (price *priceImplementation) LowFloat() float64 {
	return cast.ToFloat64(price.Low())
}

func (price *priceImplementation) SetLow(low string) PriceInterface {
	price.LowField = low
	return price
}

func (price *priceImplementation) ID() string {
	return price.ShortID.ID
}

func (price *priceImplementation) SetID(id string) PriceInterface {
	price.ShortID.ID = id
	return price
}

func (price *priceImplementation) Open() string {
	return price.OpenField
}

func (price *priceImplementation) OpenFloat() float64 {
	return cast.ToFloat64(price.Open())
}

func (price *priceImplementation) SetOpen(open string) PriceInterface {
	price.OpenField = open
	return price
}

// Time returns the time as an ISO8601 formatted string.
func (price *priceImplementation) Time() string {
	return price.TimeField
}

func (price *priceImplementation) TimeCarbon() *carbon.Carbon {
	return carbon.Parse(price.Time(), carbon.UTC)
}

// SetTime sets the time for a Price, must be in UTC.
// The time is stored as an ISO8601 formatted string.
func (price *priceImplementation) SetTime(timeUtc string) PriceInterface {
	price.TimeField = carbon.Parse(timeUtc, carbon.UTC).ToIso8601ZuluString()
	return price
}

func (price *priceImplementation) Volume() string {
	return price.VolumeField
}

func (price *priceImplementation) VolumeFloat() float64 {
	return cast.ToFloat64(price.Volume())
}

func (price *priceImplementation) SetVolume(volume string) PriceInterface {
	price.VolumeField = volume
	return price
}
