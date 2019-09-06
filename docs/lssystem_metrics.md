### Spectrum System Metrics

```
# HELP spectrum_system_compression_compressed_capacity The total used capacity for all compressed volume copies in non-data reduction pools.
# TYPE spectrum_system_compression_compressed_capacity gauge

# HELP spectrum_system_compression_uncompressed_capacity The total uncompressed used capacity for all compressed volume copies in non-data reduction pools
# TYPE spectrum_system_compression_uncompressed_capacity gauge

# HELP spectrum_system_compression_virtual_capacity The total virtual capacity for all compressed volume copies in non-data reduction pools. Compressed volumes that are in data reduction pools do not count towards this value. This value is in unsigned decimal format.
# TYPE spectrum_system_compression_virtual_capacity gauge

# HELP spectrum_system_deduplication_capacity_saving The total amount of used capacity that is saved by data deduplication. This saving is before any compression.
# TYPE spectrum_system_deduplication_capacity_saving gauge

# HELP spectrum_system_overhead_capacity The overhead capacity consumption in all storage pools that is not attributed to data.
# TYPE spectrum_system_overhead_capacity gauge

# HELP spectrum_system_physical_capacity the total physical capacity of all fully allocated and thin-provisioned storage that is managed by the storage system
# TYPE spectrum_system_physical_capacity gauge

# HELP spectrum_system_physical_free_capacity The total free physical capacity of all fully allocated and thin-provisioned storage that is managed by the storage system
# TYPE spectrum_system_physical_free_capacity gauge

# HELP spectrum_system_space_allocated_to_vdisks The sum of mdiskgrp real_capacity
# TYPE spectrum_system_space_allocated_to_vdisks gauge

# HELP spectrum_system_space_in_mdisk_grps The sum of mdiskgrp capacity
# TYPE spectrum_system_space_in_mdisk_grps gauge

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

# HELP spectrum_system_used_capacity_after_reduction The total amount of capacity that is used for thin-provisioned and compressed volume copies in the storage pool after data reduction occurs.
# TYPE spectrum_system_used_capacity_after_reduction gauge

# HELP spectrum_system_used_capacity_before_reduction The total amount of data that is written to thin-provisioned and compressed volume copies that are in data reduction storage pools - before data reduction occurs
# TYPE spectrum_system_used_capacity_before_reduction gauge

# HELP spectrum_system_physical_capacity_used_percent The physical capacity utilization.
# TYPE spectrum_system_physical_capacity_used_percent

# HELP spectrum_system_mdiskgrp_capacity_used_percent The mdiskgrp capacity utilization
# TYPE spectrum_system_mdiskgrp_capacity_used_percent

# HELP spectrum_system_volume_capacity_used_percent The volume capacity utilization.
# TYPE spectrum_system_volume_capacity_used_percent gauge
```
