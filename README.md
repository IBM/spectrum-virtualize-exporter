# spectrum-virtualize-exporter
A prometheus.io exporter for IBM Spectrum Virtualize

This exporter collects performance and metrics stats from Spectrum Virtualize and makes it available for prometheus to scrape.

## Usage

|Flag	|Description	|Default Value|	
| :---: | :---: | :---: |
| config.file | Path to configuration file | spectrumVirtualize.yml |
| web.telemetry-path | Path under which to expose metrics | /metrics |
| web.listen-address | Address on which to expose metrics and web interface | :9119 |
| web.disable-exporter-metrics | Exclude metrics about the exporter itself (promhttp_*, process_*, go_*) | false

## Building and running
Prerequisites:
* Go compiler

Building:
* binary
    
    ```go build ```
* docker image

    ```docker build -t spectrum-virtualize-exporter .```

Running
* binary

    ```./spectrum-virtualize-exporter --config.file=/etc/spectrumVirtualize/spectrumVirtualize.yml```

* docker image
    ```
    docker run -it -d -p 9119:9119 -v /etc/spectrumVirtualize/spectrumVirtualize.yml:/etc/spectrumVirtualize/spectrumVirtualize.yml --name spectrum-virtualize-exporter spectrum-virtualize-exporter --config.file=/etc/spectrumVirtualize/spectrumVirtualize.yml --log.level debug --restart always
    ```

Visit http://localhost:9119/metrics

## Configuration

The spectrum-virtualize-exporter reads from spectrumVirtualize.yml config file by default. Edit your config YAML file, Enter the IP address of the storage device, your username, and your password there. 

## Exported Metrics


### Exporter itself metrics

```
# HELP spectrum_collector_duration_seconds Duration of a collector scrape for one target
# TYPE spectrum_collector_duration_seconds gauge

# HELP spectrum_collector_success Scrape of target was sucessful
# TYPE spectrum_collector_success gauge

# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary

# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge

# HELP go_info Information about the Go environment.
# TYPE go_info gauge

# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge

# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter

# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge

# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter

# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge

# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge

# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge

# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge

# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge

# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge

# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge

# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge

# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge

# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter

# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter

# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge

# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge

# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge

# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge

# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge

# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge

# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge

# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge

# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge

# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge

# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge

# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
```

