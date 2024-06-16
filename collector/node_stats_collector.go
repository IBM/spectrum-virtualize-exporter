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
	"github.com/IBM/spectrum-virtualize-exporter/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
)

const prefix_nodeStats = "spectrum_nodestats_"

var (
	nodeStats_metrics [46]*prometheus.Desc
)

func init() {
	registerCollector("lsnodestats", defaultDisabled, NewNodeStatsCollector)
}

// nodeStatsCollector collects nodeStats metrics
type nodeStatsCollector struct {
}

func NewNodeStatsCollector() (Collector, error) {
	labelnames := []string{"resource", "node"}
	if len(utils.ExtraLabelNames) > 0 {
		labelnames = append(labelnames, utils.ExtraLabelNames...)
	}
	nodeStats_metrics = [46]*prometheus.Desc{
		prometheus.NewDesc(prefix_nodeStats+"compression_cpu_pc", "The percentage of allocated CPU capacity that is used for compression.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"cpu_pc", "The percentage of allocated CPU capacity that is used for the system.", labelnames, nil),

		prometheus.NewDesc(prefix_nodeStats+"fc_mb", "The total number of megabytes transferred per second (MBps) for Fibre Channel traffic on the system. This value includes host I/O and any bandwidth that is used for communication within the system.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"fc_io", "The total input/output (I/O) operations that are transferred per seconds for Fibre Channel traffic on the system. This value includes host I/O and any bandwidth that is used for communication within the system.", labelnames, nil),

		prometheus.NewDesc(prefix_nodeStats+"sas_mb", "The total number of megabytes transferred per second (MBps) for serial-attached SCSI (SAS) traffic on the system. This value includes host I/O and bandwidth that is used for background RAID activity.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"sas_io", "The total I/O operations that are transferred per second for SAS traffic on the system. This value includes host I/O and bandwidth that is used for background RAID activity.", labelnames, nil),

		prometheus.NewDesc(prefix_nodeStats+"iscsi_mb", "The total number of megabytes transferred per second (MBps) for iSCSI traffic on the system.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"iscsi_io", "The total I/O operations that are transferred per second for iSCSI traffic on the system.", labelnames, nil),

		prometheus.NewDesc(prefix_nodeStats+"write_cache_pc", "The percentage of the write cache usage for the node.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"total_cache_pc", "The total percentage for both the write and read cache usage for the node.", labelnames, nil),

		prometheus.NewDesc(prefix_nodeStats+"vdisk_mb", "The average number of megabytes transferred per second (MBps) for read and write operations to volumes during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"vdisk_io", "The average number of I/O operations that are transferred per second for read and write operations to volumes during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"vdisk_ms", "The average amount of time in milliseconds that the system takes to respond to read and write requests to volumes over the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_nodeStats+"mdisk_mb", "The average number of megabytes transferred per second (MBps) for read and write operations to MDisks during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"mdisk_io", "The average number of I/O operations that are transferred per second for read and write operations to MDisks during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"mdisk_ms", "The average amount of time in milliseconds that the system takes to respond to read and write requests to MDisks over the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_nodeStats+"drive_mb", "The average number of megabytes transferred per second (MBps) for read and write operations to drives during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"drive_io", "The average number of I/O operations that are transferred per second for read and write operations to drives during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"drive_ms", "The average amount of time in milliseconds that the system takes to respond to read and write requests to drives over the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_nodeStats+"vdisk_r_mb", "The average number of megabytes transferred per second (MBps) for read operations to volumes during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"vdisk_r_io", "The average number of I/O operations that are transferred per second for read operations to volumes during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"vdisk_r_ms", "The average amount of time in milliseconds that the system takes to respond to read requests to volumes over the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_nodeStats+"vdisk_w_mb", "The average number of megabytes transferred per second (MBps) for read and write operations to volumes during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"vdisk_w_io", "The average number of I/O operations that are transferred per second for write operations to volumes during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"vdisk_w_ms", "The average amount of time in milliseconds that the system takes to respond to write requests to volumes over the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_nodeStats+"mdisk_r_mb", "The average number of megabytes transferred per second (MBps) for read operations to MDisks during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"mdisk_r_io", "The average number of I/O operations that are transferred per second for read operations to MDisks during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"mdisk_r_ms", "The average amount of time in milliseconds that the system takes to respond to read requests to MDisks over the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_nodeStats+"mdisk_w_mb", "The average number of megabytes transferred per second (MBps) for write operations to MDisks during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"mdisk_w_io", "TThe average number of I/O operations that are transferred per second for write operations to MDisks during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"mdisk_w_ms", "the average amount of time in milliseconds that the system takes to respond to write requests to MDisks over the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_nodeStats+"drive_r_mb", "The average number of megabytes transferred per second (MBps) for read operations to drives during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"drive_r_io", "The average number of I/O operations that are transferred per second for read operations to drives during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"drive_r_ms", "The average amount of time in milliseconds that the system takes to respond to read requests to drives over the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_nodeStats+"drive_w_mb", "The average number of megabytes transferred per second (MBps) for write operations to drives during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"drive_w_io", "The average number of I/O operations that are transferred per second for write operations to drives during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"drive_w_ms", "The average amount of time in milliseconds that the system takes to respond write requests to drives over the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_nodeStats+"iplink_mb", "The average number of megabytes requested to be transferred per second (MBps) over the IP partnership link during the sample period. This value is calculated before any compression of the data takes place. This value does not include iSCSI host input/output (I/O) operations.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"iplink_io", "TThe total input/output (I/O) operations that are transferred per second for IP partnership traffic on the system.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"iplink_comp_mb", "The average number of compressed megabytes transferred per second (MBps) over the IP Replication link during the sample period. This value is calculated after any compression of |the data takes place. This value does not include iSCSI host I/O operations.", labelnames, nil),

		prometheus.NewDesc(prefix_nodeStats+"cloud_up_mb", "The average number of megabytes transferred per second (Mbps) for upload operations to a cloud account during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"cloud_up_ms", "The average amount of time (in milliseconds) it takes for the system to respond to upload requests to a cloud account during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"cloud_down_mb", "The average number of Mbps for download operations to a cloud account during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"cloud_down_ms", "The average amount of time (in milliseconds) it takes for the system to respond to download requests to a cloud account during the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_nodeStats+"iser_mb", "The total number of megabytes transferred per second (MBps) for iSER traffic on the system.", labelnames, nil),
		prometheus.NewDesc(prefix_nodeStats+"iser_io", "The total I/O operations that are transferred per second for iSER traffic on the system.", labelnames, nil),
	}

	return &nodeStatsCollector{}, nil
}

// Describe describes the metrics
func (*nodeStatsCollector) Describe(ch chan<- *prometheus.Desc) {

	for _, nodestat_metric := range nodeStats_metrics {
		ch <- nodestat_metric
	}

}

// Collect collects metrics from Spectrum Virtualize Restful API
func (c *nodeStatsCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	logger.Debugln("entering NodeStats collector ...")
	nodeStatsResp, err := sClient.CallSpectrumAPI("lsnodestats", true)
	if err != nil {
		logger.Errorf("Executing lsnodestats cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("response of lsnodestats: ", nodeStatsResp)
	// This is a sample output of lsnodestats
	// [
	// {
	//     "node_id": "1",
	//     "node_name": "node1",
	//     "stat_name": "compression_cpu_pc",
	//     "stat_current": "0",
	//     "stat_peak": "0",
	//     "stat_peak_time": "181217083654"
	// },
	// {
	//     "node_id": "1",
	//     "node_name": "node1",
	//     "stat_name": "cpu_pc",
	//     "stat_current": "1",
	//     "stat_peak": "1",
	//     "stat_peak_time": "181217083654"
	// },
	// {
	//     "node_id": "1",
	//     "node_name": "node1",
	//     "stat_name": "fc_mb",
	//     "stat_current": "0",
	//     "stat_peak": "0",
	//     "stat_peak_time": "181217083654"
	// },
	// ....
	// ....
	// {
	//     "node_id": "2",
	//     "node_name": "node2",
	//     "stat_name": "mdisk_io",
	//     "stat_current": "0",
	//     "stat_peak": "0",
	//     "stat_peak_time": "181217083656"
	// },
	// {
	//     "node_id": "2",
	//     "node_name": "node2",
	//     "stat_name": "mdisk_ms",
	//     "stat_current": "0",
	//     "stat_peak": "0",
	//     "stat_peak_time": "181217083656"
	// },
	// {
	//     "node_id": "2",
	//     "node_name": "node2",
	//     "stat_name": "drive_mb",
	//     "stat_current": "232",
	//     "stat_peak": "422",
	//     "stat_peak_time": "181217083651"
	// }
	// ....
	// ....
	// ]

	nodeStatsArray := gjson.Parse(nodeStatsResp).Array()
	nodesNumber := len(nodeStatsArray) / len(nodeStats_metrics)

	for i, nodeStats_metric := range nodeStats_metrics {
		for node := 0; node < nodesNumber; node++ {
			index := len(nodeStats_metrics)*node + i
			labelvalues := []string{sClient.Hostname, nodeStatsArray[index].Get("node_name").String()}
			ch <- prometheus.MustNewConstMetric(nodeStats_metric, prometheus.GaugeValue, nodeStatsArray[index].Get("stat_current").Float(), labelvalues...)
		}
	}
	logger.Debugln("exit NodeStats collector")
	return nil
}
