package service

import (
	"encoding/json"
	"fmt"
	"gio-device-driver/pkg/model"
	"log"
	"net/http"
	"net/url"
	"os"
)

type FogNode struct {
	url string
}

func (fogNode *FogNode) GetDevice(deviceId string) (*model.FogNodeDevice, error) {
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
	resp, err := http.Get(devicesUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

var fogNodeInstance *FogNode = nil

func NewFogNode() (*FogNode, error) {
	fogNodeHost := os.Getenv("FOG_NODE_HOST")
	fogNodePort := os.Getenv("FOG_NODE_PORT")

	if deviceServiceInstance == nil {
		u := fmt.Sprintf("http://%s:%s", fogNodeHost, fogNodePort)
		log.Printf("FogNode URL: %s\n", u)

		_, err := url.Parse(u)
		if err != nil {
			return nil, err
		}

		fogNodeInstance = &FogNode{
			url: u,
		}
	}

	return fogNodeInstance, nil
}
