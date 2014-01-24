package td

/*
#cgo LDFLAGS: -ltelldus-core
#include <telldus-core.h>
#include "events.h"
*/
import "C"
import (
	"time"
	"unsafe"
)

type Sensor struct {
	Protocol string
	Model    string
	Id       int
	DataType int
	Value    string
	When     time.Time
}

type sensorSubscription struct {
	cs   chan Sensor
	cbId C.int
}

func (c *sensorSubscription) Close() error {
	C.tdUnregisterCallback(c.cbId)
	close(c.cs)

	return nil
}

//export newSensorEvent
func newSensorEvent(protocol, model *C.char, sensorId, dataType int, value *C.char, ts, callbackId int, context unsafe.Pointer) {
	cs := (*sensorSubscription)(context).cs

	cs <- Sensor{
		Protocol: C.GoString(protocol),
		Model:    C.GoString(model),
		Id:       sensorId,
		DataType: dataType,
		Value:    C.GoString(value),
		When:     time.Unix(int64(ts), 0),
	}
}

func (conn *Connection) SubscribeSensorEvents() <-chan Sensor {
	s := &sensorSubscription{cs: make(chan Sensor)}

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
