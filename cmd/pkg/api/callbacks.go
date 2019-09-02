package api

import (
	"encoding/json"
	"fmt"
	"gio-device-driver/cmd/pkg/devices"
	"gio-device-driver/cmd/pkg/model"
	"gio-device-driver/cmd/pkg/service"
	"log"
	"net/http"
)

// This is the implementation of the webhook for readings notifications
func OnReadingCreated(w http.ResponseWriter, r *http.Request) {
	log.Println("Callback Called!")

	var data service.CallbackResponseData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, "invalid data")
		return
	}

	log.Printf("Data received: %v\n", data)

	// Process data
	processed := processData(data)

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

	log.Println("Data sent successfully")
}

func processData(data service.CallbackResponseData) *model.Reading {
	// Check SmartVase characteristics
	for _, char := range devices.SmartVaseCharacteristics {
		if char.UUID == data.Reading.ID {
			return char.Process(&data.Reading)
		}
	}

	return nil
}
