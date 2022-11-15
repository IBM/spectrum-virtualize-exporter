# Portfc Metrics

## Metrics Definition

```txt
# HELP spectrum_portfc_status Indicates whether the port is configured to a device of Fibre Channel (FC) port. 0-active; 1-inactive_configured; 2-inactive_unconfigured.
# TYPE spectrum_portfc_status gauge
# HELP spectrum_portfc_attachment Indicates if the port is attached to a FC switch. 0-yes; 1-no.
# TYPE spectrum_portfc_attachment gauge
```

## Metrics Value

### spectrum_portfc_status

- 0: active
- 1: inactive_configured
- 2: inactive_unconfigured

### spectrum_portfc_attachment

- 0: yes
- 1: no

## Sample Metrics

```txt
spectrum_portfc_status{node_name="node1",port_id="1",resource="SARA-wdc04-03",target="172.16.64.20",wwpn="500507681011038D"} 0
spectrum_portfc_status{node_name="node1",port_id="2",resource="SARA-wdc04-03",target="172.16.64.20",wwpn="500507681012038D"} 0
spectrum_portfc_status{node_name="node1",port_id="5",resource="SARA-wdc04-03",target="172.16.64.20",wwpn="500507681022038D"} 0
spectrum_portfc_status{node_name="node1",port_id="6",resource="SARA-wdc04-03",target="172.16.64.20",wwpn="500507681023038D"} 0
spectrum_portfc_status{node_name="node2",port_id="1",resource="SARA-wdc04-03",target="172.16.64.20",wwpn="500507681011039F"} 0
spectrum_portfc_status{node_name="node2",port_id="2",resource="SARA-wdc04-03",target="172.16.64.20",wwpn="500507681012039F"} 0
spectrum_portfc_status{node_name="node2",port_id="5",resource="SARA-wdc04-03",target="172.16.64.20",wwpn="500507681021039F"} 0
spectrum_portfc_status{node_name="node2",port_id="6",resource="SARA-wdc04-03",target="172.16.64.20",wwpn="500507681022039F"} 0

spectrum_portfc_attachment{node_name="node1",port_id="1",resource="SARA-wdc04-03",target="172.16.64.20",wwpn="500507681011038D"} 0
spectrum_portfc_attachment{node_name="node1",port_id="2",resource="SARA-wdc04-03",target="172.16.64.20",wwpn="500507681012038D"} 0
spectrum_portfc_attachment{node_name="node1",port_id="5",resource="SARA-wdc04-03",target="172.16.64.20",wwpn="500507681022038D"} 0
spectrum_portfc_attachment{node_name="node1",port_id="6",resource="SARA-wdc04-03",target="172.16.64.20",wwpn="500507681023038D"} 0
spectrum_portfc_attachment{node_name="node2",port_id="1",resource="SARA-wdc04-03",target="172.16.64.20",wwpn="500507681011039F"} 0
spectrum_portfc_attachment{node_name="node2",port_id="2",resource="SARA-wdc04-03",target="172.16.64.20",wwpn="500507681012039F"} 0
spectrum_portfc_attachment{node_name="node2",port_id="5",resource="SARA-wdc04-03",target="172.16.64.20",wwpn="500507681021039F"} 0
spectrum_portfc_attachment{node_name="node2",port_id="6",resource="SARA-wdc04-03",target="172.16.64.20",wwpn="500507681022039F"} 0
```
