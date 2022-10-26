package collector_s

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_enclosurebattery = "spectrum_enclosurebattery_"

var (
	battery_status              *prometheus.Desc
	battery_end_of_life_warning *prometheus.Desc
)

func init() {
	registerCollector("lsenclosurebattery", defaultEnabled, NewEnclosureBatteryCollector)
	labelnames_status := []string{"target", "resource", "enclosure_id", "battery_id"}
	labelnames_eolw := []string{"target", "resource", "enclosure_id", "battery_id"}
	battery_status = prometheus.NewDesc(prefix_enclosurebattery+"status", "Identifies the status of the battery. 0-online; 1-offline; 2-degraded.", labelnames_status, nil)
	battery_end_of_life_warning = prometheus.NewDesc(prefix_enclosurebattery+"end_of_life_warning", "Identifies the battery's end of life. Replace the battery if yes. 0-no; 1-yes.", labelnames_eolw, nil)
}

//enclosureBatteryCollector collects enclosurebattery setting metrics
type enclosureBatteryCollector struct {
}

func NewEnclosureBatteryCollector() (Collector, error) {
	return &enclosureBatteryCollector{}, nil
}

//Describe describes the metrics
func (*enclosureBatteryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- battery_status
	ch <- battery_end_of_life_warning
}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *enclosureBatteryCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	logger.Debugln("Entering enclosurebattery collector ...")
	respData, err := sClient.CallSpectrumAPI("lsenclosurebattery", true)
	if err != nil {
		logger.Errorf("Executing lsenclosurebattery cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("Response of lsenclosurebattery: ", respData)
	/* This is a sample output of lsenclosurebattery
	[
		{
			"enclosure_id": "1",
			"battery_id": "1",
			"status": "online",
			"charging_status": "idle",
			"recondition_needed": "no",
			"percent_charged": "100",
			"end_of_life_warning": "no"
		},
		{
			"enclosure_id": "1",
			"battery_id": "2",
			"status": "online",
			"charging_status": "idle",
			"recondition_needed": "no",
			"percent_charged": "100",
			"end_of_life_warning": "no"
		}
	] */
	if !gjson.Valid(respData) {
		return fmt.Errorf("invalid json for lsenclosurebattery:\n%v", respData)
	}
	jsonBatteries := gjson.Parse(respData)
	jsonBatteries.ForEach(func(key, battery gjson.Result) bool {
		enclosure_id := battery.Get("enclosure_id").String()
		battery_id := battery.Get("battery_id").String()
		status := battery.Get("status").String()                           // ["online", "offline", "degraded"]
		end_of_life_warning := battery.Get("end_of_life_warning").String() // ["yes", "no"]

		v_status := 0
		switch status {
		case "online":
			v_status = 0
		case "offline":
			v_status = 1
		case "degraded":
			v_status = 2
		}
		ch <- prometheus.MustNewConstMetric(battery_status, prometheus.GaugeValue, float64(v_status), sClient.IpAddress, sClient.Hostname, enclosure_id, battery_id)

		v_eolw := 0
		switch end_of_life_warning {
		case "no":
			v_eolw = 0
		case "yes":
			v_eolw = 1
		}
		ch <- prometheus.MustNewConstMetric(battery_end_of_life_warning, prometheus.GaugeValue, float64(v_eolw), sClient.IpAddress, sClient.Hostname, enclosure_id, battery_id)
		return true
	})

	logger.Debugln("Leaving enclosurebattery collector.")
	return nil
}
