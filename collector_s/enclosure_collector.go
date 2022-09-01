package collector_s

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_enclosure = "spectrum_enclosure_"

var (
	enclosure_status   *prometheus.Desc
	enclosure_canister *prometheus.Desc
	enclosure_psu      *prometheus.Desc
)

func init() {
	registerCollector("lsenclosure", defaultEnabled, NewEnclosureCollector)
	labelnames_status := []string{"target", "resource", "enclosure_id"}
	labelnames_canister := []string{"target", "resource", "enclosure_id", "total_canisters"}
	labelnames_psu := []string{"target", "resource", "enclosure_id", "total_PSUs"}
	enclosure_status = prometheus.NewDesc(prefix_enclosure+"status", "Indicates whether an enclosure is visible to the SAS network. 0-online; 1-offline; 2-degraded.", labelnames_status, nil)
	enclosure_canister = prometheus.NewDesc(prefix_enclosure+"canister_offline", "Indicates the number of canisters that are contained in this enclosure that are offline. .", labelnames_canister, nil)
	enclosure_psu = prometheus.NewDesc(prefix_enclosure+"psu_offline", "Indicates the number of power-supply units (PSUs) contained in this enclosure that are offline.", labelnames_psu, nil)
}

//enclosureCollector collects enclosure setting metrics
type enclosureCollector struct {
}

func NewEnclosureCollector() (Collector, error) {
	return &enclosureCollector{}, nil
}

//Describe describes the metrics
func (*enclosureCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- enclosure_status
	ch <- enclosure_canister
	ch <- enclosure_psu
}

//Collect collects metrics from Spectrum Virtualize Restful API
func (c *enclosureCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	log.Debugln("Entering enclosure collector ...")
	respData, err := sClient.CallSpectrumAPI("lsenclosure", true)
	if err != nil {
		log.Errorf("Executing lsenclosure cmd failed: %s", err.Error())
		return err
	}
	log.Debugln("Response of lsenclosure: ", respData)
	/* This is a sample output of lsenclosure
	[
	    {
	        "id": "1",
	        "status": "online",
	        "type": "control",
	        "managed": "yes",
	        "IO_group_id": "0",
	        "IO_group_name": "io_grp0",
	        "product_MTM": "9846-AF8",
	        "serial_number": "78E008V",
	        "total_canisters": "2",
	        "online_canisters": "2",
	        "total_PSUs": "2",
	        "online_PSUs": "2",
	        "drive_slots": "24",
	        "total_fan_modules": "0",
	        "online_fan_modules": "0",
	        "total_sems": "0",
	        "online_sems": "0"
	    }
	] */
	if !gjson.Valid(respData) {
		return fmt.Errorf("invalid json for lsenclosure:\n%v", respData)
	}
	jsonEnclosures := gjson.Parse(respData)
	jsonEnclosures.ForEach(func(key, enclosure gjson.Result) bool {
		enclosure_id := enclosure.Get("id").String()
		status := enclosure.Get("status").String() // ["online", "offline", "degraded"]
		canister_total := enclosure.Get("total_canisters").String()
		canister_online := enclosure.Get("online_canisters").String()
		psu_total := enclosure.Get("total_PSUs").String()
		psu_online := enclosure.Get("online_PSUs").String()

		v_status := 0
		switch status {
		case "online":
			v_status = 0
		case "offline":
			v_status = 1
		case "degraded":
			v_status = 2
		}
		ch <- prometheus.MustNewConstMetric(enclosure_status, prometheus.GaugeValue, float64(v_status), sClient.IpAddress, sClient.Hostname, enclosure_id)

		i_canister_total, err := strconv.Atoi(canister_total)
		if err != nil {
			log.Errorf("Parsing total_canister as int failed: %s", err.Error())
		}
		i_canister_online, err := strconv.Atoi(canister_online)
		if err != nil {
			log.Errorf("Parsing online_canister as int failed: %s", err.Error())
		}
		i_canister_offline := i_canister_total - i_canister_online
		ch <- prometheus.MustNewConstMetric(enclosure_canister, prometheus.GaugeValue, float64(i_canister_offline), sClient.IpAddress, sClient.Hostname, enclosure_id, canister_total)

		i_psu_total, err := strconv.Atoi(psu_total)
		if err != nil {
			log.Errorf("Parsing total_psu as int failed: %s", err.Error())
		}
		i_psu_online, err := strconv.Atoi(psu_online)
		if err != nil {
			log.Errorf("Parsing online_psu as int failed: %s", err.Error())
		}
		i_psu_offline := i_psu_total - i_psu_online
		ch <- prometheus.MustNewConstMetric(enclosure_psu, prometheus.GaugeValue, float64(i_psu_offline), sClient.IpAddress, sClient.Hostname, enclosure_id, psu_total)
		return true
	})

	log.Debugln("Leaving enclosure collector.")
	return nil
}
