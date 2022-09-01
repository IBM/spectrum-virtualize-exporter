# Enclosure Metrics

## Metrics Definition

```txt
# HELP spectrum_enclosure_status Indicates whether an enclosure is visible to the SAS network. 0-online; 1-offline; 2-degraded.
# TYPE spectrum_enclosure_status gauge

# HELP spectrum_enclosure_canister_offline Indicates the number of canisters that are contained in this enclosure that are offline.
# TYPE spectrum_enclosure_canister_offline gauge

# HELP spectrum_enclosure_psu_offline Indicates the number of power-supply units (PSUs) contained in this enclosure that are offline.
# TYPE spectrum_enclosure_psu_offline gauge
```

## Metrics Value

### spectrum_enclosure_status

- 0: online
- 1: offline
- 2: degraded

## Sample Metrics

```txt
spectrum_enclosure_status{enclosure_id="1",resource="SARA-wdc04-03",target="172.16.64.20"} 0

spectrum_enclosure_canister_offline{enclosure_id="1",resource="SARA-wdc04-03",target="172.16.64.20",total_canisters="2"} 0

spectrum_enclosure_psu_offline{enclosure_id="1",resource="SARA-wdc04-03",target="172.16.64.20",total_PSUs="2"} 0
```
