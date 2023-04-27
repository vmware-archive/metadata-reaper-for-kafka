package finder

import (
	"context"

	"github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/vdp/kafka-metadata-cleaner/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sLister struct {
	k8sConfig *config.K8sConfig
	k8sClient *dynamic.DynamicClient
}

var _ Finder = &K8sLister{}

func NewK8sLister(config config.K8sConfig) (Finder, error) {
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
	return &K8sLister{
		k8sClient: k8sClient,
		k8sConfig: &config,
	}, nil

}

func (lister *K8sLister) GetTopics() ([]string, error) {

	resourceId := schema.GroupVersionResource{
		Group:    "kafka.strimzi.io",
		Version:  "v1beta2",
		Resource: "kafkatopics",
	}
	list, err := lister.k8sClient.Resource(resourceId).Namespace(lister.k8sConfig.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var kafkaTopics []string
	for _, d := range list.Items {
		topicName, found, err := unstructured.NestedString(d.Object, "spec", "topicName")
		if err != nil || !found {
			logrus.Error("Topic name not found for topic: " + d.GetName() + ": " + err.Error())
			continue
		}
		logrus.Debug("Found topic: " + topicName)
		kafkaTopics = append(kafkaTopics, topicName)
	}
	return kafkaTopics, nil
}
