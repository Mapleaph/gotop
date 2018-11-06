package widgets

import (
	psHost "github.com/shirou/gopsutil/host"
)

func (self *Temp) update() {
	sensors, err := psHost.SensorsTemperatures()
	if err != nil {
		// return err
	}
	for _, sensor := range sensors {
		self.Data[sensor.SensorKey] = int(sensor.Temperature)
	}
}
