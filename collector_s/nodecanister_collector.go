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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_nodecanister = "spectrum_nodecanister_"

var (
	nodecanister_status *prometheus.Desc
)

func init() {
	registerCollector("lsnodecanister", defaultEnabled, NewNodecanisterCollector)
}

// nodecanisterCollector collects nodecanister setting metrics
type nodecanisterCollector struct {
}

func NewNodecanisterCollector() (Collector, error) {
	labelnames := []string{"resource", "node_name"}
	if len(utils.ExtraLabelNames) > 0 {
		labelnames = append(labelnames, utils.ExtraLabelNames...)
	}
	nodecanister_status = prometheus.NewDesc(prefix_nodecanister+"status", "Status of nodes that are part of the system. 0-online; 1-offline; 2-service; 3-flushing; 4-pending; 5-adding; 6-deleting.", labelnames, nil)
	return &nodecanisterCollector{}, nil
}

// Describe describes the metrics
func (*nodecanisterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- nodecanister_status
}

// Collect collects metrics from Spectrum Virtualize Restful API
func (c *nodecanisterCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	logger.Debugln("entering nodecanister collector ...")
	respData, err := sClient.CallSpectrumAPI("lsnodecanister", true)
	if err != nil {
		logger.Errorf("executing lsnodecanister cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("response of lsnodecanister: ", respData)
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

		labelvalues := []string{sClient.Hostname, node_name}
		if len(utils.ExtraLabelValues) > 0 {
			labelvalues = append(labelvalues, utils.ExtraLabelValues...)
		}

		ch <- prometheus.MustNewConstMetric(nodecanister_status, prometheus.GaugeValue, float64(v_status), labelvalues...)
		return true
	})

	logger.Debugln("exit nodecanister exit")
	return nil
}
