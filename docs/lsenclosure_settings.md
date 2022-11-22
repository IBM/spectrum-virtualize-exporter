# Enclosure Metrics

## Metrics Definition

```txt
# HELP spectrum_enclosure_status Indicates whether an enclosure is visible to the SAS network. 0-online; 1-offline; 2-degraded.
# TYPE spectrum_enclosure_status gauge
```

## Metrics Value

### spectrum_enclosure_status

- 0: online
- 1: offline
- 2: degraded

## Sample Metrics

```txt
spectrum_enclosure_status{enclosure_id="1",resource="SARA-wdc04-03",target="172.16.64.20"} 0
```
