#ifndef _EVENTS_H
#define _EVENTS_H

#include <telldus-core.h>

void WINAPI sensorEvent(const char *protocol, const char *model, int sensorId, int dataType, const char *value, int ts, int callbackId, void *context);
int registerSensorEvent(void *context);

void WINAPI deviceEvent(int deviceId, int method, const char *data, int callbackId, void *context);
int registerDeviceEvent(void *context);

#endif
