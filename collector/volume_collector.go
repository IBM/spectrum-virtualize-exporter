package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_volume = "spectrum_volume_"

var (
	volumeCapacity *prometheus.Desc
)

func init() {
	registerCollector("lsvdisk", defaultDisabled, NewVolumeCollector)
}

//volumeCollector collects vdisk metrics
type volumeCollector struct {
}

func NewVolumeCollector() (Collector, error) {
	labelnames := []string{"resource", "volume_id", "volume_name", "mdisk_grp_name"}
	if len(utils.ExtraLabelNames) > 0 {
		labelnames = append(labelnames, utils.ExtraLabelNames...)
	}
	volumeCapacity = prometheus.NewDesc(prefix_volume+"capacity", "The virtual capacity of the volume that is the size of the volume as seen by the host.", labelnames, nil)

	return &volumeCollector{}, nil
}

//Describe describes the metrics
func (*volumeCollector) Describe(ch chan<- *prometheus.Desc) {

	ch <- volumeCapacity

}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *volumeCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {
	logger.Debugln("entering volume collector ...")
	volumeResp, err := sClient.CallSpectrumAPI("lsvdisk", true)
	if err != nil {
		logger.Errorf("Executing lsvdisk cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("response of lsvdisk: ", volumeResp)
	// This is a sample output of lsvdisk
	// 	[
	//     {
	//         "id": "0",
	//         "name": "MGMT1_MGMT1-boot",
	//         "IO_group_id": "0",
	//         "IO_group_name": "io_grp0",
	//         "status": "online",
	//         "mdisk_grp_id": "0",
	//         "mdisk_grp_name": "Pool0",
	//         "capacity": "128.00GB",
	//         "type": "striped",
	//         "FC_id": "",
	//         "FC_name": "",
	//         "RC_id": "",
	//         "RC_name": "",
	//         "vdisk_UID": "600507681081001D4800000000000001",
	//         "fc_map_count": "0",
	//         "copy_count": "1",
	//         "fast_write_state": "empty",
	//         "se_copy_count": "0",
	//         "RC_change": "no",
	//         "compressed_copy_count": "0",
	//         "parent_mdisk_grp_id": "0",
	//         "parent_mdisk_grp_name": "Pool0",
	//         "formatting": "no",
	//         "encrypt": "no",
	//         "volume_id": "0",
	//         "volume_name": "MGMT1_MGMT1-boot",
	//         "function": ""
	//     }
	// ]

	volumeArray := gjson.Parse(volumeResp).Array()
	for _, volume := range volumeArray {
		capacity_bytes, err := utils.ToBytes(volume.Get("capacity").String())
		if err != nil {
			logger.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		labelvalues := []string{sClient.Hostname, volume.Get("volume_id").String(), volume.Get("volume_name").String(), volume.Get("mdisk_grp_name").String()}
		if len(utils.ExtraLabelValues) > 0 {
			labelvalues = append(labelvalues, utils.ExtraLabelValues...)
		}
		ch <- prometheus.MustNewConstMetric(volumeCapacity, prometheus.GaugeValue, float64(capacity_bytes), labelvalues...)
	}
	logger.Debugln("exit volume collector")
	return nil
}
