#include "events.h"

void WINAPI sensorEvent(const char *protocol, const char *model, int sensorId, int dataType, const char *value, int ts, int callbackId, void *context) {
	newSensorEvent(protocol, model,	sensorId, dataType,	value, ts, callbackId, context);
}

int registerSensorEvent(void *context) {
	return tdRegisterSensorEvent( (TDSensorEvent)&sensorEvent, context);
}

void WINAPI deviceEvent(int deviceId, int method, const char *data, int callbackId, void *context) {
	newDeviceEvent(deviceId, method, data, callbackId, context);
}

int registerDeviceEvent(void *context) {
	return tdRegisterDeviceEvent( (TDDeviceEvent)&deviceEvent, context);
}

