package api

import (
	"encoding/json"
	"gio-device-driver/pkg/logging"
	"gio-device-driver/pkg/model"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      []string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

const (
	CallbackEndpointPath = "/callbacks/readings"
)

func errorHandler(w http.ResponseWriter, status int, message string) {
	r := model.ApiResponse{Code: status, Message: message}
	w.WriteHeader(int(status))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(r); err != nil {
		log.Println(err)
	}
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logging.Logger(handler, route.Name)

		router.
			Methods(route.Method...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	// Callback for SmartVases
	Route{
		"OnReadingCreated",
		[]string{http.MethodPost},
		CallbackEndpointPath,
		OnReadingCreated,
	},

	// Trigger for actions
	Route{
		"TriggerAction",
		[]string{http.MethodPost},
		"/devices/{deviceId}/actions/{actionName}",
		TriggerAction,
	},
}
