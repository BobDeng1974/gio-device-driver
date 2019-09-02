package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gio-device-driver/pkg/model"
	"log"
	"net/http"
	"net/url"
	"os"
)

type CallbackResponseData struct {
	PeripheralID string        `json:"peripheral_id"`
	Reading      model.Reading `json:"reading"`
}
type DeviceService struct {
	url *url.URL
}

func (ds *DeviceService) Register(id string, roomName string) (*model.GioDevice, error) {
	// Create the room
	roomData := model.Room{
		Name: roomName,
	}

	roomBody, _ := json.Marshal(roomData)

	roomUrl := fmt.Sprintf("%s/rooms", ds.url)
	roomResponse, err := http.Post(roomUrl, "application/json", bytes.NewBuffer(roomBody))

	if err != nil {
		return nil, err
	}

	defer roomResponse.Body.Close()

	if roomResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cannot perform the requested operation: (%d) %s", roomResponse.StatusCode, roomResponse.Status)
	}

	var room model.Room
	if err := json.NewDecoder(roomResponse.Body).Decode(&room); err != nil {
		return nil, err
	}

	// Register the Device
	deviceData := model.GioDevice{
		Name: "device" + id,
		Mac:  id,
	}

	deviceBody, _ := json.Marshal(deviceData)

	devicesUrl := fmt.Sprintf("%s/rooms/%s/devices", ds.url, room.ID)
	deviceResponse, err := http.Post(devicesUrl, "application/json", bytes.NewBuffer(deviceBody))

	if err != nil {
		return nil, err
	}

	defer deviceResponse.Body.Close()

	if deviceResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cannot perform the requested operation: (%d) %s", deviceResponse.StatusCode, deviceResponse.Status)
	}

	// Take the id from the response
	var device model.GioDevice
	_ = json.NewDecoder(deviceResponse.Body).Decode(&device)

	return &device, nil
}

func (ds *DeviceService) SendData(device *model.GioDevice, reading *model.Reading) error {
	body, err := json.Marshal(reading)
	if err != nil {
		return err
	}

	readingsUrl := fmt.Sprintf("%s/rooms/%s/devices/%s/readings", ds.url, device.Room, device.ID)
	res, err := http.Post(readingsUrl, "application/json", bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot perform the requested operation: (%d) %s", res.StatusCode, res.Status)
	}

	return nil
}

var instance *DeviceService = nil

func NewDeviceService() (*DeviceService, error) {
	serviceHost := os.Getenv("DEVICE_SERVICE_HOST")
	servicePort := os.Getenv("DEVICE_SERVICE_PORT")

	if instance == nil {
		u := fmt.Sprintf("http://%s:%s", serviceHost, servicePort)
		log.Printf("DeviceService URL: %s\n", u)

		serviceUrl, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		instance = &DeviceService{serviceUrl}
	}

	return instance, nil
}