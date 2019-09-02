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

	host := flag.String("host", "localhost", "IP address of the current host")
	port := flag.Int("port", 8080, "port to be used")

	flag.Parse()

	go registerService(*host, *port)

	log.Printf("Server started on port %d\n", *port)

	router := api.NewRouter()

	p := fmt.Sprintf(":%d", *port)

	log.Fatal(http.ListenAndServe(p, router))
}

func registerService(host string, port int) {
	retries := 10

	for retries > 0 {
		log.Printf("Try to register to FogNode (trial: %d)\n", retries)
		callbackUuid, err := registerCallback(host, port)
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

func registerCallback(host string, port int) (string, error) {
	fogNodeHost := os.Getenv("FOG_NODE_HOST")
	fogNodePort := os.Getenv("FOG_NODE_PORT")

	fogNodeUrl := fmt.Sprintf("http://%s:%s", fogNodeHost, fogNodePort)

	_, err := url.Parse(fogNodeUrl)
	if err != nil {
		return "", err
	}

	callbackUrl := fmt.Sprintf("http://%s:%d%s", host, port, api.CallbackEndpointPath)
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
	if fogNodeHost := os.Getenv("FOG_NODE_HOST"); fogNodeHost == "" {
		panic("FOG_NODE_HOST not set.")
	}
	if fogNodePort := os.Getenv("FOG_NODE_PORT"); fogNodePort == "" {
		panic("FOG_NODE_PORT not set.")
	}

	if DeviceServiceHost := os.Getenv("DEVICE_SERVICE_HOST"); DeviceServiceHost == "" {
		panic("DEVICE_SERVICE_HOST not set.")
	}
	if DeviceServicePort := os.Getenv("DEVICE_SERVICE_PORT"); DeviceServicePort == "" {
		panic("DEVICE_SERVICE_PORT not set.")
	}
}
