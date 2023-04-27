package deleter

import (
	"context"

	"github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/vdp/kafka-metadata-cleaner/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sDeleter struct {
	k8sConfig *config.K8sConfig
	k8sClient *dynamic.DynamicClient
}

var _ Deleter = &K8sDeleter{}

func NewK8sDeleter(config config.K8sConfig) (Deleter, error) {
	k8sconfig, err := clientcmd.BuildConfigFromFlags("", config.KubeConfigPath)
	if err != nil {
		logrus.Error("Error reading kubeconfigpath: " + config.KubeConfigPath + " " + err.Error())
		return nil, err
	}
	k8sClient, err := dynamic.NewForConfig(k8sconfig)
	if err != nil {
		logrus.Error("Unable to create k8s client: " + err.Error())
		return nil, err
	}
	return &K8sDeleter{
		k8sClient: k8sClient,
		k8sConfig: &config,
	}, nil

}

func (deleter *K8sDeleter) DeleteTopics(topics []string) error {
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	resourceId := schema.GroupVersionResource{
		Group:    "kafka.strimzi.io",
		Version:  "v1beta2",
		Resource: "kafkatopics",
	}
	for _, topicName := range topics {
		logrus.Info("Topic deleted: " + topicName)
		if err := deleter.k8sClient.Resource(resourceId).Namespace(deleter.k8sConfig.Namespace).Delete(context.TODO(), topicName, deleteOptions); err != nil {
			logrus.Error("Error deleting topic " + topicName + " " + err.Error())
		}
	}
	return nil
}
