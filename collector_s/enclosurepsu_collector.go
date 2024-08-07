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

const prefix_enclosurepsu = "spectrum_enclosurepsu_"

var (
	psu_status *prometheus.Desc
)

func init() {
	registerCollector("lsenclosurepsu", defaultEnabled, NewEnclosurePsuCollector)
}

// enclosurePsuCollector collects enclosurepsu setting metrics
type enclosurePsuCollector struct {
}

func NewEnclosurePsuCollector() (Collector, error) {
	labelnames := []string{"resource", "enclosure_id", "psu_id"}
	if len(utils.ExtraLabelNames) > 0 {
		labelnames = append(labelnames, utils.ExtraLabelNames...)
	}
	psu_status = prometheus.NewDesc(prefix_enclosurepsu+"status", "Indicates status of each power-supply unit (PSU) in enclosures.", labelnames, nil)
	return &enclosurePsuCollector{}, nil
}

// Describe describes the metrics
func (*enclosurePsuCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- psu_status
}

// Collect collects metrics from Spectrum Virtualize Restful API
func (c *enclosurePsuCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	logger.Debugln("entering enclosurepsu collector ...")
	respData, err := sClient.CallSpectrumAPI("lsenclosurepsu", true)
	if err != nil {
		logger.Errorf("executing lsenclosurepsu cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("response of lsenclosurepsu: ", respData)
	/* This is a sample output of lsenclosurepsu
	[
		{
			"enclosure_id": "1",
			"PSU_id": "1",
			"status": "online",
			"input_power": "ac"
		},
		{
			"enclosure_id": "1",
			"PSU_id": "2",
			"status": "online",
			"input_power": "ac"
		}
	] */
	if !gjson.Valid(respData) {
		return fmt.Errorf("invalid json for lsenclosurepsu:\n%v", respData)
	}
	jsonEnclosures := gjson.Parse(respData)
	jsonEnclosures.ForEach(func(key, enclosure gjson.Result) bool {
		enclosure_id := enclosure.Get("enclosure_id").String()
		psu_id := enclosure.Get("PSU_id").String()
		status := enclosure.Get("status").String() // ["online", "offline", "degraded"]

		v_status := 0
		switch status {
		case "online":
			v_status = 0
		case "offline":
			v_status = 1
		case "degraded":
			v_status = 2
		}

		labelvalues := []string{sClient.Hostname, enclosure_id, psu_id}
		if len(utils.ExtraLabelValues) > 0 {
			labelvalues = append(labelvalues, utils.ExtraLabelValues...)
		}

		ch <- prometheus.MustNewConstMetric(psu_status, prometheus.GaugeValue, float64(v_status), labelvalues...)
		return true
	})

	logger.Debugln("exit enclosurepsu exit")
	return nil
}
