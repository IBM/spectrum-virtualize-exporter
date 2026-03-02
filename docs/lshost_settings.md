# Host Metrics

## Metrics Definition

```txt
# HELP host_status Host connection status. 0-online; 1-offline; 2-degraded.
# TYPE host_status gauge
```

## Metrics Value

### host_status

- 0: online
- 1: offline
- 2: degraded

## Sample Metrics

```txt
host_status{host_name="dal1-qz2-sr3-rk196-m01",resource="SARA",target="192.168.196.120"} 0
host_status{host_name="dal1-qz2-sr3-rk196-a01",resource="SARA",target="192.168.196.120"} 0
host_status{host_name="dal1-qz2-sr3-rk196-m04",resource="SARA",target="192.168.196.120"} 2
...
```
