package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/collector"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/s_collector"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	configFile             = kingpin.Flag("config.file", "Path to configuration file.").Default("spectrumVirtualize.yml").String()
	metricsContext         = kingpin.Flag("web.metrics-context", "Context under which to expose metrics.").Default("/metrics").String()
	settingsContext        = kingpin.Flag("web.settings-context", "Context under which to expose settings.").Default("/settings").String()
	listenAddress          = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9119").String()
	disableExporterMetrics = kingpin.Flag("web.disable-exporter-metrics", "Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).").Default("true").Bool()
	// maxRequests            = kingpin.Flag("web.max-requests", "Maximum number of parallel scrape requests. Use 0 to disable.").Default("40").Int()
	cfg                     *utils.Config
	enableMetricsCollector  bool = true
	enableSettingsCollector bool = true
	authTokenCache          sync.Map
)

type handler struct {
	// exporterMetricsRegistry is a separate registry for the metrics about the exporter itself.
	exporterMetricsRegistry *prometheus.Registry
	includeExporterMetrics  bool
	// maxRequests             int
}

func main() {
	r := mux.NewRouter()
	CSRF := csrf.Protect([]byte("spectrum-expor-32-bytes-auth-key"))
	// Parse flags.
	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("spectrum_virtualize_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	//Bail early if the config is bad.
	log.Infoln("Loading config from", *configFile)
	// var err error
	c, err := utils.GetConfig(*configFile)
	if err != nil {
		log.Fatalf("Error parsing config file: %s", err)
	}
	cfg = c

	log.Infoln("Starting Spectrum_Virtualize_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())
	//Launch http services
	// http.HandleFunc(*metricsContext, handlerMetricRequest)
	r.Handle(*metricsContext, newHandler(!*disableExporterMetrics))
	r.Handle(*settingsContext, newHandler(!*disableExporterMetrics))
	//	r.HandleFunc(*settingsContext, testFunc)
	r.HandleFunc("/", rootFunc)
	// http.Handle(*metricsContext, prometheus.Handler()) // Normal metrics endpoint for Spectrum Virtualize exporter itself.

	log.Infof("Listening for %s on %s\n", *metricsContext, *listenAddress)
	log.Infof("Listening for %s on %s\n", *settingsContext, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, CSRF(r)))

}

func rootFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write([]byte(`<html>
		<head><title>Spectrum Virtualize exporter</title></head>
		<body>
			<h1>Spectrum Virtualize exporter</h1>
			<p><a href='` + *metricsContext + `'>Metrics</a></p>
		</body>
	</html>`))
	} else {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
	}
}

func targetsForRequest(r *http.Request) ([]utils.Target, error) {
	reqTarget := r.URL.Query().Get("target")
	if reqTarget == "" {
		return cfg.Targets, nil
	}
	for _, t := range cfg.Targets {
		if t.IpAddress == reqTarget {
			return []utils.Target{t}, nil
		}
	}
	return nil, fmt.Errorf("the target '%s' not defined in the configuration file", reqTarget)
}

func newHandler(includeExporterMetrics bool) *handler {
	h := &handler{
		exporterMetricsRegistry: prometheus.NewRegistry(),
		includeExporterMetrics:  includeExporterMetrics,
		// maxRequests:             maxRequests,
	}
	if h.includeExporterMetrics {
		h.exporterMetricsRegistry.MustRegister(
			prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
			prometheus.NewGoCollector(),
		)
	}
	return h
}

// ServeHTTP implements http.Handler.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var handler http.Handler
		targets, err := targetsForRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		//handler, err := h.settingsHandler(targets...)
		if r.RequestURI == "/metrics" {
			handler, err = h.metricsHandler(targets...)
		}
		if r.RequestURI == "/settings" {
			handler, err = h.settingsHandler(targets...)
		}

		if err != nil {
			log.Warnln("Couldn't create handler:", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Couldn't create handler: %s", err)))
			return
		}
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
	}
}

func (h *handler) metricsHandler(targets ...utils.Target) (http.Handler, error) {

	registry := prometheus.NewRegistry()
	sc, err := collector.NewSVCCollector(targets, &authTokenCache) //new a Spectrum Virtualize Collector
	// registry.MustRegister(version.NewCollector("Spectrum-Virtualize-Exporter"))

	if err != nil {
		log.Fatalf("Couldn't create metrics collector: %s", err)
	}
	if enableMetricsCollector {
		log.Infof("Enabled metrics collectors:")
		for n := range sc.Collectors {
			log.Infof(" - %s", n)
		}
		enableMetricsCollector = false
	}

	if err := registry.Register(sc); err != nil {
		return nil, fmt.Errorf("couldn't register metrics SVC collector: %s", err)
	}
	handler := promhttp.HandlerFor(
		prometheus.Gatherers{h.exporterMetricsRegistry, registry},
		promhttp.HandlerOpts{
			ErrorLog:      log.NewErrorLogger(),
			ErrorHandling: promhttp.ContinueOnError,
			// MaxRequestsInFlight: h.maxRequests,
		},
	)
	if h.includeExporterMetrics {
		// Note that we have to use h.exporterMetricsRegistry here to
		// use the same promhttp metrics for all expositions.
		handler = promhttp.InstrumentMetricHandler(
			h.exporterMetricsRegistry, handler,
		)
	}
	return handler, nil
}

func (h *handler) settingsHandler(targets ...utils.Target) (http.Handler, error) {

	registry := prometheus.NewRegistry()
	sc, err := s_collector.NewSVCCollector(targets, &authTokenCache) //new a Spectrum Virtualize Collector
	// registry.MustRegister(version.NewCollector("Spectrum-Virtualize-Exporter"))

	if err != nil {
		log.Fatalf("Couldn't create setting collector: %s", err)
	}
	if enableSettingsCollector {
		log.Infof("Enabled setting collectors:")
		for n := range sc.Collectors {
			log.Infof(" - %s", n)
		}
		enableSettingsCollector = false
	}

	if err := registry.Register(sc); err != nil {
		return nil, fmt.Errorf("couldn't register settings SVC collector: %s", err)
	}
	handler := promhttp.HandlerFor(
		prometheus.Gatherers{h.exporterMetricsRegistry, registry},
		promhttp.HandlerOpts{
			ErrorLog:      log.NewErrorLogger(),
			ErrorHandling: promhttp.ContinueOnError,
			// MaxRequestsInFlight: h.maxRequests,
		},
	)
	if h.includeExporterMetrics {
		// Note that we have to use h.exporterMetricsRegistry here to
		// use the same promhttp metrics for all expositions.
		handler = promhttp.InstrumentMetricHandler(
			h.exporterMetricsRegistry, handler,
		)
	}
	return handler, nil
}
