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

	srv, _ := service.NewFogNode()

	device, err := srv.GetDevice(deviceId)
	if err != nil {
		errorHandler(w, http.StatusNotFound, err.Error())
		return
	}

	if device == nil {
		errorHandler(w, http.StatusNotFound, "device not found")
		return
	}

	err = srv.TriggerAction(device, actionName)
	if err != nil {
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

	log.Println("Data sent successfully")
}
