package cmd

import (
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/vdp/kafka-metadata-cleaner/config"
	"gitlab.eng.vmware.com/vdp/kafka-metadata-cleaner/pkg/deleter"
	"gitlab.eng.vmware.com/vdp/kafka-metadata-cleaner/pkg/filter"
	"gitlab.eng.vmware.com/vdp/kafka-metadata-cleaner/pkg/finder"
)

var topicsCmd = &cobra.Command{
	Use:     "topics",
	Aliases: []string{"topics"},
	Short:   "topics",
	Run:     startTopicDeletion,
}

func init() {
	rootCmd.AddCommand(topicsCmd)
}

func startTopicDeletion(cmd *cobra.Command, args []string) {
	cfg := config.Instance
	var err error
	logrus.Info("Starting lister.....")
	for _, cluster := range cfg.KafkaClusters {
		var lister finder.Finder
		if cluster.ListerMode == "kafka" {
			lister, err = finder.NewKafkaLister(cluster.KafkaConfig)
			if err != nil {
				logrus.Error("Error creating kafka admin: " + cluster.KafkaConfig.BootstrapServer + err.Error())
			}
		} else {
			lister, err = finder.NewK8sLister(cluster.K8sConfig)
			if err != nil {
				logrus.Error("Error creating k8s context: " + err.Error())
			}
		}
		topics, err := lister.GetTopics()
		if err != nil {
			logrus.Error("Error listing kafka topics: " + cluster.KafkaConfig.BootstrapServer)
		}
		kafkaMetadataLister, _ := filter.NewKafkaFilter(cluster.KafkaConfig, topics)
		topicsFiltered, _ := kafkaMetadataLister.FilterTopics()
		if cfg.DryRun {
			logrus.Info("Running in dry-run mode....")
			for _, topicName := range topicsFiltered {
				logrus.Info("Dry-run enabled..... Topic has mark for deletion: " + topicName)
			}
		} else {
			if cluster.ListerMode == "kafka" {
				topicsDeleter, _ := deleter.NewKafkaDeleter(cluster.KafkaConfig)
				topicsDeleter.DeleteTopics(topicsFiltered)
			} else {
				topicsDeleter, _ := deleter.NewK8sDeleter(cluster.K8sConfig)
				topicsDeleter.DeleteTopics(topicsFiltered)
			}
		}
		logrus.Info("Total topics: " + strconv.Itoa(len(topics)) + " empty topics: " + strconv.Itoa(len(topicsFiltered)))
	}
}
