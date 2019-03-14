### MdiskGrp(Pool) Metrics 
```
# HELP spectrum_mdiskgrp_capacity The total amount of MDisk storage that is assigned to the storage pool..
# TYPE spectrum_mdiskgrp_capacity gauge

# HELP spectrum_mdiskgrp_compression_active Indicates whether any compressed volume copies are in the storage pool.
# TYPE spectrum_mdiskgrp_compression_active gauge

# HELP spectrum_mdiskgrp_compression_compressed_capacity The total used capacity for all compressed volume copies in regular storage pools.
# TYPE spectrum_mdiskgrp_compression_compressed_capacity gauge

# HELP spectrum_mdiskgrp_compression_uncompressed_capacity the total uncompressed used capacity for all compressed volume copies in regular storage pools
# TYPE spectrum_mdiskgrp_compression_uncompressed_capacity gauge

# HELP spectrum_mdiskgrp_compression_virtual_capacity The total virtual capacity for all compressed volume copies in regular storage pools. 
# TYPE spectrum_mdiskgrp_compression_virtual_capacity gauge

# HELP spectrum_mdiskgrp_deduplication_capcacity_saving The capacity that is saved by deduplication before compression in a data reduction pool.
# TYPE spectrum_mdiskgrp_deduplication_capcacity_saving gauge

# HELP spectrum_mdiskgrp_extent_size The sizes of the extents for this group
# TYPE spectrum_mdiskgrp_extent_size gauge

# HELP spectrum_mdiskgrp_free_capacity The amount of MDisk storage that is immediately available. Additionally, reclaimable_capacity can eventually become available
# TYPE spectrum_mdiskgrp_free_capacity gauge

# HELP spectrum_mdiskgrp_overallocation The ratio of the virtual_capacity value to the capacity
# TYPE spectrum_mdiskgrp_overallocation gauge

# HELP spectrum_mdiskgrp_overhead_capacity The MDisk capacity that is reserved for internal usage.
# TYPE spectrum_mdiskgrp_overhead_capacity gauge

# HELP spectrum_mdiskgrp_real_capacity The total MDisk storage capacity assigned to volume copies.
# TYPE spectrum_mdiskgrp_real_capacity gauge

# HELP spectrum_mdiskgrp_reclaimable_capacity The MDisk capacity that is reserved for internal usage.
# TYPE spectrum_mdiskgrp_reclaimable_capacity gauge

# HELP spectrum_mdiskgrp_used_capacity The amount of data that is stored on MDisks.
# TYPE spectrum_mdiskgrp_used_capacity gauge

# HELP spectrum_mdiskgrp_used_capacity_after_reduction The data that is stored on MDisks for non-fully-allocated volume copies in a data reduction pool.
# TYPE spectrum_mdiskgrp_used_capacity_after_reduction gauge

# HELP spectrum_mdiskgrp_used_capacity_before_reduction The data that is stored on non-fully-allocated volume copies in a data reduction pool.
# TYPE spectrum_mdiskgrp_used_capacity_before_reduction gauge

# HELP spectrum_mdiskgrp_virtual_capacity The total host mappable capacity of all volume copies in the storage pool.
# TYPE spectrum_mdiskgrp_virtual_capacity gauge

```