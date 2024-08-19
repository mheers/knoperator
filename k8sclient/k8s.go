package k8sclient

import (
	"flag"
	"path/filepath"
	"sync"

	"github.com/mheers/knoperator/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	K8sClient *kubernetes.Clientset
	once      sync.Once
)

// Init initializes a message queue client
func Init(appConfig *config.Config) (*kubernetes.Clientset, error) {
	var err error
	once.Do(func() {

		clusterConfig, err := getClusterConfig(appConfig)
		if err != nil {
			panic(err.Error())
		}

		// create the clientset
		clientset, err := kubernetes.NewForConfig(clusterConfig)
		if err != nil {
			panic(err.Error())
		}
		K8sClient = clientset
	})
	return K8sClient, err
}

func getClusterConfig(appConfig *config.Config) (*rest.Config, error) {
	if appConfig.K8sInCluster {
		return getInClusterConfig()
	}
	return getOutOfClusterConfig()
}

func getInClusterConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return config, err
}

func getOutOfClusterConfig() (*rest.Config, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}
