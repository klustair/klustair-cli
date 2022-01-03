package kubectl

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	// all available auth providers for go-client
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	_ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
	"k8s.io/client-go/rest"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	clientset *kubernetes.Clientset
	inCluster bool
}

func (c *Client) Init() {
	c.clientset, c.inCluster = GetClientset()
}

func NewKubectlClient() *Client {
	c := new(Client)
	c.Init()
	return c
}

func GetClientset() (*kubernetes.Clientset, bool) {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err == nil {
		log.Debug("kubectl: in-cluster config")
		// creates the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

		return clientset, true
	} else {
		log.Debug("kubectl: local config")
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()

		configOverrides := &clientcmd.ConfigOverrides{}

		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

		config, err := kubeConfig.ClientConfig()
		if err != nil {
			panic(err.Error())
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		return clientset, false
	}
}

func (c *Client) GetNamespaces() (*v1.NamespaceList, error) {
	return c.clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
}

func (c *Client) GetPods(namespace string) (*v1.PodList, error) {
	return c.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
}
