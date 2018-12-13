package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/spectrum-virtualize-exporter/utils"
	"github.com/tidwall/gjson"
)

const prefix_volume = "spectrum_volume"

var (
	volumeCapacity *prometheus.Desc
)

func init() {
	labelnames := []string{"target", "volume_id", "volume_name"}
	volumeCapacity = prometheus.NewDesc(prefix_volume+"capacity", "The virtual capacity of the volume that is the size of the volume as seen by the host.", labelnames, nil)

}

//nodeStatsCollector collects vdisk metrics
type volumeCollector struct {
}

func NewVolumeCollector() Collector {
	return &volumeCollector{}
}

//Describe describes the metrics
func (*volumeCollector) Describe(ch chan<- *prometheus.Desc) {

	ch <- volumeCapacity

}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *volumeCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	log.Debugln("volume collector is starting")

	reqSystemURL := "https://" + sClient.IpAddress + ":7443/rest/lsvdisk"
	volumeRes, err := sClient.CallSpectrumAPI(reqSystemURL)
	volumeArray := gjson.Parse(volumeRes).Array()
	// nodeStats_metrics = make([]*prometheus.Desc, len(nodeStatsArray), len(nodeStatsArray))
	for _, volume := range volumeArray {

		capacity_bytes, errors := utils.ToBytes(volume.Get("capacity").String())
		ch <- prometheus.MustNewConstMetric(volumeCapacity, prometheus.GaugeValue, float64(capacity_bytes), sClient.IpAddress, volume.Get("volume_id").String(), volume.Get("volume_name").String())
		err = errors
	}

	return err

}
