package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_mdiskgrp = "spectrum_mdiskgrp_"

var (
	mdiskgrp_capacity                          *prometheus.Desc
	extent_size                                *prometheus.Desc
	free_capacity                              *prometheus.Desc
	virtual_capacity                           *prometheus.Desc
	used_capacity                              *prometheus.Desc
	real_capacity                              *prometheus.Desc
	overallocation                             *prometheus.Desc
	mdiskgrp_compression_active                *prometheus.Desc
	mdiskgrp_compression_virtual_capacity      *prometheus.Desc
	mdiskgrp_compression_compressed_capacity   *prometheus.Desc
	mdiskgrp_compression_uncompressed_capacity *prometheus.Desc
	mdiskgrp_used_capacity_before_reduction    *prometheus.Desc
	mdiskgrp_used_capacity_after_reduction     *prometheus.Desc
	mdiskgrp_overhead_capacity                 *prometheus.Desc
	mdiskgrp_deduplication_capcacity_saving    *prometheus.Desc
	reclaimable_capacity                       *prometheus.Desc
)

func init() {
	registerCollector("lsmdiskgrp", defaultDisabled, NewMdiskgrpCollector)
	labelnames := []string{"target", "resource", "name", "status"}
	mdiskgrp_capacity = prometheus.NewDesc(prefix_mdiskgrp+"capacity", "The total amount of MDisk storage that is assigned to the storage pool..", labelnames, nil)
	extent_size = prometheus.NewDesc(prefix_mdiskgrp+"extent_size", "The sizes of the extents for this group", labelnames, nil)
	free_capacity = prometheus.NewDesc(prefix_mdiskgrp+"free_capacity", "The amount of MDisk storage that is immediately available. Additionally, reclaimable_capacity can eventually become available", labelnames, nil)
	virtual_capacity = prometheus.NewDesc(prefix_mdiskgrp+"virtual_capacity", "The total host mappable capacity of all volume copies in the storage pool.", labelnames, nil)
	used_capacity = prometheus.NewDesc(prefix_mdiskgrp+"used_capacity", "The amount of data that is stored on MDisks.", labelnames, nil)
	real_capacity = prometheus.NewDesc(prefix_mdiskgrp+"real_capacity", "The total MDisk storage capacity assigned to volume copies.", labelnames, nil)
	overallocation = prometheus.NewDesc(prefix_mdiskgrp+"overallocation", "The ratio of the virtual_capacity value to the capacity", labelnames, nil)
	mdiskgrp_compression_active = prometheus.NewDesc(prefix_mdiskgrp+"compression_active", "Indicates whether any compressed volume copies are in the storage pool.", labelnames, nil)
	mdiskgrp_compression_virtual_capacity = prometheus.NewDesc(prefix_mdiskgrp+"compression_virtual_capacity", "The total virtual capacity for all compressed volume copies in regular storage pools. ", labelnames, nil)
	mdiskgrp_compression_compressed_capacity = prometheus.NewDesc(prefix_mdiskgrp+"compression_compressed_capacity", "The total used capacity for all compressed volume copies in regular storage pools.", labelnames, nil)
	mdiskgrp_compression_uncompressed_capacity = prometheus.NewDesc(prefix_mdiskgrp+"compression_uncompressed_capacity", "the total uncompressed used capacity for all compressed volume copies in regular storage pools", labelnames, nil)
	mdiskgrp_used_capacity_before_reduction = prometheus.NewDesc(prefix_mdiskgrp+"used_capacity_before_reduction", "The data that is stored on non-fully-allocated volume copies in a data reduction pool.", labelnames, nil)
	mdiskgrp_used_capacity_after_reduction = prometheus.NewDesc(prefix_mdiskgrp+"used_capacity_after_reduction", "The data that is stored on MDisks for non-fully-allocated volume copies in a data reduction pool.", labelnames, nil)
	mdiskgrp_overhead_capacity = prometheus.NewDesc(prefix_mdiskgrp+"overhead_capacity", "The MDisk capacity that is reserved for internal usage.", labelnames, nil)
	mdiskgrp_deduplication_capcacity_saving = prometheus.NewDesc(prefix_mdiskgrp+"deduplication_capcacity_saving", "The capacity that is saved by deduplication before compression in a data reduction pool.", labelnames, nil)
	reclaimable_capacity = prometheus.NewDesc(prefix_mdiskgrp+"reclaimable_capacity", "The MDisk capacity that is reserved for internal usage.", labelnames, nil)

}

//mdiskgrpCollector collects mdisk metrics
type mdiskgrpCollector struct {
}

func NewMdiskgrpCollector() (Collector, error) {
	return &mdiskgrpCollector{}, nil
}

