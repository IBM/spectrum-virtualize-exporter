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
}

//enclosureCanisterCollector collects enclosurecanister setting metrics
type enclosureCanisterCollector struct {
}

func NewEnclosureCanisterCollector() (Collector, error) {
	labelnames := []string{"resource", "enclosure_id", "canister_id", "node_name"}
	if len(utils.ExtraLabelNames) > 0 {
		labelnames = append(labelnames, utils.ExtraLabelNames...)
	}
	canister_status = prometheus.NewDesc(prefix_enclosurecanister+"status", "Identifies status of each canister in enclosures.", labelnames, nil)
	return &enclosureCanisterCollector{}, nil
}

//Describe describes the metrics
func (*enclosureCanisterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- canister_status
}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *enclosureCanisterCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	logger.Debugln("entering enclosurecanister collector ...")
	respData, err := sClient.CallSpectrumAPI("lsenclosurecanister", true)
	if err != nil {
		logger.Errorf("executing lsenclosurecanister cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("response of lsenclosurecanister: ", respData)
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

		labelvalues := []string{sClient.Hostname, enclosure_id, canister_id, node_name}
		if len(utils.ExtraLabelValues) > 0 {
			labelvalues = append(labelvalues, utils.ExtraLabelValues...)
		}

		ch <- prometheus.MustNewConstMetric(canister_status, prometheus.GaugeValue, float64(v_status), labelvalues...)
		return true
	})

	logger.Debugln("exit enclosurecanister exit")
	return nil
}
