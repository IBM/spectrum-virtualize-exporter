# Node Canister Metrics

## Metrics Definition

```txt
# HELP nodecanister_status Status of nodes that are part of the system. 0-online; 1-offline; 2-service; 3-flushing; 4-pending; 5-adding; 6-deleting.
# TYPE nodecanister_status gauge
```

## Metrics Value

### nodecanister_status

- 0: online
- 1: offline
- 2: service
- 3: flushing
- 4: pending
- 5: adding
- 6: deleting

## Sample Metrics

```txt
nodecanister_status{node_name="node1",resource="SARA-wdc04-03",target="172.16.64.20"} 0
nodecanister_status{node_name="node2",resource="SARA-wdc04-03",target="172.16.64.20"} 1
```
