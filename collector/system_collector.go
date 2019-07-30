package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/spectrum-virtualize-exporter/utils"
	"github.com/tidwall/gjson"
)

const prefix_sys = "spectrum_system_"

var (
	total_mdisk_capacity                 *prometheus.Desc
	space_in_mdisk_grps                  *prometheus.Desc
	space_allocated_to_vdisks            *prometheus.Desc
	total_free_space                     *prometheus.Desc
	total_vdiskcopy_capacity             *prometheus.Desc
	total_used_capacity                  *prometheus.Desc
	total_overallocation                 *prometheus.Desc
	total_vdisk_capacity                 *prometheus.Desc
	total_allocated_extent_capacity      *prometheus.Desc
	compression_virtual_capacity         *prometheus.Desc
	compression_compressed_capacity      *prometheus.Desc
	compression_uncompressed_capacity    *prometheus.Desc
	total_drive_raw_capacity             *prometheus.Desc
	tier0_flash_compressed_data_used     *prometheus.Desc
	tier1_flash_compressed_data_used     *prometheus.Desc
	tier_enterprise_compressed_data_used *prometheus.Desc
	tier_nearline_compressed_data_used   *prometheus.Desc
	total_reclaimable_capacity           *prometheus.Desc
	physical_capacity                    *prometheus.Desc
	physical_free_capacity               *prometheus.Desc
	used_capacity_before_reduction       *prometheus.Desc
	used_capacity_after_reduction        *prometheus.Desc
	overhead_capacity                    *prometheus.Desc
	deduplication_capacity_saving        *prometheus.Desc

	tier_capacity      *prometheus.Desc
	tier_free_capacity *prometheus.Desc

	physical_capacity_usage *prometheus.Desc
	volume_capacity_usage   *prometheus.Desc
	mdiskgrp_capacity_usage *prometheus.Desc
	hostname                string
)

