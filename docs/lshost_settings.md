# Host Metrics

## Metrics Definition

```txt
# HELP spectrum_host_status LPAR connection status. 0-online/active; 1-inactive; 2-offline; 3-degraded.
# TYPE spectrum_host_status gauge
```

## Metrics Value

### spectrum_host_status

- 0: online/active
- 1: inactive
- 2: offline
- 3: degraded

## Sample Metrics

```txt
spectrum_host_status{lpar_name="dal1-qz2-sr3-rk196-m01",resource="SARA",target="192.168.196.120"} 0
spectrum_host_status{lpar_name="dal1-qz2-sr3-rk196-a01",resource="SARA",target="192.168.196.120"} 0
spectrum_host_status{lpar_name="dal1-qz2-sr3-rk196-m04",resource="SARA",target="192.168.196.120"} 3
...
```