package collector_s

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_drive = "spectrum_drive_"

var (
	drive_status                     *prometheus.Desc
	drive_firmware_level             *prometheus.Desc
	drive_firmware_level_consistency *prometheus.Desc
)

func init() {
	registerCollector("lsdrive", defaultEnabled, NewDriveCollector)
}

//driveCollector collects drive setting metrics
type DriveCollector struct {
}

func NewDriveCollector() (Collector, error) {
	labelnames_drive := []string{"resource", "drive_id", "enclosure_id", "slot_id"}
	labelnames_firmware := []string{"resource", "drive_id", "firmware_level"}
	labelnames_firmware_consistency := []string{"resource"}
	if len(utils.ExtraLabelNames) > 0 {
		labelnames_drive = append(labelnames_drive, utils.ExtraLabelNames...)
		labelnames_firmware = append(labelnames_firmware, utils.ExtraLabelNames...)
		labelnames_firmware_consistency = append(labelnames_firmware_consistency, utils.ExtraLabelNames...)
	}
	drive_status = prometheus.NewDesc(prefix_drive+"status", "Indicates the status of the drive. 0-online; 1-offline; 2-degraded.", labelnames_drive, nil)
	drive_firmware_level = prometheus.NewDesc(prefix_drive+"firmware_level", "Indicates the firmware level consistency of disks. 0-consistent; 1-inconsistent.", labelnames_firmware, nil)
	drive_firmware_level_consistency = prometheus.NewDesc(prefix_drive+"firmware_level_consistency", "Indicates the firmware level consistency of disks. 0-consistent; 1-inconsistent.", labelnames_firmware_consistency, nil)
	return &DriveCollector{}, nil
}

//Describe describes the metrics
func (*DriveCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- drive_status
	ch <- drive_firmware_level
	ch <- drive_firmware_level_consistency
}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *DriveCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	logger.Debugln("entering drive collector ...")
	respData, err := sClient.CallSpectrumAPI("lsdrive", true)
	if err != nil {
		logger.Errorf("executing lsdrive cmd failed: %s", err.Error())
		return err
	}
	logger.Debugln("response of lsdrive: ", respData)
	/* This is a sample output of lsdrive
	[
	    {
	        "id": "0",
	        "status": "online",
	        "error_sequence_number": "",
	        "use": "member",
	        "tech_type": "tier0_flash",
	        "capacity": "20.0TB",
	        "mdisk_id": "0",
	        "mdisk_name": "mdisk0",
	        "member_id": "0",
	        "enclosure_id": "1",
	        "slot_id": "1",
	        "node_id": "",
	        "node_name": "",
	        "auto_manage": "inactive",
	        "drive_class_id": "0"
	    },
		...
	] */
	if !gjson.Valid(respData) {
		return fmt.Errorf("invalid json for lsdrive:\n%v", respData)
	}
	jsonDrives := gjson.Parse(respData)
	var drives []string
	jsonDrives.ForEach(func(key, drive gjson.Result) bool {
		drive_id := drive.Get("id").String()
		enclosure_id := drive.Get("enclosure_id").String()
		slot_id := drive.Get("slot_id").String()
		status := drive.Get("status").String() // ["online", "offline", "degraded"]
		drives = append(drives, drive_id)

		v_status := 0
		switch status {
		case "online":
			v_status = 0
		case "offline":
			v_status = 1
		case "degraded":
			v_status = 2
		}
		labelvalues_drive := []string{sClient.Hostname, drive_id, enclosure_id, slot_id}
		if len(utils.ExtraLabelValues) > 0 {
			labelvalues_drive = append(labelvalues_drive, utils.ExtraLabelValues...)
		}
		ch <- prometheus.MustNewConstMetric(drive_status, prometheus.GaugeValue, float64(v_status), labelvalues_drive...)
		return true
	})
	v_firmware_consistency_total := 0
	v_firmware_consistency := 0
	base_level := ""
	for _, drive_id := range drives {
		resp, err := sClient.CallSpectrumAPI("lsdrive/"+drive_id, true)
		if err != nil {
			logger.Errorf("executing lsdrive/%s cmd failed: %s", drive_id, err.Error())
			return err
		}
		logger.Debugf("response of lsdrive/%s: %s", drive_id, resp)
		/* This is a sample output of lsdrive/<id>
		{
		    "id": "0",
		    "status": "online",
		    "error_sequence_number": "",
		    "use": "member",
		    "UID": "4e564d653031454b323332595331374247383752303558001014101406480000",
		    "tech_type": "tier0_flash",
		    "capacity": "20.0TB",
		    "block_size": "512",
		    "vendor_id": "IBM-C062",
		    "product_id": "10140648",
		    "FRU_part_number": "01YM583",
		    "FRU_identity": "11S01EK232YS17BG87R05X",
		    "RPM": "",
		    "firmware_level": "1_2_11  ",
		    "FPGA_level": "",
		    "date_of_manufacture": "",
		    "mdisk_id": "0",
		    "mdisk_name": "mdisk0",
		    "member_id": "0",
		    "enclosure_id": "1",
		    "slot_id": "1",
		    "node_id": "",
		    "node_name": "",
		    "quorum_id": "0",
		    "port_1_status": "online",
		    "port_2_status": "online",
		    "interface_speed": "",
		    "protection_enabled": "yes",
		    "auto_manage": "inactive",
		    "drive_class_id": "0",
		    "write_endurance_used": "0",
		    "write_endurance_usage_rate": "",
		    "replacement_date": "",
		    "transport_protocol": "nvme",
		    "compressed": "yes",
		    "physical_capacity": "8.73TB",
		    "physical_used_capacity": "4.30TB",
		    "effective_used_capacity": "4.35TB"
		} */
		if !gjson.Valid(resp) {
			return fmt.Errorf("invalid json for lsdrive/%s:\n%v", drive_id, resp)
		}
		jsonDrive := gjson.Parse(resp)
		firmware_level := strings.TrimSpace(jsonDrive.Get("firmware_level").String())
		if firmware_level == "" {
			firmware_level = "unknown"
		}

		if base_level == "" {
			base_level = firmware_level
		} else {
			if firmware_level == base_level {
				v_firmware_consistency = 0
			} else {
				v_firmware_consistency = 1
				v_firmware_consistency_total = 1
			}
		}
		labelvalues_firmware := []string{sClient.Hostname, drive_id, firmware_level}
		if len(utils.ExtraLabelValues) > 0 {
			labelvalues_firmware = append(labelvalues_firmware, utils.ExtraLabelValues...)
		}
		ch <- prometheus.MustNewConstMetric(drive_firmware_level, prometheus.GaugeValue, float64(v_firmware_consistency), labelvalues_firmware...)
	}
	labelvalues_firmware_consistency := []string{sClient.Hostname}
	if len(utils.ExtraLabelValues) > 0 {
		labelvalues_firmware_consistency = append(labelvalues_firmware_consistency, utils.ExtraLabelValues...)
	}
	ch <- prometheus.MustNewConstMetric(drive_firmware_level_consistency, prometheus.GaugeValue, float64(v_firmware_consistency_total), labelvalues_firmware_consistency...)

	logger.Debugln("exit drive exit")
	return nil
}
