
# Alerts

## Alert Filter

- agent_tag_service
- agent_tag_region
- agent_tag_zone
- agent_tag_pod
- agent_tag_cpc
- agent_tag_lpar

e.g:

```txt
agent_tag_service = paas
agent_tag_region = eu-gb
agent_tag_zone = lon06z1
agent_tag_pod = lon06z1
agent_tag_cpc = CPCA
agent_tag_lpar = lon3-qz1-sr2-rk006-m01
```

## Alert Definition

| Name | Severity | Condition | Value | Segments | Description |
| --- | --- | --- | --- | --- | --- |
| FS9K Callhome Status Alert | Low | `avg(avg(spectrum_callhome_info)) > 0.0` | `0`: status --enabled, connection --active<br>`1`: status --disabled<br>`2`: status --enabled, connection in ["error", "untried"] | resource | Alert when Callhome status is disabled/inctive. |
| FS9K Enclosure Status Alert | High | `avg(avg(spectrum_enclosure_status)) > 0.0` | `0`: online<br>`1`: offline<br>`2`: degraded | resource<br>enclosure_id | Alert when the enclosure status is offline/degraded. |
| FS9K Canister Status Alert | High | `avg(avg(spectrum_enclosurecanister_status)) > 0.0` | `0`: online<br>`1`: offline<br>`2`: degraded | resource<br>enclosure_id<br>canister_id<br>node_name | Alert when canister status is offline/degraded. |
| FS9K PSU Status Alert | High | `avg(avg(spectrum_enclosurepsu_status)) > 0.0` | `0`: online<br>`1`: offline<br>`2`: degraded | resource<br>enclosure_id<br>psu_id | Alert when PSU status is offline/degraded. |
| FS9K Battery Status Alert | High | `avg(avg(spectrum_enclosurebattery_status)) > 0.0` | `0`: online<br>`1`: offline<br>`2`: degraded | resource<br>enclosure_id<br>battery_id | Alert when the battery status is offline/degraded. |
| FS9K Battery End of Life Alert | High | `avg(avg(spectrum_enclosurebattery_end_of_life_warning)) > 0.0` | `0`: no<br>`1`: yes | resource<br>enclosure_id<br>battery_id | Alert when the battery end of life warning is on. |
| FS9K Drive Status Alert | High | `avg(avg(spectrum_drive_status)) > 0.0` | `0`: online<br>`1`: offline<br>`2`: degraded | resource<br>drive_id | Alert when drive status is offline/degraded. |
| FS9K Disk Firmware Level Inconsistent Alert | High | `avg(avg(spectrum_drive_firmware_level_consistency)) > 0.0` | `0`: consistent<br>`1`: inconsistent | resource | Alert when disk drive firmware level is inconsistent. |
| FS9K Port Status Alert | High | `avg(avg(spectrum_portfc_status)) > 0.0` | `0`: active<br>`1`: inactive_configured<br>`2`: inactive_unconfigured | resource<br>node_name<br>port_id<br>wwpn | Alert when port status is not active. |
| FS9K Port Attachment Alert | High | `avg(avg(spectrum_portfc_attachment)) > 0.0` | `0`: yes<br>`1`: no | resource<br>node_name<br>port_id<br>wwpn | Alert when port is not attached to a FC switch. |
| FS9K Host Connection Status Alert | High | `avg(avg(spectrum_host_status)) >= 1.0` | `0`: online/active<br>`1`: inactive<br>`2`: offline<br>`3`: degraded | resource<br>host_name | Alert when host connection status is inctive/offline/degraded. |
| FS9K Node Status Alert | High | `avg(avg(spectrum_nodecanister_status)) > 0.0` | `0`: online<br>`1`: offline<br>`2`: service<br>`3`: flushing<br>`4`: pending<br>`5`: adding<br>`6`: deleting | resource<br>node_name | Alert when node status is not online. |
| FS9K Managed Disks Status Alert | High | `avg(avg(spectrum_mdisk_status)) > 0.0` | `0`: online<br>`1`: offline<br>`2`: excluded<br>`3`: degraded_paths<br>`4`: degraded_ports<br>`5`: degraded | resource<br>pod_name<br>mdisk_name | Alert when managed disks status is not online. |
| FS9K Storage Pool Status Alert | High | `avg(avg(spectrum_mdiskgrp_status)) > 0.0` | `0`: online<br>`1`: offline<br>`2`: others | resource<br>pool_name | Alert when storage pool status is not online. |
| FS9K IP Status Alert | Low | `avg(avg(spectrum_ip_status)) > 0.0` | `0`: connectable<br>`1`: unreachable | resource<br>ip_name<br>ip_address | Alert when PSYS/SSYS/SVC1/SVC2 IP is unreachable. |
