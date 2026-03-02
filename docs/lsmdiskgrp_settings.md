# MDiskgrp Metrics

## Metrics Definition

```txt
# HELP mdiskgrp_status Status of storage pools that are visible to the system. 0-online; 1-offline; 2-others.
# TYPE mdiskgrp_status gauge
```

## Metrics Value

### mdiskgrp_status

- 0: online
- 1: offline
- 2: others

## Sample Metrics

```txt
mdiskgrp_status{pool_name="Pool0",resource="SARA-wdc04-03",target="172.16.64.20"} 0
```
