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

	statistics_status                 *prometheus.Desc
	statistics_frequency              *prometheus.Desc
	gm_link_tolerance                 *prometheus.Desc
	gm_inter_cluster_delay_simulation *prometheus.Desc
	gm_intra_cluster_dalay_simulation *prometheus.Desc
	gm_max_host_delay                 *prometheus.Desc
	inventory_mail_interval           *prometheus.Desc
	auth_service_configured           *prometheus.Desc
	auth_service_enabled              *prometheus.Desc
	auth_service_pwd_set              *prometheus.Desc
	auth_service_cert_set             *prometheus.Desc
	relationship_bandwidth_limit      *prometheus.Desc
	easy_tier_acceleration            *prometheus.Desc
	has_nas_key                       *prometheus.Desc
	rc_buffer_size                    *prometheus.Desc
	compression_active                *prometheus.Desc
	cache_prefetch                    *prometheus.Desc
	compression_destage_mode          *prometheus.Desc
	high_temp_mode                    *prometheus.Desc
	vdisk_protection_time             *prometheus.Desc
	vdisk_protection_enabled          *prometheus.Desc
	odx                               *prometheus.Desc
	max_replication_delay             *prometheus.Desc
	partnership_exclusion_threshold   *prometheus.Desc
	gen1_compatibility_mode_enabled   *prometheus.Desc
	unmap                             *prometheus.Desc
	enhanced_callhome                 *prometheus.Desc
	censor_callhome                   *prometheus.Desc
	physical_capacity_usage           *prometheus.Desc
	volume_capacity_usage             *prometheus.Desc
	mdiskgrp_capacity_usage           *prometheus.Desc
)

