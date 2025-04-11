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

package collector

import (
	"fmt"

	"github.com/IBM/spectrum-virtualize-exporter/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
)

const prefix_stats = "spectrum_systemstats_"

var (
	metricsPromDescMap map[string]*prometheus.Desc
	metricsDescMap map[string]string //update InitmetricsDescMap function for any new stat added in lssystemstats api rsp
)

type systemStatsCollector struct {
}

func init() {
	registerCollector("lssystemstats", defaultEnabled, NewSystemStatsCollector)
}

func InitmetricsDescMap() {
	metricsDescMap = make(map[string]string)
	metricsDescMap["compression_cpu_pc"] = "The percentage of allocated CPU capacity that is used for compression."
	metricsDescMap["cpu_pc"] = "The percentage of allocated CPU capacity that is used for the system."

	metricsDescMap["fc_mb"] = "The total number of megabytes transferred per second (MBps) for Fibre Channel traffic on the system. This value includes host I/O and any bandwidth that is used for communication within the system."
	metricsDescMap["fc_io"] = "The total input/output (I/O) operations that are transferred per seconds for Fibre Channel traffic on the system. This value includes host I/O and any bandwidth that is used for communication within the system."

	metricsDescMap["sas_mb"] = "The total number of megabytes transferred per second (MBps) for serial-attached SCSI (SAS) traffic on the system. This value includes host I/O and bandwidth that is used for background RAID activity."
	metricsDescMap["sas_io"] = "The total I/O operations that are transferred per second for SAS traffic on the system. This value includes host I/O and bandwidth that is used for background RAID activity."

	metricsDescMap["iscsi_mb"] = "The total number of megabytes transferred per second (MBps) for iSCSI traffic on the system."
	metricsDescMap["iscsi_io"] = "The total I/O operations that are transferred per second for iSCSI traffic on the system."

	metricsDescMap["write_cache_pc"] = "The percentage of the write cache usage for the node."
	metricsDescMap["total_cache_pc"] = "The total percentage for both the write and read cache usage for the node."

	metricsDescMap["vdisk_mb"] = "The average number of megabytes transferred per second (MBps) for read and write operations to volumes during the sample period."
	metricsDescMap["vdisk_io"] = "The average number of I/O operations that are transferred per second for read and write operations to volumes during the sample period."
	metricsDescMap["vdisk_ms"] = "The average amount of time in milliseconds that the system takes to respond to read and write requests to volumes over the sample period."

	metricsDescMap["mdisk_mb"] = "The average number of megabytes transferred per second (MBps) for read and write operations to MDisks during the sample period."
	metricsDescMap["mdisk_io"] = "The average number of I/O operations that are transferred per second for read and write operations to MDisks during the sample period."
	metricsDescMap["mdisk_ms"] = "The average amount of time in milliseconds that the system takes to respond to read and write requests to MDisks over the sample period."

	metricsDescMap["drive_mb"] = "The average number of megabytes transferred per second (MBps) for read and write operations to drives during the sample period."
	metricsDescMap["drive_io"] = "The average number of I/O operations that are transferred per second for read and write operations to drives during the sample period."
	metricsDescMap["drive_ms"] = "The average amount of time in milliseconds that the system takes to respond to read and write requests to drives over the sample period."

	metricsDescMap["vdisk_r_mb"] = "The average number of megabytes transferred per second (MBps) for read operations to volumes during the sample period."
	metricsDescMap["vdisk_r_io"] = "The average number of I/O operations that are transferred per second for read operations to volumes during the sample period."
	metricsDescMap["vdisk_r_ms"] = "The average amount of time in milliseconds that the system takes to respond to read requests to volumes over the sample period."

	metricsDescMap["vdisk_w_mb"] = "The average number of megabytes transferred per second (MBps) for read and write operations to volumes during the sample period."
	metricsDescMap["vdisk_w_io"] = "The average number of I/O operations that are transferred per second for write operations to volumes during the sample period."
	metricsDescMap["vdisk_w_ms"] = "The average amount of time in milliseconds that the system takes to respond to write requests to volumes over the sample period."

	metricsDescMap["mdisk_r_mb"] = "The average number of megabytes transferred per second (MBps) for read operations to MDisks during the sample period."
	metricsDescMap["mdisk_r_io"] = "The average number of I/O operations that are transferred per second for read operations to MDisks during the sample period."
	metricsDescMap["mdisk_r_ms"] = "The average amount of time in milliseconds that the system takes to respond to read requests to MDisks over the sample period."

	metricsDescMap["mdisk_w_mb"] = "The average number of megabytes transferred per second (MBps) for write operations to MDisks during the sample period."
	metricsDescMap["mdisk_w_io"] = "TThe average number of I/O operations that are transferred per second for write operations to MDisks during the sample period."
	metricsDescMap["mdisk_w_ms"] = "the average amount of time in milliseconds that the system takes to respond to write requests to MDisks over the sample period."

	metricsDescMap["drive_r_mb"] = "The average number of megabytes transferred per second (MBps) for read operations to drives during the sample period."
	metricsDescMap["drive_r_io"] = "The average number of I/O operations that are transferred per second for read operations to drives during the sample period."
	metricsDescMap["drive_r_ms"] = "The average amount of time in milliseconds that the system takes to respond to read requests to drives over the sample period."

	metricsDescMap["drive_w_mb"] = "The average number of megabytes transferred per second (MBps) for write operations to drives during the sample period."
	metricsDescMap["drive_w_io"] = "The average number of I/O operations that are transferred per second for write operations to drives during the sample period."
	metricsDescMap["drive_w_ms"] = "The average amount of time in milliseconds that the system takes to respond write requests to drives over the sample period."

	metricsDescMap["power_w"] = "the power that is consumed in watts."
	metricsDescMap["temp_c"] = " the ambient temperature in Celsius."
	metricsDescMap["temp_f"] = "the ambient temperature in Fahrenheit."

	metricsDescMap["iplink_mb"] = "The average number of megabytes requested to be transferred per second (MBps) over the IP partnership link during the sample period. This value is calculated before any compression of the data takes place. This value does not include iSCSI host input/output (I/O) operations."
	metricsDescMap["iplink_io"] = "TThe total input/output (I/O) operations that are transferred per second for IP partnership traffic on the system."
	metricsDescMap["iplink_comp_mb"] = "The average number of compressed megabytes transferred per second (MBps) over the IP Replication link during the sample period. This value is calculated after any compression of |the data takes place. This value does not include iSCSI host I/O operations."

	metricsDescMap["cloud_up_mb"] = "The average number of megabytes transferred per second (Mbps) for upload operations to a cloud account during the sample period."
	metricsDescMap["cloud_up_ms"] = "The average amount of time (in milliseconds) it takes for the system to respond to upload requests to a cloud account during the sample period."
	metricsDescMap["cloud_down_mb"] = "The average number of Mbps for download operations to a cloud account during the sample period."
	metricsDescMap["cloud_down_ms"] = "The average amount of time (in milliseconds) it takes for the system to respond to download requests to a cloud account during the sample period."

	metricsDescMap["iser_mb"] = "The total number of megabytes transferred per second (MBps) for iSER traffic on the system."
	metricsDescMap["iser_io"] = "The total I/O operations that are transferred per second for iSER traffic on the system."
	metricsDescMap["nvme_rdma_mb"] = "The total number of megabytes transferred per second (MBps) for NVMe over RDMA traffic on the system."
	metricsDescMap["nvme_rdma_io"] = "The total I/O operations that are transferred per second for NVMe over RDMA traffic on the system."
	metricsDescMap["nvme_tcp_mb"] = "The total number of megabytes transferred per second (MBps) for NVMe over TCP traffic on the system."
	metricsDescMap["nvme_tcp_io"] = "The total I/O operations that are transferred per second for NVMe over TCP traffic on the system."
}

