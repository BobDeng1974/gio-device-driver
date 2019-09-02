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
	actionUUID := findActionUUID(device, actionName)
	if actionUUID == "" {
		log.Printf("Action not recognised %s\n", actionName)
		errorHandler(w, http.StatusBadRequest, err.Error())
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

	log.Println("Action performed successfully")
}

// Get the UUID of the action with a given name
func findActionUUID(device *model.FogNodeDevice, actionName string) string {
	for _, char := range device.Characteristics {
		if actionName == char.Name {
			return char.UUID
		}
	}

	return ""
}
