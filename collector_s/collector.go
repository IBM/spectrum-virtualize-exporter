package collector_s

import (
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	prefix          = "spectrum_collector_s_"
	defaultEnabled  = true
	defaultDisabled = false
)

var (
	mySC                       *svcCollector = nil
	once                       sync.Once
	scrapeDurationDesc         *prometheus.Desc
	authTokenRenewIntervalDesc *prometheus.Desc
	authTokenRenewSuccessDesc  *prometheus.Desc
	authTokenRenewFailureDesc  *prometheus.Desc
	hosts                      []utils.Target
	collectors                 map[string]Collector
	sClients                   map[string]*utils.SpectrumClient
	factories                  = make(map[string]func() (Collector, error))
	collectorState             = make(map[string]*bool)
)

type SVCCollector interface {
	Describe(ch chan<- *prometheus.Desc)
	Collect(ch chan<- prometheus.Metric)
}

// SVCollector implements the prometheus.Collector interface
type svcCollector struct{}

func init() {
	labelnames := []string{"target", "resource"}
	// metric name, help information, Array of defined label names, defined labels
	scrapeDurationDesc = prometheus.NewDesc("spectrum_collector_scrape_duration_seconds", "Duration of a collector scraping for one host", labelnames, nil)
	authTokenRenewIntervalDesc = prometheus.NewDesc("spectrum_collector_authtoken_renew_interval_seconds", "Interval of the last renewing auth token", labelnames, nil)
	authTokenRenewSuccessDesc = prometheus.NewDesc("spectrum_collector_authtoken_renew_success_total", "Cumulative count of successful verification of renewed auth token", labelnames, nil)
	authTokenRenewFailureDesc = prometheus.NewDesc("spectrum_collector_authtoken_renew_failure_total", "Cumulative count of failed verification of renewed auth token", labelnames, nil)
}

func registerCollector(collector string, isDefaultEnabled bool, factory func() (Collector, error)) {
	var helpDefaultState string
	if isDefaultEnabled {
		helpDefaultState = "enabled"
	} else {
		helpDefaultState = "disabled"
	}

	flagName := fmt.Sprintf("collector.%s", collector)
	flagHelp := fmt.Sprintf("Enable the %s collector (default: %s).", collector, helpDefaultState)
	defaultValue := fmt.Sprintf("%v", isDefaultEnabled)

	flag := kingpin.Flag(flagName, flagHelp).Default(defaultValue).Bool()
	collectorState[collector] = flag

	factories[collector] = factory
}

// NewSVCCollector creates a new Spectrum Virtualize Collector.
func NewSVCCollector(targets []utils.Target, tokenCaches map[string]*utils.AuthToken, tokenMutexes map[string]*sync.Mutex, colCounters map[string]*utils.Counter) (SVCCollector, error) {
	var (
		err       error = nil
		collector Collector
	)
	once.Do(func() {
		hosts = targets
		collectors = make(map[string]Collector)
		log.Infof("Enabled setting collectors:")
		for key, enabled := range collectorState {
			if *enabled {
				collector, err = factories[key]()
				if err != nil {
					log.Errorln("Failed to load setting collector: ", key)
					return
				}
				collectors[key] = collector
				log.Infof(" - %s", key)
			}
		}
		sClients = make(map[string]*utils.SpectrumClient)
		for _, t := range targets {
			sClients[t.IpAddress] = &utils.SpectrumClient{
				UserName:       t.Userid,
				Password:       t.Password,
				IpAddress:      t.IpAddress,
				VerifyCert:     t.VerifyCert,
				AuthTokenCache: tokenCaches[t.IpAddress],
				AuthTokenMutex: tokenMutexes[t.IpAddress],
				ColCounter:     colCounters[t.IpAddress],
			}
		}
		mySC = &svcCollector{}
	})
	return *mySC, err
}

// Describe implements the Prometheus.Collector interface.
func (c svcCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- authTokenRenewSuccessDesc
	ch <- authTokenRenewFailureDesc
	ch <- authTokenRenewIntervalDesc
	for _, col := range collectors {
		col.Describe(ch)
	}
}

// Collect implements the Prometheus.Collector interface.
func (c svcCollector) Collect(ch chan<- prometheus.Metric) {
	wg := &sync.WaitGroup{}
	wg.Add(len(hosts))
	for _, h := range hosts {
		go c.collectForHost(h, ch, wg)
	}
	wg.Wait()
}

func (c *svcCollector) collectForHost(host utils.Target, ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	success := 0
	var counter utils.Counter
	spectrumClient := sClients[host.IpAddress]
	defer func() {
		ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(start).Seconds(), spectrumClient.IpAddress, spectrumClient.Hostname)
		ch <- prometheus.MustNewConstMetric(authTokenRenewSuccessDesc, prometheus.CounterValue, float64(counter.AuthTokenRenewSuccessCount), spectrumClient.IpAddress, spectrumClient.Hostname)
		ch <- prometheus.MustNewConstMetric(authTokenRenewFailureDesc, prometheus.CounterValue, float64(counter.AuthTokenRenewFailureCount), spectrumClient.IpAddress, spectrumClient.Hostname)
		ch <- prometheus.MustNewConstMetric(authTokenRenewIntervalDesc, prometheus.GaugeValue, float64(counter.AuthTokenRenewIntervalSeconds), spectrumClient.IpAddress, spectrumClient.Hostname)
	}()

	counter, success = sClients[host.IpAddress].RenewAuthToken(true)

	if success == 0 {
		log.Errorln("No valid auth token, skip executing setting collectors")
	} else {
		for k, col := range collectors {
			err := col.Collect(*spectrumClient, ch)
			if err != nil && err.Error() != "EOF" {
				log.Errorln(k + ": " + err.Error())
			}
		}
	}
}

// Collector is the interface a collector has to implement.
type Collector interface {
	// Describe metrics
	Describe(ch chan<- *prometheus.Desc)

	// Collect metrics
	Collect(client utils.SpectrumClient, ch chan<- prometheus.Metric) error
}
