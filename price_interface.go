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

	GetClose() string
	GetCloseFloat() float64
	SetClose(close string) PriceInterface

	GetHigh() string
	GetHighFloat() float64
	SetHigh(high string) PriceInterface

	GetLow() string
	GetLowFloat() float64
	SetLow(low string) PriceInterface

	GetOpen() string
	GetOpenFloat() float64
	SetOpen(open string) PriceInterface

	GetTime() string
	GetTimeCarbon() carbon.Carbon
	SetTime(time string) PriceInterface

	GetVolume() string
	GetVolumeFloat() float64
	SetVolume(volume string) PriceInterface
}
