# Enclosure Battery Metrics

## Metrics Definition

```txt
# HELP enclosurebattery_status Identifies status of each battery in enclosures. 0-online; 1-offline; 2-degraded.
# TYPE enclosurebattery_status gauge

# HELP enclosurebattery_end_of_life_warning Identifies the battery's end of life. Replace the battery if yes. 0-no; 1-yes.
# TYPE enclosurebattery_end_of_life_warning gauge
```

## Metrics Value

### enclosurebattery_status

- 0: online;
- 1: offline
- 2: degraded

### enclosurebattery_end_of_life_warning

- 0: no
- 1: yes

## Sample Metrics

```txt
enclosurebattery_status{battery_id="1",enclosure_id="1",resource="SARA-wdc04-03",target="172.16.64.20"} 0
enclosurebattery_status{battery_id="2",enclosure_id="1",resource="SARA-wdc04-03",target="172.16.64.20"} 0

enclosurebattery_end_of_life_warning{battery_id="1",enclosure_id="1",resource="SARA-wdc04-03",target="172.16.64.20"} 0
enclosurebattery_end_of_life_warning{battery_id="2",enclosure_id="1",resource="SARA-wdc04-03",target="172.16.64.20"} 1
```
