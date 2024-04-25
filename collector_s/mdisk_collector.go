package collector_s

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_mdisk = "spectrum_mdisk_"

var (
	mdisk_status *prometheus.Desc
)

func init() {
	registerCollector("lsmdisk_s", defaultEnabled, NewMdiskCollector)
}

// mdiskCollector collects mdisk metrics
type mdiskCollector struct {
}

func NewMdiskCollector() (Collector, error) {
	labelnames := []string{"resource", "pool_name", "mdisk_name"}
	if len(utils.ExtraLabelNames) > 0 {
		labelnames = append(labelnames, utils.ExtraLabelNames...)
	}
	mdisk_status = prometheus.NewDesc(prefix_mdisk+"status", "Status of managed disks (MDisks) visible to the system. 0-online; 1-offline; 2-excluded; 3-degraded_paths; 4-degraded_ports; 5-degraded.", labelnames, nil)
	return &mdiskCollector{}, nil
}

// Describe() describes the metrics
func (*mdiskCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- mdisk_status
}

// Collect() collects metrics from Spectrum Virtualize Restful API
func (c *mdiskCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	logger.Debugln("entering MDisk collector ...")
	respData, err := sClient.CallSpectrumAPI("lsmdisk", true)
	if err != nil {
		logger.Errorf("executing lsmdisk cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("response of lsmdisk: ", respData)
	/* This is a sample output of lsmdisk
	[
		{
			"id": "0",
			"name": "mdisk0",
			"status": "online",
			"mode": "array",
			"mdisk_grp_id": "0",
			"mdisk_grp_name": "Pool0",
			"capacity": "99.1TB",
			"ctrl_LUN_#": "",
			"controller_name": "",
			"UID": "",
			"tier": "tier0_flash",
			"encrypt": "no",
			"site_id": "",
			"site_name": "",
			"distributed": "yes",
			"dedupe": "no",
			"over_provisioned": "yes",
			"supports_unmap": "yes"
		}
	] */
	if !gjson.Valid(respData) {
		return fmt.Errorf("invalid json for lsmdisk:\n%v", respData)
	}
	jsonMDisks := gjson.Parse(respData)
	jsonMDisks.ForEach(func(key, port gjson.Result) bool {
		pool_name := port.Get("mdisk_grp_name").String()
		mdisk_name := port.Get("name").String()
		status := port.Get("status").String() // ["online", "offline", ...]

		for _, pool := range online_pools {
			if pool == pool_name {
				v_status := 0
				switch status {
				case "online":
					v_status = 0
				case "offline":
					v_status = 1
				case "excluded":
					v_status = 2
				case "degraded_paths":
					v_status = 3
				case "degraded_ports":
					v_status = 4
				case "degraded":
					v_status = 5
				}
				labelvalues := []string{sClient.Hostname, pool_name, mdisk_name}
				if len(utils.ExtraLabelValues) > 0 {
					labelvalues = append(labelvalues, utils.ExtraLabelValues...)
				}
				ch <- prometheus.MustNewConstMetric(mdisk_status, prometheus.GaugeValue, float64(v_status), labelvalues...)
			}
		}
		return true
	})

	logger.Debugln("exit MDisk exit")
	return nil
}
