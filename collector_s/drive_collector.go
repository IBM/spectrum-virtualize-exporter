package collector_s

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_drive = "spectrum_drive_"

var (
	drive_status         *prometheus.Desc
	drive_firmware_level *prometheus.Desc
)

func init() {
	registerCollector("lsdrive", defaultEnabled, NewDriveCollector)
	labelnames_drive := []string{"target", "resource", "drive_id"}
	labelnames_firmware := []string{"target", "resource", "drive_id", "firmware_level"}
	drive_status = prometheus.NewDesc(prefix_drive+"status", "Indicates the summary status of the drive. 0-online; 1-offline; 2-degraded.", labelnames_drive, nil)
	drive_firmware_level = prometheus.NewDesc(prefix_drive+"firmware_level", "Indicates the firmware level consistency of disks. 0-consistent; 1-inconsistent.", labelnames_firmware, nil)
}

//driveCollector collects drive setting metrics
type DriveCollector struct {
}

func NewDriveCollector() (Collector, error) {
	return &DriveCollector{}, nil
}

//Describe describes the metrics
func (*DriveCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- drive_status
	ch <- drive_firmware_level
}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *DriveCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	log.Debugln("Entering drive collector ...")
	respData, err := sClient.CallSpectrumAPI("lsdrive", true)
	if err != nil {
		log.Errorf("Executing lsdrive cmd failed: %s", err.Error())
		return err
	}
	log.Debugln("Response of lsdrive: ", respData)
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
		ch <- prometheus.MustNewConstMetric(drive_status, prometheus.GaugeValue, float64(v_status), sClient.IpAddress, sClient.Hostname, drive_id)
		return true
	})
	v_firmware_consistency := 0
	base_level := ""
	for _, drive_id := range drives {
		resp, err := sClient.CallSpectrumAPI("lsdrive/"+drive_id, true)
		if err != nil {
			log.Errorf("Executing lsdrive/%s cmd failed: %s", drive_id, err.Error())
			return err
		}
		log.Debugf("Response of lsdrive/%s: %s", drive_id, resp)
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
			}
		}
		ch <- prometheus.MustNewConstMetric(drive_firmware_level, prometheus.GaugeValue, float64(v_firmware_consistency), sClient.IpAddress, sClient.Hostname, drive_id, firmware_level)
	}

	log.Debugln("Leaving drive collector.")
	return nil
}
