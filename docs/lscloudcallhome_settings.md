# Call Home Information Metrics

## Metrics Definition

```txt
# HELP callhome_info The status of the Call Home information.
# TYPE callhome_info gauge
```

## Metrics Value

- 0: status --enabled, connection --active;
- 1: status --disabled
- 2: status --enabled, connection in ["error", "untried"]

## Sample Metrics

```txt
callhome_info{connection="active",resource="SARA-wdc04-03",status="enabled",target="172.16.64.20"} 0
```
