# MDisk Metrics

## Metrics Definition

```txt
# HELP spectrum_mdisk_status Status of managed disks (MDisks) visible to the system. 0-online; 1-offline; 2-excluded; 3-degraded_paths; 4-degraded_ports; 5-degraded.
# TYPE spectrum_mdisk_status gauge
```

## Metrics Value

### spectrum_mdisk_status

- 0: online
- 1: offline
- 2: excluded
- 3: degraded_paths
- 4: degraded_ports
- 5: degraded

## Sample Metrics

```txt
spectrum_mdisk_status{mdisk_name="mdisk0",pool_name="Pool0",resource="SARA-wdc04-03",target="172.16.64.20"} 0
```
