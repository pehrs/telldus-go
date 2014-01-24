// A simple library for communication with Telldus devices.
//
// td is a golang library for talking to Telldus
// products (home automation devices). It wraps
// around libtelldus-core to provide a simple way
// of controlling a TellStick/TellStick Duo in a
// go-ish way through the Telldus service.
package td

/*
#cgo LDFLAGS: -ltelldus-core
#include <telldus-core.h>
*/
import "C"

import (
	"errors"
	"io"
)

// Connection is the representation of a
// connection to the Telldus service.
type Connection struct {
	supportedMethods Methods
	subscriptions    []io.Closer
}

// Returns a new "connection" to the local Telldus service.
//
// supportedMethods should be all device methods the client
// implements.
func NewConnection(supportedMethods Methods) (*Connection, error) {
	C.tdInit()
	return &Connection{supportedMethods, make([]io.Closer, 0)}, nil
}

// Closes the connection to the local Telldus service.
func (conn *Connection) Close() {
	for _, s := range conn.subscriptions {
		s.Close()
	}

	C.tdClose()
}

// Returns all the devices registered with the Telldus service.
func (conn *Connection) Devices() ([]*Device, error) {
	devs := make([]*Device, 0)
	numDevs := C.tdGetNumberOfDevices()

	for i := C.int(0); i < numDevs; i++ {
		dev, err := conn.GetDevice(int(C.tdGetDeviceId(i)))
		if err != nil {
			return devs, err
		}
		devs = append(devs, dev)
	}

	return devs, getError(numDevs)
}

func (conn *Connection) addSubscription(s io.Closer) {
	conn.subscriptions = append(conn.subscriptions, s)
}

func getError(errval C.int) error {
	if errval >= C.TELLSTICK_SUCCESS {
		return nil
	}

	cStr := C.tdGetErrorString(errval)
	defer C.tdReleaseString(cStr)

	return errors.New(C.GoString(cStr))
}
