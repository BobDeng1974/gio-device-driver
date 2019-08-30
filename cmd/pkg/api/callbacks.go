package api

import (
	"encoding/json"
	"fmt"
	"gio-device-driver/cmd/pkg/service"
	"log"
	"net/http"
)

// This is the implementation of the webhook for readings notifications
func OnReadingCreated(w http.ResponseWriter, r *http.Request) {
	log.Printf("Callback Called!")

	var data service.CallbackResponseData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		errorHandler(w, http.StatusOK, err.Error())
	}

	log.Printf("Data received: %s\n", data)

	// Send data to Device Service
	srv, _ := service.NewDeviceService()

	device, err := srv.Register(data.PeripheralID, "default")
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, fmt.Sprintf("error while registering device: %s", err))
		return
	}

	err = srv.SendData(device, &data.Reading)
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, fmt.Sprintf("cannot send data to DeviceService: %s", err))
		return
	}

	log.Println("Data sent successfully")
	w.WriteHeader(http.StatusOK)
}
