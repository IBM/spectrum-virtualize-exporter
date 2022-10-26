package collector_s

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_portfc = "spectrum_portfc_"

var (
	portfc_status *prometheus.Desc
)

func init() {
	registerCollector("lsportfc", defaultEnabled, NewPortfcCollector)
	labelnames_status := []string{"target", "resource", "node_name", "port_id", "attachment"}
	portfc_status = prometheus.NewDesc(prefix_portfc+"status", "Indicates whether the port is configured to a device of Fibre Channel (FC) port. 0-active; 1-inactive_configured; 2-inactive_unconfigured.", labelnames_status, nil)
}

//portfcCollector collects portfc setting metrics
type portfcCollector struct {
}

func NewPortfcCollector() (Collector, error) {
	return &portfcCollector{}, nil
}

//Describe describes the metrics
func (*portfcCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- portfc_status
}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *portfcCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	logger.Debugln("Entering portfc collector ...")
	respData, err := sClient.CallSpectrumAPI("lsportfc", true)
	if err != nil {
		logger.Errorf("Executing lsportfc cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("Response of lsportfc: ", respData)
	/* This is a sample output of lsportfc
	[
		{
			"id": "0",
			"fc_io_port_id": "1",
			"port_id": "1",
			"type": "fc",
			"port_speed": "16Gb",
			"node_id": "1",
			"node_name": "node1",
			"WWPN": "500507681011038D",
			"nportid": "010400",
			"status": "active",
			"attachment": "switch",
			"cluster_use": "local_partner",
			"adapter_location": "1",
			"adapter_port_id": "1"
		},
		...
		{
			"id": "16",
			"fc_io_port_id": "1",
			"port_id": "1",
			"type": "fc",
			"port_speed": "16Gb",
			"node_id": "2",
			"node_name": "node2",
			"WWPN": "500507681011039F",
			"nportid": "010600",
			"status": "active",
			"attachment": "switch",
			"cluster_use": "local_partner",
			"adapter_location": "1",
			"adapter_port_id": "1"
		},
		...
	] */
	if !gjson.Valid(respData) {
		return fmt.Errorf("invalid json for lsportfc:\n%v", respData)
	}
	jsonPorts := gjson.Parse(respData)
	jsonPorts.ForEach(func(key, port gjson.Result) bool {
		port_id := port.Get("port_id").String()
		node_name := port.Get("node_name").String()
		status := port.Get("status").String() // ["active", "inactive_configured", "inactive_unconfigured"]
		attachment := port.Get("attachment").String()

		v_status := 0
		switch status {
		case "active":
			v_status = 0
		case "inactive_configured":
			v_status = 1
		case "inactive_unconfigured":
			v_status = 2
		}
		ch <- prometheus.MustNewConstMetric(portfc_status, prometheus.GaugeValue, float64(v_status), sClient.IpAddress, sClient.Hostname, node_name, port_id, attachment)
		return true
	})

	logger.Debugln("Leaving portfc collector.")
	return nil
}
