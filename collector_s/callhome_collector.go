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
