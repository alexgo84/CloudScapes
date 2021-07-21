package deployer

import (
	"CloudScapes/pkg/logger"
	"CloudScapes/pkg/shared"
	"context"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8s_errors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/client-go/kubernetes"
)

func applyNamespace(ctx context.Context, api *kubernetes.Clientset, spec *shared.K8sDeployment) error {
	logger.Log(logger.DEBUG, "DEPLOYING", logger.Any("k8s_deployment", *spec))

	ns := makeNamespace(spec)

	if _, err := api.CoreV1().Namespaces().Create(ctx, &ns, metaV1.CreateOptions{}); err == nil {
		return nil
	} else if err != nil && !k8s_errors.IsAlreadyExists(err) {
		return err
	}

	// namespace already exists. get the existing namespace and update it with labels
	existingNamespace, err := api.CoreV1().Namespaces().Get(ctx, ns.Name, metaV1.GetOptions{})
	if err != nil {
		return err
	}

	updateMeta(&ns.ObjectMeta, &existingNamespace.ObjectMeta)

	if _, err := api.CoreV1().Namespaces().Update(ctx, existingNamespace, metaV1.UpdateOptions{}); err != nil {
		return err
	}
	return nil
}

func makeNamespace(k8sDeploy *shared.K8sDeployment) v1.Namespace {
	return v1.Namespace{ObjectMeta: metaV1.ObjectMeta{Name: k8sDeploy.Name, Labels: buildLabels(k8sDeploy)}}
}
