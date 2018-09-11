package appdconfig

import (
	"github.com/kelseyhightower/envconfig"
	"os"
)

type AppDConfig struct {
	Sink           string `default:"Controller" envconfig:"sink"`
	SamplingRate   uint16 `default:"2" envconfig:"sampling_rate"`
	ControllerHost string `envconfig:"controller_host"`
	ControllerPort uint16 `default:"8090" envconfig:"controller_port"`
	AccessKey      string `envconfig:"access_key"`
	Account        string `envconfig:"account"`
	SslEnabled     bool   `default:"false" envconfig:"ssl_enabled"`
	NozzleAppName  string `ignored:"true"`
	NozzleTierName string `ignored:"true"`
	NozzleTierId   string `envconfig:"nozzle_tier_id"`
	NozzleNodeName string `ignored:"true"`
}

func Parse() (*AppDConfig, error) {
	config := &AppDConfig{}
	err := envconfig.Process("appd", config)
	if err != nil {
		return nil, err
	}
	config.NozzleAppName, err = GetEnvWithDefault("APPD_NOZZLE_APP", "appd-nozzle")
	if err != nil {
		return nil, err
	}
	config.NozzleTierName, err = GetEnvWithDefault("APPD_NOZZLE_TIER", "appd-nozzle-tier")
	if err != nil {
		return nil, err
	}
	config.NozzleNodeName, err = GetEnvWithDefault("APPD_NOZZLE_NODE", "appd-nozzle-node")
	if err != nil {
		return nil, err
	}
	return config, nil
}

func GetEnvWithDefault(envVariable, defaultValue string) (string, error) {
	envValue := os.Getenv(envVariable)
	if envValue == "" {
		return defaultValue, nil
	}
	return envValue, nil
}