func init() {
	registerCollector("lssystem", defaultEnabled, NewSystemCollector)
	labelnames := []string{"target", "resource"}
	total_mdisk_capacity = prometheus.NewDesc(prefix_sys+"total_mdisk_capacity", "The sum of mdiskgrp capacity plus the capacity of all unmanaged MDisks", labelnames, nil)
	space_in_mdisk_grps = prometheus.NewDesc(prefix_sys+"space_in_mdisk_grps", "The sum of mdiskgrp capacity", labelnames, nil)
	space_allocated_to_vdisks = prometheus.NewDesc(prefix_sys+"space_allocated_to_vdisks", "The sum of mdiskgrp real_capacity", labelnames, nil)
	total_free_space = prometheus.NewDesc(prefix_sys+"total_free_space", "The sum of mdiskgrp free_capacity", labelnames, nil)
	total_vdiskcopy_capacity = prometheus.NewDesc(prefix_sys+"total_vdiskcopy_capacity", "The total virtual capacity of all volume copies in the cluster", labelnames, nil)
	total_used_capacity = prometheus.NewDesc(prefix_sys+"total_used_capacity", "The sum of mdiskgrp used_capacity", labelnames, nil)
	total_overallocation = prometheus.NewDesc(prefix_sys+"total_overallocation_percent", "The total_vdiskcopy_capacity as a percentage of total_mdisk_capacity. If total_mdisk_capacity is zero, then total_overallocation should display 100", labelnames, nil)
	total_vdisk_capacity = prometheus.NewDesc(prefix_sys+"total_vdisk_capacity", "The total virtual capacity of volumes in the cluster", labelnames, nil)
	total_allocated_extent_capacity = prometheus.NewDesc(prefix_sys+"total_allocated_extent_capacity", "The total size of all extents that are allocated to VDisks or otherwise in use by the system.", labelnames, nil)
	compression_virtual_capacity = prometheus.NewDesc(prefix_sys+"compression_virtual_capacity", "The total virtual capacity for all compressed volume copies in non-data reduction pools. Compressed volumes that are in data reduction pools do not count towards this value. This value is in unsigned decimal format.", labelnames, nil)
	compression_compressed_capacity = prometheus.NewDesc(prefix_sys+"compression_compressed_capacity", "The total used capacity for all compressed volume copies in non-data reduction pools.", labelnames, nil)
	compression_uncompressed_capacity = prometheus.NewDesc(prefix_sys+"compression_uncompressed_capacity", "The total uncompressed used capacity for all compressed volume copies in non-data reduction pools", labelnames, nil)
	total_drive_raw_capacity = prometheus.NewDesc(prefix_sys+"total_drive_raw_capacity", "The total known capacity of all discovered drives (regardless of drive use)", labelnames, nil)
	tier0_flash_compressed_data_used = prometheus.NewDesc(prefix_sys+"tier0_flash_compressed_data_used", "The capacity of compressed data used on the flash tier 0 storage tier", labelnames, nil)
	tier1_flash_compressed_data_used = prometheus.NewDesc(prefix_sys+"tier1_flash_compressed_data_used", "The capacity of compressed data used on the flash tier 1 storage tier.", labelnames, nil)
	tier_enterprise_compressed_data_used = prometheus.NewDesc(prefix_sys+"tier_enterprise_compressed_data_used", "The capacity of compressed data that is used on the tier 2 enterprise storage tier.", labelnames, nil)
	tier_nearline_compressed_data_used = prometheus.NewDesc(prefix_sys+"tier_nearline_compressed_data_used", "The capacity of compressed data that is used on the tier 3 nearline storage tier.", labelnames, nil)
	total_reclaimable_capacity = prometheus.NewDesc(prefix_sys+"total_reclaimable_capacity", "The unused (free) capacity that will be available after data is reduced", labelnames, nil)
	physical_capacity = prometheus.NewDesc(prefix_sys+"physical_capacity", "the total physical capacity of all fully allocated and thin-provisioned storage that is managed by the storage system", labelnames, nil)
	physical_free_capacity = prometheus.NewDesc(prefix_sys+"physical_free_capacity", "The total free physical capacity of all fully allocated and thin-provisioned storage that is managed by the storage system", labelnames, nil)
	used_capacity_before_reduction = prometheus.NewDesc(prefix_sys+"used_capacity_before_reduction", "The total amount of data that is written to thin-provisioned and compressed volume copies that are in data reduction storage pools - before data reduction occurs", labelnames, nil)
	used_capacity_after_reduction = prometheus.NewDesc(prefix_sys+"used_capacity_after_reduction", "The total amount of capacity that is used for thin-provisioned and compressed volume copies in the storage pool after data reduction occurs.", labelnames, nil)
	overhead_capacity = prometheus.NewDesc(prefix_sys+"overhead_capacity", "The overhead capacity consumption in all storage pools that is not attributed to data.", labelnames, nil)
	deduplication_capacity_saving = prometheus.NewDesc(prefix_sys+"deduplication_capacity_saving", "The total amount of used capacity that is saved by data deduplication. This saving is before any compression.", labelnames, nil)

	tier_capacity = prometheus.NewDesc(prefix_sys+"tier_capacity", "The total MDisk storage in the tier.", []string{"target", "resource", "tier"}, nil)
	tier_free_capacity = prometheus.NewDesc(prefix_sys+"tier_free_capacity", "The amount of MDisk storage in the tier that is unused.", []string{"target", "resource", "tier"}, nil)

	physical_capacity_usage = prometheus.NewDesc(prefix_sys+"physical_capacity_used_percent", "The physical capacity utilization", labelnames, nil)
	volume_capacity_usage = prometheus.NewDesc(prefix_sys+"volume_capacity_used_percent", "The volume capacity utilization", labelnames, nil)
	mdiskgrp_capacity_usage = prometheus.NewDesc(prefix_sys+"mdiskgrp_capacity_used_percent", "The mdiskgrp capacity utilization", labelnames, nil)

}

// systemCollector collects system metrics
type systemCollector struct {
}

