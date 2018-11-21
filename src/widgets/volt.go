package widgets

import (
    "errors"
    "strconv"
    "regexp"
    "os/exec"
    "strings"
    "time"
    "sort"
    "fmt"
    ui "github.com/cjbassi/termui"
)

type Volt struct {
    *ui.Block
    Data map[string]float32
    Max  map[string]float32
    Min  map[string]float32
    interval time.Duration
    VoltHigh ui.Color
    VoltLow  ui.Color
}

func NewVolt() *Volt {
    self := &Volt{
        Block: ui.NewBlock(),
        Data: make(map[string]float32),
        Max: make(map[string]float32),
        Min: make(map[string]float32),
        interval: time.Second,
    }

    self.Label = "Voltages"

    self.update()

    ticker := time.NewTicker(self.interval)

    go func() {
        for range ticker.C {
            self.update()
        }
    } ()

    return self

}

func (self *Volt) update() {

    output, _ := exec.Command("sensors", "-Au").Output()
    sectionStrings := strings.Split(string(output), "\n\n")
    v := sectionExamine(sectionStrings)
    self.Data = v.Data
    self.Max = v.Max
    self.Min = v.Min
}

func sectionExamine(sec []string) Volt {

    var matched bool
    var sectionVolt string
    //var sectionCoretemp string
    var in_val_group []string
    var in_name_group []string
    var filtered_in_name string

    volt := Volt{
        Data: make(map[string]float32),
        Max: make(map[string]float32),
        Min: make(map[string]float32),
    }

    for i:=0; i<len(sec); i++ {
        matched, _ = regexp.MatchString("in\\d", sec[i])
        if (matched) {
            sectionVolt = sec[i]
        }
        //matched, _ = regexp.MatchString("coretemp", sec[i])
        //if (matched) {
        //    sectionCoretemp = sec[i]
        //    fmt.Println(sectionCoretemp)
        //}
    }

    if (len(sectionVolt) != 0) {

        r := regexp.MustCompile("\n[\\w+\\d\\.]+:");
        val_group := r.Split(sectionVolt, -1)
        name_group := r.FindAllString(sectionVolt, -1)

        for i:=0; i<len(val_group); i++ {
            matched = regexp.MustCompile(`in\d`).MatchString(val_group[i])
            if (matched) {
                in_val_group = append(in_val_group, val_group[i])
                filtered_in_name = strings.TrimSpace(strings.Split(name_group[i-1], ":")[0])
                in_name_group = append(in_name_group, filtered_in_name)
            }
        }

        var inx_input_val float32
        var inx_max_val float32
        var inx_min_val float32
        var r_input = regexp.MustCompile(`.*_input:\s*[\d\.]+\n`)
        var r_max = regexp.MustCompile(`.*_max:\s*[\d\.]+\n`)
        var r_min = regexp.MustCompile(`.*_min:\s*[\d\.]+\n`)

        for i:=0; i<len(in_val_group); i++ {
            if tmp_inx_input_val, err := get_inx_x_val(in_val_group[i], r_input); err == nil {
                inx_input_val = tmp_inx_input_val
            }
            if tmp_inx_max_val, err := get_inx_x_val(in_val_group[i], r_max); err == nil {
                inx_max_val = tmp_inx_max_val
            }
            if tmp_inx_min_val, err := get_inx_x_val(in_val_group[i], r_min); err == nil {
                inx_min_val = tmp_inx_min_val
            }

    //        fmt.Printf("name:%s,input:%.2f,min:%.2f,max:%.2f\n", in_name_group[i],
    //                                                       inx_input_val,
    //                                                       inx_min_val,
    //                                                       inx_max_val)
            volt.Data[in_name_group[i]] = inx_input_val
            volt.Max[in_name_group[i]] = inx_max_val
            volt.Min[in_name_group[i]] = inx_min_val
        }
    }
    return volt
}

func get_inx_x_val(inx_val_string string, r *regexp.Regexp) (float32, error) {
    tmp := r.FindString(inx_val_string)
    if (tmp != "") {
        tmp_val_string := strings.TrimSpace(regexp.MustCompile(`:`).Split(tmp, -1)[1])
        if tmp_val, err := strconv.ParseFloat(tmp_val_string, 32); err == nil {
            val := float32(tmp_val)
            return val, nil
        }
    }

    return 0, errors.New("get_inx_x_val failed")
}

// Buffer implements ui.Bufferer interface and renders the widget.
func (self *Volt) Buffer() *ui.Buffer {
        buf := self.Block.Buffer()

        var keys []string
        for key := range self.Data {
                keys = append(keys, key)
        }
        sort.Strings(keys)

        for y, key := range keys {
                if y+1 > self.Y {
                        break
                }

		fg := self.VoltLow

                if self.Data[key] > self.Max[key] {
		    fg = self.VoltHigh
		}
		if self.Data[key] < self.Min[key] {
		    fg = self.VoltHigh
		}
                s := ui.MaxString(key, (self.X - 4))
                buf.SetString(1, y+1, s, self.Fg, self.Bg)
                buf.SetString(self.X-6, y+1, fmt.Sprintf("%05.2fV", self.Data[key]), fg, self.Bg)
        }

        return buf
}
