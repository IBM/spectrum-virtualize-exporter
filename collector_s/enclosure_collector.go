package collector_s

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_enclosure = "spectrum_enclosure_"

var (
	enclosure_status *prometheus.Desc
)

func init() {
	registerCollector("lsenclosure", defaultEnabled, NewEnclosureCollector)
}

//enclosureCollector collects enclosure setting metrics
type enclosureCollector struct {
}

func NewEnclosureCollector() (Collector, error) {
	labelnames := []string{"resource", "enclosure_id"}
	if len(utils.ExtraLabelNames) > 0 {
		labelnames = append(labelnames, utils.ExtraLabelNames...)
	}
	enclosure_status = prometheus.NewDesc(prefix_enclosure+"status", "Indicates whether an enclosure is visible to the SAS network. 0-online; 1-offline; 2-degraded.", labelnames, nil)
	return &enclosureCollector{}, nil
}

//Describe describes the metrics
func (*enclosureCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- enclosure_status
}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *enclosureCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	logger.Debugln("entering enclosure collector ...")
	respData, err := sClient.CallSpectrumAPI("lsenclosure", true)
	if err != nil {
		logger.Errorf("executing lsenclosure cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("response of lsenclosure: ", respData)
	/* This is a sample output of lsenclosure
	[
	    {
	        "id": "1",
	        "status": "online",
	        "type": "control",
	        "managed": "yes",
	        "IO_group_id": "0",
	        "IO_group_name": "io_grp0",
	        "product_MTM": "9846-AF8",
	        "serial_number": "78E008V",
	        "total_canisters": "2",
	        "online_canisters": "2",
	        "total_PSUs": "2",
	        "online_PSUs": "2",
	        "drive_slots": "24",
	        "total_fan_modules": "0",
	        "online_fan_modules": "0",
	        "total_sems": "0",
	        "online_sems": "0"
	    }
	] */
	if !gjson.Valid(respData) {
		return fmt.Errorf("invalid json for lsenclosure:\n%v", respData)
	}
	jsonEnclosures := gjson.Parse(respData)
	jsonEnclosures.ForEach(func(key, enclosure gjson.Result) bool {
		enclosure_id := enclosure.Get("id").String()
		status := enclosure.Get("status").String() // ["online", "offline", "degraded"]

		v_status := 0
		switch status {
		case "online":
			v_status = 0
		case "offline":
			v_status = 1
		case "degraded":
			v_status = 2
		}
		labelvalues := []string{sClient.Hostname, enclosure_id}
		if len(utils.ExtraLabelValues) > 0 {
			labelvalues = append(labelvalues, utils.ExtraLabelValues...)
		}
		ch <- prometheus.MustNewConstMetric(enclosure_status, prometheus.GaugeValue, float64(v_status), labelvalues...)
		return true
	})

	logger.Debugln("exit enclosure exit")
	return nil
}
