### MdiskGrp(Pool) Metrics 
```
# HELP mdiskgrp_capacity The total amount of MDisk storage that is assigned to the storage pool..
# TYPE mdiskgrp_capacity gauge

# HELP mdiskgrp_compression_active Indicates whether any compressed volume copies are in the storage pool.
# TYPE mdiskgrp_compression_active gauge

# HELP mdiskgrp_compression_compressed_capacity The total used capacity for all compressed volume copies in regular storage pools.
# TYPE mdiskgrp_compression_compressed_capacity gauge

# HELP mdiskgrp_compression_uncompressed_capacity the total uncompressed used capacity for all compressed volume copies in regular storage pools
# TYPE mdiskgrp_compression_uncompressed_capacity gauge

# HELP mdiskgrp_compression_virtual_capacity The total virtual capacity for all compressed volume copies in regular storage pools. 
# TYPE mdiskgrp_compression_virtual_capacity gauge

# HELP mdiskgrp_deduplication_capcacity_saving The capacity that is saved by deduplication before compression in a data reduction pool.
# TYPE mdiskgrp_deduplication_capcacity_saving gauge

# HELP mdiskgrp_extent_size The sizes of the extents for this group
# TYPE mdiskgrp_extent_size gauge

# HELP mdiskgrp_free_capacity The amount of MDisk storage that is immediately available. Additionally, reclaimable_capacity can eventually become available
# TYPE mdiskgrp_free_capacity gauge

# HELP mdiskgrp_overallocation The ratio of the virtual_capacity value to the capacity
# TYPE mdiskgrp_overallocation gauge

# HELP mdiskgrp_overhead_capacity The MDisk capacity that is reserved for internal usage.
# TYPE mdiskgrp_overhead_capacity gauge

# HELP mdiskgrp_real_capacity The total MDisk storage capacity assigned to volume copies.
# TYPE mdiskgrp_real_capacity gauge

# HELP mdiskgrp_reclaimable_capacity The MDisk capacity that is reserved for internal usage.
# TYPE mdiskgrp_reclaimable_capacity gauge

# HELP mdiskgrp_used_capacity The amount of data that is stored on MDisks.
# TYPE mdiskgrp_used_capacity gauge

# HELP mdiskgrp_used_capacity_after_reduction The data that is stored on MDisks for non-fully-allocated volume copies in a data reduction pool.
# TYPE mdiskgrp_used_capacity_after_reduction gauge

# HELP mdiskgrp_used_capacity_before_reduction The data that is stored on non-fully-allocated volume copies in a data reduction pool.
# TYPE mdiskgrp_used_capacity_before_reduction gauge

# HELP mdiskgrp_virtual_capacity The total host mappable capacity of all volume copies in the storage pool.
# TYPE mdiskgrp_virtual_capacity gauge

```