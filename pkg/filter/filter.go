package filter

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/sirupsen/logrus"

	"github.com/Shopify/sarama"
	"gitlab.eng.vmware.com/vdp/kafka-metadata-cleaner/config"
	"gitlab.eng.vmware.com/vdp/kafka-metadata-cleaner/pkg/common"
)

var SYSTEM_TOPICS = []string{
	"__consumer_offsets",
	"^strimzi",
	"^__strimzi",
	"schemas",
}

type Filter interface {
	FilterTopics() ([]string, error)
	filterEmptyTopics() error
	filterSystemTopics() error
	filterSchemaTopics() error
}

type KafkaFilter struct {
	kafkaTopicsFiltered map[string]bool
	BootstrapServer     string
	KafkaConfig         config.KafkaConfig
	KafkaClient         sarama.ClusterAdmin
}

var _ Filter = &KafkaFilter{}

func NewKafkaFilter(config config.KafkaConfig, topics []string) (Filter, error) {

	KafkaConfig := sarama.NewConfig()
	KafkaConfig.Version = sarama.V2_4_0_0
	if config.Tls.Enabled == true {
		tlsConfig, err := common.NewTLSConfig(config.Tls.ClientCert,
			config.Tls.ClientKey,
			config.Tls.ClusterCert)
		if err != nil {
			log.Fatal(err)
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
		logrus.Fatal("Error while creating cluster admin for filtering: ", err.Error())
		return nil, err
	}
	kafkaTopicsFiltered := make(map[string]bool, len(topics))

	for _, topicName := range topics {
		kafkaTopicsFiltered[topicName] = true
	}

	return &KafkaFilter{
		BootstrapServer:     config.BootstrapServer,
		KafkaConfig:         config,
		KafkaClient:         admin,
		kafkaTopicsFiltered: kafkaTopicsFiltered,
	}, nil

}

func (filter *KafkaFilter) FilterTopics() ([]string, error) {

	err := filter.filterSystemTopics()
	if err != nil {
		logrus.Error("Unable to filter system topics: ", err.Error())
		return nil, err
	}

	err = filter.filterEmptyTopics()
	if err != nil {
		logrus.Error("Unable to filter empty topics: ", err.Error())
		return nil, err
	}

	if filter.KafkaConfig.SchemaRegistryConfig.Enable {
		err = filter.filterSchemaTopics()
		if err != nil {
			logrus.Error("Unable to filter empty topics: ", err.Error())
			return nil, err
		}
	}

	counter := 0
	var kafkaTopics []string
	for topic, value := range filter.kafkaTopicsFiltered {
		if value {
			counter += 1
			kafkaTopics = append(kafkaTopics, topic)
		}
	}
	return kafkaTopics, nil
}

func (filter *KafkaFilter) filterEmptyTopics() error {
	brokerLogDirs, err := filter.KafkaClient.DescribeLogDirs([]int32{0, 1, 2, 3, 4, 5, 6, 7, 8})
	if err != nil {
		logrus.Error("Error conection to kafka cluster to describe brokers: ", err.Error())
		return err
	}
	for index, logDirs := range brokerLogDirs {
		logrus.Debug("Checking broker number " + string(index))
		for _, logdir := range logDirs {
			logrus.Debug("Checking log dir " + logdir.Path + " for broker " + string(index))
			for _, topic := range logdir.Topics {
				if filter.kafkaTopicsFiltered[topic.Topic] != false {
					for _, partition := range topic.Partitions {
						if partition.Size != 0 || partition.OffsetLag > 0 {
							logrus.Debug("Topic " + topic.Topic + " is not empty filtering out")
							filter.kafkaTopicsFiltered[topic.Topic] = false
							break
						}
					}
				}
			}
		}
	}
	return nil
}

func (filter *KafkaFilter) filterSystemTopics() error {
	for topic, _ := range filter.kafkaTopicsFiltered {
		for _, systemTopic := range SYSTEM_TOPICS {
			match, err := regexp.MatchString(systemTopic, topic)
			if err != nil {
				logrus.Error("Error compiling systemTopic "+systemTopic+" against topic ", topic, err.Error())
				return err
			}
			if match {
				logrus.Debug("Filtering out system topic " + topic)
				filter.kafkaTopicsFiltered[topic] = false
			}

		}
	}
	return nil
}

// filterSchemaTopics talks with schema registry and check if we have and schema created for the topic
func (filter *KafkaFilter) filterSchemaTopics() error {
	for topic, value := range filter.kafkaTopicsFiltered {
		if value {
			logrus.Debug("Checking if topic " + topic + " has an schema")
			schemaRegistryUrl := fmt.Sprintf("https://%s/subjects/%s-value/versions\n", filter.KafkaConfig.SchemaRegistryConfig.SchemaRegistryUrl, topic)
			resp, err := http.Get(schemaRegistryUrl)
			if err != nil {
				defer resp.Body.Close()
				logrus.Errorf("Unable to connect to schema registry: " + " " + err.Error())
				filter.kafkaTopicsFiltered[topic] = false
				return err
			} else if resp.StatusCode == 200 {
				filter.kafkaTopicsFiltered[topic] = false
				logrus.Debug("Filtering out topic " + topic + " with schema")
			}
		}
	}
	return nil
}
