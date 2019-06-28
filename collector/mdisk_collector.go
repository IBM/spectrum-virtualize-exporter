package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/spectrum-virtualize-exporter/utils"
	"github.com/tidwall/gjson"
)

const prefix_mdisk = "spectrum_mdisk_"

var mdiskCapacity *prometheus.Desc

func init() {
	registerCollector("lsmdisk", defaultDisabled, NewMdiskCollector)
	labelnames := []string{"target", "resource", "name", "status", "mdisk_grp_name", "tier"}
	mdiskCapacity = prometheus.NewDesc(prefix_mdisk+"capacity", "The capacity of the MDisk by pool.", labelnames, nil)

}

//mdiskCollector collects mdisk metrics
type mdiskCollector struct {
}

func NewMdiskCollector() (Collector, error) {
	return &mdiskCollector{}, nil
}

//Describe describes the metrics
func (*mdiskCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- mdiskCapacity
}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *mdiskCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	log.Debugln("MDisk collector is starting")
	reqSystemURL := "https://" + sClient.IpAddress + ":7443/rest/lsmdisk"
	mDiskRes, err := sClient.CallSpectrumAPI(reqSystemURL)
	mDiskArray := gjson.Parse(mDiskRes).Array()
	for _, mdisk := range mDiskArray {
		capacity_bytes, err := utils.ToBytes(mdisk.Get("capacity").String())
		ch <- prometheus.MustNewConstMetric(mdiskCapacity, prometheus.GaugeValue, float64(capacity_bytes), sClient.IpAddress, sClient.Hostname, mdisk.Get("name").String(), mdisk.Get("status").String(), mdisk.Get("mdisk_grp_name").String(), mdisk.Get("tier").String())
		if err != nil {
			return err
		}

	}
	return err

}
