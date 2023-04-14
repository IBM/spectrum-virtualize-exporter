package collector_s

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_mdiskgrp = "spectrum_mdiskgrp_"

var (
	mdiskgrp_status *prometheus.Desc
	online_pools    []string
)

func init() {
	registerCollector("lsmdiskgrp_s", defaultEnabled, NewMdiskgrpCollector)
}

//mdiskgrpCollector collects mdisk metrics
type mdiskgrpCollector struct {
}

func NewMdiskgrpCollector() (Collector, error) {
	labelnames := []string{"resource", "pool_name"}
	if len(utils.ExtraLabelNames) > 0 {
		labelnames = append(labelnames, utils.ExtraLabelNames...)
	}
	mdiskgrp_status = prometheus.NewDesc(prefix_mdiskgrp+"status", "Status of storage pools that are visible to the system. 0-online; 1-offline; 2-others.", labelnames, nil)
	return &mdiskgrpCollector{}, nil
}

//Describe() describes the metrics
func (*mdiskgrpCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- mdiskgrp_status
}

//Collect() collects metrics from Spectrum Virtualize Restful API
func (c *mdiskgrpCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	logger.Debugln("Entering MDiskgrp collector ...")
	respData, err := sClient.CallSpectrumAPI("lsmdiskgrp", true)
	if err != nil {
		logger.Errorf("Executing lsmdiskgrp cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("Response of lsmdiskgrp: ", respData)
	/* This is a sample output of lsmdiskgrp
	[
		{
			"id": "0",
			"name": "Pool0",
			"status": "online",
			"mdisk_count": "1",
			"vdisk_count": "114",
			"capacity": "99.01TB",
			"extent_size": "1024",
			"free_capacity": "36.66TB",
			"virtual_capacity": "62.35TB",
			"used_capacity": "62.35TB",
			"real_capacity": "62.35TB",
			"overallocation": "62",
			"warning": "80",
			"easy_tier": "auto",
			"easy_tier_status": "balanced",
			"compression_active": "no",
			"compression_virtual_capacity": "0.00MB",
			"compression_compressed_capacity": "0.00MB",
			"compression_uncompressed_capacity": "0.00MB",
			"parent_mdisk_grp_id": "0",
			"parent_mdisk_grp_name": "Pool0",
			"child_mdisk_grp_count": "0",
			"child_mdisk_grp_capacity": "0.00MB",
			"type": "parent",
			"encrypt": "no",
			"owner_type": "none",
			"owner_id": "",
			"owner_name": "",
			"site_id": "",
			"site_name": "",
			"data_reduction": "no",
			"used_capacity_before_reduction": "0.00MB",
			"used_capacity_after_reduction": "0.00MB",
			"overhead_capacity": "0.00MB",
			"deduplication_capacity_saving": "0.00MB",
			"reclaimable_capacity": "0.00MB",
			"easy_tier_fcm_over_allocation_max": ""
		}
	] */
	if !gjson.Valid(respData) {
		return fmt.Errorf("invalid json for lsmdiskgrp:\n%v", respData)
	}
	jsonPools := gjson.Parse(respData)
	jsonPools.ForEach(func(key, port gjson.Result) bool {
		pool_name := port.Get("name").String()
		status := port.Get("status").String() // ["online", "offline", ...]

		v_status := 0
		switch status {
		case "online":
			v_status = 0
			found := false
			for _, pool := range online_pools {
				if pool == pool_name {
					found = true
				}
			}
			if !found {
				online_pools = append(online_pools, pool_name)
				logger.Debugf("Appended the '%s' into the online_pools.", pool_name)
			} else {
				logger.Debugf("The '%s' already exists in the online_pools.", pool_name)
			}
		case "offline":
			v_status = 1
		default:
			v_status = 2
		}

		labelvalues := []string{sClient.Hostname, pool_name}
		if len(utils.ExtraLabelValues) > 0 {
			labelvalues = append(labelvalues, utils.ExtraLabelValues...)
		}

		ch <- prometheus.MustNewConstMetric(mdiskgrp_status, prometheus.GaugeValue, float64(v_status), labelvalues...)
		return true
	})

	logger.Debugln("Leaving MDiskgrp collector.")
	return nil
}
