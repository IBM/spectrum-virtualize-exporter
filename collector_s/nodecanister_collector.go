package collector_s

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_nodecanister = "spectrum_nodecanister_"

var (
	nodecanister_status *prometheus.Desc
)

func init() {
	registerCollector("lsnodecanister", defaultEnabled, NewNodecanisterCollector)
	labelnames_status := []string{"target", "resource", "node_name"}
	nodecanister_status = prometheus.NewDesc(prefix_nodecanister+"status", "Status of nodes that are part of the system. 0-online; 1-offline; 2-service; 3-flushing; 4-pending; 5-adding; 6-deleting.", labelnames_status, nil)
}

//nodecanisterCollector collects nodecanister setting metrics
type nodecanisterCollector struct {
}

func NewNodecanisterCollector() (Collector, error) {
	return &nodecanisterCollector{}, nil
}

//Describe describes the metrics
func (*nodecanisterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- nodecanister_status
}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *nodecanisterCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	log.Debugln("Entering nodecanister collector ...")
	respData, err := sClient.CallSpectrumAPI("lsnodecanister", true)
	if err != nil {
		log.Errorf("Executing lsnodecanister cmd failed: %s", err.Error())
		return err
	}
	log.Debugln("Response of lsnodecanister: ", respData)
	/* This is a sample output of lsnodecanister
	[
		{
			"id": "1",
			"name": "node1",
			"UPS_serial_number": "",
			"WWNN": "500507681000038D",
			"status": "online",
			"IO_group_id": "0",
			"IO_group_name": "io_grp0",
			"config_node": "no",
			"UPS_unique_id": "",
			"hardware": "AF8",
			"iscsi_name": "iqn.1986-03.com.ibm:2145.sara-wdc04-03.node1",
			"iscsi_alias": "",
			"panel_name": "01-1",
			"enclosure_id": "1",
			"canister_id": "1",
			"enclosure_serial_number": "78E008V",
			"site_id": "",
			"site_name": ""
		},
			...
	] */
	if !gjson.Valid(respData) {
		return fmt.Errorf("invalid json for lsnodecanister:\n%v", respData)
	}
	jsonNodes := gjson.Parse(respData)
	jsonNodes.ForEach(func(key, port gjson.Result) bool {
		node_name := port.Get("name").String()
		status := port.Get("status").String() // ["online", "offline", "degraded"]

		v_status := 0
		switch status {
		case "online":
			v_status = 0
		case "offline":
			v_status = 1
		case "service":
			v_status = 2
		case "flushing":
			v_status = 3
		case "pending":
			v_status = 4
		case "adding":
			v_status = 5
		case "deleting":
			v_status = 6
		}

		ch <- prometheus.MustNewConstMetric(nodecanister_status, prometheus.GaugeValue, float64(v_status), sClient.IpAddress, sClient.Hostname, node_name)
		return true
	})

	log.Debugln("Leaving nodecanister collector.")
	return nil
}