//Describe describes the metrics
func (*mdiskgrpCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- mdiskgrp_capacity
	ch <- extent_size
	ch <- free_capacity
	ch <- virtual_capacity
	ch <- used_capacity
	ch <- real_capacity
	ch <- overallocation
	ch <- mdiskgrp_compression_active
	ch <- mdiskgrp_compression_virtual_capacity
	ch <- mdiskgrp_compression_compressed_capacity
	ch <- mdiskgrp_compression_uncompressed_capacity
	ch <- mdiskgrp_used_capacity_before_reduction
	ch <- mdiskgrp_used_capacity_after_reduction
	ch <- mdiskgrp_overhead_capacity
	ch <- mdiskgrp_deduplication_capcacity_saving
	ch <- reclaimable_capacity
}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *mdiskgrpCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	logger.Debugln("Entering MDiskgrp collector ...")
	mDiskGrpResp, err := sClient.CallSpectrumAPI("lsmdiskgrp", true)
	if err != nil {
		logger.Errorf("Executing lsmdiskgrp cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("Response of lsmdiskgrp: ", mDiskGrpResp)
	// This is a sample output of lsmdiskgrp
	// 	[
	//     {
	//         "id": "0",
	//         "name": "Pool0",
	//         "status": "online",
	//         "mdisk_count": "1",
	//         "vdisk_count": "5",
	//         "capacity": "99.01TB",
	//         "extent_size": "1024",
	//         "free_capacity": "98.46TB",
	//         "virtual_capacity": "656.00GB",
	//         "used_capacity": "556.00GB",
	//         "real_capacity": "558.02GB",
	//         "overallocation": "0",
	//         "warning": "80",
	//         "easy_tier": "auto",
	//         "easy_tier_status": "balanced",
	//         "compression_active": "no",
	//         "compression_virtual_capacity": "0.00MB",
	//         "compression_compressed_capacity": "0.00MB",
	//         "compression_uncompressed_capacity": "0.00MB",
	//         "parent_mdisk_grp_id": "0",
	//         "parent_mdisk_grp_name": "Pool0",
	//         "child_mdisk_grp_count": "0",
	//         "child_mdisk_grp_capacity": "0.00MB",
	//         "type": "parent",
	//         "encrypt": "no",
	//         "owner_type": "none",
	//         "site_id": "",
	//         "site_name": "",
	//         "data_reduction": "no",
	//         "used_capacity_before_reduction": "0.00MB",
	//         "used_capacity_after_reduction": "0.00MB",
	//         "overhead_capacity": "0.00MB",
	//         "deduplication_capacity_saving": "0.00MB",
	//         "reclaimable_capacity": "0.00MB"
	//     }
	// ]

	mDiskGrpArray := gjson.Parse(mDiskGrpResp).Array()
	for _, mdiskgrp := range mDiskGrpArray {
		mdiskgrp_capacity_bytes, err := utils.ToBytes(mdiskgrp.Get("capacity").String())
		if err != nil {
			logger.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		ch <- prometheus.MustNewConstMetric(mdiskgrp_capacity, prometheus.GaugeValue, float64(mdiskgrp_capacity_bytes), sClient.IpAddress, sClient.Hostname, mdiskgrp.Get("name").String(), mdiskgrp.Get("status").String())

		extent_size_bytes, err := utils.ToBytes(mdiskgrp.Get("extent_size").String() + "MB")
		if err != nil {
			logger.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		ch <- prometheus.MustNewConstMetric(extent_size, prometheus.GaugeValue, float64(extent_size_bytes), sClient.IpAddress, sClient.Hostname, mdiskgrp.Get("name").String(), mdiskgrp.Get("status").String())

		free_capacity_bytes, err := utils.ToBytes(mdiskgrp.Get("free_capacity").String())
		if err != nil {
			logger.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		ch <- prometheus.MustNewConstMetric(free_capacity, prometheus.GaugeValue, float64(free_capacity_bytes), sClient.IpAddress, sClient.Hostname, mdiskgrp.Get("name").String(), mdiskgrp.Get("status").String())

		virtual_capacity_bytes, err := utils.ToBytes(mdiskgrp.Get("virtual_capacity").String())
		if err != nil {
			logger.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		ch <- prometheus.MustNewConstMetric(virtual_capacity, prometheus.GaugeValue, float64(virtual_capacity_bytes), sClient.IpAddress, sClient.Hostname, mdiskgrp.Get("name").String(), mdiskgrp.Get("status").String())

		used_capacity_bytes, err := utils.ToBytes(mdiskgrp.Get("used_capacity").String())
		if err != nil {
			logger.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		ch <- prometheus.MustNewConstMetric(used_capacity, prometheus.GaugeValue, float64(used_capacity_bytes), sClient.IpAddress, sClient.Hostname, mdiskgrp.Get("name").String(), mdiskgrp.Get("status").String())

		real_capacity_bytes, err := utils.ToBytes(mdiskgrp.Get("real_capacity").String())
		if err != nil {
			logger.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		ch <- prometheus.MustNewConstMetric(real_capacity, prometheus.GaugeValue, float64(real_capacity_bytes), sClient.IpAddress, sClient.Hostname, mdiskgrp.Get("name").String(), mdiskgrp.Get("status").String())

		overallocation_pc, err := strconv.ParseFloat(mdiskgrp.Get("overallocation").String(), 64)
		if err != nil {
			logger.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		ch <- prometheus.MustNewConstMetric(overallocation, prometheus.GaugeValue, float64(overallocation_pc), sClient.IpAddress, sClient.Hostname, mdiskgrp.Get("name").String(), mdiskgrp.Get("status").String())

		mdiskgrp_compression_active_value, err := utils.ToBool(mdiskgrp.Get("compression_active").String())
		if err != nil {
			logger.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		ch <- prometheus.MustNewConstMetric(mdiskgrp_compression_active, prometheus.GaugeValue, mdiskgrp_compression_active_value, sClient.IpAddress, sClient.Hostname, mdiskgrp.Get("name").String(), mdiskgrp.Get("status").String())

		mdiskgrp_compression_virtual_capacity_bytes, err := utils.ToBytes(mdiskgrp.Get("compression_virtual_capacity").String())
		if err != nil {
			logger.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		ch <- prometheus.MustNewConstMetric(mdiskgrp_compression_virtual_capacity, prometheus.GaugeValue, float64(mdiskgrp_compression_virtual_capacity_bytes), sClient.IpAddress, sClient.Hostname, mdiskgrp.Get("name").String(), mdiskgrp.Get("status").String())

		mdiskgrp_compression_compressed_capacity_bytes, err := utils.ToBytes(mdiskgrp.Get("compression_compressed_capacity").String())
		if err != nil {
			logger.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		ch <- prometheus.MustNewConstMetric(mdiskgrp_compression_compressed_capacity, prometheus.GaugeValue, float64(mdiskgrp_compression_compressed_capacity_bytes), sClient.IpAddress, sClient.Hostname, mdiskgrp.Get("name").String(), mdiskgrp.Get("status").String())

		mdiskgrp_compression_uncompressed_capacity_bytes, err := utils.ToBytes(mdiskgrp.Get("compression_uncompressed_capacity").String())
		if err != nil {
			logger.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		ch <- prometheus.MustNewConstMetric(mdiskgrp_compression_uncompressed_capacity, prometheus.GaugeValue, float64(mdiskgrp_compression_uncompressed_capacity_bytes), sClient.IpAddress, sClient.Hostname, mdiskgrp.Get("name").String(), mdiskgrp.Get("status").String())

		mdiskgrp_used_capacity_before_reduction_bytes, err := utils.ToBytes(mdiskgrp.Get("used_capacity_before_reduction").String())
		if err != nil {
			logger.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		ch <- prometheus.MustNewConstMetric(mdiskgrp_used_capacity_before_reduction, prometheus.GaugeValue, float64(mdiskgrp_used_capacity_before_reduction_bytes), sClient.IpAddress, sClient.Hostname, mdiskgrp.Get("name").String(), mdiskgrp.Get("status").String())

		mdiskgrp_used_capacity_after_reduction_bytes, err := utils.ToBytes(mdiskgrp.Get("used_capacity_after_reduction").String())
		if err != nil {
			logger.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		ch <- prometheus.MustNewConstMetric(mdiskgrp_used_capacity_after_reduction, prometheus.GaugeValue, float64(mdiskgrp_used_capacity_after_reduction_bytes), sClient.IpAddress, sClient.Hostname, mdiskgrp.Get("name").String(), mdiskgrp.Get("status").String())

		mdiskgrp_overhead_capacity_bytes, err := utils.ToBytes(mdiskgrp.Get("overhead_capacity").String())
		if err != nil {
			logger.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		ch <- prometheus.MustNewConstMetric(mdiskgrp_overhead_capacity, prometheus.GaugeValue, float64(mdiskgrp_overhead_capacity_bytes), sClient.IpAddress, sClient.Hostname, mdiskgrp.Get("name").String(), mdiskgrp.Get("status").String())

		mdiskgrp_deduplication_capcacity_saving_bytes, err := utils.ToBytes(mdiskgrp.Get("deduplication_capacity_saving").String())
		if err != nil {
			logger.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		ch <- prometheus.MustNewConstMetric(mdiskgrp_deduplication_capcacity_saving, prometheus.GaugeValue, float64(mdiskgrp_deduplication_capcacity_saving_bytes), sClient.IpAddress, sClient.Hostname, mdiskgrp.Get("name").String(), mdiskgrp.Get("status").String())

		reclaimable_capacity_bytes, err := utils.ToBytes(mdiskgrp.Get("reclaimable_capacity").String())
		if err != nil {
			logger.Errorf("Converting capacity unit failed: %s", err.Error())
		}
		ch <- prometheus.MustNewConstMetric(reclaimable_capacity, prometheus.GaugeValue, float64(reclaimable_capacity_bytes), sClient.IpAddress, sClient.Hostname, mdiskgrp.Get("name").String(), mdiskgrp.Get("status").String())
	}
	logger.Debugln("Leaving MDiskgrp collector.")
	return err
}