### Spectrum System Metrics
```
# HELP spectrum_system_auth_service_cert_set 
# TYPE spectrum_system_auth_service_cert_set gauge

# HELP spectrum_system_auth_service_configured 
# TYPE spectrum_system_auth_service_configured gauge

# HELP spectrum_system_auth_service_enabled 
# TYPE spectrum_system_auth_service_enabled gauge

# HELP spectrum_system_auth_service_pwd_set 
# TYPE spectrum_system_auth_service_pwd_set gauge

# HELP spectrum_system_cache_prefetch 
# TYPE spectrum_system_cache_prefetch gauge

# HELP spectrum_system_censor_callhome 
# TYPE spectrum_system_censor_callhome gauge

# HELP spectrum_system_compression_active 
# TYPE spectrum_system_compression_active gauge

# HELP spectrum_system_compression_compressed_capacity The total used capacity for all compressed volume copies in non-data reduction pools.
# TYPE spectrum_system_compression_compressed_capacity gauge

# HELP spectrum_system_compression_destage_mode 
# TYPE spectrum_system_compression_destage_mode gauge

# HELP spectrum_system_compression_uncompressed_capacity The total uncompressed used capacity for all compressed volume copies in non-data reduction pools
# TYPE spectrum_system_compression_uncompressed_capacity gauge

# HELP spectrum_system_compression_virtual_capacity The total virtual capacity for all compressed volume copies in non-data reduction pools. Compressed volumes that are in data reduction pools do not count towards this value. This value is in unsigned decimal format.
# TYPE spectrum_system_compression_virtual_capacity gauge

# HELP spectrum_system_deduplication_capacity_saving The total amount of used capacity that is saved by data deduplication. This saving is before any compression.
# TYPE spectrum_system_deduplication_capacity_saving gauge

# HELP spectrum_system_easy_tier_acceleration 
# TYPE spectrum_system_easy_tier_acceleration gauge

# HELP spectrum_system_enhanced_callhome 
# TYPE spectrum_system_enhanced_callhome gauge

# HELP spectrum_system_gen1_compatibility_mode_enabled 
# TYPE spectrum_system_gen1_compatibility_mode_enabled gauge

# HELP spectrum_system_gm_inter_cluster_delay_simulation 
# TYPE spectrum_system_gm_inter_cluster_delay_simulation gauge

# HELP spectrum_system_gm_intra_cluster_dalay_simulation 
# TYPE spectrum_system_gm_intra_cluster_dalay_simulation gauge

# HELP spectrum_system_gm_link_tolerance 
# TYPE spectrum_system_gm_link_tolerance gauge

# HELP spectrum_system_gm_max_host_delay 
# TYPE spectrum_system_gm_max_host_delay gauge

# HELP spectrum_system_has_nas_key 
# TYPE spectrum_system_has_nas_key gauge

# HELP spectrum_system_high_temp_mode 
# TYPE spectrum_system_high_temp_mode gauge

# HELP spectrum_system_inventory_mail_interval 
# TYPE spectrum_system_inventory_mail_interval gauge

# HELP spectrum_system_max_replication_delay 
# TYPE spectrum_system_max_replication_delay gauge

# HELP spectrum_system_odx 
# TYPE spectrum_system_odx gauge

# HELP spectrum_system_overhead_capacity The overhead capacity consumption in all storage pools that is not attributed to data.
# TYPE spectrum_system_overhead_capacity gauge

# HELP spectrum_system_partnership_exclusion_threshold 
# TYPE spectrum_system_partnership_exclusion_threshold gauge

# HELP spectrum_system_physical_capacity the total physical capacity of all fully allocated and thin-provisioned storage that is managed by the storage system
# TYPE spectrum_system_physical_capacity gauge

# HELP spectrum_system_physical_free_capacity The total free physical capacity of all fully allocated and thin-provisioned storage that is managed by the storage system
# TYPE spectrum_system_physical_free_capacity gauge

# HELP spectrum_system_rc_buffer_size 
# TYPE spectrum_system_rc_buffer_size gauge

# HELP spectrum_system_relationship_bandwidth_limit 
# TYPE spectrum_system_relationship_bandwidth_limit gauge

# HELP spectrum_system_space_allocated_to_vdisks The sum of mdiskgrp real_capacity
# TYPE spectrum_system_space_allocated_to_vdisks gauge

# HELP spectrum_system_space_in_mdisk_grps The sum of mdiskgrp capacity
# TYPE spectrum_system_space_in_mdisk_grps gauge

# HELP spectrum_system_statistics_frequency 
# TYPE spectrum_system_statistics_frequency gauge

# HELP spectrum_system_statistics_status 
# TYPE spectrum_system_statistics_status gauge

# HELP spectrum_system_tier0_flash_compressed_data_used The capacity of compressed data used on the flash tier 0 storage tier
# TYPE spectrum_system_tier0_flash_compressed_data_used gauge

# HELP spectrum_system_tier1_flash_compressed_data_used The capacity of compressed data used on the flash tier 1 storage tier.
# TYPE spectrum_system_tier1_flash_compressed_data_used gauge

# HELP spectrum_system_tier_capacity The total MDisk storage in the tier.
# TYPE spectrum_system_tier_capacity gauge

# HELP spectrum_system_tier_enterprise_compressed_data_used The capacity of compressed data that is used on the tier 2 enterprise storage tier.
# TYPE spectrum_system_tier_enterprise_compressed_data_used gauge

# HELP spectrum_system_tier_free_capacity The amount of MDisk storage in the tier that is unused.
# TYPE spectrum_system_tier_free_capacity gauge

# HELP spectrum_system_tier_nearline_compressed_data_used The capacity of compressed data that is used on the tier 3 nearline storage tier.
# TYPE spectrum_system_tier_nearline_compressed_data_used gauge

# HELP spectrum_system_total_allocated_extent_capacity The total size of all extents that are allocated to VDisks or otherwise in use by the system.
# TYPE spectrum_system_total_allocated_extent_capacity gauge

# HELP spectrum_system_total_drive_raw_capacity The total known capacity of all discovered drives (regardless of drive use)
# TYPE spectrum_system_total_drive_raw_capacity gauge

# HELP spectrum_system_total_free_space The sum of mdiskgrp free_capacity
# TYPE spectrum_system_total_free_space gauge

# HELP spectrum_system_total_mdisk_capacity The sum of mdiskgrp capacity plus the capacity of all unmanaged MDisks
# TYPE spectrum_system_total_mdisk_capacity gauge

# HELP spectrum_system_total_overallocation_pc The total_vdiskcopy_capacity as a percentage of total_mdisk_capacity. If total_mdisk_capacity is zero, then total_overallocation should display 100
# TYPE spectrum_system_total_overallocation_pc gauge

# HELP spectrum_system_total_reclaimable_capacity The unused (free) capacity that will be available after data is reduced
# TYPE spectrum_system_total_reclaimable_capacity gauge

# HELP spectrum_system_total_used_capacity The sum of mdiskgrp used_capacity
# TYPE spectrum_system_total_used_capacity gauge

# HELP spectrum_system_total_vdisk_capacity The total virtual capacity of volumes in the cluster
# TYPE spectrum_system_total_vdisk_capacity gauge

# HELP spectrum_system_total_vdiskcopy_capacity The total virtual capacity of all volume copies in the cluster
# TYPE spectrum_system_total_vdiskcopy_capacity gauge

# HELP spectrum_system_unmap 
# TYPE spectrum_system_unmap gauge

# HELP spectrum_system_used_capacity_after_reduction The total amount of capacity that is used for thin-provisioned and compressed volume copies in the storage pool after data reduction occurs.
# TYPE spectrum_system_used_capacity_after_reduction gauge

# HELP spectrum_system_used_capacity_before_reduction The total amount of data that is written to thin-provisioned and compressed volume copies that are in data reduction storage pools - before data reduction occurs
# TYPE spectrum_system_used_capacity_before_reduction gauge

# HELP spectrum_system_vdisk_protection_enabled 
# TYPE spectrum_system_vdisk_protection_enabled gauge

# HELP spectrum_system_vdisk_protection_time 
# TYPE spectrum_system_vdisk_protection_time gauge
```
### Spectrum System Performance Stats
```
# HELP spectrum_systemStats_cloud_down_mb The average number of Mbps for download operations to a cloud account during the sample period.
# TYPE spectrum_systemStats_cloud_down_mb gauge

# HELP spectrum_systemStats_cloud_down_ms The average amount of time (in milliseconds) it takes for the system to respond to download requests to a cloud account during the sample period.
# TYPE spectrum_systemStats_cloud_down_ms gauge

# HELP spectrum_systemStats_cloud_up_mb The average number of megabytes transferred per second (Mbps) for upload operations to a cloud account during the sample period.
# TYPE spectrum_systemStats_cloud_up_mb gauge

# HELP spectrum_systemStats_cloud_up_ms The average amount of time (in milliseconds) it takes for the system to respond to upload requests to a cloud account during the sample period.
# TYPE spectrum_systemStats_cloud_up_ms gauge

# HELP spectrum_systemStats_compression_cpu_pc The percentage of allocated CPU capacity that is used for compression.
# TYPE spectrum_systemStats_compression_cpu_pc gauge

# HELP spectrum_systemStats_cpu_pc The percentage of allocated CPU capacity that is used for the system.
# TYPE spectrum_systemStats_cpu_pc gauge

# HELP spectrum_systemStats_drive_io The average number of I/O operations that are transferred per second for read and write operations to drives during the sample period.
# TYPE spectrum_systemStats_drive_io gauge

# HELP spectrum_systemStats_drive_mb The average number of megabytes transferred per second (MBps) for read and write operations to drives during the sample period.
# TYPE spectrum_systemStats_drive_mb gauge

# HELP spectrum_systemStats_drive_ms The average amount of time in milliseconds that the system takes to respond to read and write requests to drives over the sample period.
# TYPE spectrum_systemStats_drive_ms gauge

# HELP spectrum_systemStats_drive_r_io The average number of I/O operations that are transferred per second for read operations to drives during the sample period.
# TYPE spectrum_systemStats_drive_r_io gauge

# HELP spectrum_systemStats_drive_r_mb The average number of megabytes transferred per second (MBps) for read operations to drives during the sample period.
# TYPE spectrum_systemStats_drive_r_mb gauge

# HELP spectrum_systemStats_drive_r_ms The average amount of time in milliseconds that the system takes to respond to read requests to drives over the sample period.
# TYPE spectrum_systemStats_drive_r_ms gauge

# HELP spectrum_systemStats_drive_w_io The average number of I/O operations that are transferred per second for write operations to drives during the sample period.
# TYPE spectrum_systemStats_drive_w_io gauge

# HELP spectrum_systemStats_drive_w_mb The average number of megabytes transferred per second (MBps) for write operations to drives during the sample period.
# TYPE spectrum_systemStats_drive_w_mb gauge

# HELP spectrum_systemStats_drive_w_ms The average amount of time in milliseconds that the system takes to respond write requests to drives over the sample period.
# TYPE spectrum_systemStats_drive_w_ms gauge

# HELP spectrum_systemStats_fc_io The total input/output (I/O) operations that are transferred per seconds for Fibre Channel traffic on the system. This value includes host I/O and any bandwidth that is used for communication within the system.
# TYPE spectrum_systemStats_fc_io gauge

# HELP spectrum_systemStats_fc_mb The total number of megabytes transferred per second (MBps) for Fibre Channel traffic on the system. This value includes host I/O and any bandwidth that is used for communication within the system.
# TYPE spectrum_systemStats_fc_mb gauge

# HELP spectrum_systemStats_iplink_comp_mb The average number of compressed megabytes transferred per second (MBps) over the IP Replication link during the sample period. This value is calculated after any compression of |the data takes place. This value does not include iSCSI host I/O operations.
# TYPE spectrum_systemStats_iplink_comp_mb gauge

# HELP spectrum_systemStats_iplink_io TThe total input/output (I/O) operations that are transferred per second for IP partnership traffic on the system.
# TYPE spectrum_systemStats_iplink_io gauge

# HELP spectrum_systemStats_iplink_mb The average number of megabytes requested to be transferred per second (MBps) over the IP partnership link during the sample period. This value is calculated before any compression of the data takes place. This value does not include iSCSI host input/output (I/O) operations.
# TYPE spectrum_systemStats_iplink_mb gauge

# HELP spectrum_systemStats_iscsi_io The total I/O operations that are transferred per second for iSCSI traffic on the system.
# TYPE spectrum_systemStats_iscsi_io gauge

# HELP spectrum_systemStats_iscsi_mb The total number of megabytes transferred per second (MBps) for iSCSI traffic on the system.
# TYPE spectrum_systemStats_iscsi_mb gauge

# HELP spectrum_systemStats_iser_io The total I/O operations that are transferred per second for iSER traffic on the system.
# TYPE spectrum_systemStats_iser_io gauge

# HELP spectrum_systemStats_iser_mb The total number of megabytes transferred per second (MBps) for iSER traffic on the system.
# TYPE spectrum_systemStats_iser_mb gauge

# HELP spectrum_systemStats_mdisk_io The average number of I/O operations that are transferred per second for read and write operations to MDisks during the sample period.
# TYPE spectrum_systemStats_mdisk_io gauge

# HELP spectrum_systemStats_mdisk_mb The average number of megabytes transferred per second (MBps) for read and write operations to MDisks during the sample period.
# TYPE spectrum_systemStats_mdisk_mb gauge

# HELP spectrum_systemStats_mdisk_ms The average amount of time in milliseconds that the system takes to respond to read and write requests to MDisks over the sample period.
# TYPE spectrum_systemStats_mdisk_ms gauge

# HELP spectrum_systemStats_mdisk_r_io The average number of I/O operations that are transferred per second for read operations to MDisks during the sample period.
# TYPE spectrum_systemStats_mdisk_r_io gauge

# HELP spectrum_systemStats_mdisk_r_mb The average number of megabytes transferred per second (MBps) for read operations to MDisks during the sample period.
# TYPE spectrum_systemStats_mdisk_r_mb gauge

# HELP spectrum_systemStats_mdisk_r_ms The average amount of time in milliseconds that the system takes to respond to read requests to MDisks over the sample period.
# TYPE spectrum_systemStats_mdisk_r_ms gauge

# HELP spectrum_systemStats_mdisk_w_io TThe average number of I/O operations that are transferred per second for write operations to MDisks during the sample period.
# TYPE spectrum_systemStats_mdisk_w_io gauge

# HELP spectrum_systemStats_mdisk_w_mb The average number of megabytes transferred per second (MBps) for write operations to MDisks during the sample period.
# TYPE spectrum_systemStats_mdisk_w_mb gauge

# HELP spectrum_systemStats_mdisk_w_ms the average amount of time in milliseconds that the system takes to respond to write requests to MDisks over the sample period.
# TYPE spectrum_systemStats_mdisk_w_ms gauge

# HELP spectrum_systemStats_power_w the power that is consumed in watts.
# TYPE spectrum_systemStats_power_w gauge

# HELP spectrum_systemStats_sas_io The total I/O operations that are transferred per second for SAS traffic on the system. This value includes host I/O and bandwidth that is used for background RAID activity.
# TYPE spectrum_systemStats_sas_io gauge

# HELP spectrum_systemStats_sas_mb The total number of megabytes transferred per second (MBps) for serial-attached SCSI (SAS) traffic on the system. This value includes host I/O and bandwidth that is used for background RAID activity.
# TYPE spectrum_systemStats_sas_mb gauge

# HELP spectrum_systemStats_temp_c  the ambient temperature in Celsius.
# TYPE spectrum_systemStats_temp_c gauge

# HELP spectrum_systemStats_temp_f the ambient temperature in Fahrenheit.
# TYPE spectrum_systemStats_temp_f gauge

# HELP spectrum_systemStats_total_cache_pc The total percentage for both the write and read cache usage for the node.
# TYPE spectrum_systemStats_total_cache_pc gauge

# HELP spectrum_systemStats_vdisk_io The average number of I/O operations that are transferred per second for read and write operations to volumes during the sample period.
# TYPE spectrum_systemStats_vdisk_io gauge

# HELP spectrum_systemStats_vdisk_mb The average number of megabytes transferred per second (MBps) for read and write operations to volumes during the sample period.
# TYPE spectrum_systemStats_vdisk_mb gauge

# HELP spectrum_systemStats_vdisk_ms The average amount of time in milliseconds that the system takes to respond to read and write requests to volumes over the sample period.
# TYPE spectrum_systemStats_vdisk_ms gauge

# HELP spectrum_systemStats_vdisk_r_io The average number of I/O operations that are transferred per second for read operations to volumes during the sample period.
# TYPE spectrum_systemStats_vdisk_r_io gauge

# HELP spectrum_systemStats_vdisk_r_mb The average number of megabytes transferred per second (MBps) for read operations to volumes during the sample period.
# TYPE spectrum_systemStats_vdisk_r_mb gauge

# HELP spectrum_systemStats_vdisk_r_ms The average amount of time in milliseconds that the system takes to respond to read requests to volumes over the sample period.
# TYPE spectrum_systemStats_vdisk_r_ms gauge

# HELP spectrum_systemStats_vdisk_w_io The average number of I/O operations that are transferred per second for write operations to volumes during the sample period.
# TYPE spectrum_systemStats_vdisk_w_io gauge

# HELP spectrum_systemStats_vdisk_w_mb The average number of megabytes transferred per second (MBps) for read and write operations to volumes during the sample period.
# TYPE spectrum_systemStats_vdisk_w_mb gauge

# HELP spectrum_systemStats_vdisk_w_ms The average amount of time in milliseconds that the system takes to respond to write requests to volumes over the sample period.
# TYPE spectrum_systemStats_vdisk_w_ms gauge

# HELP spectrum_systemStats_write_cache_pc The percentage of the write cache usage for the node.
# TYPE spectrum_systemStats_write_cache_pc gauge

```
### Spectrum Node Performance Stats
```
# HELP spectrum_nodeStats_cloud_down_mb The average number of Mbps for download operations to a cloud account during the sample period.
# TYPE spectrum_nodeStats_cloud_down_mb gauge

# HELP spectrum_nodeStats_cloud_down_ms The average amount of time (in milliseconds) it takes for the system to respond to download requests to a cloud account during the sample period.
# TYPE spectrum_nodeStats_cloud_down_ms gauge

# HELP spectrum_nodeStats_cloud_up_mb The average number of megabytes transferred per second (Mbps) for upload operations to a cloud account during the sample period.
# TYPE spectrum_nodeStats_cloud_up_mb gauge

# HELP spectrum_nodeStats_cloud_up_ms The average amount of time (in milliseconds) it takes for the system to respond to upload requests to a cloud account during the sample period.
# TYPE spectrum_nodeStats_cloud_up_ms gauge

# HELP spectrum_nodeStats_compression_cpu_pc The percentage of allocated CPU capacity that is used for compression.
# TYPE spectrum_nodeStats_compression_cpu_pc gauge

# HELP spectrum_nodeStats_cpu_pc The percentage of allocated CPU capacity that is used for the system.
# TYPE spectrum_nodeStats_cpu_pc gauge

# HELP spectrum_nodeStats_drive_io The average number of I/O operations that are transferred per second for read and write operations to drives during the sample period.
# TYPE spectrum_nodeStats_drive_io gauge

# HELP spectrum_nodeStats_drive_mb The average number of megabytes transferred per second (MBps) for read and write operations to drives during the sample period.
# TYPE spectrum_nodeStats_drive_mb gauge

# HELP spectrum_nodeStats_drive_ms The average amount of time in milliseconds that the system takes to respond to read and write requests to drives over the sample period.
# TYPE spectrum_nodeStats_drive_ms gauge

# HELP spectrum_nodeStats_drive_r_io The average number of I/O operations that are transferred per second for read operations to drives during the sample period.
# TYPE spectrum_nodeStats_drive_r_io gauge

# HELP spectrum_nodeStats_drive_r_mb The average number of megabytes transferred per second (MBps) for read operations to drives during the sample period.
# TYPE spectrum_nodeStats_drive_r_mb gauge

# HELP spectrum_nodeStats_drive_r_ms The average amount of time in milliseconds that the system takes to respond to read requests to drives over the sample period.
# TYPE spectrum_nodeStats_drive_r_ms gauge

# HELP spectrum_nodeStats_drive_w_io The average number of I/O operations that are transferred per second for write operations to drives during the sample period.
# TYPE spectrum_nodeStats_drive_w_io gauge

# HELP spectrum_nodeStats_drive_w_mb The average number of megabytes transferred per second (MBps) for write operations to drives during the sample period.
# TYPE spectrum_nodeStats_drive_w_mb gauge

# HELP spectrum_nodeStats_drive_w_ms The average amount of time in milliseconds that the system takes to respond write requests to drives over the sample period.
# TYPE spectrum_nodeStats_drive_w_ms gauge

# HELP spectrum_nodeStats_fc_io The total input/output (I/O) operations that are transferred per seconds for Fibre Channel traffic on the system. This value includes host I/O and any bandwidth that is used for communication within the system.
# TYPE spectrum_nodeStats_fc_io gauge

# HELP spectrum_nodeStats_fc_mb The total number of megabytes transferred per second (MBps) for Fibre Channel traffic on the system. This value includes host I/O and any bandwidth that is used for communication within the system.
# TYPE spectrum_nodeStats_fc_mb gauge

# HELP spectrum_nodeStats_iplink_comp_mb The average number of compressed megabytes transferred per second (MBps) over the IP Replication link during the sample period. This value is calculated after any compression of |the data takes place. This value does not include iSCSI host I/O operations.
# TYPE spectrum_nodeStats_iplink_comp_mb gauge

# HELP spectrum_nodeStats_iplink_io TThe total input/output (I/O) operations that are transferred per second for IP partnership traffic on the system.
# TYPE spectrum_nodeStats_iplink_io gauge

# HELP spectrum_nodeStats_iplink_mb The average number of megabytes requested to be transferred per second (MBps) over the IP partnership link during the sample period. This value is calculated before any compression of the data takes place. This value does not include iSCSI host input/output (I/O) operations.
# TYPE spectrum_nodeStats_iplink_mb gauge

# HELP spectrum_nodeStats_iscsi_io The total I/O operations that are transferred per second for iSCSI traffic on the system.
# TYPE spectrum_nodeStats_iscsi_io gauge

# HELP spectrum_nodeStats_iscsi_mb The total number of megabytes transferred per second (MBps) for iSCSI traffic on the system.
# TYPE spectrum_nodeStats_iscsi_mb gauge

# HELP spectrum_nodeStats_iser_io The total I/O operations that are transferred per second for iSER traffic on the system.
# TYPE spectrum_nodeStats_iser_io gauge

# HELP spectrum_nodeStats_iser_mb The total number of megabytes transferred per second (MBps) for iSER traffic on the system.
# TYPE spectrum_nodeStats_iser_mb gauge

# HELP spectrum_nodeStats_mdisk_io The average number of I/O operations that are transferred per second for read and write operations to MDisks during the sample period.
# TYPE spectrum_nodeStats_mdisk_io gauge

# HELP spectrum_nodeStats_mdisk_mb The average number of megabytes transferred per second (MBps) for read and write operations to MDisks during the sample period.
# TYPE spectrum_nodeStats_mdisk_mb gauge

# HELP spectrum_nodeStats_mdisk_ms The average amount of time in milliseconds that the system takes to respond to read and write requests to MDisks over the sample period.
# TYPE spectrum_nodeStats_mdisk_ms gauge

# HELP spectrum_nodeStats_mdisk_r_io The average number of I/O operations that are transferred per second for read operations to MDisks during the sample period.
# TYPE spectrum_nodeStats_mdisk_r_io gauge

# HELP spectrum_nodeStats_mdisk_r_mb The average number of megabytes transferred per second (MBps) for read operations to MDisks during the sample period.
# TYPE spectrum_nodeStats_mdisk_r_mb gauge

# HELP spectrum_nodeStats_mdisk_r_ms The average amount of time in milliseconds that the system takes to respond to read requests to MDisks over the sample period.
# TYPE spectrum_nodeStats_mdisk_r_ms gauge

# HELP spectrum_nodeStats_mdisk_w_io TThe average number of I/O operations that are transferred per second for write operations to MDisks during the sample period.
# TYPE spectrum_nodeStats_mdisk_w_io gauge

# HELP spectrum_nodeStats_mdisk_w_mb The average number of megabytes transferred per second (MBps) for write operations to MDisks during the sample period.
# TYPE spectrum_nodeStats_mdisk_w_mb gauge

# HELP spectrum_nodeStats_mdisk_w_ms the average amount of time in milliseconds that the system takes to respond to write requests to MDisks over the sample period.
# TYPE spectrum_nodeStats_mdisk_w_ms gauge

# HELP spectrum_nodeStats_sas_io The total I/O operations that are transferred per second for SAS traffic on the system. This value includes host I/O and bandwidth that is used for background RAID activity.
# TYPE spectrum_nodeStats_sas_io gauge

# HELP spectrum_nodeStats_sas_mb The total number of megabytes transferred per second (MBps) for serial-attached SCSI (SAS) traffic on the system. This value includes host I/O and bandwidth that is used for background RAID activity.
# TYPE spectrum_nodeStats_sas_mb gauge

# HELP spectrum_nodeStats_total_cache_pc The total percentage for both the write and read cache usage for the node.
# TYPE spectrum_nodeStats_total_cache_pc gauge

# HELP spectrum_nodeStats_vdisk_io The average number of I/O operations that are transferred per second for read and write operations to volumes during the sample period.
# TYPE spectrum_nodeStats_vdisk_io gauge

# HELP spectrum_nodeStats_vdisk_mb The average number of megabytes transferred per second (MBps) for read and write operations to volumes during the sample period.
# TYPE spectrum_nodeStats_vdisk_mb gauge

# HELP spectrum_nodeStats_vdisk_ms The average amount of time in milliseconds that the system takes to respond to read and write requests to volumes over the sample period.
# TYPE spectrum_nodeStats_vdisk_ms gauge

# HELP spectrum_nodeStats_vdisk_r_io The average number of I/O operations that are transferred per second for read operations to volumes during the sample period.
# TYPE spectrum_nodeStats_vdisk_r_io gauge

# HELP spectrum_nodeStats_vdisk_r_mb The average number of megabytes transferred per second (MBps) for read operations to volumes during the sample period.
# TYPE spectrum_nodeStats_vdisk_r_mb gauge

# HELP spectrum_nodeStats_vdisk_r_ms The average amount of time in milliseconds that the system takes to respond to read requests to volumes over the sample period.
# TYPE spectrum_nodeStats_vdisk_r_ms gauge

# HELP spectrum_nodeStats_vdisk_w_io The average number of I/O operations that are transferred per second for write operations to volumes during the sample period.
# TYPE spectrum_nodeStats_vdisk_w_io gauge

# HELP spectrum_nodeStats_vdisk_w_mb The average number of megabytes transferred per second (MBps) for read and write operations to volumes during the sample period.
# TYPE spectrum_nodeStats_vdisk_w_mb gauge

# HELP spectrum_nodeStats_vdisk_w_ms The average amount of time in milliseconds that the system takes to respond to write requests to volumes over the sample period.
# TYPE spectrum_nodeStats_vdisk_w_ms gauge

# HELP spectrum_nodeStats_write_cache_pc The percentage of the write cache usage for the node.
# TYPE spectrum_nodeStats_write_cache_pc gauge
```

## Refrences

* Spectrum Virtualize RESTful API For FS9XXX:https://www.ibm.com/support/knowledgecenter/en/STSLR9_8.2.0/com.ibm.fs9100_820.doc/rest_api_overview.html

* IBM Spectrum Virtualize Interfacing Using the RESTful API:https://www.ibm.com/support/knowledgecenter/STVLF4_8.1.3/spectrum.virtualize.813.doc/Spectrum_Virtualize_API_8.1.3.pdf