package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_volumeCopy = "spectrum_volume_copy"

var (
	volumeCopy_Capacity *prometheus.Desc
)

func init() {
	registerCollector("lsvdiskcopy", defaultDisabled, NewVolumeCopyCollector)
}

//volumeCopyCollector collects volume cpoy metrics
type volumeCopyCollector struct {
}

func NewVolumeCopyCollector() (Collector, error) {
	labelnames := []string{"resource", "volume_id", "volume_name", "copy_id", "mdisk_grp_name"}
	if len(utils.ExtraLabelNames) > 0 {
		labelnames = append(labelnames, utils.ExtraLabelNames...)
	}
	volumeCopy_Capacity = prometheus.NewDesc(prefix_volumeCopy+"_capacity", "The capacity of the volume copy.", labelnames, nil)

	return &volumeCopyCollector{}, nil
}

//Describe describes the metrics
func (*volumeCopyCollector) Describe(ch chan<- *prometheus.Desc) {

	ch <- volumeCopy_Capacity

}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *volumeCopyCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {
	logger.Debugln("Entering volumeCopy collector ...")
	volumeCopyResp, err := sClient.CallSpectrumAPI("lsvdiskcopy", true)
	if err != nil {
		logger.Errorf("Executing lsvdiskcopy cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("Response of lsvdiskcopy: ", volumeCopyResp)
	// This is a sample output of lsvdiskcopy
	// [
	// {
	//     "vdisk_id": "0",
	//     "vdisk_name": "MGMT1_MGMT1-boot",
	//     "copy_id": "0",
	//     "status": "online",
	//     "sync": "yes",
	//     "primary": "yes",
	//     "mdisk_grp_id": "0",
	//     "mdisk_grp_name": "Pool0",
	//     "capacity": "128.00GB",
	//     "type": "striped",
	//     "se_copy": "no",
	//     "easy_tier": "on",
	//     "easy_tier_status": "balanced",
	//     "compressed_copy": "no",
	//     "parent_mdisk_grp_id": "0",
	//     "parent_mdisk_grp_name": "Pool0",
	//     "encrypt": "no",
	//     "deduplicated_copy": "no"
	// }
	// ]
	volumeCopyArray := gjson.Parse(volumeCopyResp).Array()
	for _, volumeCopy := range volumeCopyArray {
		volumeCopy_capacity_bytes, err := utils.ToBytes(volumeCopy.Get("capacity").String())
		if err != nil {
			logger.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		labelvalues := []string{sClient.Hostname, volumeCopy.Get("vdisk_id").String(), volumeCopy.Get("vdisk_name").String(), volumeCopy.Get("copy_id").String(), volumeCopy.Get("mdisk_grp_name").String()}
		if len(utils.ExtraLabelValues) > 0 {
			labelvalues = append(labelvalues, utils.ExtraLabelValues...)
		}
		ch <- prometheus.MustNewConstMetric(volumeCopy_Capacity, prometheus.GaugeValue, float64(volumeCopy_capacity_bytes), labelvalues...)
	}
	logger.Debugln("Leaving volumeCopy collector.")
	return nil
}
