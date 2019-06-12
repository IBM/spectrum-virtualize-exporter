package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/spectrum-virtualize-exporter/utils"
	"github.com/tidwall/gjson"
)

const prefix_volumeCopy = "spectrum_volume_copy"

var (
	volumeCopy_Capacity *prometheus.Desc
)

func init() {
	registerCollector("lsvdiskcopy", defaultDisabled, NewVolumeCopyCollector)
	labelnames := []string{"resource", "volume_id", "volume_name", "copy_id", "mdisk_grp_name"}
	volumeCopy_Capacity = prometheus.NewDesc(prefix_volumeCopy+"_capacity", "The capacity of the volume copy.", labelnames, nil)

}

//volumeCopyCollector collects volume cpoy metrics
type volumeCopyCollector struct {
}

func NewVolumeCopyCollector() (Collector, error) {
	return &volumeCopyCollector{}, nil
}

//Describe describes the metrics
func (*volumeCopyCollector) Describe(ch chan<- *prometheus.Desc) {

	ch <- volumeCopy_Capacity

}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *volumeCopyCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {
	log.Debugln("volume copy collector is starting")
	reqSystemURL := "https://" + sClient.IpAddress + ":7443/rest/lsvdiskcopy"
	volumeCopyRes, err := sClient.CallSpectrumAPI(reqSystemURL)
	volumeCopyArray := gjson.Parse(volumeCopyRes).Array()
	for _, volumeCopy := range volumeCopyArray {
		volumeCopy_capacity_bytes, err := utils.ToBytes(volumeCopy.Get("capacity").String())
		ch <- prometheus.MustNewConstMetric(volumeCopy_Capacity, prometheus.GaugeValue, float64(volumeCopy_capacity_bytes), sClient.Hostname, volumeCopy.Get("vdisk_id").String(), volumeCopy.Get("vdisk_name").String(), volumeCopy.Get("copy_id").String(), volumeCopy.Get("mdisk_grp_name").String())
		if err != nil {
			return err
		}
	}

	return err

}
