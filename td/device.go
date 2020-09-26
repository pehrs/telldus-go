package td

/*
#include <telldus-core.h>
*/
import "C"
import (
	"errors"
	"fmt"
)

// Device methods
const (
	TURNON  Methods = C.TELLSTICK_TURNON
	TURNOFF Methods = C.TELLSTICK_TURNOFF
	BELL    Methods = C.TELLSTICK_BELL
	TOGGLE  Methods = C.TELLSTICK_TOGGLE
	DIM     Methods = C.TELLSTICK_DIM
	LEARN   Methods = C.TELLSTICK_LEARN
	EXECUTE Methods = C.TELLSTICK_EXECUTE
	UP      Methods = C.TELLSTICK_UP
	DOWN    Methods = C.TELLSTICK_DOWN
	STOP    Methods = C.TELLSTICK_STOP
)

// Device states
const (
	ON     State = C.TELLSTICK_TURNON
	OFF    State = C.TELLSTICK_TURNOFF
	DIMMED State = C.TELLSTICK_DIM
)

// Device is the representation of a
// device controllable by the TellStick.
type Device struct {
	Id    int
	Name  string
	State State // Current State

	conn             *Connection
	supportedMethods Methods
	cId              C.int
}

// State is the type for representing
// a Device's state.
type State C.int

// Method is a bitmask of commands supported by the
// client and/or the device.
type Methods C.int

// Returns the device identified by the given id
func (conn *Connection) GetDevice(id int) (*Device, error) {
	dev := &Device{Id: id, cId: C.int(id), conn: conn}
	dev.Name = dev.getName()
	dev.State = dev.getLastSentCommand()

	if methods, err := dev.getMethods(); err != nil {
		return dev, err
	} else {
		dev.supportedMethods = methods
	}

	go func() {
		cs := dev.SubscribeEvents()
		for newState := range cs {
			dev.State = newState
		}
	}()

	return dev, nil
}

// Turn on the Device.
func (d *Device) TurnOn() error {
	err := getError(C.tdTurnOn(d.cId))
	if err == nil {
		d.State = ON
	}

	return err
}

// Turn off the Device.
func (d *Device) TurnOff() error {
	err := getError(C.tdTurnOff(d.cId))
	if err == nil {
		d.State = OFF
	}

	return err
}

// Dim Device to given level.
func (d *Device) Dim(level uint) error {
	err := getError(C.tdDim(d.cId, C.uchar(level)))
	if err == nil {
		d.State = DIMMED
	}

	return err
}

// Sends bell command to the Device.
func (d *Device) Bell() error {
	return getError(C.tdBell(d.cId))
}

// Sends up command to the Device.
func (d *Device) Up() error {
	return getError(C.tdUp(d.cId))
}

// Sends down command to the Device.
func (d *Device) Down() error {
	return getError(C.tdDown(d.cId))
}

// Sends stop command to the Device.
func (d *Device) Stop() error {
	return getError(C.tdStop(d.cId))
}

// Sends a special learn command to the Device.
// This is normaly devices of 'selflearning' type.
func (d *Device) Learn() error {
	return getError(C.tdLearn(d.cId))
}

// Send the given string to the device as a raw command.
func (d *Device) RawCommand(command string) error {
	return errors.New("Not implemented")
}

// Returns true if the Device supports the given Methods.
func (d *Device) Supports(m Methods) bool {
	return d.supportedMethods|m == d.supportedMethods
}

// Device's Stringer method.
func (d *Device) String() string {
	return fmt.Sprintf("%s (Id %d): %s", d.Name, d.Id, d.State)
}

// State's Stringer method.
func (s State) String() string {
	switch s {
	case ON:
		return "On"
	case OFF:
		return "Off"
	case DIMMED:
		return "Dimmed"
	}

	return ""
}

func (d *Device) getName() string {
	cStr := C.tdGetName(d.cId)
	defer C.tdReleaseString(cStr)

	return C.GoString(cStr)
}

func (d *Device) getMethods() (Methods, error) {
	methods := C.tdMethods(d.cId, C.int(d.conn.supportedMethods))

	return Methods(methods), getError(methods)
}

func (d *Device) getLastSentCommand() State {
	return State(C.tdLastSentCommand(d.cId, C.int(d.conn.supportedMethods)))
}
