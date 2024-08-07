// Copyright 2021-2024 IBM Corp. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	metricsCollector "github.com/IBM/spectrum-virtualize-exporter/collector"
	settingsCollector "github.com/IBM/spectrum-virtualize-exporter/collector_s"
	"github.com/IBM/spectrum-virtualize-exporter/utils"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	configFile             = kingpin.Flag("config.file", "Path to configuration file.").Default("spectrumVirtualize.yml").String()
	metricsContext         = kingpin.Flag("web.metrics-context", "Context under which to expose metrics.").Default("/metrics").String()
	settingsContext        = kingpin.Flag("web.settings-context", "Context under which to expose settings.").Default("/settings").String()
	listenAddress          = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9119").String()
	disableExporterMetrics = kingpin.Flag("web.disable-exporter-metrics", "Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).").Default("true").Bool()
	// maxRequests            = kingpin.Flag("web.max-requests", "Maximum number of parallel scrape requests. Use 0 to disable.").Default("40").Int()
	cfg *utils.Config
	//enableSettingCollectors bool                        = true
	authTokenCaches  map[string]*utils.AuthToken = make(map[string]*utils.AuthToken)
	authTokenMutexes map[string]*sync.Mutex      = make(map[string]*sync.Mutex)
	colCounters      map[string]*utils.Counter   = make(map[string]*utils.Counter)
	logger           log.Logger                  = *utils.SpectrumLogger()
	https            bool                        = true
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
	logger.Infoln("Loading config from", *configFile)
	// var err error
	c, err := utils.GetConfig(*configFile)
	if err != nil {
		logger.Fatalf("Error parsing config file: %s", err.Error())
		return
	}
	cfg = c
	for _, t := range cfg.Targets {
		authTokenCaches[t.IpAddress] = &utils.AuthToken{}
		authTokenMutexes[t.IpAddress] = &sync.Mutex{}
		colCounters[t.IpAddress] = &utils.Counter{}
	}
	for _, l := range cfg.ExtraLabels {
		utils.ExtraLabelNames = append(utils.ExtraLabelNames, l.Name)
		utils.ExtraLabelValues = append(utils.ExtraLabelValues, l.Value)
	}
	logger.Infoln("Starting Spectrum_Virtualize_exporter", version.Info())
	logger.Infoln("Build context", version.BuildContext())

	if len(utils.ExtraLabelNames) > 0 {
		msg := "Extra labels: ["
		for idx, item := range utils.ExtraLabelNames {
			msg += "    " + item + " => " + utils.ExtraLabelValues[idx] + ";"
		}
		logger.Infoln(msg, "]")
	}
	//Launch http services
	r.Handle(*metricsContext, newHandler(!*disableExporterMetrics))
	r.Handle(*settingsContext, newHandler(!*disableExporterMetrics))
	r.HandleFunc("/", rootFunc)

	if cfg.TlsServerConfig.CaCert != "" && cfg.TlsServerConfig.ServerCert != "" && cfg.TlsServerConfig.ServerKey != "" {
		startHTTPS(CSRF(r))
	} else {
		https = false
		startHTTP(CSRF(r))
	}
}

func startHTTP(handler http.Handler) {
	server := http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  45 * time.Second,
		Addr:         *listenAddress,
		Handler:      handler,
	}
	logger.Infof("Listening(HTTP) for %s on %s\n", *metricsContext, *listenAddress)
	logger.Infof("Listening(HTTP) for %s on %s\n", *settingsContext, *listenAddress)
	logger.Fatal(server.ListenAndServe())
}

func startHTTPS(handler http.Handler) {
	// load CA certificate file and add it to list of client CAs
	caCertFile, err := os.ReadFile(cfg.TlsServerConfig.CaCert)
	if err != nil {
		logger.Fatalf("error reading CA certificate: %s", err.Error())
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertFile)

	// Create the TLS Config with the CA pool and enable Client certificate validation
	tlsConfig := &tls.Config{
		ClientCAs:                caCertPool,
		ClientAuth:               tls.RequireAndVerifyClientCert,
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,       //tls1.2 FIPS/IBM cloud approved
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,         //tls1.2 FIPS/IBM cloud approved
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,       //tls1.2 FIPS/IBM cloud approved
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,         //tls1.2 FIPS/IBM cloud approved
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256, //tls1.2 IBM cloud approved
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,   //tls1.2 IBM cloud approved
			tls.TLS_AES_256_GCM_SHA384,                        //tls1.3 IBM cloud approved
			tls.TLS_AES_128_GCM_SHA256,                        //tls1.3 IBM cloud approved
			tls.TLS_CHACHA20_POLY1305_SHA256,                  //tls1.3 IBM cloud approved
		},
	}

	server := http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  45 * time.Second,
		Addr:         *listenAddress,
		Handler:      handler,
		TLSConfig:    tlsConfig,
	}

	logger.Infof("Listening(HTTPS) for %s on %s\n", *metricsContext, *listenAddress)
	logger.Infof("Listening(HTTPS) for %s on %s\n", *settingsContext, *listenAddress)
	logger.Fatal(server.ListenAndServeTLS(cfg.TlsServerConfig.ServerCert, cfg.TlsServerConfig.ServerKey))
}

func rootFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if https {
			w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}
		_, _ = w.Write([]byte(`<html>
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
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
			collectors.NewGoCollector(),
			//prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
			//prometheus.NewGoCollector(),
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
		if r.RequestURI == "/metrics" {
			handler, err = h.metricsHandler(targets...)
		}
		if r.RequestURI == "/settings" {
			handler, err = h.settingsHandler(targets...)
		}

		if err != nil {
			logger.Warnln("Couldn't create handler:", err)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(fmt.Sprintf("Couldn't create handler: %s", err.Error())))
			return
		}
		if https {
			w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
	}
}

func (h *handler) metricsHandler(targets ...utils.Target) (http.Handler, error) {

	registry := prometheus.NewRegistry()
	sc, err := metricsCollector.NewSVCCollector(targets, authTokenCaches, authTokenMutexes, colCounters) //new a Spectrum Virtualize Collector
	// registry.MustRegister(version.NewCollector("Spectrum-Virtualize-Exporter"))

	if err != nil {
		logger.Fatalf("Couldn't create metrics collector: %s", err.Error())
	}

	if err := registry.Register(sc); err != nil {
		return nil, fmt.Errorf("couldn't register metrics SVC collector: %s", err.Error())
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
	sc, err := settingsCollector.NewSVCCollector(targets, authTokenCaches, authTokenMutexes, colCounters) //new a Spectrum Virtualize Collector
	// registry.MustRegister(version.NewCollector("Spectrum-Virtualize-Exporter"))

	if err != nil {
		logger.Fatalf("Couldn't create setting collector: %s", err.Error())
	}

	if err := registry.Register(sc); err != nil {
		return nil, fmt.Errorf("couldn't register settings SVC collector: %s", err.Error())
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
