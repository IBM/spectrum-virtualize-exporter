package s_collector

import (
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/tidwall/gjson"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	prefix          = "spectrum_s_collector_"
	defaultEnabled  = true
	defaultDisabled = false
)

var (
	scrapeDurationDesc         *prometheus.Desc
	scrapeSuccessDesc          *prometheus.Desc
	authTokenRequestErrorDesc  *prometheus.Desc
	authTokenRenewSuccessDesc  *prometheus.Desc
	authTokenRenewFailureDesc  *prometheus.Desc
	authTokenCache             *sync.Map
	authTokenRequestErrorCount int = 0
	authTokenRenewSuccessCount int = 0
	authTokenRenewFailureCount int = 0
	factories                      = make(map[string]func() (Collector, error))
	collectorState                 = make(map[string]*bool)
)

// SVCollector implements the prometheus.Collector interface
type SVCCollector struct {
	targets    []utils.Target
	Collectors map[string]Collector
}

func init() {
	labelnames := []string{"target", "resource"}
	// metric name, help information, Array of defined label names, defined labels
	scrapeDurationDesc = prometheus.NewDesc(prefix+"scrape_duration_seconds", "Duration of a collector scrape for one resource", labelnames, nil)
	scrapeSuccessDesc = prometheus.NewDesc(prefix+"scrape_success", "Scrape of resource is successful or not", labelnames, nil)
	authTokenRequestErrorDesc = prometheus.NewDesc(prefix+"authtoken_request_error_total", "Cumulative error count of requesting auth token", labelnames, nil)
	authTokenRenewSuccessDesc = prometheus.NewDesc(prefix+"authtoken_renew_success_total", "Cumulative count of success verification of renewed auth token", labelnames, nil)
	authTokenRenewFailureDesc = prometheus.NewDesc(prefix+"authtoken_renew_failure_total", "Cumulative count of failed verification of renewed auth token", labelnames, nil)
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

// newSVCCollector creates a new Spectrum Virtualize Collector.
func NewSVCCollector(targets []utils.Target, authToken *sync.Map) (*SVCCollector, error) {
	authTokenCache = authToken
	collectors := make(map[string]Collector)
	// log.Infof("Enabled collectors:")
	for key, enabled := range collectorState {
		if *enabled {
			// log.Infof(" - %s", key)
			collector, err := factories[key]()
			if err != nil {
				return nil, err
			}
			collectors[key] = collector
		}
	}
	return &SVCCollector{targets, collectors}, nil
}

// Describe implements the Prometheus.Collector interface.
func (c SVCCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeSuccessDesc
	ch <- scrapeDurationDesc
	ch <- authTokenRequestErrorDesc
	ch <- authTokenRenewSuccessDesc
	ch <- authTokenRenewFailureDesc

	for _, col := range c.Collectors {
		col.Describe(ch)
	}
}

// Collect implements the Prometheus.Collector interface.
func (c SVCCollector) Collect(ch chan<- prometheus.Metric) {

	hosts := c.targets
	wg := &sync.WaitGroup{}
	wg.Add(len(hosts))
	for _, h := range hosts {
		go c.collectForHost(h, ch, wg)
	}
	wg.Wait()
}

func (c *SVCCollector) collectForHost(host utils.Target, ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	success := 0
	spectrumClient := utils.SpectrumClient{
		UserName:   host.Userid,
		Password:   host.Password,
		IpAddress:  host.IpAddress,
		VerifyCert: host.VerifyCert,
	}
	defer func() {
		ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(start).Seconds(), spectrumClient.IpAddress, spectrumClient.Hostname)
		ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, float64(success), spectrumClient.IpAddress, spectrumClient.Hostname)
		ch <- prometheus.MustNewConstMetric(authTokenRequestErrorDesc, prometheus.CounterValue, float64(authTokenRequestErrorCount), spectrumClient.IpAddress, spectrumClient.Hostname)
		ch <- prometheus.MustNewConstMetric(authTokenRenewFailureDesc, prometheus.CounterValue, float64(authTokenRenewFailureCount), spectrumClient.IpAddress, spectrumClient.Hostname)
		ch <- prometheus.MustNewConstMetric(authTokenRenewSuccessDesc, prometheus.CounterValue, float64(authTokenRenewSuccessCount), spectrumClient.IpAddress, spectrumClient.Hostname)
	}()

	lc := 1
	for lc < 4 {
		log.Debugf("Looking for cached Auth Token for %s", host.IpAddress)
		cacheHit := false
		result, ok := authTokenCache.Load(host.IpAddress)
		if !ok {
			log.Debug("Authtoken not found in cache.")
			log.Debugf("Retrieving authToken for %s", host.IpAddress)
			// get our authtoken for future interactions
			authtoken, err := spectrumClient.RetriveAuthToken()
			if err != nil {
				log.Errorf("Error getting auth token for %s, the error was %v.", host.IpAddress, err)
				authTokenRequestErrorCount++
				return
			}
			authTokenCache.Store(host.IpAddress, authtoken)
			spectrumClient.AuthToken = authtoken
		} else {
			cacheHit = true
			log.Debugf("Authtoken pulled from cache for %s", host.IpAddress)
			spectrumClient.AuthToken = result.(string)
		}
		//test to make sure that our auth token is good
		// if not delete it and loop back
		i := 0
		for i < 2 {
			validateURL := "https://" + host.IpAddress + ":7443/rest/lssystem"
			systemMetrics, err := spectrumClient.CallSpectrumAPI(validateURL)
			if err != nil {
				i++
				time.Sleep(2 * time.Second)
				continue
			} else {
				spectrumClient.Hostname = gjson.Get(systemMetrics, "name").String()
				//We have a valid auth token, we can break out of this loop
				break
			}
		}
		if i > 1 { //new auth token verification failed
			if !cacheHit {
				authTokenRenewFailureCount++
			}
			authTokenCache.Delete(host.IpAddress)
			log.Infof("\nInvalid authToken for %s, re-requesting authtoken....", host.IpAddress)
			lc++
		} else { //new auth token verification succeeded
			if !cacheHit {
				authTokenRenewSuccessCount++
			}
			break
		}
	}
	if lc > 3 {
		log.Errorf("Error getting auth token for %s, please check network or username and password", host.IpAddress)
		return
	}
	success = 1
	for k, col := range c.Collectors {
		err := col.Collect(spectrumClient, ch)
		if err != nil && err.Error() != "EOF" {
			log.Errorln(k + ": " + err.Error())
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
