* [Enhance] Enhance err logging and response examples
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
