package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"gio-device-driver/pkg/api"
	"gio-device-driver/pkg/model"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	checkVariables()

	port := flag.Int("port", 8080, "port to be used")

	flag.Parse()

	go registerService()

	log.Printf("Server started on port %d\n", *port)

	router := api.NewRouter()

	p := fmt.Sprintf(":%d", *port)

	log.Fatal(http.ListenAndServe(p, router))
}

func registerService() {
	retries := 10

	for retries > 0 {
		log.Printf("Try to register to FogNode (trial: %d)\n", retries)
		callbackUuid, err := registerCallback()
		if err == nil {
			log.Printf("Callback UUID: %s\n", callbackUuid)
			return
		}

		log.Println(err)

		retries--

		// Sleep before try again
		time.Sleep(5 * time.Second)
	}

	panic("Cannot register to FogNode!")
}

func registerCallback() (string, error) {
	callbackHost := os.Getenv("CALLBACK_HOST")
	callbackPort := os.Getenv("CALLBACK_PORT")

	fogNodeHost := os.Getenv("FOG_NODE_HOST")
	fogNodePort := os.Getenv("FOG_NODE_PORT")

	fogNodeUrl := fmt.Sprintf("http://%s:%s", fogNodeHost, fogNodePort)

	if _, err := url.Parse(fogNodeUrl); err != nil {
		return "", err
	}

	callbackUrl := fmt.Sprintf("http://%s:%s%s", callbackHost, callbackPort, api.CallbackEndpointPath)
	if _, err := url.Parse(callbackUrl); err != nil {
		return "", err
	}

	callbackData := struct {
		Url string `json:"url"`
	}{
		Url: callbackUrl,
	}

	dataJson, _ := json.Marshal(callbackData)

	registrationUrl := fmt.Sprintf("%s/callbacks", fogNodeUrl)

	log.Printf("FogNode URL: %s\n", fogNodeUrl)
	log.Printf("callbackUrl: %s\n", callbackUrl)
	log.Printf("registrationUrl: %s\n", registrationUrl)

	registrationResp, err := http.Post(registrationUrl, "application/json", bytes.NewBuffer(dataJson))
	if err != nil {
		log.Printf("Error while registrering device: %s\n", err)
		return "", err
	}

	var message model.ApiResponse
	err = json.NewDecoder(registrationResp.Body).Decode(&message)
	if err != nil {
		log.Printf("Error while decoding: %s\n", err)
		return "", err
	}

	// Return the UUID
	return message.Message, nil
}

func checkVariables() {
	varNames := []string{"FOG_NODE_HOST", "FOG_NODE_PORT", "DEVICE_SERVICE_HOST", "DEVICE_SERVICE_PORT", "CALLBACK_HOST", "CALLBACK_PORT"}
	for _, name := range varNames {
		if v := os.Getenv(name); v == "" {
			panic(fmt.Sprintf("%s not set", name))
		}
	}
}