func NewSystemCollector() (Collector, error) {
	return &systemCollector{}, nil
}

//Describe describes the metrics
func (*systemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- total_mdisk_capacity
	ch <- space_in_mdisk_grps
	ch <- space_allocated_to_vdisks
	ch <- total_free_space
	ch <- total_vdiskcopy_capacity
	ch <- total_used_capacity
	ch <- total_overallocation
	ch <- total_vdisk_capacity
	ch <- total_allocated_extent_capacity
	ch <- compression_virtual_capacity
	ch <- compression_compressed_capacity
	ch <- compression_uncompressed_capacity
	ch <- total_drive_raw_capacity
	ch <- tier0_flash_compressed_data_used
	ch <- tier1_flash_compressed_data_used
	ch <- tier_enterprise_compressed_data_used
	ch <- tier_nearline_compressed_data_used
	ch <- total_reclaimable_capacity
	ch <- physical_capacity
	ch <- physical_free_capacity
	ch <- used_capacity_before_reduction
	ch <- used_capacity_after_reduction
	ch <- overhead_capacity
	ch <- deduplication_capacity_saving

	ch <- tier_capacity
	ch <- tier_free_capacity

	ch <- physical_capacity_usage
	ch <- volume_capacity_usage
	ch <- mdiskgrp_capacity_usage
}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *systemCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {
	log.Debugln("Entering System collector ...")
	// labelvalues := []string{sClient.IpAddress}
	labelvalues := []string{sClient.IpAddress, sClient.Hostname}
	reqSystemURL := "https://" + sClient.IpAddress + ":7443/rest/lssystem"
	systemMetrics, err := sClient.CallSpectrumAPI(reqSystemURL)
	// This is a sample output of lssystem
	// {
	// 	"id": "0000020420400752",
	// 	"name": "SARA",
	// 	"location": "local",
	// 	"partnership": "",
	// 	"total_mdisk_capacity": "99.0TB",
	// 	"space_in_mdisk_grps": "99.0TB",
	// 	"space_allocated_to_vdisks": "558.02GB",
	// 	"total_free_space": "98.5TB",
	// 	"total_vdiskcopy_capacity": "656.00GB",
	// 	"total_used_capacity": "556.00GB",
	// 	"total_overallocation": "0",
	// 	"total_vdisk_capacity": "656.00GB",
	// 	"total_allocated_extent_capacity": "559.00GB",
	// 	"statistics_status": "on",
	// 	"statistics_frequency": "15",
	// 	"cluster_locale": "en_US",
	// 	"time_zone": "410 GMT",
	// 	"code_level": "8.2.0.2 (build 145.23.1811141325000)",
	// 	"console_IP": "172.16.192.20:443",
	// 	"id_alias": "0000020420400752",
	// 	"gm_link_tolerance": "300",
	// 	"gm_inter_cluster_delay_simulation": "0",
	// 	"gm_intra_cluster_delay_simulation": "0",
	// 	"gm_max_host_delay": "5",
	// 	"cluster_ntp_IP_address": "172.16.192.15",
	// 	"cluster_isns_IP_address": "",
	// 	"iscsi_auth_method": "none",
	// 	"iscsi_chap_secret": "",
	// 	"relationship_bandwidth_limit": "25",
	// 	"tiers": [
	// 		{
	// 			"tier": "tier0_flash",
	// 			"tier_capacity": "99.01TB",
	// 			"tier_free_capacity": "98.46TB"
	// 		},
	// 		{
	// 			"tier": "tier1_flash",
	// 			"tier_capacity": "0.00MB",
	// 			"tier_free_capacity": "0.00MB"
	// 		},
	// 		{
	// 			"tier": "tier_enterprise",
	// 			"tier_capacity": "0.00MB",
	// 			"tier_free_capacity": "0.00MB"
	// 		},
	// 		{
	// 			"tier": "tier_nearline",
	// 			"tier_capacity": "0.00MB",
	// 			"tier_free_capacity": "0.00MB"
	// 		}
	// 	],
	// 	"easy_tier_acceleration": "off",
	// 	"has_nas_key": "no",
	// 	"layer": "storage",
	// 	"rc_buffer_size": "48",
	// 	"compression_active": "no",
	// 	"compression_virtual_capacity": "0.00MB",
	// 	"compression_compressed_capacity": "0.00MB",
	// 	"compression_uncompressed_capacity": "0.00MB",
	// 	"cache_prefetch": "on",
	// 	"email_organization": "IBM Blockchain",
	// 	"email_machine_address": "44060 Digital Loundoun Plaza",
	// 	"email_machine_city": "Ashburn",
	// 	"email_machine_state": "VA",
	// 	"email_machine_zip": "20147",
	// 	"email_machine_country": "US",
	// 	"total_drive_raw_capacity": "0",
	// 	"compression_destage_mode": "off",
	// 	"rc_auth_method": "none",
	// 	"vdisk_protection_time": "15",
	// 	"vdisk_protection_enabled": "no",
	// 	"product_name": "IBM FlashSystem 9100",
	// 	"max_replication_delay": "0",
	// 	"partnership_exclusion_threshold": "315",
	// 	"tier0_flash_compressed_data_used": "0.00MB",
	// 	"tier1_flash_compressed_data_used": "0.00MB",
	// 	"tier_enterprise_compressed_data_used": "0.00MB",
	// 	"tier_nearline_compressed_data_used": "0.00MB",
	// 	"total_reclaimable_capacity": "0.00MB",
	// 	"physical_capacity": "42.90TB",
	// 	"physical_free_capacity": "42.90TB",
	// 	"used_capacity_before_reduction": "0.00MB",
	// 	"used_capacity_after_reduction": "0.00MB",
	// 	"overhead_capacity": "0.00MB",
	// 	"deduplication_capacity_saving": "0.00MB",
	// }

	if err != nil {
		log.Errorf("Executing lssystem cmd failed: %s", err)
	}
	log.Debugln("Response of lssystem: ", systemMetrics)
	total_mdisk_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "total_mdisk_capacity").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(total_mdisk_capacity, prometheus.GaugeValue, float64(total_mdisk_capacity_bytes), labelvalues...)

	space_in_mdisk_grps_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "space_in_mdisk_grps").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(space_in_mdisk_grps, prometheus.GaugeValue, float64(space_in_mdisk_grps_bytes), labelvalues...)

	space_allocated_to_vdisks_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "space_allocated_to_vdisks").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(space_allocated_to_vdisks, prometheus.GaugeValue, float64(space_allocated_to_vdisks_bytes), labelvalues...)

	total_free_space_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "total_free_space").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(total_free_space, prometheus.GaugeValue, float64(total_free_space_bytes), labelvalues...)

	total_vdiskcopy_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "total_vdiskcopy_capacity").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(total_vdiskcopy_capacity, prometheus.GaugeValue, float64(total_vdiskcopy_capacity_bytes), labelvalues...)

	total_used_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "total_used_capacity").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(total_used_capacity, prometheus.GaugeValue, float64(total_used_capacity_bytes), labelvalues...)

	total_overallocation_pc, err := strconv.ParseFloat(gjson.Get(systemMetrics, "total_overallocation").String(), 64)
	if err != nil {
		log.Errorf("Parsing string as float failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(total_overallocation, prometheus.GaugeValue, total_overallocation_pc, labelvalues...)

	total_vdisk_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "total_vdisk_capacity").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(total_vdisk_capacity, prometheus.GaugeValue, float64(total_vdisk_capacity_bytes), labelvalues...)

	total_allocated_extent_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "total_allocated_extent_capacity").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(total_allocated_extent_capacity, prometheus.GaugeValue, float64(total_allocated_extent_capacity_bytes), labelvalues...)

	compression_virtual_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "compression_virtual_capacity").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(compression_virtual_capacity, prometheus.GaugeValue, float64(compression_virtual_capacity_bytes), labelvalues...)

	compression_compressed_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "compression_compressed_capacity").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(compression_compressed_capacity, prometheus.GaugeValue, float64(compression_compressed_capacity_bytes), labelvalues...)

	compression_uncompressed_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "compression_uncompressed_capacity").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(compression_uncompressed_capacity, prometheus.GaugeValue, float64(compression_uncompressed_capacity_bytes), labelvalues...)

	total_drive_raw_capacity_bytes, err := strconv.ParseFloat(gjson.Get(systemMetrics, "total_drive_raw_capacity").String(), 64)
	if err != nil {
		log.Errorf("Parsing string as float failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(total_drive_raw_capacity, prometheus.GaugeValue, float64(total_drive_raw_capacity_bytes), labelvalues...)

	tier0_flash_compressed_data_used_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "tier0_flash_compressed_data_used").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(tier0_flash_compressed_data_used, prometheus.GaugeValue, float64(tier0_flash_compressed_data_used_bytes), labelvalues...)

	tier1_flash_compressed_data_used_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "tier1_flash_compressed_data_used").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(tier1_flash_compressed_data_used, prometheus.GaugeValue, float64(tier1_flash_compressed_data_used_bytes), labelvalues...)

	tier_enterprise_compressed_data_used_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "tier_enterprise_compressed_data_used").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(tier_enterprise_compressed_data_used, prometheus.GaugeValue, float64(tier_enterprise_compressed_data_used_bytes), labelvalues...)

	tier_nearline_compressed_data_used_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "tier_nearline_compressed_data_used").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(tier_nearline_compressed_data_used, prometheus.GaugeValue, float64(tier_nearline_compressed_data_used_bytes), labelvalues...)

	total_reclaimable_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "total_reclaimable_capacity").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(total_reclaimable_capacity, prometheus.GaugeValue, float64(total_reclaimable_capacity_bytes), labelvalues...)

	physical_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "physical_capacity").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(physical_capacity, prometheus.GaugeValue, float64(physical_capacity_bytes), labelvalues...)

	physical_free_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "physical_free_capacity").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(physical_free_capacity, prometheus.GaugeValue, float64(physical_free_capacity_bytes), labelvalues...)

	used_capacity_before_reduction_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "used_capacity_before_reduction").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(used_capacity_before_reduction, prometheus.GaugeValue, float64(used_capacity_before_reduction_bytes), labelvalues...)

	used_capacity_after_reduction_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "used_capacity_after_reduction").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(used_capacity_after_reduction, prometheus.GaugeValue, float64(used_capacity_after_reduction_bytes), labelvalues...)

	overhead_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "overhead_capacity").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(overhead_capacity, prometheus.GaugeValue, float64(overhead_capacity_bytes), labelvalues...)

	deduplication_capacity_saving_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "deduplication_capacity_saving").String())
	if err != nil {
		log.Errorf("Converting capacity unit failed: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(deduplication_capacity_saving, prometheus.GaugeValue, float64(deduplication_capacity_saving_bytes), labelvalues...)

	tierArray := gjson.Get(systemMetrics, "tiers").Array()
	for _, tier := range tierArray {
		tier_capacity_bytes, err := utils.ToBytes(tier.Get("tier_capacity").String())
		if err != nil {
			log.Errorf("Converting capacity unit failed: %s", err)
		}
		ch <- prometheus.MustNewConstMetric(tier_capacity, prometheus.GaugeValue, float64(tier_capacity_bytes), sClient.IpAddress, sClient.Hostname, tier.Get("tier").String())

		tier_free_capacity_bytes, err := utils.ToBytes(tier.Get("tier_free_capacity").String())
		if err != nil {
			log.Errorf("Converting capacity unit failed: %s", err)
		}
		ch <- prometheus.MustNewConstMetric(tier_free_capacity, prometheus.GaugeValue, float64(tier_free_capacity_bytes), sClient.IpAddress, sClient.Hostname, tier.Get("tier").String())
	}

	physical_capacity_usage_value := float64(physical_capacity_bytes-physical_free_capacity_bytes-total_reclaimable_capacity_bytes) / float64(physical_capacity_bytes) * 100
	ch <- prometheus.MustNewConstMetric(physical_capacity_usage, prometheus.GaugeValue, float64(physical_capacity_usage_value), labelvalues...)

	stored_capacity_logical := space_allocated_to_vdisks_bytes - overhead_capacity_bytes - total_reclaimable_capacity_bytes
	compression_savings := compression_uncompressed_capacity_bytes - compression_compressed_capacity_bytes + used_capacity_before_reduction_bytes - used_capacity_after_reduction_bytes + total_reclaimable_capacity_bytes
	deduplication_savings := deduplication_capacity_saving_bytes
	total_provisioned := total_vdiskcopy_capacity_bytes
	reqMdiskURL := "https://" + sClient.IpAddress + ":7443/rest/lsmdisk"
	mDiskResp, err := sClient.CallSpectrumAPI(reqMdiskURL)
	if err != nil {
		log.Errorf("Executing lsmdisk cmd failed: %s", err)
	}
	log.Debugln("Response of lsmdisk: ", mDiskResp)
	mDisks := gjson.Parse(mDiskResp).Array()
	var drive_thin_savings uint64
	for _, mdisk := range mDisks {
		mdisk_name := mdisk.Get("name").String()
		reqMdiskDeatilURL := "https://" + sClient.IpAddress + ":7443/rest/lsmdisk/" + mdisk_name
		mDiskDetailResp, err := sClient.CallSpectrumAPI(reqMdiskDeatilURL)
		if err != nil {
			log.Errorf("Executing lsmdisk/%s cmd failed: %s", mdisk_name, err)
		}
		log.Debugln("Response of lsmdisk/%s: ", mdisk_name, mDiskResp)

		// {
		// 	"id": "0",
		// 	"name": "mdisk0",
		// 	"status": "online",
		// 	"mode": "array",
		// 	"mdisk_grp_id": "0",
		// 	"mdisk_grp_name": "Pool0",
		// 	"capacity": "99.1TB",
		// 	"redundancy": "2",
		// 	"distributed": "yes",
		// 	"drive_class_id": "0",
		// 	"drive_count": "8",
		// 	"dedupe": "no",
		// 	"over_provisioned": "yes",
		// 	"provisioning_group_id": "0",
		// 	"physical_capacity": "42.90TB",
		// 	"physical_free_capacity": "42.73TB",
		// 	"write_protected": "no",
		// 	"allocated_capacity": "7.13TB",
		// 	"effective_used_capacity": "181.33GB"
		// }

		allocated_capapcity_bytes, err := utils.ToBytes(gjson.Get(mDiskDetailResp, "allocated_capacity").String())
		effective_used_capacity_bytes, err := utils.ToBytes(gjson.Get(mDiskDetailResp, "effective_used_capacity").String())
		thin_saving := allocated_capapcity_bytes - effective_used_capacity_bytes
		drive_thin_savings += thin_saving
	}

	written_capacity_with_FCMs := stored_capacity_logical + compression_savings + deduplication_savings - drive_thin_savings
	log.Debugln("written_capacity_with_FCMs", written_capacity_with_FCMs)
	volume_capacity_usage_value := float64(written_capacity_with_FCMs) / float64(total_provisioned) * 100
	log.Debugln("volume_capacity_usage_value", volume_capacity_usage_value)
	ch <- prometheus.MustNewConstMetric(volume_capacity_usage, prometheus.GaugeValue, volume_capacity_usage_value, labelvalues...)

	mdiskgrp_capacity_usage_value := float64(total_mdisk_capacity_bytes-total_free_space_bytes-total_reclaimable_capacity_bytes) / float64(total_mdisk_capacity_bytes) * 100
	ch <- prometheus.MustNewConstMetric(mdiskgrp_capacity_usage, prometheus.GaugeValue, mdiskgrp_capacity_usage_value, labelvalues...)

	log.Debugln("Leaving System collector.")
	return err
}
