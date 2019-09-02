package api

import (
	"encoding/json"
	"fmt"
	"gio-device-driver/pkg/devices"
	"gio-device-driver/pkg/model"
	"gio-device-driver/pkg/service"
	"log"
	"net/http"
)

// This is the implementation of the webhook for readings notifications
// This function is called when the FogNode creates a new Reading and notify this driver by an http call.
func OnReadingCreated(w http.ResponseWriter, r *http.Request) {
	var data service.CallbackResponseData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, "invalid data")
		return
	}

	// Process data
	processed := processData(data)

	if processed == nil {
		errorHandler(w, http.StatusBadRequest, fmt.Sprintf("characteristic not recognised: %s", data.Reading.Name))
		return
	}

	// Send data to Device Service
	srv, _ := service.NewDeviceService()

	device, err := srv.Register(data.PeripheralID, "default")
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, fmt.Sprintf("error while registering device: %s", err))
		return
	}

	err = srv.SendData(device, processed)
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, fmt.Sprintf("cannot send data to DeviceService: %s", err))
		return
	}

	w.WriteHeader(http.StatusOK)

	m := model.ApiResponse{
		Code:    http.StatusOK,
		Message: "Done",
	}

	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		log.Println(err)
	}
}

func processData(data service.CallbackResponseData) *model.Reading {
	// Check SmartVase characteristics
	for _, char := range devices.SmartVaseCharacteristics {
		if char.UUID == data.Reading.Name {
			return char.Process(&data.Reading)
		}
	}

	return nil
}
