package deployer

import (
	"CloudScapes/pkg/logger"
	"CloudScapes/pkg/shared"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Deployer struct {
	k8sClient *kubernetes.Clientset
	crdClient *rest.RESTClient
	asyncMode bool
}

func NewDeployer(asyncMode bool) (*Deployer, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// Custom REST client to handle customer resource definitions
	// crdConfig := *config
	// crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{
	// 	Group:   "v1",
	// 	Version: "v1",
	// }
	// crdConfig.APIPath = "/apis"
	// // remember to add the CRDs to schema by calling "my_crd.AddToScheme(scheme.Scheme)"" at init
	// crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	// crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	// crdClient, err := rest.RESTClientFor(&crdConfig)

	return &Deployer{
		k8sClient: k8sClient,
		crdClient: nil,
		asyncMode: asyncMode,
	}, nil
}

func (d *Deployer) ApplySpec(k8sDeploy *shared.K8sDeployment) error {
	logger.Log(logger.DEBUG, "DEPLOYING", logger.Any("k8s_deployment", *k8sDeploy))
	return nil
}
