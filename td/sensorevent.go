package td

/*
#cgo LDFLAGS: -ltelldus-core
#include <telldus-core.h>
#include "events.h"
*/
import "C"
import (
	"strconv"
	"time"
	"unsafe"
)

const (
	TEMPERATURE    = C.TELLSTICK_TEMPERATURE
	HUMIDITY       = C.TELLSTICK_HUMIDITY
	RAIN_RATE      = C.TELLSTICK_RAINRATE
	RAIN_TOTAL     = C.TELLSTICK_RAINTOTAL
	WIND_DIRECTION = C.TELLSTICK_WINDDIRECTION
	WIND_AVERAGE   = C.TELLSTICK_WINDAVERAGE
	WIND_GUST      = C.TELLSTICK_WINDGUST
)

type SensorEvent struct {
	Protocol string
	Model    string
	Id       int
	When     time.Time
	Value    SensorValue
}

type SensorValue interface {
	RawValue() string
}

type TempValue string

func (v TempValue) RawValue() string {
	return string(v)
}

func (v TempValue) Float64() float64 {
	f, _ := strconv.ParseFloat(string(v), 64)
	return f
}

type HumidityValue string

func (v HumidityValue) RawValue() string {
	return string(v)
}

func (v HumidityValue) Int() int {
	i, _ := strconv.Atoi(string(v))
	return i
}

type GenericValue struct {
	Value string
	Type  int
}

func (v GenericValue) RawValue() string {
	return v.Value
}

type sensorSubscription struct {
	cs   chan SensorEvent
	cbId C.int
}

func (c *sensorSubscription) Close() error {
	C.tdUnregisterCallback(c.cbId)
	close(c.cs)

	return nil
}

func toSensorValue(value string, dataType int) SensorValue {
	switch dataType {
	case TEMPERATURE:
		return TempValue(value)
	case HUMIDITY:
		return HumidityValue(value)
	}
	return GenericValue{value, dataType}
}

//export newSensorEvent
func newSensorEvent(
	protocol *C.char,
	model *C.char,
	sensorId C.int,
	dataType C.int,
	value *C.char,
	ts C.int,
	callbackId C.int,
	context unsafe.Pointer) {

	cs := (*sensorSubscription)(context).cs

	cs <- SensorEvent{
		Protocol: C.GoString(protocol),
		Model:    C.GoString(model),
		Id:       int(sensorId),
		Value:    toSensorValue(C.GoString(value), int(dataType)),
		When:     time.Unix(int64(ts), 0),
	}
}

func (conn *Connection) SubscribeSensorEvents() <-chan SensorEvent {
	s := &sensorSubscription{cs: make(chan SensorEvent)}

	s.cbId = C.registerSensorEvent(unsafe.Pointer(s))
	conn.addSubscription(s)

	return s.cs
}

type deviceStateSubscription struct {
	cbId C.int
	cs   chan State
	dev  *Device
}

func (s *deviceStateSubscription) Close() error {
	C.tdUnregisterCallback(s.cbId)
	close(s.cs)

	return nil
}

//export newDeviceEvent
func newDeviceEvent(deviceId, method int, data *C.char, callbackId int, context unsafe.Pointer) {
	s := (*deviceStateSubscription)(context)
	if deviceId != s.dev.Id {
		return
	}

	s.dev.State = State(method)
}

//TODO: We should close the subscription when Device goes out of scope!
func (d *Device) SubscribeEvents() <-chan State {
	s := &deviceStateSubscription{cs: make(chan State), dev: d}

	s.cbId = C.registerDeviceEvent(unsafe.Pointer(s))
	d.conn.addSubscription(s)

	return s.cs
}
