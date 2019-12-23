package k8s

import (
	"os"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Clientset todo
type Clientset interface {
	GetConfig() *rest.Config
	GetClient() (*kubernetes.Clientset, error)
	GetExtClient() (*apiextensionsclient.Clientset, error)
}

// ClientsetImpl todo
type ClientsetImpl struct {
	Clientset
	Config *rest.Config
}

// GetConfig todo
func (s *ClientsetImpl) GetConfig() *rest.Config {
	return s.Config
}

// GetClient todo
func (s *ClientsetImpl) GetClient() (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(s.Config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// GetExtClient todo
func (s *ClientsetImpl) GetExtClient() (*apiextensionsclient.Clientset, error) {
	clientset, err := apiextensionsclient.NewForConfig(s.Config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// GetK8sClient todo
func GetK8sClient() (Clientset, error) {
	var config *rest.Config
	var err error

	// get the config, either from kubeconfig or using our
	// in-cluster service account
	kubeConfig := os.Getenv("KUBECONFIG")
	if len(kubeConfig) != 0 {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			return nil, err
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	return &ClientsetImpl{
		Config: config,
	}, nil
}
