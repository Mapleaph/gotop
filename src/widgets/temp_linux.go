package widgets

import (
	"strings"

	psHost "github.com/shirou/gopsutil/host"
)

func (self *Temp) update() {
	sensors, err := psHost.SensorsTemperatures()
	if err != nil {
		// return err
	}
	for _, sensor := range sensors {
		// only sensors with input in their name are giving us live temp info
		if strings.Contains(sensor.SensorKey, "input") {
			// removes '_input' from the end of the sensor name
			label := sensor.SensorKey[:strings.Index(sensor.SensorKey, "_input")]
			self.Data[label] = int(sensor.Temperature)
		}
	}
}
