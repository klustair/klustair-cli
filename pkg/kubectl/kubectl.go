package kubectl

import (
	"context"
	"fmt"

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
)

type Client struct {
	clientset *kubernetes.Clientset
	inCluster bool
}

func (c *Client) Init(inCluster bool) {
	c.inCluster = inCluster
	c.clientset = GetClientset(inCluster)
}

func NewKubectlClient(inCluster bool) *Client {
	c := new(Client)
	c.Init(inCluster)
	return c
}

func GetClientset(inCluster bool) *kubernetes.Clientset {

	if inCluster {

		// creates the in-cluster config
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		// creates the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		// TODO remove me
		fmt.Printf("clientset: %+v\n", clientset)
		return clientset
	} else {
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
		return clientset
	}
}

func (c *Client) GetNamespaces() (*v1.NamespaceList, error) {
	return c.clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
}

func (c *Client) GetPods(namespace string) (*v1.PodList, error) {
	return c.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
}
