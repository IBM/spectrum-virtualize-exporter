package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/spectrum-virtualize-exporter/utils"
	"github.com/tidwall/gjson"
)

const prefix_stats = "spectrum_systemstats_"

var (
	metrics    [49]*prometheus.Desc
	metricHelp []string
)

type systemStatsCollector struct {
}

func init() {
	registerCollector("lssystemstats", defaultEnabled, NewSystemStatsCollector)
	labelnames := []string{"target"}
	metrics = [49]*prometheus.Desc{
		prometheus.NewDesc(prefix_stats+"compression_cpu_pc", "The percentage of allocated CPU capacity that is used for compression.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"cpu_pc", "The percentage of allocated CPU capacity that is used for the system.", labelnames, nil),

		prometheus.NewDesc(prefix_stats+"fc_mb", "The total number of megabytes transferred per second (MBps) for Fibre Channel traffic on the system. This value includes host I/O and any bandwidth that is used for communication within the system.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"fc_io", "The total input/output (I/O) operations that are transferred per seconds for Fibre Channel traffic on the system. This value includes host I/O and any bandwidth that is used for communication within the system.", labelnames, nil),

		prometheus.NewDesc(prefix_stats+"sas_mb", "The total number of megabytes transferred per second (MBps) for serial-attached SCSI (SAS) traffic on the system. This value includes host I/O and bandwidth that is used for background RAID activity.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"sas_io", "The total I/O operations that are transferred per second for SAS traffic on the system. This value includes host I/O and bandwidth that is used for background RAID activity.", labelnames, nil),

		prometheus.NewDesc(prefix_stats+"iscsi_mb", "The total number of megabytes transferred per second (MBps) for iSCSI traffic on the system.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"iscsi_io", "The total I/O operations that are transferred per second for iSCSI traffic on the system.", labelnames, nil),

		prometheus.NewDesc(prefix_stats+"write_cache_pc", "The percentage of the write cache usage for the node.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"total_cache_pc", "The total percentage for both the write and read cache usage for the node.", labelnames, nil),

		prometheus.NewDesc(prefix_stats+"vdisk_mb", "The average number of megabytes transferred per second (MBps) for read and write operations to volumes during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"vdisk_io", "The average number of I/O operations that are transferred per second for read and write operations to volumes during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"vdisk_ms", "The average amount of time in milliseconds that the system takes to respond to read and write requests to volumes over the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_stats+"mdisk_mb", "The average number of megabytes transferred per second (MBps) for read and write operations to MDisks during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"mdisk_io", "The average number of I/O operations that are transferred per second for read and write operations to MDisks during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"mdisk_ms", "The average amount of time in milliseconds that the system takes to respond to read and write requests to MDisks over the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_stats+"drive_mb", "The average number of megabytes transferred per second (MBps) for read and write operations to drives during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"drive_io", "The average number of I/O operations that are transferred per second for read and write operations to drives during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"drive_ms", "The average amount of time in milliseconds that the system takes to respond to read and write requests to drives over the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_stats+"vdisk_r_mb", "The average number of megabytes transferred per second (MBps) for read operations to volumes during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"vdisk_r_io", "The average number of I/O operations that are transferred per second for read operations to volumes during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"vdisk_r_ms", "The average amount of time in milliseconds that the system takes to respond to read requests to volumes over the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_stats+"vdisk_w_mb", "The average number of megabytes transferred per second (MBps) for read and write operations to volumes during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"vdisk_w_io", "The average number of I/O operations that are transferred per second for write operations to volumes during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"vdisk_w_ms", "The average amount of time in milliseconds that the system takes to respond to write requests to volumes over the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_stats+"mdisk_r_mb", "The average number of megabytes transferred per second (MBps) for read operations to MDisks during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"mdisk_r_io", "The average number of I/O operations that are transferred per second for read operations to MDisks during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"mdisk_r_ms", "The average amount of time in milliseconds that the system takes to respond to read requests to MDisks over the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_stats+"mdisk_w_mb", "The average number of megabytes transferred per second (MBps) for write operations to MDisks during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"mdisk_w_io", "TThe average number of I/O operations that are transferred per second for write operations to MDisks during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"mdisk_w_ms", "the average amount of time in milliseconds that the system takes to respond to write requests to MDisks over the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_stats+"drive_r_mb", "The average number of megabytes transferred per second (MBps) for read operations to drives during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"drive_r_io", "The average number of I/O operations that are transferred per second for read operations to drives during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"drive_r_ms", "The average amount of time in milliseconds that the system takes to respond to read requests to drives over the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_stats+"drive_w_mb", "The average number of megabytes transferred per second (MBps) for write operations to drives during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"drive_w_io", "The average number of I/O operations that are transferred per second for write operations to drives during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"drive_w_ms", "The average amount of time in milliseconds that the system takes to respond write requests to drives over the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_stats+"power_w", "the power that is consumed in watts.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"temp_c", " the ambient temperature in Celsius.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"temp_f", "the ambient temperature in Fahrenheit.", labelnames, nil),

		prometheus.NewDesc(prefix_stats+"iplink_mb", "The average number of megabytes requested to be transferred per second (MBps) over the IP partnership link during the sample period. This value is calculated before any compression of the data takes place. This value does not include iSCSI host input/output (I/O) operations.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"iplink_io", "TThe total input/output (I/O) operations that are transferred per second for IP partnership traffic on the system.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"iplink_comp_mb", "The average number of compressed megabytes transferred per second (MBps) over the IP Replication link during the sample period. This value is calculated after any compression of |the data takes place. This value does not include iSCSI host I/O operations.", labelnames, nil),

		prometheus.NewDesc(prefix_stats+"cloud_up_mb", "The average number of megabytes transferred per second (Mbps) for upload operations to a cloud account during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"cloud_up_ms", "The average amount of time (in milliseconds) it takes for the system to respond to upload requests to a cloud account during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"cloud_down_mb", "The average number of Mbps for download operations to a cloud account during the sample period.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"cloud_down_ms", "The average amount of time (in milliseconds) it takes for the system to respond to download requests to a cloud account during the sample period.", labelnames, nil),

		prometheus.NewDesc(prefix_stats+"iser_mb", "The total number of megabytes transferred per second (MBps) for iSER traffic on the system.", labelnames, nil),
		prometheus.NewDesc(prefix_stats+"iser_io", "The total I/O operations that are transferred per second for iSER traffic on the system.", labelnames, nil),
	}

}
func NewSystemStatsCollector() (Collector, error) {
	return &systemStatsCollector{}, nil
}

//Describe describes the metrics
func (*systemStatsCollector) Describe(ch chan<- *prometheus.Desc) {

	for _, metric := range metrics {
		ch <- metric
	}

}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *systemStatsCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {
	log.Debugln("SystemStats collector is starting")
	labelvalues := []string{sClient.IpAddress}
	reqSystemURL := "https://" + sClient.IpAddress + ":7443/rest/lssystemstats"
	systemStats, err := sClient.CallSpectrumAPI(reqSystemURL)
	systemStatsMetrics := gjson.Parse(systemStats).Array()
	for i, systemStatsMetric := range systemStatsMetrics {
		ch <- prometheus.MustNewConstMetric(metrics[i], prometheus.GaugeValue, systemStatsMetric.Get("stat_current").Float(), labelvalues...)

	}
	return err

}
