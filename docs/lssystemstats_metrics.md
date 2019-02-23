### Spectrum System Performance Stats Metrics

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
