package collector_s

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_callhome = "spectrum_callhome_"

var callhomeInfo *prometheus.Desc

func init() {
	registerCollector("lscloudcallhome", defaultEnabled, NewCallhomeInfoCollector)
	labelnames := []string{"target", "resource", "status", "connection"}
	callhomeInfo = prometheus.NewDesc(prefix_callhome+"info", "The status of the Call Home information.", labelnames, nil)

}

//callhomeInfoCollector collects callhome setting metrics
type callhomeInfoCollector struct {
}

func NewCallhomeInfoCollector() (Collector, error) {
	return &callhomeInfoCollector{}, nil
}

//Describe describes the metrics
func (*callhomeInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- callhomeInfo
}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *callhomeInfoCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	log.Debugln("Entering Callhome collector ...")
	respData, err := sClient.CallSpectrumAPI("lscloudcallhome", true)
	if err != nil {
		log.Errorf("Executing lscloudcallhome cmd failed: %s", err.Error())
	}
	log.Debugln("Response of lscloudcallhome: ", respData)
	// This is a sample output of lscloudcallhome
	//	{
	//		"status": "disabled",          ["disabled", "enabled"]
	//		"connection": "",              ["active", "error", "untried"]
	//		"error_sequence_number": "",
	//		"last_success": "220308065924",
	//		"last_failure": "220308065307"
	//	}

	jsonCallhome := gjson.Parse(respData)

	status := jsonCallhome.Get("status").String()
	connection := jsonCallhome.Get("connection").String()

	value := 0
	// 0: status --enabled, connection --active;
	// 1: status --disabled
	// 2: status --enabled, connection in ["error", "untried"]
	if status != "enabled" {
		value ^= 1
	} else {
		if connection != "active" {
			value ^= 2
		}
	}

	ch <- prometheus.MustNewConstMetric(callhomeInfo, prometheus.GaugeValue, float64(value), sClient.IpAddress, sClient.Hostname, status, connection)

	log.Debugln("Leaving Callhome collector.")
	return err
}
