# Spectrum Exporter Metrics

## Metrics Definition

```txt
# HELP spectrum_collector_authtoken_renew_failure_total Cumulative count of failed verification of renewed auth token
# TYPE spectrum_collector_authtoken_renew_failure_total counter

# HELP spectrum_collector_authtoken_renew_interval_seconds Interval of the last renewing auth token
# TYPE spectrum_collector_authtoken_renew_interval_seconds gauge

# HELP spectrum_collector_authtoken_renew_success_total Cumulative count of successful verification of renewed auth token
# TYPE spectrum_collector_authtoken_renew_success_total counter

# HELP spectrum_collector_scrape_duration_seconds Duration of a collector scraping for one host
# TYPE spectrum_collector_scrape_duration_seconds gauge
```
