package collector

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/spectrum-virtualize-exporter/utils"
)

const prefix = "spectrum_"

var (
	scrapeDurationDesc *prometheus.Desc
	scrapeSuccessDesc  *prometheus.Desc

	authTokenCache sync.Map
)

// SVCollector implements the prometheus.Collecotor interface
type SVCCollector struct {
	targets    []utils.Targets
	collectors map[string]Collector
}

func init() {

	scrapeDurationDesc = prometheus.NewDesc(prefix+"collector_duration_seconds", "Duration of a collector scrape for one target", []string{"target"}, nil) // metric name, help information, Arrar of defined label names, defined labels
	scrapeSuccessDesc = prometheus.NewDesc(prefix+"collector_success", "Scrape of target was sucessful", []string{"target"}, nil)

}

// newSVCCollector creates a new Spectrum Virtualize Collector.
func NewSVCCollector(targets []utils.Targets) *SVCCollector {
	collectors := collectors()
	// systemexporter := NewSystemCollector()
	return &SVCCollector{targets, collectors}
}

func collectors() map[string]Collector {
	m := map[string]Collector{}
	m["system"] = NewSystemCollector()
	m["system_stats"] = NewSystemStatsCollector()
	m["node_stats"] = NewNodeStatsCollector()
	m["node_stats"] = NewVolumeCollector()
	m["node_stats"] = NewVolumeCopyCollector()
	return m
}

// Describe implements the Prometheus.Collector interface.
func (c SVCCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeSuccessDesc
	ch <- scrapeDurationDesc

	for _, col := range c.collectors {
		col.Describe(ch)
	}
}

// Collect implements the Prometheus.Collector interface.
func (c SVCCollector) Collect(ch chan<- prometheus.Metric) {

	hosts := c.targets
	wg := &sync.WaitGroup{}
	wg.Add(len(hosts))

	// for _, h := range hosts {
	// 	// go c.collectForHost(strings.Trim(h, " "), ch, wg)
	// }

	for _, h := range hosts {
		go c.collectForHost(h, ch, wg)
	}

	wg.Wait()
}

func (c *SVCCollector) collectForHost(host utils.Targets, ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	defer wg.Done()
	l := []string{host.IpAddress}
	start := time.Now()
	var success float64
	spectrumClient := utils.SpectrumClient{
		UserName:  host.Userid,
		Password:  host.Password,
		IpAddress: host.IpAddress,
	}
	defer func() {
		ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(start).Seconds(), l...)

	}()
	result, ok := authTokenCache.Load(host.IpAddress)
	if !ok {
		log.Debug("Authtoken not found in cache.")
		log.Debugf("Retrieving authToken for %s", host.IpAddress)
		// get our authtoken for future interactions
		a, err := spectrumClient.RetriveAuthToken()
		if err != nil {
			log.Debugf("Error getting auth token for %s", host.IpAddress)
			success = 0

		} else {
			authTokenCache.Store(host.IpAddress, a)
			result, _ := authTokenCache.Load(host.IpAddress)
			spectrumClient.AuthToken = result.(string)

			success = 1
		}
	} else {
		log.Debugf("Authtoken pulled from cache for %s", host.IpAddress)
		spectrumClient.AuthToken = result.(string)
		success = 1

	}

	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, l...)

	for k, col := range c.collectors {
		// err = col.Collect(spectrumClient, ch, l)
		err := col.Collect(spectrumClient, ch)
		if err != nil && err.Error() != "EOF" {
			log.Errorln(k + ": " + err.Error())
		}
	}
}

// Collector is the interface a collector has to implement.
//Collector collects metrics from spectrum virtual using rest api
type Collector interface {
	//Describe describes the metrics
	Describe(ch chan<- *prometheus.Desc)

	//Collect collects metrics from spectrum virtual
	// Collect(client utils.SpectrumClient, ch chan<- prometheus.Metric, labelvalues []string) error
	Collect(client utils.SpectrumClient, ch chan<- prometheus.Metric) error
}
