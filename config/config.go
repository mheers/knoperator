package config

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"time"

	"dario.cat/mergo"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/mheers/knoperator/helpers"
	"github.com/sethvargo/go-envconfig"
	"github.com/sirupsen/logrus"
)

// Config describes the config
type Config struct {
	MQJWT       string `env:"KNOPERATOR_MQ_JWT"`
	MQURI       string `env:"KNOPERATOR_MQ_URI"`
	MQUSeed     string `env:"KNOPERATOR_MQ_USEED"`
	MQCredsPath string `env:"KNOPERATOR_MQ_CREDS_PATH"`

	K8sInCluster              bool   `env:"KNOPERATOR_K8S_INCLUSTER"`
	K8sNamespace              string `env:"KNOPERATOR_K8S_NAMESPACE"`
	K8sPodName                string `env:"KNOPERATOR_K8S_POD_NAME"`
	K8sDefaultImagePullPolicy string `env:"KNOPERATOR_K8S_DEFAULT_IMAGE_PULL_POLICY"`

	BaseHostPath string `env:"KNOPERATOR_BASE_HOSTPATH"`

	DataDir string `env:"KNOPERATOR_DATA_DIR"`
}

// OverlayConfigWithEnv overlays the config with values from the env
func (cfg *Config) OverlayConfigWithEnv(print bool) error {
	ctx := context.Background()
	overlayCfg := &Config{}
	err := envconfig.Process(ctx, overlayCfg)
	if err != nil {
		return err
	}

	if print {
		overlayCfg.Print()
	}

	err = mergo.Merge(cfg, overlayCfg, mergo.WithOverride)
	if err != nil {
		return err
	}
	return nil
}

// Print prints the config if log-level is debug
func (cfg *Config) Print() {
	s := reflect.ValueOf(cfg).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if typeOfT.Field(i).Name != "FeatureTelemetryEnabled" && typeOfT.Field(i).Name != "FeatureTelemetryURL" {
			fmt.Printf("Config: %s %s = %v\n",
				typeOfT.Field(i).Name, f.Type(), f.Interface())
		}
	}
}

// GetFakeConfig creates a config for testing purposes only
func GetFakeConfig() *Config {
	gofakeit.Seed(time.Now().UTC().UnixNano())

	cfg := &Config{
		MQJWT:        "eyJ0eXAiOiJKV1QiLCJhbGciOiJlZDI1NTE5LW5rZXkifQ.eyJqdGkiOiI2MlRRVlU0WkQ1RlVDV0NaS09ZWkg3TU1YM1Q3SzROVk80VkJEQVE2UTc0S0JZN01IRVpRIiwiaWF0IjoxNjUxMzE1Mzc4LCJpc3MiOiJBQjVVVU1XRTdMQkVQVjNSSjVUQVRIT1Y1RjJXT0xRREtZQ0VTVFBNQ1VISVAyM0VQQVpNMlBDRSIsIm5hbWUiOiJ0ZXN0Iiwic3ViIjoiVUJJSkhONlpDT0g0M0pCVFJLVzQyRk5aMjJZQ1JBM09DSFM3SUg3QVBaNjJLMzVVN1E2VzNYSVIiLCJuYXRzIjp7InB1YiI6eyJhbGxvdyI6WyJmbnhwLioiLCJfSU5CT1guXHUwMDNlIl19LCJzdWIiOnsiYWxsb3ciOlsiZm54cC4qIiwiX0lOQk9YLlx1MDAzZSJdfSwic3VicyI6LTEsImRhdGEiOi0xLCJwYXlsb2FkIjotMSwiYmVhcmVyX3Rva2VuIjp0cnVlLCJ0eXBlIjoidXNlciIsInZlcnNpb24iOjJ9fQ.t4nO6cJuTumTTF0mIzw64iYnTBjR_2DGPcFS-hYq2dQn5KS1Tuk5cfpJsqPXreuWCuQJjnM3-QzGrLLyhqeHDA",
		MQUSeed:      "SUAIX7MNFMX7G7LB2P3EB53HDCAEMNRM4HLPGMJ7U4OZU55VWDXTT655EU",
		MQURI:        "ws://localhost:9222",
		K8sInCluster: false,
		DataDir:      "/tmp/data",
	}

	err := cfg.OverlayConfigWithEnv(false)
	if err != nil {
		return nil
	}

	logLevel := os.Getenv("LOGLEVEL")
	if logLevel != "" {
		helpers.SetLogLevel(logLevel)
	}

	return cfg
}

var configInstance *Config

func GetConfig(print bool) *Config {
	if configInstance == nil {
		configInstance = &Config{}
		err := configInstance.OverlayConfigWithEnv(false)
		if err != nil {
			logrus.Fatalf("config: %s", err)
		}
		if print {
			configInstance.Print()
		}
	}
	return configInstance
}
