package main

import (
	"flag"
	"fmt"
	"gio-device-driver/pkg/api"
	"gio-device-driver/pkg/service"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	fogNodeRegistrationRetries = 20
	callbackRegistrationDelay  = 5 * time.Second
	heartbeatDelay             = 2 * time.Minute
)

func main() {
	checkVariables()

	port := flag.Int("port", 8080, "port to be used")

	flag.Parse()

	err := registerService()
	if err != nil {
		panic(err)
	}

	go startHeartbeat()

	log.Printf("Server started on port %d\n", *port)

	router := api.NewRouter()

	p := fmt.Sprintf(":%d", *port)

	log.Fatal(http.ListenAndServe(p, router))
}

// Register the service to the FogNode for notifications
func registerService() error {
	retries := fogNodeRegistrationRetries

	for retries > 0 {
		log.Printf("Try to register to FogNode (trial: %d)\n", retries)
		callbackUuid, err := registerCallback()
		if err == nil {
			log.Printf("Callback UUID: %s\n", callbackUuid)
			return nil
		}

		log.Println(err)

		retries--

		// Sleep before try again
		time.Sleep(callbackRegistrationDelay)
	}

	return fmt.Errorf("cannot register to FogNode")
}

// This function periodically checks the status of the FogNode and tries to register again the callback in case of service unavailable.
func startHeartbeat() {
	log.Println("Heartbeat started")

	ticker := time.NewTicker(heartbeatDelay)
	go func() {
		for {
			// Wait for tick
			<-ticker.C

			log.Println("Checking FogNode")
			err := registerService()
			if err != nil {
				panic(err)
			}
		}
	}()
}

// Register callback to FogNode
func registerCallback() (string, error) {
	callbackHost := os.Getenv("CALLBACK_HOST")
	callbackPort := os.Getenv("CALLBACK_PORT")

	callbackUrl := fmt.Sprintf("http://%s:%s%s", callbackHost, callbackPort, api.CallbackEndpointPath)
	if _, err := url.Parse(callbackUrl); err != nil {
		return "", err
	}

	fogNode, err := service.NewFogNode()
	if err != nil {
		return "", err
	}

	uuid, err := fogNode.RegisterCallback(callbackUrl)
	if err != nil {
		return "", err
	}

	return uuid, nil
}

func checkVariables() {
	varNames := []string{"FOG_NODE_PORT", "DEVICE_SERVICE_HOST", "DEVICE_SERVICE_PORT", "CALLBACK_HOST", "CALLBACK_PORT"}
	for _, name := range varNames {
		if v := os.Getenv(name); v == "" {
			panic(fmt.Sprintf("%s not set", name))
		}
	}
}
