package service

import "github.com/chiangwanyi/iot_test/model"

type DeviceService interface {
	CreateDevice(device *model.Device) (int64, error)
}

type DeviceServiceImpl struct{}

func (s *DeviceServiceImpl) CreateDevice(device *model.Device) (int64, error) {
	if id, err := device.Create(); err != nil {
		return 0, err
	} else {
		return id, nil
	}
}
