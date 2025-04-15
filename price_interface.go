package tradingstore

import (
	"github.com/dromara/carbon/v2"
)

type PriceInterface interface {
	// from dataobject

	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	// methods

	// setters and getters
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
