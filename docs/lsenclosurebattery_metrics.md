# Enclosure Battery Metrics

## Metrics Definition

```txt
# HELP spectrum_enclosurebattery_status Identifies the status of the battery. 0-online; 1-offline; 2-degraded.
# TYPE spectrum_enclosurebattery_status gauge

# HELP spectrum_enclosurebattery_end_of_life_warning Identifies the battery's end of life. Replace the battery if yes. 0-no; 1-yes.
# TYPE spectrum_enclosurebattery_end_of_life_warning gauge
```

## Metrics Value

### spectrum_enclosurebattery_status

- 0: online;
- 1: offline
- 2: degraded

### spectrum_enclosurebattery_end_of_life_warning

- 0: no
- 1: yes

## Sample Metrics

```txt
spectrum_enclosurebattery_status{battery_id="1",enclosure_id="1",resource="SARA-wdc04-03",target="172.16.64.20"} 0
spectrum_enclosurebattery_status{battery_id="2",enclosure_id="1",resource="SARA-wdc04-03",target="172.16.64.20"} 0

spectrum_enclosurebattery_end_of_life_warning{battery_id="1",enclosure_id="1",resource="SARA-wdc04-03",target="172.16.64.20"} 0
spectrum_enclosurebattery_end_of_life_warning{battery_id="2",enclosure_id="1",resource="SARA-wdc04-03",target="172.16.64.20"} 1
```
