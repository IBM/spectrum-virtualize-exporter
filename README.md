# spectrum-virtualize-exporter

This [Prometheus](https://prometheus.io) [Exporter](https://prometheus.io/docs/instrumenting/exporters)
collects metrics from storage solutions which are built with the
 [IBM Spectrum Virtualize software](https://www.ibm.com/support/home/product/10000647/IBM_Spectrum_Virtualize_software).
Storage solutions built with the IBM Spectrum Virtualize software are the
[IBM FlashSystem V9000 system](https://www.ibm.com/support/knowledgecenter/STKMQV_8.2.1/com.ibm.storage.vflashsystem9000.8.2.1.doc/svc_svcovr_1bcfiq.html),
the [IBM SAN Volume Controller](https://www.ibm.com/us-en/marketplace/san-volume-controller) and
 the [IBM Storwize Family](https://www.ibm.com/it-infrastructure/storage/storwize<Paste>)

## Usage

| Flag | Description | Default Value |
| --- | --- | --- |
| config.file | Path to configuration file | spectrumVirtualize.yml |
| web.telemetry-path | Path under which to expose metrics | /metrics |
| web.listen-address | Address on which to expose metrics and web interface | :9119 |
| web.disable-exporter-metrics | Exclude metrics about the exporter itself (promhttp_*, process_*, go_*) | true |
| --collector.name | Collector are enabled, the name means name of CLI Command | By default enabled collectors: lssystem and lssystemstats. |
| --no-collector.name | Collectors that are enabled by default can be disabled, the name means name of CLI Command | By default disabled collectors: lsnodestats, lsmdisk, lsmdiskgrp, lsvdisk and lsvdiskcopy. |

## Building and running

* Prerequisites:
  * Go compiler

* Building:
  
  * binary

    ```bash
    export GOPATH=your_gopath
    cd your_gopath
    git clone git@github.ibm.com:ZaaS/spectrum-virtualize-exporter.git
    cd spectrum-virtualize-exporter
    go build
    go install (Optional but recommended. This step will copy spectrum-virtualize-exporter binary package into $GOPATH/bin directory. It will be connvenient to copy the package to Monitoring docker image)
    ```

  * docker image

    ```bash
    docker build -t spectrum-virtualize-exporter .
    ```

* Running:
  * Run Locally

    ```bash
    ./spectrum-virtualize-exporter --config.file=/etc/spectrumVirtualize/spectrumVirtualize.yml
    ```

  * Run as docker image

    ```bash
    docker run -it -d -p 9119:9119 -v /etc/spectrumVirtualize/spectrumVirtualize.yml:/etc/spectrumVirtualize/spectrumVirtualize.yml --name spectrum-virtualize-exporter spectrum-virtualize-exporter --config.file=/etc/spectrumVirtualize/spectrumVirtualize.yml --log.level debug
    ```

  * Visit http://localhost:9119/metrics

## Configuration

The spectrum-virtualize-exporter reads from [spectrumVirtualize.yml](spectrumVirtualize.yml) config file by default. Edit your config YAML file, Enter the IP address of the storage device, your username, and your password there. 

```bash
targets:
  - ipAddress: IP address
    userid: user
    password: password
    verifyCert: true
```

## Exported Metrics

| CLI Command | Description | Default | Metrics | Total number of metrics |
| --- | --- | --- | --- | --- |
| - | Metrics from the exporter itself. | Disabled | [List](docs/exporter_metrics.md) | 35 |
| lssystem | Get a detailed view of a clustered system (system). | Enabled | [List](docs/lssystem_metrics.md) | 57 |
| lssystemstats | Get the most recent values of all node statistics in a system. | Enabled | [List](docs/lssystemstats_metrics.md) | 49 |
| lsnodestats | Ge the most recent values of statistics for all nodes. | Disabled | [List](docs/lsnodestats_metrics.md)| 46 |
| lsmdisk | Get a detailed view of managed disks (MDisks) visible to the clustered system. | Disabled | [List](docs/lsmdisk_metrics.md) | 1 |
| lsmdiskgrp | Get a detailed view of storage pools that are visible to the clustered system. | Disabled | [List](docs/lsmdiskgrp_metrics.md) | 16 |
| lsvdisk | Get detailed view of volumes that are recognized by the system. | Disabled | [List](docs/lsvdisk_metrics.md) | 1 |
| lsvdiskcopy | Get volume copy information. | Disabled | [List](docs/lsvdiskcopy_metrics.md) | 1 |

## Exported Settings

| CLI Command | Description | Default | Metrics | Total number of metrics |
| --- | --- | --- | --- | --- |
| - | Metrics from the exporter itself. | Disabled | [List](docs/exporter_metrics.md) | 35 |
| lscloudcallhome | The status of the Call Home information. | Enabled | [List](docs/lscloudcallhome_metrics.md) | 1 |
| lsenclosurebattery | The information about the batteries. | Enabled | [List](docs/lsenclosurebattery_metrics.md) | 2 |

## References

* [IBM Spectrum Virtualize RESTful API For FS9xxx](https://www.ibm.com/support/knowledgecenter/en/STSLR9_8.2.0/com.ibm.fs9100_820.doc/rest_api_overview.html)

* [IBM Spectrum Virtualize Interfacing Using the RESTful API](https://www.ibm.com/support/knowledgecenter/STVLF4_8.1.3/spectrum.virtualize.813.doc/Spectrum_Virtualize_API_8.1.3.pdf)
