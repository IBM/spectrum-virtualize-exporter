# spectrum-virtualize-exporter
A prometheus.io exporter for IBM Spectrum Virtualize

This exporter collects performance and metrics stats from Spectrum Virtualize and makes it available for prometheus to scrape.

## Usage

|Flag	|Description	|Default Value|	
| :---: | :---: | :---: |
| config.file | Path to configuration file | spectrumVirtualize.yml |
| web.telemetry-path | Path under which to expose metrics | /metrics |
| web.listen-address | Address on which to expose metrics and web interface | :9119 |
| web.disable-exporter-metrics | Exclude metrics about the exporter itself (promhttp_*, process_*, go_*) | false

## Building and running
Prerequisites:
* Go compiler

Building:
* binary
    
    ```go build ```
* docker image

    ```docker build -t spectrum-virtualize-exporter .```

Running
* binary

    ```./spectrum-virtualize-exporter --config.file=/etc/spectrumVirtualize/spectrumVirtualize.yml```

* docker image
    ```
    docker run -it -d -p 9119:9119 -v /etc/spectrumVirtualize/spectrumVirtualize.yml:/etc/spectrumVirtualize/spectrumVirtualize.yml --name spectrum-virtualize-exporter spectrum-virtualize-exporter --config.file=/etc/spectrumVirtualize/spectrumVirtualize.yml --log.level debug --restart always
    ```

Visit http://localhost:9119/metrics

## Configuration

The spectrum-virtualize-exporter reads from spectrumVirtualize.yml config file by default. Edit your config YAML file, Enter the IP address of the storage device, your username, and your password there. 

## Exported Metrics

| CLI Command | Description | Metrics |
| --- | --- | --- |
| - | Metrics from the exporter itself. | [Metrics List](docs/exporter_metrics.md) |
| lssystem | Get a detailed view of a clustered system (system). | [Metrics List](docs/lssystem_metrics.md) |
| lssystemstats | Get the most recent values of all node statistics in a system | [Metrics List](docs/lssystemstats_metrics.md) |
| lsnodestats | Ge the most recent values of statistics for all nodes. | [Metrics List](docs/lsnodestats_metrics.md)|


## References

* Spectrum Virtualize RESTful API For FS9XXX:https://www.ibm.com/support/knowledgecenter/en/STSLR9_8.2.0/com.ibm.fs9100_820.doc/rest_api_overview.html

* IBM Spectrum Virtualize Interfacing Using the RESTful API:https://www.ibm.com/support/knowledgecenter/STVLF4_8.1.3/spectrum.virtualize.813.doc/Spectrum_Virtualize_API_8.1.3.pdf
