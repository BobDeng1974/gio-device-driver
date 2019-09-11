package smartvase

import (
	"gio-device-driver/pkg/model"
)

type BLEService struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

func (bles BLEService) String() string {
	return bles.UUID
}

type ProcessFunc func(reading *model.Reading) *model.Reading

type BLECharacteristic struct {
	UUID    string      `json:"uuid"`
	Name    string      `json:"name"`
	Process ProcessFunc `json:"-"`
}

func (blec BLECharacteristic) String() string {
	return blec.UUID
}

var SmartVaseCharacteristics = [...]BLECharacteristic{
	{
		UUID: "02759250523e493b8f941765effa1b20",
		Name: "light",
		Process: func(reading *model.Reading) *model.Reading {
			reading.Value = reading.Value[1 : len(reading.Value)-1]
			reading.Name = "light"
			return reading
		},
	},
	{
		UUID: "e95d9250251d470aa062fa1922dfa9a8",
		Name: "temperature",
		Process: func(reading *model.Reading) *model.Reading {
			reading.Value = reading.Value[1 : len(reading.Value)-1]
			reading.Name = "temperature"
			return reading
		},
	},
	{
		UUID: "73cd7350d32c4345a543487435c70c48",
		Name: "moisture",
		Process: func(reading *model.Reading) *model.Reading {
			reading.Value = reading.Value[1 : len(reading.Value)-1]
			reading.Name = "moisture"
			return reading
		},
	},
	{
		UUID: "ce9e7625c44341db9cb581e567f3ba93",
		Name: "watering",
		Process: func(reading *model.Reading) *model.Reading {
			return nil
		},
	},
}
