## 0.9.4 / 2022-04-18

* [FIXBUG] Fix vulunerabilities for CVE-2022-21698

## 0.9.3 / 2021-06-21

* [FEATURE] Support to verify the certification of Spectrum device
* [CHANGE] Add a parameter 'verifyCert' into spectrumVirtualize.yml file to turn on/off the ability of certification verification

## 0.9.2 / 2020-04-30

* [CHANGE] Disable metrics from the exporter itself by default

## 0.9.1 / 2020-01-17

* [CHANGE] Use go module to organise the dependency modules

## 0.9.0 / 2019-09-11

* [CHANGE] Disable http methods other than GET

## 0.8.0 / 2019-07-30

* [BUGFIX] Fix formulas of volume usage.
* [FEATURE] Enhance error logging and response examples

## 0.7.0 / 2019-06-28

* [FEATURE] Add 'target' label, its value is ipaddress

## 0.6.0 / 2019-06-21

* [CHANGE] Change type of token log form INFO to DEBUG.

## 0.5.0 / 2019-06-15

### **Breaking changes**

* [CHANGE] Change the label 'target' from ip to hostname.

### Changes

* [BUGFIX] Fix value of Capacity Usage metrics from decimal to percent.
* [BUGFIX] Fix issue of showing always the "Enabled collectors" logs.
* [CHANGE] Add the column of 'Total number of metrics' to 'Exported Metrics' table.
* [CHANGE] Change the label name from 'target' to 'resource'.

## 0.4.0 / 2019-03-16

### **Breaking changes**

* [CHANGE] By default following collectors are disabled: lsnodestats, lsmdisk,
           lsmdiskgrp, lsvdisk and lsvdiskcopy. Following collectors are
           enabled by default: lssystem and lssystemstats.

### Changes

* [FEATURE] Add disable/enable collector flag.

## 0.3.0 / 2019-03-14

* [BUGFIX] Fix issue of "Connection aborted".

## 0.2.0 / 2019-03-14

* [FEATURE] Add lsmdisk collector.
* [FEATURE] Add lsmdiskgrp collector.
* [FEATURE] Add lsvdisk collector.
* [FEATURE] Add lsvdiskgrp collector.
* [Feature] Add capacity usage metrics.

## 0.1.0 / 2019-03-14

* [CLEANUP] Introduced semantic versioning and changelog. From now on,
  changes will be reported in this file.
