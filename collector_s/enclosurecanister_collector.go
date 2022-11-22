package collector_s

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_enclosurecanister = "spectrum_enclosurecanister_"

var (
	canister_status *prometheus.Desc
)

func init() {
	registerCollector("lsenclosurecanister", defaultEnabled, NewEnclosureCanisterCollector)
	labelnames_status := []string{"target", "resource", "enclosure_id", "canister_id", "node_name"}
	canister_status = prometheus.NewDesc(prefix_enclosurecanister+"status", "Identifies status of each canister in enclosures.", labelnames_status, nil)
}

//enclosureCanisterCollector collects enclosurecanister setting metrics
type enclosureCanisterCollector struct {
}

func NewEnclosureCanisterCollector() (Collector, error) {
	return &enclosureCanisterCollector{}, nil
}

//Describe describes the metrics
func (*enclosureCanisterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- canister_status
}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *enclosureCanisterCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	logger.Debugln("Entering enclosurecanister collector ...")
	respData, err := sClient.CallSpectrumAPI("lsenclosurecanister", true)
	if err != nil {
		logger.Errorf("Executing lsenclosurecanister cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("Response of lsenclosurecanister: ", respData)
	/* This is a sample output of lsenclosurecanister
	[
		{
			"enclosure_id": "1",
			"canister_id": "1",
			"status": "online",
			"type": "node",
			"node_id": "1",
			"node_name": "node1"
		},
		{
			"enclosure_id": "1",
			"canister_id": "2",
			"status": "online",
			"type": "node",
			"node_id": "2",
			"node_name": "node2"
		}
	] */
	if !gjson.Valid(respData) {
		return fmt.Errorf("invalid json for lsenclosurecanister:\n%v", respData)
	}
	jsonCanisters := gjson.Parse(respData)
	jsonCanisters.ForEach(func(key, canister gjson.Result) bool {
		enclosure_id := canister.Get("enclosure_id").String()
		canister_id := canister.Get("canister_id").String()
		node_name := canister.Get("node_name").String()
		status := canister.Get("status").String() // ["online", "offline", "degraded"]

		v_status := 0
		switch status {
		case "online":
			v_status = 0
		case "offline":
			v_status = 1
		case "degraded":
			v_status = 2
		}
		ch <- prometheus.MustNewConstMetric(canister_status, prometheus.GaugeValue, float64(v_status), sClient.IpAddress, sClient.Hostname, enclosure_id, canister_id, node_name)
		return true
	})

	logger.Debugln("Leaving enclosurecanister collector.")
	return nil
}
