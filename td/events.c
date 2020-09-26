#include "events.h"
#include "_cgo_export.h"

void WINAPI sensorEvent(const char *protocol, const char *model, int sensorId, int dataType, const char *value, int ts, int callbackId, void *context) {
	newSensorEvent((char*)protocol, (char*)model,	sensorId, dataType,	(char*)value, ts, callbackId, context);
}

int registerSensorEvent(void *context) {
	return tdRegisterSensorEvent( (TDSensorEvent)&sensorEvent, context);
}

void WINAPI deviceEvent(int deviceId, int method, const char *data, int callbackId, void *context) {
	newDeviceEvent(deviceId, method, (char*)data, callbackId, context);
}

int registerDeviceEvent(void *context) {
	return tdRegisterDeviceEvent( (TDDeviceEvent)&deviceEvent, context);
}

