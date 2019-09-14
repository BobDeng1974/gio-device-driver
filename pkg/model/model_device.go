package model

type GioDevice struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
	Mac  string `json:"mac"`
	Room string `json:"room"`
}

type Reading struct {
	ID                string `json:"id,omitempty"`
	Name              string `json:"name"`
	Value             string `json:"value"`
	Unit              string `json:"unit"`
	CreationTimestamp string `json:"creation_timestamp,omitempty"`
}

type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type FogNodeDevice struct {
	ID              string              `json:"id"`
	Name            string              `json:"name"`
	Characteristics []BLECharacteristic `json:"characteristics"`
}

type BLECharacteristic struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type ActionData struct {
	Value int `json:"value"`
}
