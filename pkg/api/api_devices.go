package api

import (
	"encoding/json"
	"gio-device-driver/pkg/model"
	"gio-device-driver/pkg/service"
	"gio-device-driver/pkg/smartvase"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func GetDevices(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting all connected devices")

	fogNode, _ := service.NewFogNode()

	devices, err := fogNode.GetDevices()
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}

	code := http.StatusOK
	w.WriteHeader(code)

	err = json.NewEncoder(w).Encode(devices)
	if err != nil {
		log.Println(err)
	}
}

func TriggerAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceId := vars["deviceId"]
	actionName := vars["actionName"]

	log.Printf("Requesting triggering for action %s on device %s\n", actionName, deviceId)

	fogNode, _ := service.NewFogNode()

	device, err := fogNode.GetDevice(deviceId)
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

	// get action UUID
	actionUUID := findActionUUID(actionName)
	if actionUUID == "" {
		log.Printf("Action not recognised %s\n", actionName)
		errorHandler(w, http.StatusBadRequest, "action not recognised")
		return
	}

	log.Printf("Action UUID found: %s\n", actionUUID)

	err = fogNode.TriggerAction(device, actionUUID)
	if err != nil {
		log.Printf("Cannot perform action %s (UUID: %s)\n", actionName, actionUUID)
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
}

// Get the UUID of the action with a given name
func findActionUUID(actionName string) string {
	for _, char := range smartvase.SmartVaseCharacteristics {
		if actionName == char.Name {
			return char.UUID
		}
	}

	return ""
}
