// Copyright 2021-2024 IBM Corp. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector_s

import (
	"fmt"

	"github.com/IBM/spectrum-virtualize-exporter/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
)

const prefix_host = "spectrum_host_"

var (
	host_status *prometheus.Desc
)

func init() {
	registerCollector("lshost", defaultEnabled, NewHostCollector)
}

// hostCollector collects host setting metrics
type hostCollector struct {
}

func NewHostCollector() (Collector, error) {
	labelnames := []string{"resource", "host_name"}
	if len(utils.ExtraLabelNames) > 0 {
		labelnames = append(labelnames, utils.ExtraLabelNames...)
	}
	host_status = prometheus.NewDesc(prefix_host+"status", "Host connection status. 0-online; 1-offline; 2-degraded.", labelnames, nil)
	return &hostCollector{}, nil
}

// Describe() describes the metrics
func (*hostCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- host_status
}

// Collect() collects metrics from Spectrum Virtualize Restful API
func (c *hostCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	logger.Debugln("entering host collector ...")
	respData, err := sClient.CallSpectrumAPI("lshost", true)
	if err != nil {
		logger.Errorf("executing lshost cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("response of lshost: ", respData)
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

		labelvalues := []string{sClient.Hostname, host_name}
		if len(utils.ExtraLabelValues) > 0 {
			labelvalues = append(labelvalues, utils.ExtraLabelValues...)
		}

		ch <- prometheus.MustNewConstMetric(host_status, prometheus.GaugeValue, float64(v_status), labelvalues...)
		return true
	})

	logger.Debugln("exit host exit")
	return nil
}
