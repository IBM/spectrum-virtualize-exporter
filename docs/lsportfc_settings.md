# Portfc Metrics

## Metrics Definition

```txt
# HELP spectrum_portfc_status Indicates whether the port is configured to a device of Fibre Channel (FC) port. 0-active; 1-inactive_configured; 2-inactive_unconfigured.
# TYPE spectrum_portfc_status gauge
```

## Metrics Value

### spectrum_portfc_status

- 0: active
- 1: inactive_configured
- 2: inactive_unconfigured

## Sample Metrics

```txt
spectrum_portfc_status{attachment="none",node_name="node1",port_id="3",resource="SARA-wdc04-03",target="172.16.64.20"} 2
spectrum_portfc_status{attachment="none",node_name="node1",port_id="4",resource="SARA-wdc04-03",target="172.16.64.20"} 2
spectrum_portfc_status{attachment="none",node_name="node1",port_id="7",resource="SARA-wdc04-03",target="172.16.64.20"} 2
spectrum_portfc_status{attachment="none",node_name="node1",port_id="8",resource="SARA-wdc04-03",target="172.16.64.20"} 2
spectrum_portfc_status{attachment="none",node_name="node2",port_id="3",resource="SARA-wdc04-03",target="172.16.64.20"} 2
spectrum_portfc_status{attachment="none",node_name="node2",port_id="4",resource="SARA-wdc04-03",target="172.16.64.20"} 2
spectrum_portfc_status{attachment="none",node_name="node2",port_id="7",resource="SARA-wdc04-03",target="172.16.64.20"} 2
spectrum_portfc_status{attachment="none",node_name="node2",port_id="8",resource="SARA-wdc04-03",target="172.16.64.20"} 2
spectrum_portfc_status{attachment="switch",node_name="node1",port_id="1",resource="SARA-wdc04-03",target="172.16.64.20"} 0
spectrum_portfc_status{attachment="switch",node_name="node1",port_id="2",resource="SARA-wdc04-03",target="172.16.64.20"} 0
spectrum_portfc_status{attachment="switch",node_name="node1",port_id="5",resource="SARA-wdc04-03",target="172.16.64.20"} 0
spectrum_portfc_status{attachment="switch",node_name="node1",port_id="6",resource="SARA-wdc04-03",target="172.16.64.20"} 0
spectrum_portfc_status{attachment="switch",node_name="node2",port_id="1",resource="SARA-wdc04-03",target="172.16.64.20"} 0
spectrum_portfc_status{attachment="switch",node_name="node2",port_id="2",resource="SARA-wdc04-03",target="172.16.64.20"} 0
spectrum_portfc_status{attachment="switch",node_name="node2",port_id="5",resource="SARA-wdc04-03",target="172.16.64.20"} 0
spectrum_portfc_status{attachment="switch",node_name="node2",port_id="6",resource="SARA-wdc04-03",target="172.16.64.20"} 0
```
