package sinks

// Allowed sinks
const Stdout string = "stdout"
const MachineAgent string = "machineagent"
const Controller string = "controller"

// SSL Certificates Location
const SystemSSLCert string = "/etc/ssl/certs/ca-certificates.crt"
const AgentSSLCert string = "/tmp/ca-certificates.crt"

// Metric Types

const ValueMetricEvent string = "ValueMetricEvent"
const CounterEvent string = "CounterEvent"
const HttpStartStopEvent string = "HttpStartStopMetricEvent"
const LogMessageEvent string = "LogMessageEvent"
const ErrorEvent string = "ErrorEvent"
const ContainerEvent string = "ContainerEvent"


// Metrics we are interested in ...  {<metric origin>: [<name of the metric>]}
var MetricFilter = map[string][]string{
	"gorouter":      []string{"file_descriptors", "backend_exhausted_conns", "latency", "ms_since_last_registry_update", "bad_gateways", "responses.5xx", "registry_message.route-emitter", "total_requests", "total_routes", "latency.uaa"},
	"mysql":         []string{"/mysql/available", "/mysql/galera/wsrep_ready", "/mysql/galera/wsrep_cluster_size", "/mysql/galera/wsrep_cluster_status", "/mysql/net/connections", "/mysql/performance/questions", "/mysql/performance/busy_time"},
	"route_emitter": []string{"RouteEmitterSyncDuration"},
	"locket":        []string{"ActiveLocks", "ActivePresences"},
	"bbs":           []string{"ConvergenceLRPDuration", "RequestLatency", "Domain.cf-apps", "LRPsExtra", "LRPsMissing", "CrashedActualLRPs", "LRPsRunning", "LockHeld"},
	"auctioneer":    []string{"AuctioneerLRPAuctionsFailed", "AuctioneerFetchStatesDuration", "AuctioneerLRPAuctionsStarted", "LockHeld", "AuctioneerTaskAuctionsFailed"},
	"rep":           []string{"CapacityRemainingMemory", "CapacityRemainingDisk", "RepBulkSyncDuration", "UnhealthyCell", "CapacityTotalMemory", "CapacityRemainingContainers", "CapacityTotalContainers", "CapacityTotalDisk"},
	"uaa":           []string{"requests.global.completed.count", "server.inflight.count"},
	"bosh-system-metrics-forwarder": []string{"system.healthy", "system.mem.percent", "system.disk.system.percent", "system.disk.ephemeral.percent", "system.disk.persistent.percent", "system.cpu.user"},
	"loggregator.doppler":           []string{"dropped", "ingress"},
	"loggregator.rlp":               []string{"dropped", "ingress"},
	"cf-syslog-drain.scheduler":     []string{"drains"},
}

// Alias for Origins - used in the metric path Application Infrastructure Performance|...|PCF Firehose Monitor|<alias>
var MetricAlias = map[string]string{
	"auctioneer":    "Diego Auctioneer Metrics",
	"bbs":           "Diego BBS Metrics",
	"rep":           "Diego Cell Metrics",
	"locket":        "Diego Locket Metrics",
	"route_emitter": "Diego Route Emitter Metrics",
	"mysql":         "PAS MySQL KPIs",
	"gorouter":      "Gorouter Metrics",
	"uaa":           "UAA Metrics",
	"bosh-system-metrics-forwarder": "System (BOSH) Metrics",
	"loggregator.doppler":           "Loggregator Doppler Metrics",
	"cf-syslog-drain.adapter":       "CF Syslog Drain Metrics",
	"loggregator.rlp":               "Loggregator RLP Metrics",
	"cf-syslog-drain.scheduler":     "CF Syslog Drain Bindings",
}
