package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/collector"
	"github.ibm.com/ZaaS/spectrum-virtualize-exporter/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	configFile             = kingpin.Flag("config.file", "Path to configuration file.").Default("spectrumVirtualize.yml").String()
	metricsPath            = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	listenAddress          = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9119").String()
	disableExporterMetrics = kingpin.Flag("web.disable-exporter-metrics", "Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).").Bool()
	// maxRequests            = kingpin.Flag("web.max-requests", "Maximum number of parallel scrape requests. Use 0 to disable.").Default("40").Int()
	cfg             *utils.Config
	enableCollector bool = true
)

type handler struct {
	// exporterMetricsRegistry is a separate registry for the metrics about the exporter itself.
	exporterMetricsRegistry *prometheus.Registry
	includeExporterMetrics  bool
	// maxRequests             int
}

func main() {

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
	// http.HandleFunc(*metricsPath, handlerMetricRequest)
	http.Handle(*metricsPath, newHandler(!*disableExporterMetrics))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Write([]byte(`<html>
			<head><title>Spectrum Virtualize exporter</title></head>
			<body>
				<h1>Spectrum Virtualize exporter</h1>
				<p><a href='` + *metricsPath + `'>Metrics</a></p>
			</body>
		</html>`))
		} else {
			http.Error(w, "403 Forbidden", 403)
		}
	})
	// http.Handle(*metricsPath, prometheus.Handler()) // Normal metrics endpoint for Spectrum Virtualize exporter itself.

	log.Infof("Listening for %s on %s\n", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))

}

func targetsForRequest(r *http.Request) ([]utils.Targets, error) {
	reqTarget := r.URL.Query().Get("target")
	if reqTarget == "" {
		return cfg.Targets, nil
	}
	for _, t := range cfg.Targets {
		if t.IpAddress == reqTarget {
			return []utils.Targets{t}, nil
		}
	}

	return nil, fmt.Errorf("The target '%s' os not defined in the configuration file", reqTarget)
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
		targets, err := targetsForRequest(r)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		} else {
			handler, err := h.innerHandler(targets...)
			if err != nil {
				log.Warnln("Couldn't create  metrics handler:", err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("Couldn't create  metrics handler: %s", err)))
				return
			}
			handler.ServeHTTP(w, r)
		}
	} else {
		http.Error(w, "403 Forbidden", 403)
	}
}

func (h *handler) innerHandler(targets ...utils.Targets) (http.Handler, error) {

	registry := prometheus.NewRegistry()
	sc, err := collector.NewSVCCollector(targets) //new a Spectrum Virtualize Collector
	// registry.MustRegister(version.NewCollector("Spectrum-Virtualize-Exporter"))

	if err != nil {
		log.Fatalf("Couldn't create collector: %s", err)
	}
	if enableCollector == true {
		log.Infof("Enabled collectors:")
		for n := range sc.Collectors {
			log.Infof(" - %s", n)
		}
		enableCollector = false
	}

	if err := registry.Register(sc); err != nil {
		return nil, fmt.Errorf("couldn't register SVC collector: %s", err)
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
