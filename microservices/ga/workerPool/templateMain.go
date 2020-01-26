{{define "templateMain.go"}}{{.PavedroadInfo}}

// {{.ProjectInfo}}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// Constants to build up a k8s style URL
const (
	// {{.NameExported}}APIVersion Version API URL
	{{.NameExported}}APIVersion string = "/api/{{.APIVersion}}"

	// {{.NameExported}}NamespaceID Prefix for namespaces
	{{.NameExported}}NamespaceID string = "namespace"

	// {{.NameExported}}DefaultNamespace Default namespace
	{{.NameExported}}DefaultNamespace string = "{{.Namespace}}"

	// {{.NameExported}}ResourceType CRD Type per k8s
	{{.NameExported}}ResourceType string = "{{.Name}}"

	// The email or account login used by 3rd party provider
	{{.NameExported}}Key string = "/{key}"

	// {{.NameExported}}LivenessEndPoint
	{{.NameExported}}LivenessEndPoint string = "{{.Liveness}}"

	// {{.NameExported}}ReadinessEndPoint
	{{.NameExported}}ReadinessEndPoint string = "{{.Readiness}}"

	// {{.NameExported}}MetricsEndPoint
	{{.NameExported}}MetricsEndPoint string = "{{.Metrics}}"

	// EventCollectorManagementEndPoint
	EventCollectorManagementEndPoint string = "management"

	// {{.NameExported}}JobsEndPoint
	{{.NameExported}}JobsEndPoint string = "jobs"

	// {{.NameExported}}SchedulerEndPoint
	{{.NameExported}}SchedulerEndPoint string = "scheduler"
)

// {{.NameExported}}App Top level construct containing building blocks
// for this micro service
type {{.NameExported}}App struct {
	// Router http request router, gorilla mux for this app
	Router *mux.Router

	// Dispatcher manages jobs for workers
	Dispatcher dispatcher

	// Scheduler creates and forwards jobs to dispatcher
	Scheduler  Scheduler

	// Live http server is start
	Live bool

	// Ready once dispatcher has complete initialization
	Ready bool

	httpInterruptChan chan os.Signal

	// Logs
	accessLog *os.File
}

// HTTP server configuration
type httpConfig struct {
	ip							string
	port						string
	shutdownTimeout time.Duration
	readTimeout			time.Duration
	writeTimeout		time.Duration
	listenString		string
	logPath					string
	diagnosticsFile string
	accessFile      string
}

// Set default http configuration
var httpconf = httpConfig{ip: "127.0.0.1", port: "8081", shutdownTimeout: 15, readTimeout: 60, writeTimeout: 60, listenString: "127.0.0.1:8081", logPath: "logs/", diagnosticsFile: "diagnostics.log", accessFile: "access.log"}

// shutdownTimeout will be initialized based on the default or HTTP_SHUTDOWN_TIMEOUT
var shutdowTimeout time.Duration

// GitTag contains current git tab for this repository
var GitTag string

// Version contains version specified in definitions file
var Version string

// Build holds latest git commit hash in short form
var Build string

// printVersion
func printVersion() {
	fmt.Printf("{\"Version\": \"%v\", \"Build\": \"%v\", \"GitTag\": \"%v\"}\n",
		Version, Build, GitTag)
	os.Exit(0)
}

// main entry point for server
func main() {
	a := {{.NameExported}}App{}

	versionFlag := flag.Bool("v", false, "Print version information")
	flag.Parse()

	if *versionFlag {
		printVersion()
	}

	// Setup logging
	openErrorLogFile(httpconf.logPath + httpconf.diagnosticsFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	a.accessLog = openAccessLogFile(httpconf.logPath + httpconf.accessFile)

	a.Initialize()
	a.Run(httpconf.listenString)
}{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}