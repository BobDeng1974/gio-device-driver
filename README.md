# gio-device-driver
Microservice that reads data from the Giò Plants Fog Node software installed on the machine, applying filters and
forwarding the result to the Giò Device service.

## How does it works

The service starts registering a webhook (`http://<host>:<port>/callbacks/readings`) with the gio-fog-node tools to be notified when a new reading is produced by a connected device.

When the webhook is called, it filters and cleans the data received. Then, cleaned data are sent to the device service.

## Run

The service requires a few data to successfully start:

Options:

- --port: specifies the host used by the service to expose its callbacks endpoints. Default value is 8080

Environment variables:

- FOG_NODE_PORT: specifies the port in which the gio-fog-node tool is running
- DEVICE_SERVICE_HOST: specifies the host in which the gio-device-ms service is running
- DEVICE_SERVICE_PORT: specifies the port in which the gio-device-ms service is running
- CALLBACK_PORT: specifies the port used for composing the callback url

### Go
`gio-device-driver` is developed as a Go module.
```bash
export FOG_NODE_PORT=5002
export DEVICE_SERVICE_HOST=localhost
export DEVICE_SERVICE_PORT=5001
export CALLBACK_PORT=5006

go build -o devicedriver cmd/devicedriver/main.go

./devicedriver -host localhost -port 5004
```

### Using Docker

```bash
docker build -t gio-device-driver:latest .

docker run -it --port 5004:8080 gio-device-driver:latest
```

## REST API

- POST /callbacks/readings: endpoint for device readings callbacks. Called by the Fog Node tools to notify a new reading produce by a connected device.

    Example body:
    
```json
{
    "peripheral_id": "xx:xx:xx:xx",
    "reading": {
        "name": "temperature",
        "value": "23",
        "unit": "C°",
        "creation_timestamp": "yyyyyy"  
    }
}
```

    Example response:
    
```json
{
    "status": 200,
    "message": "Done"
}
```
  
 - POST /devices/{deviceId}/actions/{actionName}: triggers an action on the specified device.
    The device is identified by its MAC address. Action name is the symbolic name associated to the action.
    The Device Driver will maps the name with the proper BLE characteristic UUID.
    It allows send a value with the request.
   
    Example body:
```json
{
    "value": "22"
}
```

    Example response:
    
```json
{
    "status": 200,
    "message": "Done"
}
```