func NewSystemStatsCollector() (Collector, error) {
	InitmetricsDescMap()
	labelnames := []string{"resource"}
	if len(utils.ExtraLabelNames) > 0 {
		labelnames = append(labelnames, utils.ExtraLabelNames...)
	}
	metricsPromDescMap = make(map[string]*prometheus.Desc)
	for statName, statDef := range metricsDescMap {
		metricsPromDescMap[statName] = prometheus.NewDesc(prefix_stats+statName, statDef, labelnames, nil)
	}
	metricsDescMap = nil //ready for cleanup
	return &systemStatsCollector{}, nil
}

// Describe describes the metricsPromDescMap
func (*systemStatsCollector) Describe(ch chan<- *prometheus.Desc) {

	for _, metric := range metricsPromDescMap {
		ch <- metric
	}

}

// Collect collects metrics from Spectrum Virtualize Restful API
func (c *systemStatsCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {
	logger.Debugln("entering SystemStats collector ...")
	systemStatsResp, err := sClient.CallSpectrumAPI("lssystemstats", true)
	if err != nil {
		logger.Errorf("Executing lssystemstats cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("response of lssystemstats: ", systemStatsResp)
	if !gjson.Valid(systemStatsResp) {
		return fmt.Errorf("invalid json for lscloudcallhome: %v", systemStatsResp)
	}
	/* This is a sample output of lssystemstats
		[
	    {
	        "stat_name": "compression_cpu_pc",
	        "stat_current": "0",
	        "stat_peak": "0",
	        "stat_peak_time": "181217033223"
	    },
	    {
	        "stat_name": "cpu_pc",
	        "stat_current": "1",
	        "stat_peak": "1",
	        "stat_peak_time": "181217033223"
	    },
	    {
	        "stat_name": "fc_mb",
	        "stat_current": "0",
	        "stat_peak": "0",
	        "stat_peak_time": "181217033223"
	    },
	    .......
	    .........
	] */

	labelvalues := []string{sClient.Hostname}
	if len(utils.ExtraLabelValues) > 0 {
		labelvalues = append(labelvalues, utils.ExtraLabelValues...)
	}

	systemStats := gjson.Parse(systemStatsResp).Array()

	for _, systemStat := range systemStats {

		mrtricDesc, isExist := metricsPromDescMap[systemStat.Get("stat_name").String()]
		//new stat will be ignored if not added in metricsDescMap
		if isExist {
			ch <- prometheus.MustNewConstMetric(mrtricDesc, prometheus.GaugeValue, systemStat.Get("stat_current").Float(), labelvalues...)
		}
	}
	logger.Debugln("exit SystemStats collector")
	return nil
}
