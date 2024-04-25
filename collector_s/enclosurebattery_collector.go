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
}

// enclosureBatteryCollector collects enclosurebattery setting metrics
type enclosureBatteryCollector struct {
}

func NewEnclosureBatteryCollector() (Collector, error) {
	labelnames_status := []string{"resource", "enclosure_id", "battery_id"}
	labelnames_eolw := []string{"resource", "enclosure_id", "battery_id"}
	if len(utils.ExtraLabelNames) > 0 {
		labelnames_status = append(labelnames_status, utils.ExtraLabelNames...)
		labelnames_eolw = append(labelnames_eolw, utils.ExtraLabelNames...)
	}
	battery_status = prometheus.NewDesc(prefix_enclosurebattery+"status", "Identifies status of each battery in enclosures. 0-online; 1-offline; 2-degraded.", labelnames_status, nil)
	battery_end_of_life_warning = prometheus.NewDesc(prefix_enclosurebattery+"end_of_life_warning", "Identifies the battery's end of life. Replace the battery if yes. 0-no; 1-yes.", labelnames_eolw, nil)
	return &enclosureBatteryCollector{}, nil
}

// Describe describes the metrics
func (*enclosureBatteryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- battery_status
	ch <- battery_end_of_life_warning
}

// Collect collects metrics from Spectrum Virtualize Restful API
func (c *enclosureBatteryCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	logger.Debugln("entering enclosurebattery collector ...")
	respData, err := sClient.CallSpectrumAPI("lsenclosurebattery", true)
	if err != nil {
		logger.Errorf("executing lsenclosurebattery cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("response of lsenclosurebattery: ", respData)
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

		labelvalues := []string{sClient.Hostname, enclosure_id, battery_id}
		if len(utils.ExtraLabelValues) > 0 {
			labelvalues = append(labelvalues, utils.ExtraLabelValues...)
		}

		v_status := 0
		switch status {
		case "online":
			v_status = 0
		case "offline":
			v_status = 1
		case "degraded":
			v_status = 2
		}
		ch <- prometheus.MustNewConstMetric(battery_status, prometheus.GaugeValue, float64(v_status), labelvalues...)

		v_eolw := 0
		switch end_of_life_warning {
		case "no":
			v_eolw = 0
		case "yes":
			v_eolw = 1
		}
		ch <- prometheus.MustNewConstMetric(battery_end_of_life_warning, prometheus.GaugeValue, float64(v_eolw), labelvalues...)
		return true
	})

	logger.Debugln("exit enclosurebattery exit")
	return nil
}
