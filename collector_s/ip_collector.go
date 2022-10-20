package collector_s

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_ip = "spectrum_ip_"

var (
	ip_status *prometheus.Desc
)

func init() {
	registerCollector("ip", defaultEnabled, NewIPCollector)
	labelnames_status := []string{"target", "resource", "ip_name", "ip_address"}
	ip_status = prometheus.NewDesc(prefix_ip+"status", "IP connection status. 0-connectable; 1-unreachable.", labelnames_status, nil)
}

//ipCollector collects ip setting metrics
type ipCollector struct {
}

func NewIPCollector() (Collector, error) {
	return &ipCollector{}, nil
}

//Describe() describes the metrics
func (*ipCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- ip_status
}

//Collect() collects metrics from Spectrum Virtualize Restful API
func (c *ipCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	log.Debugln("Entering IP collector ...")
	hosts := make(map[string]string)
	str := sClient.IpAddress[:len(sClient.IpAddress)-1]
	hosts["PSYS"] = sClient.IpAddress
	hosts["SSYS"] = str + "1"
	hosts["SVC1"] = str + "2"
	hosts["SVC2"] = str + "3"

	for ip_name, ip_address := range hosts {
		cmd := fmt.Sprintf("ping -c 1 -w 2 %s> /dev/null 2>&1 && echo $? || echo $?", ip_address)
		respData, err := exec.Command("/bin/sh", "-c", cmd).Output()
		if err != nil {
			log.Errorf("Ping %s failed: %s", ip_address, err.Error())
			return err
		}
		status := strings.TrimRight(string(respData), "\n")
		log.Debugf("Ping %s: %s", ip_address, status)

		v_status := 0
		switch status {
		case "0":
			v_status = 0
		default:
			v_status = 1
		}
		ch <- prometheus.MustNewConstMetric(ip_status, prometheus.GaugeValue, float64(v_status), sClient.IpAddress, sClient.Hostname, ip_name, ip_address)
	}

	log.Debugln("Leaving IP collector.")
	return nil
}