func init() {
	registerCollector("lssystem", defaultEnabled, NewSystemCollector)
	labelnames := []string{"target"}
	total_mdisk_capacity = prometheus.NewDesc(prefix_sys+"total_mdisk_capacity", "The sum of mdiskgrp capacity plus the capacity of all unmanaged MDisks", labelnames, nil)
	space_in_mdisk_grps = prometheus.NewDesc(prefix_sys+"space_in_mdisk_grps", "The sum of mdiskgrp capacity", labelnames, nil)
	space_allocated_to_vdisks = prometheus.NewDesc(prefix_sys+"space_allocated_to_vdisks", "The sum of mdiskgrp real_capacity", labelnames, nil)
	total_free_space = prometheus.NewDesc(prefix_sys+"total_free_space", "The sum of mdiskgrp free_capacity", labelnames, nil)
	total_vdiskcopy_capacity = prometheus.NewDesc(prefix_sys+"total_vdiskcopy_capacity", "The total virtual capacity of all volume copies in the cluster", labelnames, nil)
	total_used_capacity = prometheus.NewDesc(prefix_sys+"total_used_capacity", "The sum of mdiskgrp used_capacity", labelnames, nil)
	total_overallocation = prometheus.NewDesc(prefix_sys+"total_overallocation_pc", "The total_vdiskcopy_capacity as a percentage of total_mdisk_capacity. If total_mdisk_capacity is zero, then total_overallocation should display 100", labelnames, nil)
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

	tier_capacity = prometheus.NewDesc(prefix_sys+"tier_capacity", "The total MDisk storage in the tier.", []string{"target", "tier"}, nil)
	tier_free_capacity = prometheus.NewDesc(prefix_sys+"tier_free_capacity", "The amount of MDisk storage in the tier that is unused.", []string{"target", "tier"}, nil)

	statistics_status = prometheus.NewDesc(prefix_sys+"statistics_status", "", []string{"target"}, nil)
	statistics_frequency = prometheus.NewDesc(prefix_sys+"statistics_frequency", "", []string{"target"}, nil)
	gm_link_tolerance = prometheus.NewDesc(prefix_sys+"gm_link_tolerance", "", []string{"target"}, nil)
	gm_inter_cluster_delay_simulation = prometheus.NewDesc(prefix_sys+"gm_inter_cluster_delay_simulation", "", []string{"target"}, nil)
	gm_intra_cluster_dalay_simulation = prometheus.NewDesc(prefix_sys+"gm_intra_cluster_dalay_simulation", "", []string{"target"}, nil)
	gm_max_host_delay = prometheus.NewDesc(prefix_sys+"gm_max_host_delay", "", []string{"target"}, nil)
	inventory_mail_interval = prometheus.NewDesc(prefix_sys+"inventory_mail_interval", "", []string{"target"}, nil)
	auth_service_configured = prometheus.NewDesc(prefix_sys+"auth_service_configured", "", []string{"target"}, nil)
	auth_service_enabled = prometheus.NewDesc(prefix_sys+"auth_service_enabled", "", []string{"target"}, nil)
	auth_service_pwd_set = prometheus.NewDesc(prefix_sys+"auth_service_pwd_set", "", []string{"target"}, nil)
	auth_service_cert_set = prometheus.NewDesc(prefix_sys+"auth_service_cert_set", "", []string{"target"}, nil)
	relationship_bandwidth_limit = prometheus.NewDesc(prefix_sys+"relationship_bandwidth_limit", "", []string{"target"}, nil)
	easy_tier_acceleration = prometheus.NewDesc(prefix_sys+"easy_tier_acceleration", "", []string{"target"}, nil)
	has_nas_key = prometheus.NewDesc(prefix_sys+"has_nas_key", "", []string{"target"}, nil)
	rc_buffer_size = prometheus.NewDesc(prefix_sys+"rc_buffer_size", "", []string{"target"}, nil)
	compression_active = prometheus.NewDesc(prefix_sys+"compression_active", "", []string{"target"}, nil)
	cache_prefetch = prometheus.NewDesc(prefix_sys+"cache_prefetch", "", []string{"target"}, nil)
	compression_destage_mode = prometheus.NewDesc(prefix_sys+"compression_destage_mode", "", []string{"target"}, nil)
	high_temp_mode = prometheus.NewDesc(prefix_sys+"high_temp_mode", "", []string{"target"}, nil)
	vdisk_protection_time = prometheus.NewDesc(prefix_sys+"vdisk_protection_time", "", []string{"target"}, nil)
	vdisk_protection_enabled = prometheus.NewDesc(prefix_sys+"vdisk_protection_enabled", "", []string{"target"}, nil)
	odx = prometheus.NewDesc(prefix_sys+"odx", "", []string{"target"}, nil)
	max_replication_delay = prometheus.NewDesc(prefix_sys+"max_replication_delay", "", []string{"target"}, nil)
	partnership_exclusion_threshold = prometheus.NewDesc(prefix_sys+"partnership_exclusion_threshold", "", []string{"target"}, nil)
	gen1_compatibility_mode_enabled = prometheus.NewDesc(prefix_sys+"gen1_compatibility_mode_enabled", "", []string{"target"}, nil)
	unmap = prometheus.NewDesc(prefix_sys+"unmap", "", []string{"target"}, nil)
	enhanced_callhome = prometheus.NewDesc(prefix_sys+"enhanced_callhome", "", []string{"target"}, nil)
	censor_callhome = prometheus.NewDesc(prefix_sys+"censor_callhome", "", []string{"target"}, nil)

	physical_capacity_usage = prometheus.NewDesc(prefix_sys+"physical_capacity_usage", "physical capacity utilization", []string{"target"}, nil)
	volume_capacity_usage = prometheus.NewDesc(prefix_sys+"volume_capacity_usage", "volume capacity utilization", []string{"target"}, nil)
	mdiskgrp_capacity_usage = prometheus.NewDesc(prefix_sys+"mdiskgrp_capacity_usage", "mdiskgrp capacity utilization", []string{"target"}, nil)

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

	ch <- statistics_status
	ch <- statistics_frequency
	ch <- gm_link_tolerance
	ch <- gm_inter_cluster_delay_simulation
	ch <- gm_intra_cluster_dalay_simulation
	ch <- gm_max_host_delay
	ch <- inventory_mail_interval
	ch <- auth_service_configured
	ch <- auth_service_enabled
	ch <- auth_service_pwd_set
	ch <- auth_service_cert_set
	ch <- relationship_bandwidth_limit
	ch <- easy_tier_acceleration
	ch <- has_nas_key
	ch <- rc_buffer_size
	ch <- compression_active
	ch <- cache_prefetch
	ch <- compression_destage_mode
	ch <- high_temp_mode
	ch <- vdisk_protection_time
	ch <- vdisk_protection_enabled
	ch <- odx
	ch <- max_replication_delay
	ch <- partnership_exclusion_threshold
	ch <- gen1_compatibility_mode_enabled
	ch <- unmap
	ch <- enhanced_callhome
	ch <- censor_callhome

	ch <- physical_capacity_usage
	ch <- volume_capacity_usage
	ch <- mdiskgrp_capacity_usage
}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *systemCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {
	log.Debugln("System collector is starting")
	labelvalues := []string{sClient.IpAddress}

	reqSystemURL := "https://" + sClient.IpAddress + ":7443/rest/lssystem"
	systemMetrics, _ := sClient.CallSpectrumAPI(reqSystemURL)

	total_mdisk_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "total_mdisk_capacity").String())
	ch <- prometheus.MustNewConstMetric(total_mdisk_capacity, prometheus.GaugeValue, float64(total_mdisk_capacity_bytes), labelvalues...)

	space_in_mdisk_grps_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "space_in_mdisk_grps").String())
	ch <- prometheus.MustNewConstMetric(space_in_mdisk_grps, prometheus.GaugeValue, float64(space_in_mdisk_grps_bytes), labelvalues...)

	space_allocated_to_vdisks_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "space_allocated_to_vdisks").String())
	ch <- prometheus.MustNewConstMetric(space_allocated_to_vdisks, prometheus.GaugeValue, float64(space_allocated_to_vdisks_bytes), labelvalues...)

	total_free_space_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "total_free_space").String())
	ch <- prometheus.MustNewConstMetric(total_free_space, prometheus.GaugeValue, float64(total_free_space_bytes), labelvalues...)

	total_vdiskcopy_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "total_vdiskcopy_capacity").String())
	ch <- prometheus.MustNewConstMetric(total_vdiskcopy_capacity, prometheus.GaugeValue, float64(total_vdiskcopy_capacity_bytes), labelvalues...)

	total_used_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "total_used_capacity").String())
	ch <- prometheus.MustNewConstMetric(total_used_capacity, prometheus.GaugeValue, float64(total_used_capacity_bytes), labelvalues...)

	total_overallocation_pc, err := strconv.ParseFloat(gjson.Get(systemMetrics, "total_overallocation").String(), 64)
	ch <- prometheus.MustNewConstMetric(total_overallocation, prometheus.GaugeValue, total_overallocation_pc, labelvalues...)

	total_vdisk_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "total_vdisk_capacity").String())
	ch <- prometheus.MustNewConstMetric(total_vdisk_capacity, prometheus.GaugeValue, float64(total_vdisk_capacity_bytes), labelvalues...)

	total_allocated_extent_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "total_allocated_extent_capacity").String())
	ch <- prometheus.MustNewConstMetric(total_allocated_extent_capacity, prometheus.GaugeValue, float64(total_allocated_extent_capacity_bytes), labelvalues...)

	compression_virtual_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "compression_virtual_capacity").String())
	ch <- prometheus.MustNewConstMetric(compression_virtual_capacity, prometheus.GaugeValue, float64(compression_virtual_capacity_bytes), labelvalues...)

	compression_compressed_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "compression_compressed_capacity").String())
	ch <- prometheus.MustNewConstMetric(compression_compressed_capacity, prometheus.GaugeValue, float64(compression_compressed_capacity_bytes), labelvalues...)

	compression_uncompressed_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "compression_uncompressed_capacity").String())
	ch <- prometheus.MustNewConstMetric(compression_uncompressed_capacity, prometheus.GaugeValue, float64(compression_uncompressed_capacity_bytes), labelvalues...)

	total_drive_raw_capacity_bytes, err := strconv.ParseFloat(gjson.Get(systemMetrics, "total_drive_raw_capacity").String(), 64)
	ch <- prometheus.MustNewConstMetric(total_drive_raw_capacity, prometheus.GaugeValue, float64(total_drive_raw_capacity_bytes), labelvalues...)

	tier0_flash_compressed_data_used_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "tier0_flash_compressed_data_used").String())
	ch <- prometheus.MustNewConstMetric(tier0_flash_compressed_data_used, prometheus.GaugeValue, float64(tier0_flash_compressed_data_used_bytes), labelvalues...)

	tier1_flash_compressed_data_used_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "tier1_flash_compressed_data_used").String())
	ch <- prometheus.MustNewConstMetric(tier1_flash_compressed_data_used, prometheus.GaugeValue, float64(tier1_flash_compressed_data_used_bytes), labelvalues...)

	tier_enterprise_compressed_data_used_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "tier_enterprise_compressed_data_used").String())
	ch <- prometheus.MustNewConstMetric(tier_enterprise_compressed_data_used, prometheus.GaugeValue, float64(tier_enterprise_compressed_data_used_bytes), labelvalues...)

	tier_nearline_compressed_data_used_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "tier_nearline_compressed_data_used").String())
	ch <- prometheus.MustNewConstMetric(tier_nearline_compressed_data_used, prometheus.GaugeValue, float64(tier_nearline_compressed_data_used_bytes), labelvalues...)

	total_reclaimable_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "total_reclaimable_capacity").String())
	ch <- prometheus.MustNewConstMetric(total_reclaimable_capacity, prometheus.GaugeValue, float64(total_reclaimable_capacity_bytes), labelvalues...)

	physical_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "physical_capacity").String())
	ch <- prometheus.MustNewConstMetric(physical_capacity, prometheus.GaugeValue, float64(physical_capacity_bytes), labelvalues...)

	physical_free_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "physical_free_capacity").String())
	ch <- prometheus.MustNewConstMetric(physical_free_capacity, prometheus.GaugeValue, float64(physical_free_capacity_bytes), labelvalues...)

	used_capacity_before_reduction_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "used_capacity_before_reduction").String())
	ch <- prometheus.MustNewConstMetric(used_capacity_before_reduction, prometheus.GaugeValue, float64(used_capacity_before_reduction_bytes), labelvalues...)

	used_capacity_after_reduction_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "used_capacity_after_reduction").String())
	ch <- prometheus.MustNewConstMetric(used_capacity_after_reduction, prometheus.GaugeValue, float64(used_capacity_after_reduction_bytes), labelvalues...)

	overhead_capacity_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "overhead_capacity").String())
	ch <- prometheus.MustNewConstMetric(overhead_capacity, prometheus.GaugeValue, float64(overhead_capacity_bytes), labelvalues...)

	deduplication_capacity_saving_bytes, err := utils.ToBytes(gjson.Get(systemMetrics, "deduplication_capacity_saving").String())
	ch <- prometheus.MustNewConstMetric(deduplication_capacity_saving, prometheus.GaugeValue, float64(deduplication_capacity_saving_bytes), labelvalues...)

	tierArray := gjson.Get(systemMetrics, "tiers").Array()
	for _, tier := range tierArray {
		tier_capacity_bytes, err := utils.ToBytes(tier.Get("tier_capacity").String())
		ch <- prometheus.MustNewConstMetric(tier_capacity, prometheus.GaugeValue, float64(tier_capacity_bytes), sClient.IpAddress, tier.Get("tier").String())

		tier_free_capacity_bytes, err := utils.ToBytes(tier.Get("tier_free_capacity").String())
		ch <- prometheus.MustNewConstMetric(tier_free_capacity, prometheus.GaugeValue, float64(tier_free_capacity_bytes), sClient.IpAddress, tier.Get("tier").String())

		return err
	}

	statistics_status_value, err := utils.ToBool(gjson.Get(systemMetrics, "statistics_status").String())
	ch <- prometheus.MustNewConstMetric(statistics_status, prometheus.GaugeValue, statistics_status_value, labelvalues...)

	statistics_frequency_value := gjson.Get(systemMetrics, "statistics_frequency").Float()
	ch <- prometheus.MustNewConstMetric(statistics_frequency, prometheus.GaugeValue, statistics_frequency_value, labelvalues...)

	gm_link_tolerance_value := gjson.Get(systemMetrics, "gm_link_tolerance").Float()
	ch <- prometheus.MustNewConstMetric(gm_link_tolerance, prometheus.GaugeValue, gm_link_tolerance_value, labelvalues...)

	gm_inter_cluster_delay_simulation_value := gjson.Get(systemMetrics, "gm_inter_cluster_delay_simulation").Float()
	ch <- prometheus.MustNewConstMetric(gm_inter_cluster_delay_simulation, prometheus.GaugeValue, gm_inter_cluster_delay_simulation_value, labelvalues...)

	gm_intra_cluster_dalay_simulation_value := gjson.Get(systemMetrics, "gm_intra_cluster_dalay_simulation").Float()
	ch <- prometheus.MustNewConstMetric(gm_intra_cluster_dalay_simulation, prometheus.GaugeValue, gm_intra_cluster_dalay_simulation_value, labelvalues...)

	gm_max_host_delay_value := gjson.Get(systemMetrics, "gm_max_host_delay").Float()
	ch <- prometheus.MustNewConstMetric(gm_max_host_delay, prometheus.GaugeValue, gm_max_host_delay_value, labelvalues...)

	inventory_mail_interval_value := gjson.Get(systemMetrics, "inventory_mail_interval").Float()
	ch <- prometheus.MustNewConstMetric(inventory_mail_interval, prometheus.GaugeValue, inventory_mail_interval_value, labelvalues...)

	auth_service_configured_value, err := utils.ToBool(gjson.Get(systemMetrics, "auth_service_configured").String())
	ch <- prometheus.MustNewConstMetric(auth_service_configured, prometheus.GaugeValue, auth_service_configured_value, labelvalues...)

	auth_service_enabled_value, err := utils.ToBool(gjson.Get(systemMetrics, "auth_service_enabled").String())
	ch <- prometheus.MustNewConstMetric(auth_service_enabled, prometheus.GaugeValue, auth_service_enabled_value, labelvalues...)

	auth_service_pwd_set_value, err := utils.ToBool(gjson.Get(systemMetrics, "auth_service_pwd_set").String())
	ch <- prometheus.MustNewConstMetric(auth_service_pwd_set, prometheus.GaugeValue, auth_service_pwd_set_value, labelvalues...)

	auth_service_cert_set_value, err := utils.ToBool(gjson.Get(systemMetrics, "auth_service_cert_set").String())
	ch <- prometheus.MustNewConstMetric(auth_service_cert_set, prometheus.GaugeValue, auth_service_cert_set_value, labelvalues...)

	relationship_bandwidth_limit_value, err := utils.ToBool(gjson.Get(systemMetrics, "relationship_bandwidth_limit").String())
	ch <- prometheus.MustNewConstMetric(relationship_bandwidth_limit, prometheus.GaugeValue, relationship_bandwidth_limit_value, labelvalues...)

	easy_tier_acceleration_value, err := utils.ToBool(gjson.Get(systemMetrics, "easy_tier_acceleration").String())
	ch <- prometheus.MustNewConstMetric(easy_tier_acceleration, prometheus.GaugeValue, easy_tier_acceleration_value, labelvalues...)

	has_nas_key_value, err := utils.ToBool(gjson.Get(systemMetrics, "has_nas_key").String())
	ch <- prometheus.MustNewConstMetric(has_nas_key, prometheus.GaugeValue, has_nas_key_value, labelvalues...)

	rc_buffer_size_value := gjson.Get(systemMetrics, "rc_buffer_size").Float()
	ch <- prometheus.MustNewConstMetric(rc_buffer_size, prometheus.GaugeValue, rc_buffer_size_value, labelvalues...)

	compression_active_value, err := utils.ToBool(gjson.Get(systemMetrics, "compression_active").String())
	ch <- prometheus.MustNewConstMetric(compression_active, prometheus.GaugeValue, compression_active_value, labelvalues...)

	cache_prefetch_value, err := utils.ToBool(gjson.Get(systemMetrics, "cache_prefetch").String())
	ch <- prometheus.MustNewConstMetric(cache_prefetch, prometheus.GaugeValue, cache_prefetch_value, labelvalues...)

	compression_destage_mode_value, err := utils.ToBool(gjson.Get(systemMetrics, "compression_destage_mode").String())
	ch <- prometheus.MustNewConstMetric(compression_destage_mode, prometheus.GaugeValue, compression_destage_mode_value, labelvalues...)

	high_temp_mode_value, err := utils.ToBool(gjson.Get(systemMetrics, "high_temp_mode").String())
	ch <- prometheus.MustNewConstMetric(high_temp_mode, prometheus.GaugeValue, high_temp_mode_value, labelvalues...)

	vdisk_protection_time_value := gjson.Get(systemMetrics, "vdisk_protection_time").Float()
	ch <- prometheus.MustNewConstMetric(vdisk_protection_time, prometheus.GaugeValue, vdisk_protection_time_value, labelvalues...)

	vdisk_protection_enabled_value, err := utils.ToBool(gjson.Get(systemMetrics, "vdisk_protection_enabled").String())
	ch <- prometheus.MustNewConstMetric(vdisk_protection_enabled, prometheus.GaugeValue, vdisk_protection_enabled_value, labelvalues...)

	odx_value, err := utils.ToBool(gjson.Get(systemMetrics, "odx").String())
	ch <- prometheus.MustNewConstMetric(odx, prometheus.GaugeValue, odx_value, labelvalues...)

	max_replication_delay_value := gjson.Get(systemMetrics, "max_replication_delay").Float()
	ch <- prometheus.MustNewConstMetric(max_replication_delay, prometheus.GaugeValue, max_replication_delay_value, labelvalues...)

	partnership_exclusion_threshold_value := gjson.Get(systemMetrics, "partnership_exclusion_threshold").Float()
	ch <- prometheus.MustNewConstMetric(partnership_exclusion_threshold, prometheus.GaugeValue, partnership_exclusion_threshold_value, labelvalues...)

	gen1_compatibility_mode_enabled_value, err := utils.ToBool(gjson.Get(systemMetrics, "gen1_compatibility_mode_enabled").String())
	ch <- prometheus.MustNewConstMetric(gen1_compatibility_mode_enabled, prometheus.GaugeValue, gen1_compatibility_mode_enabled_value, labelvalues...)

	unmap_value, err := utils.ToBool(gjson.Get(systemMetrics, "unmap").String())
	ch <- prometheus.MustNewConstMetric(unmap, prometheus.GaugeValue, unmap_value, labelvalues...)

	enhanced_callhome_value, err := utils.ToBool(gjson.Get(systemMetrics, "enhanced_callhome").String())
	ch <- prometheus.MustNewConstMetric(enhanced_callhome, prometheus.GaugeValue, enhanced_callhome_value, labelvalues...)

	censor_callhome_value, err := utils.ToBool(gjson.Get(systemMetrics, "censor_callhome").String())
	ch <- prometheus.MustNewConstMetric(censor_callhome, prometheus.GaugeValue, censor_callhome_value, labelvalues...)

	physical_capacity_usage_value := (float64(physical_capacity_bytes) - float64(physical_free_capacity_bytes)) / float64(physical_capacity_bytes)
	ch <- prometheus.MustNewConstMetric(physical_capacity_usage, prometheus.GaugeValue, physical_capacity_usage_value, labelvalues...)

	volume_capacity_usage_value := float64(total_used_capacity_bytes) / float64(total_vdisk_capacity_bytes)
	ch <- prometheus.MustNewConstMetric(volume_capacity_usage, prometheus.GaugeValue, float64(volume_capacity_usage_value), labelvalues...)

	mdiskgrp_capacity_usage_value := (float64(total_mdisk_capacity_bytes) - float64(total_free_space_bytes) - float64(total_reclaimable_capacity_bytes)) / float64(total_mdisk_capacity_bytes)
	ch <- prometheus.MustNewConstMetric(mdiskgrp_capacity_usage, prometheus.GaugeValue, float64(mdiskgrp_capacity_usage_value), labelvalues...)

	return err
}
