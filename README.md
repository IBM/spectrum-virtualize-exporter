# spectrum-virtualize-exporter

This [Prometheus](https://prometheus.io) [Exporter](https://prometheus.io/docs/instrumenting/exporters)
collects metrics from storage solutions which are built with the
 [IBM Spectrum Virtualize software](https://www.ibm.com/support/home/product/10000647/IBM_Spectrum_Virtualize_software).
Storage solutions built with the IBM Spectrum Virtualize software are the
[IBM FlashSystem V9000 system](https://www.ibm.com/support/knowledgecenter/STKMQV_8.2.1/com.ibm.storage.vflashsystem9000.8.2.1.doc/svc_svcovr_1bcfiq.html),
the [IBM SAN Volume Controller](https://www.ibm.com/us-en/marketplace/san-volume-controller) and
 the [IBM Storwize Family](https://www.ibm.com/it-infrastructure/storage/storwize\<Paste>)

## Usage

| Flag | Description | Default Value |
| --- | --- | --- |
| config.file | Path to configuration file | spectrumVirtualize.yml |
| web.metrics-context | Context under which to expose metrics | /metrics |
| web.settings-context | Context under which to expose setting metrics | /settings |
| web.listen-address | Address on which to expose metrics and web interface | :9119 |
| web.disable-exporter-metrics | Exclude metrics about the exporter itself (promhttp_*, process_*, go_*) | true |
| --collector.[name] | Enable or disable collector. The [name] is in the list "`lsmdisk`, `lsmdiskgrp`, `lsnodestats`, `lssystem`, `lssystemstats`, `lsvdisk`, `lsvdiskcopy`, `lscloudcallhome`, `lsdrive`, `lsenclosure`, `lsenclosurebattery`, `lsenclosurecanister`, `lsenclosurepsu`, `lshost`, `ip`, `lsmdisk_s`, `lsmdiskgrp_s`, `lsnodecanister`, `lsportfc`" | [true\|false]. <br> By default enabled collectors: `lssystem`, `lssystemstats`, `lscloudcallhome`, `lsdrive`, `lsenclosure`, `lsenclosurebattery`, `lsenclosurecanister`, `lsenclosurepsu`, `lshost`, `ip`, `lsmdisk_s`, `lsmdiskgrp_s`, `lsnodecanister`, `lsportfc`. |

## Building and running

* Prerequisites:
  * Go compiler

* Building:
  
  * binary

    ```bash
    export GOPATH=your_gopath
    cd your_gopath
    git clone https://github.com/IBM/spectrum-virtualize-exporter.git
    cd spectrum-virtualize-exporter
    make binary
    go install (Optional but recommended. This step will copy spectrum-virtualize-exporter binary package into $GOPATH/bin directory. It will be connvenient to copy the package to Monitoring docker image)
    ```

  * docker image

    ```bash
    make docker
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

  * Visit <http://localhost:9119/metrics>

## Configuration

The spectrum-virtualize-exporter loads the [./spectrumVirtualize.yml](spectrumVirtualize.yml) config file by default.

### Required settings

* `targets.[].ipAddress`: IP address of the storage device.
* `targets.[].userid`: Username to access the storage device.
* `targets.[].password`: User password to access the storage device.

### Optionally settings

* `extra_labels.[].name`: Customized label name adding to metrics.
* `extra_labels.[].value`: Value of the customized label.
* `tls_server_config.ca_cert`: The CA certificate chain file in pem format for verifying client certificate.
* `tls_server_config.server_cert`: The server's certificate chain file in pem format.
* `tls_server_config.server_key`: The server's private key file.

### Config File Sample

```yaml
targets:
  - ipAddress: IP address
    userid: user
    password: password
extra_labels:
  - name: pod_name
    value: pod_value
tls_server_config:
  ca_cert: ./certs/ca-root.crt
  server_cert: ./certs/server.crt
  server_key: ./certs/server.key
```

**If any of the "ca_cert", "server_cert" or "server_key" are not provided, the exporter http server will start without https(mTLS) enabled.**

## Exported Metrics

* It recommended to scrape every 30 seconds.

| RESTful API | Description | Default | Metrics | Total number of metrics |
| --- | --- | --- | --- | --- |
| - | Metrics from the prometheus exporter itself. | Disabled | [List](docs/exporter_prometheus_metrics.md) | 30 |
| - | Metrics from the spectrum exporter itself. | Enabled | [List](docs/exporter_spectrum_metrics.md) | 4 |
| lssystem | Get a detailed view of a clustered system (system). | Enabled | [List](docs/lssystem_metrics.md) | 57 |
| lssystemstats | Get the most recent values of all node statistics in a system. | Enabled | [List](docs/lssystemstats_metrics.md) | 49 |
| lsnodestats | Ge the most recent values of statistics for all nodes. | Disabled | [List](docs/lsnodestats_metrics.md)| 46 |
| lsmdisk | Get a detailed view of managed disks (MDisks) visible to the clustered system. | Disabled | [List](docs/lsmdisk_metrics.md) | 1 |
| lsmdiskgrp | Get a detailed view of storage pools that are visible to the clustered system. | Disabled | [List](docs/lsmdiskgrp_metrics.md) | 16 |
| lsvdisk | Get detailed view of volumes that are recognized by the system. | Disabled | [List](docs/lsvdisk_metrics.md) | 1 |
| lsvdiskcopy | Get volume copy information. | Disabled | [List](docs/lsvdiskcopy_metrics.md) | 1 |

## Exported Setting Metrics

* It recommended to scrape every >15 minutes.

| RESTful API | Description | Default | Metrics | Total number of metrics |
| --- | --- | --- | --- | --- |
| - | Metrics from the prometheus exporter itself. | Disabled | [List](docs/exporter_prometheus_metrics.md) | 30 |
| - | Metrics from the spectrum exporter itself. | Enabled | [List](docs/exporter_spectrum_metrics.md) | 4 |
| lscloudcallhome | The status of the Call Home information. | Enabled | [List](docs/lscloudcallhome_settings.md) | 1 |
| lsenclosure | The summary of the enclosures including canister and PSU. | Enabled | [List](docs/lsenclosure_settings.md) | 1 |
| lsenclosurebattery | The information about the batteries. | Enabled | [List](docs/lsenclosurebattery_settings.md) | 2 |
| lsenclosurecanister | The detailed status of each canister in enclosures. | Enabled | [List](docs/lsenclosurecanister_settings.md) | 1 |
| lsenclosurepsu | The information about each power-supply unit (PSU) in enclosures. | Enabled | [List](docs/lsenclosurepsu_settings.md) | 1 |
| lsdrive | The configuration information and drive vital product data (VPD). | Enabled | [List](docs/lsdrive_settings.md) | 3 |
| lshost | The concise information about all the hosts visible to the system. | Enabled | [List](docs/lshost_settings.md) | 1 |
| lsnodecanister | The node canisters that are part of the system. | Enabled | [List](docs/lsnodecanister_settings.md) | 1 |
| lsportfc | The status and properties of the Fibre Channel (FC) input/output (I/O) ports for the clustered system. | Enabled | [List](docs/lsportfc_settings.md) | 1 |
| lsmdisk | The info of managed disks (MDisks) visible to the system. | Enabled | [List](docs/lsmdisk_settings.md) | 1 |
| lsmdiskgrp | The info of storage pools that are visible to the system. | Enabled | [List](docs/lsmdiskgrp_settings.md) | 1 |
| - | The connection status of system IPs(PSYS, SSYS, SVC1, SVC2). | Enabled | [List](docs/lsenclosurebattery_settings.md) | 1 |

## References

* [IBM Spectrum Virtualize RESTful API For FS9xxx](https://www.ibm.com/support/knowledgecenter/en/STSLR9_8.2.0/com.ibm.fs9100_820.doc/rest_api_overview.html)

* [IBM Spectrum Virtualize Interfacing Using the RESTful API](https://www.ibm.com/support/knowledgecenter/STVLF4_8.1.3/spectrum.virtualize.813.doc/Spectrum_Virtualize_API_8.1.3.pdf)

## Contributing

Third party contributions to this project are welcome!

In order to contribute, create a [Git pull request](https://help.github.com/articles/using-pull-requests/), considering this:

* Test is required.
* Each commit should only contain one "logical" change.
* A "logical" change should be put into one commit, and not split over multiple
  commits.
* Large new features should be split into stages.
* The commit message should not only summarize what you have done, but explain
  why the change is useful.
* The commit message must follow the format explained below.

What comprises a "logical" change is subject to sound judgement. Sometimes, it
makes sense to produce a set of commits for a feature (even if not large).
For example, a first commit may introduce a (presumably) compatible API change
without exploitation of that feature. With only this commit applied, it should
be demonstrable that everything is still working as before. The next commit may
be the exploitation of the feature in other components.

For further discussion of good and bad practices regarding commits, see:

* [OpenStack Git Commit Good Practice](https://wiki.openstack.org/wiki/GitCommitMessages)

* [How to Get Your Change Into the Linux Kernel](https://www.kernel.org/doc/Documentation/process/submitting-patches.rst)

## License

The spectrum-virtualize-exporter is licensed under the [Apache 2.0 License](https://github.com/IBM/spectrum-virtualize-exporter/blob/master/LICENSE).
