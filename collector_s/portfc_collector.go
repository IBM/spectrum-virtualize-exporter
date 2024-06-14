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

const prefix_portfc = "spectrum_portfc_"

var (
	portfc_status     *prometheus.Desc
	portfc_attachment *prometheus.Desc
)

func init() {
	registerCollector("lsportfc", defaultEnabled, NewPortfcCollector)
}

// portfcCollector collects portfc setting metrics
type portfcCollector struct {
}

func NewPortfcCollector() (Collector, error) {
	labelnames_status := []string{"resource", "node_name", "port_id", "wwpn"}
	labelnames_attachment := []string{"resource", "node_name", "port_id", "wwpn"}
	if len(utils.ExtraLabelNames) > 0 {
		labelnames_status = append(labelnames_status, utils.ExtraLabelNames...)
		labelnames_attachment = append(labelnames_attachment, utils.ExtraLabelNames...)
	}
	portfc_status = prometheus.NewDesc(prefix_portfc+"status", "Indicates whether the port is configured to a device of Fibre Channel (FC) port. 0-active; 1-inactive_configured; 2-inactive_unconfigured.", labelnames_status, nil)
	portfc_attachment = prometheus.NewDesc(prefix_portfc+"attachment", "Indicates if the port is attached to a FC switch. 0-yes; 1-no.", labelnames_attachment, nil)
	return &portfcCollector{}, nil
}

// Describe describes the metrics
func (*portfcCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- portfc_status
	ch <- portfc_attachment
}

// Collect collects metrics from Spectrum Virtualize Restful API
func (c *portfcCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	logger.Debugln("entering portfc collector ...")
	respData, err := sClient.CallSpectrumAPI("lsportfc", true)
	if err != nil {
		logger.Errorf("executing lsportfc cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("response of lsportfc: ", respData)
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
		if port_id != "1" && port_id != "2" && port_id != "5" && port_id != "6" {
			return true
		}
		node_name := port.Get("node_name").String()
		wwpn := port.Get("WWPN").String()
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
		v_attachment := 0
		if attachment != "switch" {
			v_attachment = 1
		}

		labelvalues := []string{sClient.Hostname, node_name, port_id, wwpn}
		if len(utils.ExtraLabelValues) > 0 {
			labelvalues = append(labelvalues, utils.ExtraLabelValues...)
		}

		ch <- prometheus.MustNewConstMetric(portfc_status, prometheus.GaugeValue, float64(v_status), labelvalues...)
		ch <- prometheus.MustNewConstMetric(portfc_attachment, prometheus.GaugeValue, float64(v_attachment), labelvalues...)
		return true
	})

	logger.Debugln("exit portfc exit")
	return nil
}
