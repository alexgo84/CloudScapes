package deployer

import (
	"CloudScapes/pkg/shared"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func buildLabels(k8sDeploy *shared.K8sDeployment) map[string]string {
	return map[string]string{
		"company": k8sDeploy.Name,
	}
}

func updateMeta(src, dst *metaV1.ObjectMeta) {
	for k, v := range src.Annotations {
		dst.Annotations[k] = v
	}
	for k, v := range src.Labels {
		dst.Labels[k] = v
	}
}
