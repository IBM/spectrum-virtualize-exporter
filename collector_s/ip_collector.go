package collector_s

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
)

const prefix_ip = "spectrum_ip_"

var (
	ip_status *prometheus.Desc
)

func init() {
	registerCollector("ip", defaultEnabled, NewIPCollector)
}

//ipCollector collects ip setting metrics
type ipCollector struct {
}

func NewIPCollector() (Collector, error) {
	labelnames := []string{"resource", "ip_name", "ip_address"}
	if len(utils.ExtraLabelNames) > 0 {
		labelnames = append(labelnames, utils.ExtraLabelNames...)
	}
	ip_status = prometheus.NewDesc(prefix_ip+"status", "IP connection status. 0-connectable; 1-unreachable.", labelnames, nil)
	return &ipCollector{}, nil
}

//Describe() describes the metrics
func (*ipCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- ip_status
}

//Collect() collects metrics from Spectrum Virtualize Restful API
func (c *ipCollector) Collect(sClient utils.SpectrumClient, ch chan<- prometheus.Metric) error {

	logger.Debugln("entering IP collector ...")
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
			logger.Errorf("Ping %s failed: %s", ip_address, err.Error())
			return err
		}
		status := strings.TrimRight(string(respData), "\n")
		logger.Debugf("Ping %s: %s", ip_address, status)

		v_status := 0
		switch status {
		case "0":
			v_status = 0
		default:
			v_status = 1
		}

		labelvalues := []string{sClient.Hostname, ip_name, ip_address}
		if len(utils.ExtraLabelValues) > 0 {
			labelvalues = append(labelvalues, utils.ExtraLabelValues...)
		}

		ch <- prometheus.MustNewConstMetric(ip_status, prometheus.GaugeValue, float64(v_status), labelvalues...)
	}

	logger.Debugln("exit IP exit")
	return nil
}
