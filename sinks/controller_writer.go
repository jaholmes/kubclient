package sinks

import (
	appd "appdynamics"
	"log"
	"io"
	"os"
)

type ControllerClient struct {
	logger     *log.Logger
	regMetrics map[string]bool
}

func NewControllerClient(host, accessKey, account, app, tier, node string, port uint16, useSSL bool, logger *log.Logger) *ControllerClient {
	cfg := appd.Config{}

	cfg.AppName = app
	cfg.TierName = tier
	cfg.NodeName = node
	cfg.Controller.Host = host
	cfg.Controller.Port = port
	cfg.Controller.UseSSL = useSSL
	if useSSL {
		err := writeSSLFromEnv(logger)
		if err != nil {
			logger.Printf("Unable to set SSL certificates. Using system certs %s", SystemSSLCert)
			cfg.Controller.CertificateFile = SystemSSLCert

		} else {
			logger.Printf("Setting agent certs to %s", AgentSSLCert)
			cfg.Controller.CertificateFile = AgentSSLCert
		}
	}
	cfg.Controller.Account = account
	cfg.Controller.AccessKey = accessKey
	cfg.UseConfigFromEnv = true
	cfg.InitTimeoutMs = 1000
	appd.InitSDK(&cfg)
	logger.Println(&cfg.Controller)
	return &ControllerClient{logger: logger, regMetrics: make(map[string]bool)}
}

func writeSSLFromEnv(logger *log.Logger) error {
	from, err := os.Open(SystemSSLCert)
	if err != nil {
	  logger.Println(err)
	  return err
	}
	defer from.Close()
  
	logger.Printf("Copying system certificates from %s to %s \n", SystemSSLCert, AgentSSLCert)
	to, err := os.OpenFile(AgentSSLCert, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Println(err)
		return err
	}
	defer to.Close()
  
	_, err = io.Copy(to, from)
	if err != nil {
		logger.Println(err)
		return err
	}

	logger.Println("Writing Trusted certificates to %s", AgentSSLCert)
	trustedCerts  := os.Getenv("APPD_TRUSTED_CERTS")
	if trustedCerts != "NOTSET" {
		if _, err = to.Write([]byte(trustedCerts)); err != nil {
			logger.Println(err)
			return err
		}
	} else {
		logger.Println("Trusted certificates not found, skipping")
	}

	return nil
}

func (c *ControllerClient) PostBatch(events []interface{}) error {
	bt := appd.StartBT("PostBatch", "")
	for _, event := range events {
		if event != nil {
			dataPoint, ok := event.(*DataPoint)
			if !ok {
				continue
			}

			if dataPoint.Allowed {
				_, pres := c.regMetrics[dataPoint.Metric]
				if !pres {
					if dataPoint.MetricType == CounterEvent {
						c.logger.Printf("Registering Counter Metric: %v", dataPoint.Metric)
						appd.AddCustomMetric("", dataPoint.Metric,
											appd.APPD_TIMEROLLUP_TYPE_SUM,
											appd.APPD_CLUSTERROLLUP_TYPE_INDIVIDUAL,
											appd.APPD_HOLEHANDLING_TYPE_REGULAR_COUNTER)
					} else {
						c.logger.Printf("Registering Value Metric: %v", dataPoint.Metric)
						appd.AddCustomMetric("", dataPoint.Metric,
											appd.APPD_TIMEROLLUP_TYPE_AVERAGE,
											appd.APPD_CLUSTERROLLUP_TYPE_INDIVIDUAL,
											appd.APPD_HOLEHANDLING_TYPE_REGULAR_COUNTER)
					}
					c.regMetrics[dataPoint.Metric] = true
				}
				appd.ReportCustomMetric("", dataPoint.Metric, dataPoint.Value)
			}
		}
	}
	appd.EndBT(bt)

	return nil
}
