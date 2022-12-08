package collector_s

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_host = "spectrum_host_"

var (
	host_status *prometheus.Desc
)

func init() {
	registerCollector("lshost", defaultEnabled, NewHostCollector)
	labelnames_status := []string{"target", "resource", "host_name"}
	host_status = prometheus.NewDesc(prefix_host+"status", "Host connection status. 0-online; 1-offline; 2-degraded.", labelnames_status, nil)
}

//hostCollector collects host setting metrics
type hostCollector struct {
}

func NewHostCollector() (Collector, error) {
	return &hostCollector{}, nil
}

//Describe() describes the metrics
func (*hostCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- host_status
}

//Collect() collects metrics from Spectrum Virtualize Restful API
func (c *hostCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	logger.Debugln("Entering host collector ...")
	respData, err := sClient.CallSpectrumAPI("lshost", true)
	if err != nil {
		logger.Errorf("Executing lshost cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("Response of lshost: ", respData)
	/* This is a sample output of lshost
	[
	    {
	        "id": "0",
	        "name": "DBM1",
	        "port_count": "6",
	        "iogrp_count": "4",
	        "status": "degraded",
	        "site_id": "",
	        "site_name": "",
	        "host_cluster_id": "",
	        "host_cluster_name": "",
	        "protocol": "scsi",
	        "owner_id": "",
	        "owner_name": ""
	    },
		...
	] */
	if !gjson.Valid(respData) {
		return fmt.Errorf("invalid json for lshost:\n%v", respData)
	}
	jsonLpars := gjson.Parse(respData)
	jsonLpars.ForEach(func(key, port gjson.Result) bool {
		host_name := port.Get("name").String()
		status := port.Get("status").String() // ["online", "offline", "degraded"]

		v_status := 0
		switch status {
		case "online":
			v_status = 0
		case "offline":
			v_status = 1
		case "degraded":
			v_status = 2
		}
		ch <- prometheus.MustNewConstMetric(host_status, prometheus.GaugeValue, float64(v_status), sClient.IpAddress, sClient.Hostname, host_name)
		return true
	})

	logger.Debugln("Leaving host collector.")
	return nil
}
