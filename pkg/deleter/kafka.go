package deleter

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/vdp/kafka-metadata-cleaner/config"
	"gitlab.eng.vmware.com/vdp/kafka-metadata-cleaner/pkg/common"
)

type KafkaDeleter struct {
	BootstrapServer string
	KafkaConfig     *sarama.Config
	KafkaClient     sarama.ClusterAdmin
}

var _ Deleter = &KafkaDeleter{}

func NewKafkaDeleter(config config.KafkaConfig) (Deleter, error) {

	KafkaConfig := sarama.NewConfig()
	KafkaConfig.Version = sarama.V0_10_2_0
	if config.Tls.Enabled == true {
		tlsConfig, err := common.NewTLSConfig(config.Tls.ClientCert,
			config.Tls.ClientKey,
			config.Tls.ClusterCert)
		if err != nil {
			logrus.Error("Error while creating tls config: ", err)
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
		logrus.Error("Error while creating cluster admin: ", err.Error())
		return nil, err
	}

	return &KafkaDeleter{
		BootstrapServer: config.BootstrapServer,
		KafkaConfig:     KafkaConfig,
		KafkaClient:     admin,
	}, nil

}

func (deleter *KafkaDeleter) DeleteTopics(topics []string) error {
	for _, topicName := range topics {
		if err := deleter.KafkaClient.DeleteTopic(topicName); err != nil {
			logrus.Error("Error deleting topic " + topicName + " " + err.Error())
			return err
		} else {
			logrus.Info("Topic deleted: " + topicName)
		}
	}
	return nil
}
