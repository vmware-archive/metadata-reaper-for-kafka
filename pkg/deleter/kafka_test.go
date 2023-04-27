package deleter

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
	"gitlab.eng.vmware.com/vdp/kafka-metadata-cleaner/config"
)

type testConfig struct {
	kafkaConfig           config.KafkaConfig
	expectConnectedErr    bool
	expectErr             bool
	expectedKafkaResponse sarama.KError
	topics                []string
}

func TestProducerSuccess(t *testing.T) {
	assert := assert.New(t)
	seedBroker := sarama.NewMockBroker(t, 1)
	defer seedBroker.Close()
	configs := []testConfig{
		{
			kafkaConfig: config.KafkaConfig{
				BootstrapServer: "localhost:9093",
			},
			expectConnectedErr: true,
			topics:             []string{"topic1", "topic2"},
		},
		{
			kafkaConfig: config.KafkaConfig{
				BootstrapServer: seedBroker.Addr(),
			},
			expectConnectedErr: false,
			topics:             []string{"topic1", "topic2"},
		},
		{
			kafkaConfig: config.KafkaConfig{
				BootstrapServer: seedBroker.Addr(),
			},
			expectConnectedErr: false,
			expectErr:          true,
			topics:             []string{"topic1", "topic2"},
		},
	}

	for _, config := range configs {
		if config.expectErr {
			seedBroker.SetHandlerByMap(map[string]sarama.MockResponse{
				"MetadataRequest": sarama.NewMockMetadataResponse(t).
					SetController(seedBroker.BrokerID()).
					SetBroker(seedBroker.Addr(), seedBroker.BrokerID()),
				"DeleteTopicsRequest": sarama.NewMockDeleteTopicsResponse(t).SetError(sarama.ErrTopicDeletionDisabled),
			})
		} else {
			seedBroker.SetHandlerByMap(map[string]sarama.MockResponse{
				"MetadataRequest": sarama.NewMockMetadataResponse(t).
					SetController(seedBroker.BrokerID()).
					SetBroker(seedBroker.Addr(), seedBroker.BrokerID()),
				"DeleteTopicsRequest": sarama.NewMockDeleteTopicsResponse(t),
			})
		}
		admin, err := NewKafkaDeleter(config.kafkaConfig)
		if config.expectConnectedErr {
			assert.Error(err)
		} else {
			err = admin.DeleteTopics(config.topics)
			if config.expectErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		}
	}
}
