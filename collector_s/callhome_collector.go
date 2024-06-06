package collector_s

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_callhome = "spectrum_callhome_"

var callhomeInfo *prometheus.Desc

func init() {
	registerCollector("lscloudcallhome", defaultEnabled, NewCallhomeInfoCollector)
}

// callhomeInfoCollector collects callhome setting metrics
type callhomeInfoCollector struct {
}

func NewCallhomeInfoCollector() (Collector, error) {
	labelnames := []string{"resource", "status", "connection"}
	if len(utils.ExtraLabelNames) > 0 {
		labelnames = append(labelnames, utils.ExtraLabelNames...)
	}
	callhomeInfo = prometheus.NewDesc(prefix_callhome+"info", "The status of the Call Home information.", labelnames, nil)

	return &callhomeInfoCollector{}, nil
}

// Describe describes the metrics
func (*callhomeInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- callhomeInfo
}

// Collect collects metrics from Spectrum Virtualize Restful API
func (c *callhomeInfoCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	logger.Debugln("entering Callhome collector ...")
	respData, err := sClient.CallSpectrumAPI("lscloudcallhome", true)
	if err != nil {
		logger.Errorf("executing lscloudcallhome cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("response of lscloudcallhome: ", respData)
	/* This is a sample output of lscloudcallhome
	{
		"status": "disabled",          // ["disabled", "enabled"]
		"connection": "",              // ["active", "error", "untried"]
		"error_sequence_number": "",
		"last_success": "220308065924",
		"last_failure": "220308065307"
	} */
	if !gjson.Valid(respData) {
		return fmt.Errorf("invalid json for lscloudcallhome:\n%v", respData)
	}
	jsonCallhome := gjson.Parse(respData)

	status := jsonCallhome.Get("status").String()
	connection := jsonCallhome.Get("connection").String()
	if connection == "" {
		connection = "unknown"
	}

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
	labelvalues := []string{sClient.Hostname, status, connection}
	if len(utils.ExtraLabelValues) > 0 {
		labelvalues = append(labelvalues, utils.ExtraLabelValues...)
	}
	ch <- prometheus.MustNewConstMetric(callhomeInfo, prometheus.GaugeValue, float64(value), labelvalues...)

	logger.Debugln("exit Callhome exit")
	return err
}
