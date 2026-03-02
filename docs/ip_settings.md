# Host Metrics

## Metrics Definition

```txt
# HELP ip_status IP connection status. 0-connectable; 1-unreachable.
# TYPE ip_status gauge
```

## Metrics Value

### ip_status

- 0: connectable
- 1: unreachable

## Sample Metrics

```txt
ip_status{ip_address="192.168.196.120",ip_name="PSYS",resource="SARA",target="192.168.196.120"} 0
ip_status{ip_address="192.168.196.121",ip_name="SSYS",resource="SARA",target="192.168.196.120"} 0
ip_status{ip_address="192.168.196.122",ip_name="SVC1",resource="SARA",target="192.168.196.120"} 0
ip_status{ip_address="192.168.196.123",ip_name="SVC2",resource="SARA",target="192.168.196.120"} 0
```
