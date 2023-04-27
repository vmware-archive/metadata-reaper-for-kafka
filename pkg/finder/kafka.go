package finder

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/vdp/kafka-metadata-cleaner/config"
	"gitlab.eng.vmware.com/vdp/kafka-metadata-cleaner/pkg/common"
)

type KafkaLister struct {
	BootstrapServer string
	KafkaConfig     *sarama.Config
	KafkaClient     sarama.ClusterAdmin
}

var _ Finder = &KafkaLister{}

func NewKafkaLister(config config.KafkaConfig) (Finder, error) {

	KafkaConfig := sarama.NewConfig()
	KafkaConfig.Version = sarama.V2_4_0_0
	if config.Tls.Enabled == true {
		tlsConfig, err := common.NewTLSConfig(config.Tls.ClientCert,
			config.Tls.ClientKey,
			config.Tls.ClusterCert)
		if err != nil {
			logrus.Fatal("Error while creating tls config: ", err)
		}
		// This can be used on test server if domain does not match cert:
		tlsConfig.InsecureSkipVerify = true
		KafkaConfig.Net.TLS.Enable = true
		KafkaConfig.Net.TLS.Config = tlsConfig
	}

	//kafka end point

	//get broker
	admin, err := sarama.NewClusterAdmin([]string{config.BootstrapServer}, KafkaConfig)
	if err != nil {
		logrus.Fatal("Error while creating cluster admin: ", err.Error())
		return nil, err
	}

	return &KafkaLister{
		BootstrapServer: config.BootstrapServer,
		KafkaConfig:     KafkaConfig,
		KafkaClient:     admin,
	}, nil

}

func (lister *KafkaLister) GetTopics() ([]string, error) {
	topics, err := lister.KafkaClient.ListTopics()
	if err != nil {
		logrus.Error("Error listing kafka topics: " + lister.BootstrapServer)
		return nil, err
	}
	var kafkaTopics []string

	for topicName, _ := range topics {
		logrus.Debug("Found topic: " + topicName)
		kafkaTopics = append(kafkaTopics, topicName)
	}
	return kafkaTopics, nil
}
