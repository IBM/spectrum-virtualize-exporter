# Drive Metrics

## Metrics Definition

```txt
# HELP spectrum_drive_status Indicates the summary status of the drive. 0-online; 1-offline; 2-degraded.
# TYPE spectrum_drive_status gauge

# HELP spectrum_drive_firmware_level Indicates the firmware level consistency of disks. 0-consistent; 1-inconsistent.
# TYPE spectrum_drive_firmware_level gauge
```

## Metrics Value

### spectrum_drive_status

- 0: online, which indicates that the drive is available through all drive ports.
- 1: offline, which indicates that the drive is unavailable.
- 2: degraded, which indicates that the drive is available but not through all drive ports.

### spectrum_drive_firmware_level

- 0: Firmware levels are consistent across all drives
- 1: Firmware level is not consistent with other drives

## Sample Metrics

```txt
spectrum_drive_status{drive_id="0",enclosure_id="1",resource="SARA-wdc04-03",slot_id="1",target="172.16.64.20"} 0
spectrum_drive_status{drive_id="1",enclosure_id="1",resource="SARA-wdc04-03",slot_id="5",target="172.16.64.20"} 0
spectrum_drive_status{drive_id="2",enclosure_id="1",resource="SARA-wdc04-03",slot_id="7",target="172.16.64.20"} 0
spectrum_drive_status{drive_id="3",enclosure_id="1",resource="SARA-wdc04-03",slot_id="6",target="172.16.64.20"} 0
spectrum_drive_status{drive_id="4",enclosure_id="1",resource="SARA-wdc04-03",slot_id="4",target="172.16.64.20"} 0
spectrum_drive_status{drive_id="5",enclosure_id="1",resource="SARA-wdc04-03",slot_id="2",target="172.16.64.20"} 0
spectrum_drive_status{drive_id="6",enclosure_id="1",resource="SARA-wdc04-03",slot_id="3",target="172.16.64.20"} 0
spectrum_drive_status{drive_id="7",enclosure_id="1",resource="SARA-wdc04-03",slot_id="8",target="172.16.64.20"} 0

spectrum_drive_firmware_level{drive_id="0",firmware_level="1_2_11",resource="SARA-wdc04-03",target="172.16.64.20"} 0
spectrum_drive_firmware_level{drive_id="1",firmware_level="1_2_11",resource="SARA-wdc04-03",target="172.16.64.20"} 0
spectrum_drive_firmware_level{drive_id="2",firmware_level="1_2_11",resource="SARA-wdc04-03",target="172.16.64.20"} 0
spectrum_drive_firmware_level{drive_id="3",firmware_level="1_2_11",resource="SARA-wdc04-03",target="172.16.64.20"} 0
spectrum_drive_firmware_level{drive_id="4",firmware_level="1_2_11",resource="SARA-wdc04-03",target="172.16.64.20"} 0
spectrum_drive_firmware_level{drive_id="5",firmware_level="1_2_11",resource="SARA-wdc04-03",target="172.16.64.20"} 0
spectrum_drive_firmware_level{drive_id="6",firmware_level="1_2_11",resource="SARA-wdc04-03",target="172.16.64.20"} 0
spectrum_drive_firmware_level{drive_id="7",firmware_level="1_2_11",resource="SARA-wdc04-03",target="172.16.64.20"} 0
```
