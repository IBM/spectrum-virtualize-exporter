package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/spectrum-virtualize-exporter/utils"
	"github.com/tidwall/gjson"
)

const prefix_volume = "spectrum_volume_"

var (
	volumeCapacity *prometheus.Desc
)

func init() {
	registerCollector("lsvdisk", defaultDisabled, NewVolumeCollector)
	labelnames := []string{"target", "resource", "volume_id", "volume_name", "mdisk_grp_name"}
	volumeCapacity = prometheus.NewDesc(prefix_volume+"capacity", "The virtual capacity of the volume that is the size of the volume as seen by the host.", labelnames, nil)
}

//volumeCollector collects vdisk metrics
type volumeCollector struct {
}

func NewVolumeCollector() (Collector, error) {
	return &volumeCollector{}, nil
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
	for _, volume := range volumeArray {
		capacity_bytes, err := utils.ToBytes(volume.Get("capacity").String())
		ch <- prometheus.MustNewConstMetric(volumeCapacity, prometheus.GaugeValue, float64(capacity_bytes), sClient.IpAddress, sClient.Hostname, volume.Get("volume_id").String(), volume.Get("volume_name").String(), volume.Get("mdisk_grp_name").String())
		if err != nil {
			return err
		}
	}

	return err

}
