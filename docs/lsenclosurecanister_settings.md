# Enclosure Metrics

## Metrics Definition

```txt
# HELP spectrum_enclosurecanister_status Identifies status of each canister in enclosures. 0-online; 1-offline; 2-degraded.
# TYPE spectrum_enclosurecanister_status gauge
```

## Metrics Value

### spectrum_enclosurecanister_status

- 0: online
- 1: offline
- 2: degraded

## Sample Metrics

```txt
spectrum_enclosurecanister_status{canister_id="1",enclosure_id="1",node_name="node1",resource="SARA-wdc04-03",target="172.16.64.20"} 0
spectrum_enclosurecanister_status{canister_id="2",enclosure_id="1",node_name="node2",resource="SARA-wdc04-03",target="172.16.64.20"} 0
```
