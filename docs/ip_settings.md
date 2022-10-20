# Host Metrics

## Metrics Definition

```txt
# HELP spectrum_ip_status IP connection status. 0-connectable; 1-unreachable.
# TYPE spectrum_ip_status gauge
```

## Metrics Value

### spectrum_ip_status

- 0: connectable
- 1: unreachable

## Sample Metrics

```txt
spectrum_ip_status{ip_address="192.168.196.120",ip_name="PSYS",resource="SARA",target="192.168.196.120"} 0
spectrum_ip_status{ip_address="192.168.196.121",ip_name="SSYS",resource="SARA",target="192.168.196.120"} 0
spectrum_ip_status{ip_address="192.168.196.122",ip_name="SVC1",resource="SARA",target="192.168.196.120"} 0
spectrum_ip_status{ip_address="192.168.196.123",ip_name="SVC2",resource="SARA",target="192.168.196.120"} 0
```
