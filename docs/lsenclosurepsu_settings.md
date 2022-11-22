# Enclosure Metrics

## Metrics Definition

```txt
# HELP spectrum_enclosurepsu_status Indicates status of each power-supply unit (PSU) in enclosures.
# TYPE spectrum_enclosurepsu_status gauge
```

## Metrics Value

### spectrum_enclosurepsu_status

- 0: online
- 1: offline
- 2: degraded

## Sample Metrics

```txt
spectrum_enclosurepsu_status{enclosure_id="1",psu_id="1",resource="SARA-wdc04-03",target="172.16.64.20"} 0
spectrum_enclosurepsu_status{enclosure_id="1",psu_id="2",resource="SARA-wdc04-03",target="172.16.64.20"} 0
```
