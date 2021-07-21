package deployer

import (
	"CloudScapes/pkg/shared"
	"context"

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

func (d *Deployer) ApplySpec(ctx context.Context, spec *shared.K8sDeployment) error {

	if err := applyNamespace(ctx, d.k8sClient, spec); err != nil {
		return err
	}

	return nil
}
