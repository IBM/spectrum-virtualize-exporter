package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
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

	log.Debugln("Entering MDisk collector ...")
	mDiskResp, err := sClient.CallSpectrumAPI("lsmdisk", true)
	if err != nil {
		log.Errorf("Executing lsmdisk cmd failed: %s", err.Error())
		return err
	}
	log.Debugln("Response of lsmdisk: ", mDiskResp)
	//This is a sample output of lsmdisk
	// 	[
	//     {
	//         "id": "0",
	//         "name": "mdisk0",
	//         "status": "online",
	//         "mode": "array",
	//         "mdisk_grp_id": "0",
	//         "mdisk_grp_name": "Pool0",
	//         "capacity": "99.1TB",
	//         "ctrl_LUN_#": "",
	//         "controller_name": "",
	//         "UID": "",
	//         "tier": "tier0_flash",
	//         "encrypt": "no",
	//         "site_id": "",
	//         "site_name": "",
	//         "distributed": "yes",
	//         "dedupe": "no",
	//         "over_provisioned": "yes",
	//         "supports_unmap": "yes"
	//     }
	// ]
	mDisks := gjson.Parse(mDiskResp).Array()
	for _, mdisk := range mDisks {
		capacity_bytes, err := utils.ToBytes(mdisk.Get("capacity").String())
		if err != nil {
			log.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		ch <- prometheus.MustNewConstMetric(mdiskCapacity, prometheus.GaugeValue, float64(capacity_bytes), sClient.IpAddress, sClient.Hostname, mdisk.Get("name").String(), mdisk.Get("status").String(), mdisk.Get("mdisk_grp_name").String(), mdisk.Get("tier").String())

	}
	log.Debugln("Leaving MDisk collector.")
	return err

}
