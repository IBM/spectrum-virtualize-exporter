# Call Home information Metrics

## Metrics Definition

```txt
# HELP spectrum_callhome_info The status of the Call Home information.
# TYPE spectrum_callhome_info gauge
```

## Metrics Value

- 0: status --enabled, connection --active;
- 1: status --disabled
- 2: status --enabled, connection in ["error", "untried"]

## Sample Metrics

```txt
spectrum_callhome_info{connection="active",resource="SARA-wdc04-03",status="enabled",target="172.16.64.20"} 0
```
