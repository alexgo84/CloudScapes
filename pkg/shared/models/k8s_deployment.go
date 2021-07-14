package shared

import "CloudScapes/pkg/wire"

// K8sDeployment is a wire.Deployment compiled with the plan that it held a reference to in planid.
// This means that any null field value in wire.Deployment was replaced with the corresponding value
// from that deployments plan.
// This is the object that will be communicated into the CloudScapes agent.
type K8sDeployment struct {

	// Name is the name of the deployment. It will also be used as namespace and more unique configurations
	Name string `json:"name" db:"name"`

	// Replicas is the amount of pods that will be created for this deployment
	Replicas int64 `json:"replicas"`

	// ClusterName is the name of the cluster this deployment will be created on. Note that
	// the CloudScapes agent on each cluster will listen on the redis PubSub channel corresponding to
	// the cluster name. When deploying, the request will be trasmitted on a channel according to the
	// deployments cluster name.
	ClusterName int64 `json:"clusterID"`

	// ImagePath is the URL for the image that should be deployed
	ImagePath string `json:"imagePath"`
	// ImageSHA will be prepended to the image path `<imagePath>@<imageSHA>`
	ImageSHA string `json:"imageSHA"`
	// ExcludeFromUpdates means that the image will not be updated during cluster wider eleases.
	// Set this to false if you would like to retain a custom image for some deployment.
	ExcludeFromUpdates bool `json:"excludeFromUpdates"`
	// CPULimit indicates the maximum CPU that would be avaialble for each pod. Example: '100m'.
	CPULimit string `json:"CPULimit"`
	// MemLimit indicates the maximum memory that would be avaialble for each pod. Example: '100m'.
	MemLimit string `json:"memLimit"`
	// CPUReq indicates the size of CPU requests the pod will make. The pod will consume at least this amount. Example: '100m'.
	CPUReq string `json:"CPUReq"`
	// CPUMem indicates the size of memory requests the pod will make. The pod will consume at least this amount. Example: '100m'.
	MemReq string `json:"memReq"`
	// DatabaseServiceName identifies the server name this deployments should use.
	DatabaseServiceName string `json:"databaseServiceName"`
	// DatabaseServiceCloud identifies a database cloud with the service provider.
	DatabaseServiceCloud string `json:"databaseServiceCloud"`
	// DatabaseServicePlan identifies a database plan with the service provider.
	DatabaseServicePlan string `json:"databaseServicePlan"`

	// EnvVars identifies the environment variables that will be set on the deployment.
	// this is generated by merging the environment variables from the deployment with those
	// of the refernced plan. When both plan and deployment declare the same key, the value
	// from the deployment will be used.
	EnvVars map[string]interface{} `json:"envVars"`
	// CronJobs identifies the CronJobs that will be scheduled on the deployment namespace.
	// this is generated by merging CronJobs from the deployment with those
	// of the refernced plan. When both plan and deployment declare the same CronJob, only the one
	// from the deployment will be used.
	CronJobs []wire.CronJob `json:"cronJobs"`
	// ConfigMaps identifies the ConfigMaps to be applied on the deployment namespace.
	// this is generated by merging ConfigMaps from the deployment with those
	// of the refernced plan. When both plan and deployment declare the same ConfigMap name, only the one
	// from the deployment will be used.
	ConfigMaps []wire.ConfigMap `json:"configMaps"`
}

func NewK8sDeployment(deploy wire.NewDeployment, plan wire.NewPlan) *K8sDeployment {
	k8sDeploy := K8sDeployment{
		Name:     deploy.Name,
		Replicas: plan.Replicas,
		CPULimit: plan.CPULimit,
		MemLimit: plan.MemLimit,
		CPUReq:   plan.CPUReq,
		MemReq:   plan.MemReq,

		DatabaseServiceName:  plan.DatabaseServiceName,
		DatabaseServiceCloud: plan.DatabaseServiceCloud,
		DatabaseServicePlan:  plan.DatabaseServicePlan,

		EnvVars: plan.EnvVars,
	}

	if deploy.Replicas != nil {
		k8sDeploy.Replicas = *deploy.Replicas
	}
	if deploy.CPULimit != nil {
		k8sDeploy.CPULimit = *deploy.CPULimit
	}
	if deploy.MemLimit != nil {
		k8sDeploy.MemLimit = *deploy.MemLimit
	}
	if deploy.CPUReq != nil {
		k8sDeploy.CPUReq = *deploy.CPUReq
	}
	if deploy.MemReq != nil {
		k8sDeploy.MemReq = *deploy.MemReq
	}
	if deploy.DatabaseServiceName != nil {
		k8sDeploy.DatabaseServiceName = *deploy.DatabaseServiceName
	}
	if deploy.DatabaseServiceCloud != nil {
		k8sDeploy.DatabaseServiceCloud = *deploy.DatabaseServiceCloud
	}
	if deploy.DatabaseServicePlan != nil {
		k8sDeploy.DatabaseServicePlan = *deploy.DatabaseServicePlan
	}

	// override environement variables
	for k, v := range deploy.EnvVars {
		k8sDeploy.EnvVars[k] = v
	}

	cronJobs := []wire.CronJob{}
cronJobLoop:
	for i := range plan.CronJobs {
		for j := range k8sDeploy.CronJobs {
			if plan.CronJobs[i].Name == k8sDeploy.CronJobs[j].Name {
				cronJobs = append(cronJobs, k8sDeploy.CronJobs[j])
				continue cronJobLoop
			}
		}
		// if an override was not found
		cronJobs = append(cronJobs, plan.CronJobs[i])
	}
	k8sDeploy.CronJobs = cronJobs

	configMaps := []wire.ConfigMap{}
configMapLoop:
	for i := range plan.ConfigMaps {
		for j := range k8sDeploy.ConfigMaps {
			if plan.ConfigMaps[i].Name == k8sDeploy.ConfigMaps[j].Name {
				configMaps = append(configMaps, k8sDeploy.ConfigMaps[j])
				continue configMapLoop
			}
		}
		// if an override was not found
		configMaps = append(configMaps, plan.ConfigMaps[i])
	}
	k8sDeploy.ConfigMaps = configMaps

	return &k8sDeploy
}
