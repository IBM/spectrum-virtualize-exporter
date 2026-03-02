# Enclosure Metrics

## Metrics Definition

```txt
# HELP enclosurepsu_status Indicates status of each power-supply unit (PSU) in enclosures.
# TYPE enclosurepsu_status gauge
```

## Metrics Value

### enclosurepsu_status

- 0: online
- 1: offline
- 2: degraded

## Sample Metrics

```txt
enclosurepsu_status{enclosure_id="1",psu_id="1",resource="SARA-wdc04-03",target="172.16.64.20"} 0
enclosurepsu_status{enclosure_id="1",psu_id="2",resource="SARA-wdc04-03",target="172.16.64.20"} 0
```
