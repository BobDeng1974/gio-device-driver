package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gio-device-driver/pkg/model"
	"gio-device-driver/pkg/utils"
	"log"
	"net/http"
	"net/url"
	"os"
)

type FogNode struct {
	url string
}

func (fogNode *FogNode) GetDevice(deviceId string) (*model.FogNodeDevice, error) {
	log.Printf("fognode: %s- deviceid: %s", fogNode, deviceId)
	devicesUrl := fmt.Sprintf("%s/devices/%s", fogNode.url, deviceId)
	resp, err := http.Get(devicesUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var device model.FogNodeDevice
	err = json.NewDecoder(resp.Body).Decode(&device)
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (fogNode *FogNode) TriggerAction(device *model.FogNodeDevice, actionName string) error {
	devicesUrl := fmt.Sprintf("%s/devices/%s/actions/%s", fogNode.url, device.ID, actionName)
	resp, err := http.Post(devicesUrl, "application/json", nil)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

// Registers a new callback to the local Fog Node
func (fogNode *FogNode) RegisterCallback(callbackUrl string) (string, error) {
	callbackData := struct {
		Url string `json:"url"`
	}{
		Url: callbackUrl,
	}

	dataJson, _ := json.Marshal(callbackData)

	registrationUrl := fmt.Sprintf("%s/callbacks", fogNode.url)

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

// Singleton
var fogNodeInstance *FogNode = nil

func NewFogNode() (*FogNode, error) {
	if fogNodeInstance == nil {
		fogNodeHost, err := utils.GetHostIP()
		if err != nil {
			return nil, err
		}
		fogNodePort := os.Getenv("FOG_NODE_PORT")
		u := fmt.Sprintf("http://%s:%s", fogNodeHost, fogNodePort)
		log.Printf("FogNode URL: %s\n", u)

		_, err = url.Parse(u)
		if err != nil {
			return nil, err
		}

		fogNodeInstance = &FogNode{
			url: u,
		}
	}

	return fogNodeInstance, nil
}
