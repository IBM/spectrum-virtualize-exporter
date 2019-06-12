package collector

import (
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/spectrum-virtualize-exporter/utils"
	"github.com/tidwall/gjson"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	prefix          = "spectrum_"
	defaultEnabled  = true
	defaultDisabled = false
)

var (
	scrapeDurationDesc        *prometheus.Desc
	scrapeSuccessDesc         *prometheus.Desc
	requestErrors             *prometheus.Desc
	authTokenCacheCounterHit  *prometheus.Desc
	authTokenCacheCounterMiss *prometheus.Desc
	authTokenCache            sync.Map
	requestErrorCount         int = 0
	authTokenMiss             int = 0
	authTokenHit              int = 0
	factories                     = make(map[string]func() (Collector, error))
	collectorState                = make(map[string]*bool)
)

// SVCollector implements the prometheus.Collecotor interface
type SVCCollector struct {
	targets    []utils.Targets
	Collectors map[string]Collector
}

func init() {
	scrapeDurationDesc = prometheus.NewDesc(prefix+"collector_duration_seconds", "Duration of a collector scrape for one resource", []string{"resource"}, nil) // metric name, help information, Arrar of defined label names, defined labels
	scrapeSuccessDesc = prometheus.NewDesc(prefix+"collector_success", "Scrape of resource was sucessful", []string{"resource"}, nil)
	requestErrors = prometheus.NewDesc(prefix+"request_errors_total", "Errors in request to the Spectrum Virtualize Exporter", []string{"resource"}, nil)
	authTokenCacheCounterHit = prometheus.NewDesc(prefix+"authtoken_cache_counter_hit", "Count of authtoken cache hits", []string{"resource"}, nil)
	authTokenCacheCounterMiss = prometheus.NewDesc(prefix+"authtoken_cache_counter_miss", "Count of authtoken cache misses", []string{"resource"}, nil)
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
func NewSVCCollector(targets []utils.Targets) (*SVCCollector, error) {
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
	ch <- requestErrors
	ch <- authTokenCacheCounterHit
	ch <- authTokenCacheCounterMiss

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

func (c *SVCCollector) collectForHost(host utils.Targets, ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	success := 0
	spectrumClient := utils.SpectrumClient{
		UserName:  host.Userid,
		Password:  host.Password,
		IpAddress: host.IpAddress,
	}
	defer func() {
		ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(start).Seconds(), spectrumClient.Hostname)
		ch <- prometheus.MustNewConstMetric(requestErrors, prometheus.CounterValue, float64(requestErrorCount), spectrumClient.Hostname)
		ch <- prometheus.MustNewConstMetric(authTokenCacheCounterMiss, prometheus.CounterValue, float64(authTokenMiss), spectrumClient.Hostname)
		ch <- prometheus.MustNewConstMetric(authTokenCacheCounterHit, prometheus.CounterValue, float64(authTokenHit), spectrumClient.Hostname)
		ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, float64(success), spectrumClient.Hostname)

	}()
	// Need to get rid of the goto cheat, replacing with a for loop, and ensureing it has backoff and a short circuit
	lc := 1
	for lc < 4 {
		log.Debugf("Looking for cached Auth Token for %s", host.IpAddress)
		result, ok := authTokenCache.Load(host.IpAddress)
		if !ok {
			log.Debug("Authtoken not found in cache.")
			log.Debugf("Retrieving authToken for %s", host.IpAddress)
			// get our authtoken for future interactions
			authtoken, err := spectrumClient.RetriveAuthToken()
			if err != nil {
				log.Debugf("Error getting auth token for %s", host.IpAddress)
				requestErrorCount += 1
				success = 0
				return

			}
			authTokenCache.Store(host.IpAddress, authtoken)
			result, _ := authTokenCache.Load(host.IpAddress)
			spectrumClient.AuthToken = result.(string)
			authTokenMiss += 1
			success = 1
		} else {
			log.Debugf("Authtoken pulled from cache for %s", host.IpAddress)
			spectrumClient.AuthToken = result.(string)
			authTokenHit += 1
			success = 1
		}
		//test to make sure that our auth token is good
		// if not delete it and loop back
		validateURL := "https://" + host.IpAddress + ":7443/rest/lssystem"
		systemMetrics, err := spectrumClient.CallSpectrumAPI(validateURL)
		if err != nil {
			authTokenCache.Delete(host.IpAddress)
			log.Debugf("Invalidating authToken for %s, re-requesting authtoken....", host.IpAddress)
			lc += 1
		} else {
			spectrumClient.Hostname = gjson.Get(systemMetrics, "name").String()
			//We have a valid auth token, we can break out of this loop
			break
		}
	}
	if lc > 3 {
		// looped and failed multiple times, so need to go further
		log.Debugf("Error getting auth token for %s", host.IpAddress)
		requestErrorCount += 1
		success = 0
		return
	}
	for k, col := range c.Collectors {
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
