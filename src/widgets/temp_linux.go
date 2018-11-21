package widgets

import (
	"strings"
        "os/exec"
	"strconv"
        "regexp"
	psHost "github.com/shirou/gopsutil/host"
)

func (self *Temp) update() {
	sensors, _ := psHost.SensorsTemperatures()
	for _, sensor := range sensors {
		// only sensors with input in their name are giving us live temp info
		if strings.Contains(sensor.SensorKey, "input") {
			// removes '_input' from the end of the sensor name
			label := sensor.SensorKey[:strings.Index(sensor.SensorKey, "_input")]
			self.Data[label] = int(sensor.Temperature)
		}
	}

        output, _ := exec.Command("hddtemp", "/dev/sda").Output()
        tempstr := strings.Split(strings.TrimSpace(string(output)), ":")[1:]
	if (len(tempstr) >= 2) {
            hddname := strings.TrimSpace(tempstr[0])
            hddtemp := strings.TrimSpace(tempstr[1])
            tempval := regexp.MustCompile(`\d*`).FindString(hddtemp)
	    if tmptemp, err := strconv.ParseInt(tempval, 10, 32); err == nil {
	        self.Data[hddname] = int(tmptemp)
	    }
	}
}
