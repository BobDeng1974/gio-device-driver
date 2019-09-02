package api

import (
	"encoding/json"
	"gio-device-driver/pkg/model"
	"gio-device-driver/pkg/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func TriggerAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceId := vars["deviceId"]
	actionName := vars["actionName"]

	log.Printf("Requesting triggering for action %s on device %s\n", actionName, deviceId)

	srv, _ := service.NewFogNode()

	device, err := srv.GetDevice(deviceId)
	if err != nil {
		log.Printf("Cannot get device with ID %s: %s\n", deviceId, err.Error())
		errorHandler(w, http.StatusNotFound, err.Error())
		return
	}

	if device == nil {
		log.Printf("Device %s not found!", deviceId)
		errorHandler(w, http.StatusNotFound, "device not found")
		return
	}

	err = srv.TriggerAction(device, actionName)
	if err != nil {
		log.Printf("Cannot perform action %s\n", actionName)
		errorHandler(w, http.StatusBadRequest, err.Error())
		return
	}

	code := http.StatusOK
	w.WriteHeader(code)
	m := model.ApiResponse{
		Code:    code,
		Message: "Done",
	}

	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		log.Println(err)
	}

	log.Println("Action performed successfully")
}
