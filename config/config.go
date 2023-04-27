package config

import (
	"github.com/sirupsen/logrus"
)

type ConfigManager interface {
	LogConfig()
}

type PrometheusConfig struct {
	PrometheusServer string
	PrometheusPort   string
}

type KafkaClusters struct {
	KafkaConfig KafkaConfig
	ListerMode  string
	DeleterMode string
	K8sConfig   K8sConfig
}

type KafkaConfig struct {
	BootstrapServer      string
	Tls                  TlsConfig
	SchemaRegistryConfig SchemaRegistryConfig
}

type SchemaRegistryConfig struct {
	SchemaRegistryUrl string
	Enable            bool
}
type TlsConfig struct {
	Enabled     bool
	ClusterCert string
	ClientCert  string
	ClientKey   string
}
type K8sConfig struct {
	KubeConfigPath string
	Namespace      string
	ClusterName    string
}

type Config struct {
	Prometheus    PrometheusConfig
	KafkaClusters []KafkaClusters
	DryRun        bool
}

var Instance *Config

func (config *Config) LogConfig() {
	logrus.Info("Configuration loaded")
	logrus.Info(Instance)
}
